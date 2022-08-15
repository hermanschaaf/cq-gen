package codegen

import (
	"fmt"
	"go/types"
	"os"
	"path"
	"strings"

	"github.com/cloudquery/cq-gen/code"
	"github.com/cloudquery/cq-gen/codegen/config"
	"github.com/cloudquery/cq-gen/codegen/source"
	"github.com/cloudquery/cq-gen/codegen/template"
	"github.com/cloudquery/cq-gen/naming"
	"github.com/cloudquery/cq-gen/rewrite"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"
	"github.com/hashicorp/go-hclog"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
	"github.com/thoas/go-funk"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	defaultImplementation = `panic("not implemented")`
	sdkPath               = "github.com/cloudquery/cq-provider-sdk"
	MaxColumnLength       = 63
)

// BuildMeta is information passed when the TableBuilder is traversing over a source.Object to build its table
type BuildMeta struct {
	// Depth represents the logical depth traversal in the build, the depth is used to avoid infinite recursion
	Depth int
	// ColumnPath is the logical traversal on the column, the path extends if it hits embedded objects.
	ColumnPath string
	// FieldPath is the dot notated path of the resource that is expected, this is used for PathResolvers
	FieldPath string
	// FieldParts is a list of fields we traversed
	FieldParts []string
	// BaseFieldIndex stores the index of the last base relation. It is used in conjunction with
	// FieldParts to create the list of accessors when generating resolvers.
	BaseFieldIndex int
	// fullColumnPath saves full path of column regardless of prefix skips, this allows users to define
	// either only column name or full embedded path
	fullColumnPath string
	// Whether this is a user defined object relation
	userRelation bool
}

func BuildColumnMeta(field source.Object, parentMeta BuildMeta, cfg config.ColumnConfig) BuildMeta {
	meta := BuildMeta{
		Depth:          0,
		BaseFieldIndex: 0,
		ColumnPath:     field.Name(),
		FieldPath:      field.Name(),
		FieldParts:     make([]string, len(parentMeta.FieldParts)),
		fullColumnPath: fmt.Sprintf("%s%s", parentMeta.fullColumnPath, field.Name()),
	}
	if cfg.Rename != "" {
		meta.ColumnPath = cfg.Rename
	}
	if parentMeta.ColumnPath != "" {
		meta.ColumnPath = fmt.Sprintf("%s%s", parentMeta.ColumnPath, meta.ColumnPath)
	}
	if cfg.SkipPrefix {
		meta.ColumnPath = parentMeta.ColumnPath
	}
	if parentMeta.FieldPath != "" {
		meta.FieldPath = fmt.Sprintf("%s.%s", parentMeta.FieldPath, field.Name())
	}
	copy(meta.FieldParts, parentMeta.FieldParts)
	return meta
}

type TableBuilder struct {
	finder             code.Finder
	source             source.DataSource
	rewriter           *rewrite.Rewriter
	log                hclog.Logger
	descriptionSource  source.DescriptionSource
	descriptionParsers []source.DescriptionParser
}

var titler = cases.Title(language.English, cases.NoLower)

func NewTableBuilder(source source.DataSource, descriptionSource source.DescriptionSource, rewriter *rewrite.Rewriter, parsers []source.DescriptionParser) TableBuilder {
	return TableBuilder{
		source:            source,
		descriptionSource: descriptionSource,
		rewriter:          rewriter,
		log: hclog.New(&hclog.LoggerOptions{
			Level:  hclog.Debug,
			Output: os.Stdout,
		}),
		descriptionParsers: parsers,
	}
}

