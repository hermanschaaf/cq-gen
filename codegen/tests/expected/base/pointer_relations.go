package base

import (
	"context"

	"github.com/cloudquery/cq-gen/codegen/tests"
	"github.com/cloudquery/cq-provider-sdk/helpers"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"
)

func PointerRelations() *schema.Table {
	return &schema.Table{
		Name:     "test_base_pointer_relations",
		Resolver: fetchBasePointerRelations,
		Columns:  []schema.Column{},
		Relations: []*schema.Table{
			{
				Name:     "test_base_pointer_relation_inner_relations",
				Resolver: fetchBasePointerRelationInnerRelations,
				Columns: []schema.Column{
					{
						Name:        "pointer_relation_cq_id",
						Description: "Unique CloudQuery ID of test_base_pointer_relations table (FK)",
						Type:        schema.TypeUUID,
						Resolver:    schema.ParentIdResolver,
					},
				},
				Relations: []*schema.Table{
					{
						Name:     "test_base_pointer_relation_inner_relation_some_bases",
						Resolver: fetchBasePointerRelationInnerRelationSomeBases,
						Columns: []schema.Column{
							{
								Name:        "pointer_relation_inner_relation_cq_id",
								Description: "Unique CloudQuery ID of test_base_pointer_relation_inner_relations table (FK)",
								Type:        schema.TypeUUID,
								Resolver:    schema.ParentIdResolver,
							},
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
					},
				},
			},
		},
	}
}

// ====================================================================================================================
//                                               Table Resolver Functions
// ====================================================================================================================

func fetchBasePointerRelations(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
	panic("not implemented")
}
func fetchBasePointerRelationInnerRelations(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
	// This is a generated base implementation. You may edit it if necessary.
	p := helpers.ToPointer(parent.Item).(*tests.PointerRelationStruct)
	if p.Inner == nil {
		return nil
	}
	res <- p.Inner.Relations
	return nil
}
func fetchBasePointerRelationInnerRelationSomeBases(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
	// This is a generated base implementation. You may edit it if necessary.
	p := helpers.ToPointer(parent.Item).(*tests.PointerRelation)
	res <- p.SomeBases
	return nil
}
