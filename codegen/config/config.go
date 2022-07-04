package config

import (
	"fmt"
	"strings"

	"github.com/creasty/defaults"
	"github.com/zclconf/go-cty/cty"

	"github.com/cloudquery/cq-gen/code"
)

type Config struct {
	Service              string                `hcl:"service"`
	AddGenerate          bool                  `hcl:"add_generate"`
	OutputDirectory      string                `hcl:"output_directory"`
	DataSource           *DataSource           `hcl:"data_source,block"`
	DescriptionSource    *DescriptionSource    `hcl:"description_source,block"`
	DescriptionModifiers []DescriptionModifier `hcl:"description_modifier,block"`
	Resources            []ResourceConfig      `hcl:"resource,block"`

	Path string
}

type DataSource struct {
	Type string `hcl:"type,label"`
	Path string `hcl:"path"`
}

type DescriptionSource struct {
	Type string `hcl:"type,label"`
	Path string `hcl:"path"`
}

type DescriptionModifier struct {
	Name        string   `hcl:"name,label"`
	Regex       string   `hcl:"regex,optional"`
	RemoveWords []string `hcl:"words,optional"`
}

func (c Config) GetResource(resource string) (ResourceConfig, error) {
	for _, r := range c.Resources {
		if r.Name == resource {
			return r, nil
		}
	}
	return ResourceConfig{}, fmt.Errorf("didn't find resource %s in config", resource)
}

type TableOptionsConfig struct {
	PrimaryKeys []string `hcl:"primary_keys,optional"`
}

type ResourceConfig struct {
	// Name of service i.e AWS,Azure etc'
	Service string `hcl:"service,label"`
	// Domain this resource belongs too, i.e Storage, Users etc'
	Domain string `hcl:"domain,label"`
	// Name of the resource table
	Name string `hcl:"name,label"`
	// Description of the table
	Description string `hcl:"description,optional"`
	// Path to the struct/source we are generating from
	Path string `hcl:"path,optional"`
	// Allow unexported fields in struct/source we are generating from to be added as well
	AllowUnexported bool `hcl:"allow_unexported,optional"`
	// Table options in the config
	TableOptions *TableOptionsConfig `hcl:"options,block"`
	// IgnoreInTests specifies whether the table should be ignored during non-nil checks in integration tests
	IgnoreInTests bool `hcl:"ignore_in_tests,optional"`

	// Column configurations we want to modify
	Columns []ColumnConfig `hcl:"column,block"`
	// Relations configurations we want to modify
	Relations []RelationConfig `hcl:"relation,block"`
	// User Relations we want to add that aren't part of the original structure
	UserRelations []RelationConfig `hcl:"user_relation,block"`

	// UserDefinedColumns are a list of columns we add that aren't part of the original struct
	UserDefinedColumn []ColumnConfig `hcl:"userDefinedColumn,block"`

	// Function configurations will be omitted if not set
	IgnoreError          *FunctionConfig `hcl:"ignoreError,block"`
	Multiplex            *FunctionConfig `hcl:"multiplex,block"`
	DeleteFilter         *FunctionConfig `hcl:"deleteFilter,block"`
	PostResourceResolver *FunctionConfig `hcl:"postResourceResolver,block"`
	Resolver             *FunctionConfig `hcl:"resolver,block"`

	// LimitDepth limits the depth cq-gen enters the structs, this is to avoid recursive structs
	LimitDepth int `hcl:"limit_depth,optional"`

	// Disables reading the struct for description comments for each column
	DisableReadDescriptions bool `hcl:"disable_auto_descriptions,optional"`
	// Disable pluralize of the name of the resource
	NoPluralize bool `hcl:"disable_pluralize,optional"`

	DescriptionPathParts []string `hcl:"description_path_parts,optional"`
}

type RelationConfig struct {
	ResourceConfig
	Rename string `hcl:"rename,optional"`
	// Embed all the relations columns into the parent struct
	Embed bool `hcl:"embed,optional"`
	// SkipPrefix skips the embedded relation name prefix for all it's embedded columns
	SkipPrefix bool `hcl:"skip_prefix,optional"`
}

type FunctionConfig struct {
	// Name of function, usually auto generated by cq-gen
	Name string `hcl:"name,label"`
	// Body to insert when function is generated, use with care, auto importing isn't supported in user defined bodies
	Body string `hcl:"body,optional"`
	// Path to a function to use.
	Path string `hcl:"path"`
	// Generate tells cq-gen to create the function code in template, usually set automatically.
	// Setting to true will force function generation in template.
	Generate bool `hcl:"generate,optional"`
	// PathResolver defines this function to be called the FieldPath traversed, this is used by generic functions
	PathResolver bool `hcl:"path_resolver,optional"`
	// Parameters to pass
	Params []cty.Value `hcl:"params,optional"`
}

func (r ResourceConfig) GetRelationConfig(name string) *RelationConfig {
	for _, r := range r.Relations {
		if strings.EqualFold(r.Name, name) {
			return &r
		}
		if _, typeName := code.PkgAndType(r.Path); typeName == name {
			return &r
		}
	}
	return nil
}

func (r ResourceConfig) GetColumnConfig(baseName string, names ...string) ColumnConfig {
	for _, c := range r.Columns {
		if c.Name == baseName {
			return c
		}
		for _, n := range names {
			if c.Name == n {
				return c
			}
		}
	}
	var c ColumnConfig
	_ = defaults.Set(&c)
	c.Name = baseName
	return c
}

type ColumnConfig struct {
	// Name of column as defined by resource, in snake_case, be careful with abbreviations
	Name string `hcl:"name,label"`
	// Description of column to display to user
	Description string `hcl:"description,optional"`
	// SkipPrefix Whether we want to skip adding the embedded prefix to a column
	SkipPrefix bool `hcl:"skip_prefix,optional" defaults:"false"`
	// Skip makes cq-gen skip the column
	Skip bool `hcl:"skip,optional" defaults:"false"`
	// GenerateResolver whether to force a resolver creation
	GenerateResolver bool `hcl:"generate_resolver,optional"`
	// Resolver unique resolver function to use
	Resolver *FunctionConfig `hcl:"resolver,block"`
	// Type Overrides column type, use carefully, validation will fail if interface{} of value isn't the same as expected ValueType
	Type string `hcl:"type,optional"`
	// Rename column name, if no resolver is passed schema.PathResolver will be used
	Rename string `hcl:"rename,optional"`
	// ExtractDescriptionFromParentField, take column description from parent spec
	ExtractDescriptionFromParentField bool `hcl:"extract_description_from_parent_field,optional" defaults:"false"`
}