func (tb TableBuilder) BuildTable(parentTable *TableDefinition, resourceCfg *config.ResourceConfig, meta BuildMeta) (*TableDefinition, error) {
	fullName := GetResourceName(parentTable, resourceCfg)

	table := &TableDefinition{
		Name:          fullName,
		FileName:      GetFileName(resourceCfg),
		TableFuncName: template.ToGo(titler.String(fullName)),
		TableName:     GetTableName(parentTable, resourceCfg.Service, resourceCfg.Domain, resourceCfg.Name),
		parentTable:   parentTable,
		Options:       resourceCfg.TableOptions,
		Description:   resourceCfg.Description,
		IgnoreInTests: resourceCfg.IgnoreInTests,
		path:          resourceCfg.Path,
	}

	// will only mark table function as copied
	tb.rewriter.GetFunctionBody(table.TableFuncName, "")
	tb.log.Debug("Building table", "table", table.TableName)
	if err := tb.buildTableFunctions(table, resourceCfg, meta); err != nil {
		return nil, err
	}

	if err := tb.addUserDefinedColumns(table, resourceCfg); err != nil {
		return nil, err
	}

	if resourceCfg.Path == "" {
		return table, nil
	}
	obj, err := tb.source.Find(resourceCfg.Path)
	if err != nil {
		return nil, err
	}

	// base field becomes this table
	if len(meta.FieldParts) == 0 {
		meta.BaseFieldIndex = 0
	} else {
		meta.BaseFieldIndex = len(meta.FieldParts) - 1
	}

	if !resourceCfg.DisableReadDescriptions {
		if table.Description == "" {
			meta.FieldParts = append(meta.FieldParts, resourceCfg.DescriptionPathParts...)
			table.Description = tb.getDescription(obj, resourceCfg.Description, meta)
		}
	}

	if len(meta.FieldParts) == 0 {
		if resourceCfg.DescriptionPathParts != nil {
			meta.FieldParts = append(meta.FieldParts, resourceCfg.DescriptionPathParts...)
		} else {
			meta.FieldParts = append(meta.FieldParts, obj.Name())
		}
	}

	if err := tb.buildColumns(table, obj, resourceCfg, meta); err != nil {
		return nil, err
	}

	if err := tb.buildTableRelations(table, resourceCfg, meta); err != nil {
		return nil, err
	}

	return table, nil
}

func (tb TableBuilder) buildTableFunctions(table *TableDefinition, resource *config.ResourceConfig, meta BuildMeta) error {
	var err error
	hasResolver := resource.Resolver != nil
	forceCustomGeneration := hasResolver && resource.Resolver.Generate
	canGenerateResolverImplementation := table.parentTable != nil && len(meta.FieldParts) > 0
	switch {
	case hasResolver && !forceCustomGeneration:
		table.Resolver, err = tb.buildResolverDefinition(table, resource.Resolver)
	case canGenerateResolverImplementation && !forceCustomGeneration:
		table.Resolver, err = tb.getPathTableResolver(meta)
	default:
		table.Resolver, err = tb.buildResolverDefinition(table, &config.FunctionConfig{
			Name:     strcase.ToLowerCamel(fmt.Sprintf("fetch%s%s", titler.String(resource.Domain), titler.String(table.Name))),
			Body:     defaultImplementation,
			Path:     path.Join(sdkPath, "provider/schema.TableResolver"),
			Generate: true, // Table functions are always generated, setting this to true will cause duplicates
		})
	}
	if err != nil {
		return err
	}

	if resource.IgnoreError != nil {
		table.IgnoreErrorFunc, err = tb.buildResolverDefinition(table, resource.IgnoreError)
		if err != nil {
			return err
		}
	}
	if resource.Multiplex != nil {
		table.MultiplexFunc, err = tb.buildResolverDefinition(table, resource.Multiplex)
		if err != nil {
			return err
		}
	}

	if resource.DeleteFilter != nil {
		table.DeleteFilterFunc, err = tb.buildResolverDefinition(table, resource.DeleteFilter)
		if err != nil {
			return err
		}
	}
	if resource.PostResourceResolver != nil {
		table.PostResourceResolver, err = tb.buildResolverDefinition(table, resource.PostResourceResolver)
		if err != nil {
			return err
		}
	}

	return nil
}

