package resolvers

import (
	"context"

	"github.com/cloudquery/cq-provider-sdk/provider/schema"
)

func WithTemplates() *schema.Table {
	return &schema.Table{
		Name:     "test_resolvers_with_template",
		Resolver: fetchResolversWithTemplates,
		Columns: []schema.Column{
			{
				Name: "int_value",
				Type: schema.TypeBigInt,
			},
			{
				Name: "bool_value",
				Type: schema.TypeBool,
			},
			{
				Name:     "embedded_field_a",
				Type:     schema.TypeBigInt,
				Resolver: schema.PathResolver("Embedded.FieldA"),
			},
		},
	}
}

// ====================================================================================================================
//                                               Table Resolver Functions
// ====================================================================================================================

func fetchResolversWithTemplates(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
	// GENERATED. Do not edit.
	panic("hello from template params")
}
