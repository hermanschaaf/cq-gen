package config

import (
	"github.com/cloudquery/cloudquery-plugin-sdk/plugin/schema"
	"github.com/creasty/defaults"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

type Config struct {
	Service         string           `hcl:"service"`
	OutputDirectory string           `hcl:"output_directory"`
	Resources       []ResourceConfig `hcl:"resource,block"`
}

type ResourceConfig struct {
	Service           string           `hcl:"service,label"`
	Domain            string           `hcl:"domain,label"`
	Name              string           `hcl:"name,label"`
	Path              string           `hcl:"path"`
	Columns           []ColumnConfig   `hcl:"column,block"`
	Relations         []ResourceConfig `hcl:"relation,block"`
	UserDefinedColumn []ColumnConfig   `hcl:"userDefinedColumn,block"`
}

func (r ResourceConfig) GetColumnConfig(name string) ColumnConfig {
	for _, c := range r.Columns {
		if c.Name == name {
			return c
		}
	}
	var c ColumnConfig
	defaults.Set(&c)
	c.Name = name
	return c
}

type ColumnConfig struct {
	Name         string           `hcl:"name,label"`
	SkipPrefix   bool             `hcl:"skip_prefix,optional" defaults:"false"`
	Skip         bool             `hcl:"skip,optional" defaults:"false"`
	Type         schema.ValueType `hcl:"type,optional"`
	ResolverPath string           `hcl:"resolver,optional"`
}

func Parse(configPath string) (*Config, error) {
	var config Config
	if err := hclsimple.DecodeFile(configPath, nil, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