// getPathTableResolver uses PathTableResolver to generate a resolver for relations
func (tb TableBuilder) getPathTableResolver(meta BuildMeta) (*ResolverDefinition, error) {
	start, end := meta.BaseFieldIndex+1, len(meta.FieldParts)
	signatureName := "schema.PathTableResolver"
	fieldPath := strings.Join(meta.FieldParts[start:end], ".")
	return &ResolverDefinition{
		Type:      nil,
		Signature: fmt.Sprintf("%s(\"%s\")", signatureName, fieldPath),
	}, nil
}

// buildColumns iterates over every field in source.Object adding a ColumnDefinition or RelationDefinition based on the type
func (tb TableBuilder) buildColumns(table *TableDefinition, object source.Object, resourceCfg *config.ResourceConfig, meta BuildMeta) error {
	fields := object.Fields()
	for _, f := range fields {
		name := f.Name()
		if !f.Exported() && !resourceCfg.AllowUnexported {
			tb.log.Debug("skipping unexported field", "field", name, "object", object.Name(), "table", table.Name)
			continue
		}
		tb.log.Debug("building column", "field", name, "object", object.Name(), "table", table.Name)
		if err := tb.buildColumn(table, f, resourceCfg, meta); err != nil {
			return err
		}
	}
	return nil
}

func (tb TableBuilder) buildColumn(table *TableDefinition, field source.Object, resourceCfg *config.ResourceConfig, meta BuildMeta) error {
	// Build initial column definition
	fieldName := field.Name()
	colDef := ColumnDefinition{
		Name:     GetColumnName(fieldName, meta),
		Type:     0,
		Resolver: nil,
	}
	tb.log.Debug("finding column configuration", "name", colDef.Name, "resource", resourceCfg.Name, "user_relation", meta.userRelation)
	// check if configuration wants column to be skipped
	cfg := resourceCfg.GetColumnConfig(colDef.Name, meta.fullColumnPath)
	if cfg.Skip {
		return nil
	}
	// Set column description, usually source.Object contains a description, but it can also be overridden by the column
	// configuration.
	if !resourceCfg.DisableReadDescriptions {
		colDef.Description = tb.getDescription(field, cfg.Description, meta)
	}

	if funk.ContainsString(resourceCfg.IgnoreColumnsInTest, colDef.Name) {
		tb.log.Debug("adding ignore in tests to column", "table", table.TableName, "column", colDef.Name, "object", field.Name())
		colDef.IgnoreInTests = true
	}

	// Set Resolver
	if err := tb.SetColumnResolver(table, field, &colDef, cfg, meta); err != nil {
		return err
	}

	// limit max column length, this is because of postgres, we can make this configurable in the future.
	if len(colDef.Name) > MaxColumnLength {
		return fmt.Errorf("column %s name length is too long, max allowed is %d chars, consider renaming/skip_prefix", colDef.Name, MaxColumnLength)
	}

	valueType := field.Type()
	if schema.ValueTypeFromString(cfg.Type) != schema.TypeInvalid {
		valueType = source.TypeUserDefined
	}
	if valueType == schema.TypeInvalid {
		return fmt.Errorf("unsupported type %T for %s", field.Type(), GetColumnName(field.Name(), meta))
	}
	switch valueType {
	case source.TypeRelation:
		tb.log.Debug("building column relation", "table", table.TableName, "column", field.Name(), "object", field.Name())
		relationCfg := resourceCfg.GetRelationConfig(colDef.Name)
		if relationCfg == nil {
			tb.log.Debug("relation config not defined in parent resource, assuming configuration", "table", table.TableName, "column", field.Name(), "object", field.Name())
			relationCfg = &config.RelationConfig{
				ResourceConfig: config.ResourceConfig{
					Service: resourceCfg.Service,
					Domain:  resourceCfg.Domain,
					Name:    colDef.Name,
					Path:    field.Path(),
				},
			}
		}
		if relationCfg.Path == "" {
			relationCfg.Path = field.Path()
		}

		// parent LimitDepth should be transferred to rel
		if resourceCfg.LimitDepth > 0 {
			relationCfg.LimitDepth = resourceCfg.LimitDepth
		}

		// increase build depth
		meta.Depth += 1
		meta.FieldParts = append(meta.FieldParts, field.Name())
		if err := tb.buildTableRelation(table, relationCfg, meta); err != nil {
			return err
		}
	case source.TypeEmbedded:
		tb.log.Debug("Building embedded column", "table", table.TableName, "column", field.Name())
		meta.FieldParts = append(meta.FieldParts, field.Name())
		if err := tb.buildColumns(table, field, resourceCfg, BuildColumnMeta(field, meta, cfg)); err != nil {
			return err
		}
	case source.TypeUserDefined:
		tb.log.Debug("Changing column to user defined", "table", table.TableName, "column", field.Name(), "valueType", valueType, "userDefinedType", cfg.Type)
		colDef.Type = schema.ValueTypeFromString(cfg.Type)
		tb.addPathResolver(field.Name(), &colDef, nil, meta)
		table.Columns = append(table.Columns, colDef)
	default:
		colDef.Type = valueType
		tb.addPathResolver(field.Name(), &colDef, nil, meta)
		table.Columns = append(table.Columns, colDef)
	}
	return nil
}

