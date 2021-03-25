package codegen

import (
	"github.com/cloudquery/cq-gen/code"
	"github.com/cloudquery/cq-gen/codegen/config"
	"github.com/cloudquery/cq-gen/rewrite"
)

type ResourceDefinition struct {
	Config config.ResourceConfig
	Table *TableDefinition
}

func buildResources(cfg *config.Config) ([]*ResourceDefinition, error) {
	rw, err := rewrite.New(cfg.OutputDirectory)
	if err != nil {
		return nil, err
	}
	b := builder{code.NewFinder(), rw}
	return b.build(cfg)
}

// builder generates all models code based from the configuration
type builder struct {
	finder   *code.Finder
	rewriter *rewrite.Rewriter
}

func (b builder) build(cfg *config.Config) ([]*ResourceDefinition, error) {
	resources := make([]*ResourceDefinition, len(cfg.Resources))
	for i, resource := range cfg.Resources {
		t, err := b.buildTable(resource)
		if err != nil {
			return nil, err
		}
		resources[i] = &ResourceDefinition{
			Config: resource,
			Table:  t,
		}
	}
	return resources, nil
}



