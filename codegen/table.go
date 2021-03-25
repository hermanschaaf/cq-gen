package codegen

import (
	"fmt"
	"github.com/cloudquery/cloudquery-plugin-sdk/plugin/schema"
	"github.com/cloudquery/cq-gen/codegen/config"
	"github.com/cloudquery/cq-gen/codegen/templates"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
	"github.com/pkg/errors"
	"go/types"
	"strings"
)

type TableDefinition struct {
	TypeName    string
	Name        string
	Description string
	Columns     []ColumnDefinition
	Relations   []*TableDefinition
	Resolver    *ResolverDefinition
}

func (t TableDefinition) UniqueResolvers() []*ResolverDefinition {

	rd := make([]*ResolverDefinition, 0)
	rd = append(rd, t.Resolver)
	existingResolvers := make(map[string]bool)
	for _, relation := range t.Relations {
		for _, ur := range relation.UniqueResolvers() {
			if _, ok := existingResolvers[ur.Name]; ok {
				continue
			}
			rd = append(rd, ur)
			existingResolvers[ur.Name] = true
		}
	}
	return rd
}

type ColumnDefinition struct {
	Name        string
	Type        schema.ValueType
	Elem        []ColumnDefinition
	Description string
	Resolver    *ResolverDefinition
}

type ResolverDefinition struct {
	Name      string
	Signature string
	Body      string
	Type      *types.Func
}

func (b builder) buildTable(resource config.ResourceConfig) (*TableDefinition, error) {
	ro, err := b.finder.FindTypeFromName(resource.Path)
	if err != nil {
		return nil, err
	}

	named := ro.(*types.Named)
	typeName := inflection.Plural(named.Obj().Name())
	resolverName := templates.ToGo(fmt.Sprintf("fetch%s%s", strings.Title(resource.Domain), strings.Title(typeName)))

	table := &TableDefinition{
		TypeName:    strings.ToLower(resource.Domain + strings.Title(typeName)),
		Name:        strings.ToLower(fmt.Sprintf("%s_%s_%s", resource.Service, resource.Domain, strcase.ToSnake(typeName))),
		Description: "",
		Columns:     make([]ColumnDefinition, 0),
		Relations:   make([]*TableDefinition, len(resource.Relations)),
		Resolver:    &ResolverDefinition{Name: resolverName, Body: "panic(\"no implemented\")"},
	}

	body := b.rewriter.GetFunction(resolverName)
	if body != "" {
		table.Resolver.Body = strings.TrimSuffix(strings.TrimPrefix(body, "\n\n"), "\n\n")
	}

	st := named.Underlying().(*types.Struct)

	for i := 0; i < st.NumFields(); i++ {
		field, tag := st.Field(i), st.Tag(i)

		columnName := strings.ToLower(strcase.ToSnake(field.Name()))
		colCfg := resource.GetColumnConfig(columnName)

		// Skip unexported, if the original field has a "-" tag or the field was requested to be skipped via config.
		if !field.Exported() || strings.Contains(tag, "-") || colCfg.Skip {
			continue
		}

		// TODO: make cleaner
		var valueType schema.ValueType
		if colCfg.Type != schema.TypeInvalid {
			valueType = colCfg.Type
		} else {
			valueType = getValueType(field.Type())
			if valueType == schema.TypeInvalid {
				return nil, fmt.Errorf("unsupported type %T", field.Type())
			}
			if valueType == TypeRelation {
				obj := getNamedType(field.Type()).Obj()
				rel, err := b.buildTableRelation(named.Obj().Name(), config.ResourceConfig{
					Service: resource.Service,
					Domain:  resource.Domain,
					Name:    strcase.ToSnake(obj.Name()), // Add prefix?
					Path:    fmt.Sprintf("%s.%s", obj.Pkg().Path(), obj.Name()),
				})
				if err != nil {
					return nil, err
				}
				table.Relations = append(table.Relations, rel)
				continue
			}
		}

		switch valueType {
		case schema.TypeEmbedded:
			columns, err := buildEmbeddedColumns(field.Name(), getNamedType(field.Type()), colCfg)
			if err != nil {
				return nil, err
			}
			table.Columns = append(table.Columns, ColumnDefinition{
				Name: columnName,
				Type: valueType,
				Elem: columns,
			})
		default:
			table.Columns = append(table.Columns, ColumnDefinition{
				Name: columnName,
				Type: valueType,
			})

		}
	}

	for _, uc := range resource.UserDefinedColumn {
		ro, err := b.finder.FindFuncFromName(uc.ResolverPath)
		if err != nil {
			return nil, errors.Wrapf(err, "user defined column %s requires resolver definition", uc.Name)
		}
		table.Columns = append(table.Columns, ColumnDefinition{
			Name:     uc.Name,
			Type:     uc.Type,
			Resolver: &ResolverDefinition{Type: ro},
		})
	}
	for i, relation := range resource.Relations {
		table.Relations[i], err = b.buildTableRelation(named.Obj().Name(), relation)
		if err != nil {
			return nil, err
		}
	}

	return table, nil
}