func (tb TableBuilder) addUserDefinedColumns(table *TableDefinition, resource *config.ResourceConfig) error {
	for _, uc := range resource.UserDefinedColumn {
		tb.log.Debug("adding user defined column", "table", table.TableName, "column", uc.Name)
		colDef := ColumnDefinition{
			Name:        uc.Name,
			Description: uc.Description,
			Type:        schema.ValueTypeFromString(uc.Type),
		}

		if uc.Resolver != nil && uc.GenerateResolver {
			return fmt.Errorf("can't set resolver with generate resolver flag set to true")
		}
		// if we were requested to generate a resolver we create the resolver with the implementation
		if uc.GenerateResolver {
			columnResolver, err := tb.buildResolverDefinition(table, &config.FunctionConfig{
				Name:     template.ToGoPrivate(fmt.Sprintf("resolve%s%s%s", titler.String(resource.Domain), titler.String(inflection.Singular(table.Name)), strcase.ToCamel(uc.Name))),
				Body:     defaultImplementation,
				Path:     path.Join(sdkPath, "provider/schema.ColumnResolver"),
				Generate: true,
			})
			if err != nil {
				return err
			}
			colDef.Resolver = columnResolver
			table.Columns = append(table.Columns, colDef)
			continue
		}

		if uc.Resolver == nil {
			tb.log.Debug("no resolver set and generate resolver was false", "table", table.TableName, "column", uc.Name)
			table.Columns = append(table.Columns, colDef)
			continue
		}
		if uc.Resolver.Body != "" {
			columnResolver, err := tb.buildResolverDefinition(table, &config.FunctionConfig{
				Name:     template.ToGoPrivate(fmt.Sprintf("resolve%s%s%s", titler.String(resource.Domain), titler.String(inflection.Singular(table.Name)), titler.String(uc.Name))),
				Body:     uc.Resolver.Body,
				Path:     path.Join(sdkPath, "provider/schema.ColumnResolver"),
				Generate: true,
			})
			if err != nil {
				return err
			}
			colDef.Resolver = columnResolver
		} else {
			ro, err := tb.finder.FindObjectFromName(uc.Resolver.Path)
			if err != nil {
				return fmt.Errorf("user defined column %s requires resolver definition %w", uc.Name, err)
			}
			// if the resolver is a path resolver
			if uc.Resolver.PathResolver {
				tb.addPathResolver(template.ToGo(uc.Name), &colDef, ro, BuildMeta{})
			} else {
				colDef.Resolver = &ResolverDefinition{Type: ro, Params: uc.Resolver.Params}
			}
		}
		table.Columns = append(table.Columns, colDef)
	}
	return nil
}

