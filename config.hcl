

service = "aws"
output_directory="resources"

//resource "aws" "s3" "buckets" {
//  path = "github.com/aws/aws-sdk-go-v2/service/s3/types.Bucket"
//
//  userDefinedColumn "account_id" {
//    type = 5
//    resolver = "github.com/cloudquery/cq-gen/resources/resolvers.ResolveAwsAccount"
//  }
//
//  relation "aws" "s3" "grants" {
//    path = "github.com/aws/aws-sdk-go-v2/service/s3/types.Grant"
//
//    column "grantee" {
//      skip_prefix = true
//    }
//  }
//}

resource "aws" "ec2" "instances" {
  path = "github.com/aws/aws-sdk-go-v2/service/ec2/types.Instance"

  column "tags" {
    type = 10 // TypeJson
  }
}