func (b builder) buildTableRelation(parent string, cfg config.ResourceConfig) (*TableDefinition, error) {

	rel, err := b.buildTable(cfg)
	if err != nil {
		return nil, err
	}
	rel.Columns = append(rel.Columns, ColumnDefinition{
		Name:     strings.ToLower(fmt.Sprintf("%s_id", parent)),
		Type:     schema.TypeUUID,
		Resolver: &ResolverDefinition{Signature: "schema.ParentIdResolver"},
	})

	return rel, nil
}

func getNamedType(typ types.Type) *types.Named {
	switch t := typ.(type) {
	case *types.Pointer:
		return getNamedType(t.Elem())
	case *types.Named:
		return t
	case *types.Slice:
		return getNamedType(t.Elem())
	}
	panic("type ")
}

func buildEmbeddedColumns(fieldName string, named *types.Named, cfg config.ColumnConfig) ([]ColumnDefinition, error) {
	st := named.Underlying().(*types.Struct)
	columns := make([]ColumnDefinition, st.NumFields())

	for i := 0; i < st.NumFields(); i++ {
		field, tag := st.Field(i), st.Tag(i)
		// Skip unexported, if the original field has a "-" tag or the field was requested to be skipped via config.
		if !field.Exported() || strings.Contains(tag, "-") {
			continue
		}
		valueType := getValueType(field.Type())
		if valueType == schema.TypeInvalid {
			return nil, fmt.Errorf("unsupported type %T", field.Type())
		}

		columnName := strings.ToLower(fmt.Sprintf("%s_%s", strcase.ToSnake(fieldName), strcase.ToSnake(field.Name())))
		if cfg.SkipPrefix {
			columnName = strings.ToLower(strcase.ToSnake(field.Name()))
		}

		columns[i] = ColumnDefinition{
			Name:     columnName,
			Type:     valueType,
			Resolver: &ResolverDefinition{Signature: fmt.Sprintf("schema.PathResolver(\"%s\")", fmt.Sprintf("%s.%s", fieldName, field.Name()))},
		}
	}
	return columns, nil
}

const TypeRelation schema.ValueType = -1

func getValueType(typ types.Type) schema.ValueType {
	if vt := getUniqueStructs(typ); vt != schema.TypeInvalid {
		return vt
	}
	switch t := typ.(type) {
	case *types.Map:
		return schema.TypeJSON
	case *types.Basic:
		return getBasicType(t)
	case *types.Named:
		return getValueType(t.Underlying())
	case *types.Struct:
		return schema.TypeEmbedded
	case *types.Slice:
		valueType := getValueType(t.Elem())
		switch valueType {
		case schema.TypeInt:
			return schema.TypeIntArray
		case schema.TypeString:
			return schema.TypeString
		case schema.TypeEmbedded:
			return TypeRelation
		default:
			return schema.TypeInvalid
		}
	case *types.Pointer:
		return getValueType(t.Elem())
	}
	return schema.TypeInvalid
}

func getBasicType(typ *types.Basic) schema.ValueType {
	switch typ.Kind() {
	case types.Bool:
		return schema.TypeBool
	// TODO: more specific conversions for int
	case types.Int, types.Int8, types.Int16, types.Int32, types.Int64, types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64, types.Uintptr:
		return schema.TypeInt
	case types.Float32, types.Float64:
		return schema.TypeFloat
	case types.String:
		return schema.TypeString
	}
	return schema.TypeInvalid
}

func getUniqueStructs(typ types.Type) schema.ValueType {
	switch typ.String() {
	case "time.Time":
		return schema.TypeTimestamp
	case "uuid.UUID":
		return schema.TypeUUID
	default:
		return schema.TypeInvalid
	}
}