func (tb TableBuilder) buildResolverDefinition(table *TableDefinition, cfg *config.FunctionConfig) (*ResolverDefinition, error) {
	ro, err := tb.finder.FindObjectFromName(cfg.Path)
	if err != nil {
		return nil, err
	}

	var signature *types.Signature
	switch t := ro.Type().(type) {
	case *types.Signature:
		signature = t
	case *types.Named:
		if _, ok := t.Underlying().(*types.Signature); !ok {
			return nil, fmt.Errorf("%s not a function", cfg.Path)
		}
		signature = t.Underlying().(*types.Signature)
	default:
		return nil, fmt.Errorf("%s not a function", cfg.Path)
	}

	body := defaultImplementation
	if cfg.Body != "" {
		body = cfg.Body
	}
	funcBody := tb.rewriter.GetFunctionBody(cfg.Name, body)
	if funcBody == defaultImplementation {
		tb.log.Debug("Using default implementation for function", "function", cfg.Name)
	}
	def := &ResolverDefinition{
		Name:      cfg.Name,
		Body:      funcBody,
		Type:      ro,
		Arguments: GetFunctionParams(signature),
		Generate:  cfg.Generate,
		Params:    cfg.Params,
	}
	// if user requested to generate or gave us a body in the configuration
	if cfg.Generate || (cfg.Body != "") {
		// Set signature of function as the generated resolver name
		def.Signature = cfg.Name
		def.Generate = true
		table.Functions = append(table.Functions, def)
	}
	return def, nil
}

func (tb TableBuilder) SetColumnResolver(tableDef *TableDefinition, field source.Object, colDef *ColumnDefinition, cfg config.ColumnConfig, meta BuildMeta) error {
	if cfg.Resolver != nil {
		ro, err := tb.finder.FindObjectFromName(cfg.Resolver.Path)
		if err != nil {
			return err
		}
		if cfg.Resolver.Path != "" && cfg.Resolver.PathResolver {
			tb.addPathResolver(field.Name(), colDef, ro, meta)
		} else {
			colDef.Resolver = &ResolverDefinition{
				Type:   ro,
				Params: cfg.Resolver.Params,
			}
		}
	}
	if cfg.Rename != "" {
		colDef.Name = cfg.Rename
		tb.addPathResolver(field.Name(), colDef, nil, meta)
	}
	if cfg.GenerateResolver {
		if colDef.Resolver != nil {
			tb.log.Warn("overriding already defined column resolver", "column", colDef.Name, "resolver", colDef.Resolver.Name)
		}
		columnResolver, err := tb.buildResolverDefinition(tableDef, &config.FunctionConfig{
			Name:     fmt.Sprintf("resolve%s%s", titler.String(tableDef.TableFuncName), titler.String(strcase.ToCamel(colDef.Name))),
			Body:     defaultImplementation,
			Path:     path.Join(sdkPath, "provider/schema.ColumnResolver"),
			Generate: true,
		})
		if err != nil {
			return err
		}
		colDef.Resolver = columnResolver
		// Set signature of function as the generated resolver name
		colDef.Resolver.Signature = colDef.Resolver.Name
		colDef.Resolver.Signature = colDef.Resolver.Name
	}
	return nil
}

func (tb TableBuilder) addPathResolver(fieldName string, definition *ColumnDefinition, funcObj types.Object, meta BuildMeta) {
	if definition.Resolver != nil {
		return
	}
	signatureName := "schema.PathResolver"
	if funcObj != nil {
		signatureName = fmt.Sprintf("%s.%s", funcObj.Pkg().Name(), funcObj.Name())
	}
	if meta.FieldPath != "" {
		tb.log.Debug("Adding embedded resolver path", "column", strcase.ToCamel(definition.Name), "field", fieldName)
		definition.Resolver = &ResolverDefinition{
			Type:      funcObj,
			Signature: fmt.Sprintf("%s(\"%s.%s\")", signatureName, meta.FieldPath, fieldName),
		}
		return
	}
	// use strcase here since sdk uses it and we didn't get a user defined path resolver
	if strcase.ToCamel(definition.Name) == fieldName && funcObj == nil {
		return
	}
	tb.log.Debug("Adding path resolver column name. camelCase is not same as original field name", "column", strcase.ToCamel(definition.Name), "field", fieldName)
	definition.Resolver = &ResolverDefinition{
		Type:      funcObj,
		Signature: fmt.Sprintf("%s(\"%s\")", signatureName, fieldName),
	}
}

