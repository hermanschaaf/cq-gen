package base

import (
	"context"

	"github.com/cloudquery/cq-provider-sdk/provider/schema"
)

func Descriptions() *schema.Table {
	return &schema.Table{
		Name:        "test_base_descriptions",
		Description: "DescriptionStruct describes itself as a struct",
		Resolver:    fetchBaseDescriptions,
		Columns: []schema.Column{
			{
				Name:        "simple",
				Description: "Simple is a simple description",
				Type:        schema.TypeString,
			},
			{
				Name:        "complex",
				Description: "Value is a string",
				Type:        schema.TypeString,
			},
			{
				Name:        "created_at",
				Description: "A sentence with a full stop that should get removed",
				Type:        schema.TypeTimestamp,
			},
		},
	}
}

// ====================================================================================================================
//                                               Table Resolver Functions
// ====================================================================================================================

func fetchBaseDescriptions(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
	panic("not implemented")
}
