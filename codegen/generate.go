package codegen

import (
	"fmt"
	"github.com/cloudquery/cq-gen/codegen/config"
	"github.com/cloudquery/cq-gen/codegen/templates"
)

func Generate(configPath string) error {

	cfg, err := config.Parse(configPath)
	if err != nil {
		return err
	}
	resources, err := buildResources(cfg)
	if err != nil {
		return err
	}

	for _, resource := range resources {
		err = templates.Render(templates.Options{
			Template: "codegen/table.gotpl",
			Filename: fmt.Sprintf("resources/%s_%s.go", resource.Config.Domain, resource.Config.Name),
			PackageName: cfg.OutputDirectory,
			Data:     resource,
			Funcs:    nil,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
