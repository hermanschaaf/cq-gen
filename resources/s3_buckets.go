package resources

import (
	"context"
	"github.com/cloudquery/cloudquery-plugin-sdk/plugin/schema"
	"github.com/cloudquery/cq-gen/resources/resolvers"
)

func S3buckets() *schema.Table {
	return &schema.Table{
		Name:     "aws_s3_buckets",
		Resolver: FetchAWSS3Buckets,
		Columns: []schema.Column{
			{
				Name: "creation_date",
				Type: schema.TypeTimestamp,
			},
			{
				Name: "name",
				Type: schema.TypeString,
			},
			{
				Name:     "account_id",
				Type:     schema.TypeString,
				Resolver: resolvers.ResolveAwsAccount,
			},
		},
		Relations: []*schema.Table{
			{
				Name:     "aws_s3_grants",
				Resolver: FetchAWSS3Grants,
				Columns: []schema.Column{
					{
						Name: "grantee",
						Type: schema.TypeEmbedded,
						Elem: []schema.Column{
							{
								Name:     "type",
								Type:     schema.TypeString,
								Resolver: schema.PathResolver("Grantee.Type"),
							},
							{
								Name:     "display_name",
								Type:     schema.TypeString,
								Resolver: schema.PathResolver("Grantee.DisplayName"),
							},
							{
								Name:     "email_address",
								Type:     schema.TypeString,
								Resolver: schema.PathResolver("Grantee.EmailAddress"),
							},
							{
								Name:     "id",
								Type:     schema.TypeString,
								Resolver: schema.PathResolver("Grantee.ID"),
							},
							{
								Name:     "uri",
								Type:     schema.TypeString,
								Resolver: schema.PathResolver("Grantee.URI"),
							},
						},
					},
					{
						Name: "permission",
						Type: schema.TypeString,
					},
					{
						Name:     "bucket_id",
						Type:     schema.TypeUUID,
						Resolver: schema.ParentIdResolver,
					},
				},
			},
		},
	}
}

// ====================================================================================================================
//                                               Table Resolver Functions
// ====================================================================================================================
func FetchAWSS3Buckets(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan interface{}) error {
	panic("no implemented")
}

func FetchAWSS3Grants(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan interface{}) error {
	panic("no implemented")
}