func (tb TableBuilder) buildTableRelations(parentTable *TableDefinition, cfg *config.ResourceConfig, _ BuildMeta) error {

	for _, relCfg := range cfg.UserRelations {
		// if relation already exists i.e was built from one of the columns we skip it
		if parentTable.RelationExists(relCfg) {
			tb.log.Debug("Relation already exists skipping", "relation", relCfg.Name, "path", relCfg.Path)
			continue
		}

		// parent LimitDepth should be transferred to rel
		if cfg.LimitDepth > 0 {
			relCfg.LimitDepth = cfg.LimitDepth
		}
		if err := tb.buildTableRelation(parentTable, &relCfg, BuildMeta{userRelation: true}); err != nil {
			return err
		}
	}
	return nil
}

func (tb TableBuilder) buildTableRelation(parentTable *TableDefinition, cfg *config.RelationConfig, meta BuildMeta) error {
	if cfg.LimitDepth != 0 && meta.Depth >= cfg.LimitDepth {
		tb.log.Warn("depth level exceeded", "parent_table", parentTable.TableName, "table", cfg.Name)
		return nil
	}
	tb.log.Debug("building column relation", "parent_table", parentTable.TableName, "table", cfg.Name)

	if cfg.Rename != "" {
		tb.log.Debug("rename relation", "from", cfg.Name, "to", cfg.Rename)
		cfg.Name = cfg.Rename
	}

	if cfg.Embed {
		return tb.buildEmbeddedRelation(parentTable, cfg, meta)
	}

	innerBuilder := TableBuilder{
		finder:             tb.finder,
		source:             tb.source,
		descriptionSource:  tb.descriptionSource,
		rewriter:           tb.rewriter,
		log:                tb.log,
		descriptionParsers: tb.descriptionParsers,
	}

	newMeta := BuildMeta{Depth: meta.Depth, BaseFieldIndex: meta.BaseFieldIndex, ColumnPath: "", FieldPath: "", FieldParts: meta.FieldParts}
	rel, err := innerBuilder.BuildTable(parentTable, &cfg.ResourceConfig, newMeta)
	if err != nil {
		return err
	}
	rel.Columns = append([]ColumnDefinition{{
		Name:        strings.ToLower(fmt.Sprintf("%s_cq_id", naming.CamelToSnake(inflection.Singular(parentTable.Name)))),
		Type:        schema.TypeUUID,
		Description: fmt.Sprintf("Unique CloudQuery ID of %s table (FK)", parentTable.TableName),
		Resolver:    &ResolverDefinition{Signature: "schema.ParentIdResolver"}},
	}, rel.Columns...)
	parentTable.Relations = append(parentTable.Relations, rel)
	return nil
}

func (tb TableBuilder) buildEmbeddedRelation(parentTable *TableDefinition, cfg *config.RelationConfig, parentMeta BuildMeta) error {
	obj, err := tb.source.Find(cfg.Path)
	if err != nil {
		return err
	}
	meta := parentMeta
	if cfg.SkipPrefix {
		meta.ColumnPath = ""
	}
	return tb.buildColumns(parentTable, obj, &cfg.ResourceConfig, meta)
}

func (tb TableBuilder) getDescription(obj source.Object, description string, meta BuildMeta) string {
	if description != "" {
		return tb.parseDescription(description)
	}
	// if alternative description source is defined
	if tb.descriptionSource != nil { //nolint
		parts := make([]string, 0, len(meta.FieldParts)+1)
		copy(parts, meta.FieldParts)
		parts = append(parts, obj.Name()) //nolint
		d, err := tb.descriptionSource.FindDescription(parts...)
		if err != nil {
			return ""
		}
		if d == "" && obj.Description() != "" {
			return obj.Description()
		}
		return d
	}
	return tb.parseDescription(obj.Description())
}

func (tb TableBuilder) parseDescription(description string) string {
	for _, dp := range tb.descriptionParsers {
		description = dp.Parse(description)
	}
	return description
}
