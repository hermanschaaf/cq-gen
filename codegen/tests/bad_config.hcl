service          = "test"
output_directory = "./tests/output"

resource "test" "bad" "duplicate" {
  path = "github.com/cloudquery/cq-gen/codegen/tests.BaseStruct"
}

resource "test" "bad" "duplicate" {
  path = "github.com/cloudquery/cq-gen/codegen/tests.BaseStruct"
}

resource "test" "bad" "long_table_name_with_more_than_63_characters_which_is_too_long" {
  path = "github.com/cloudquery/cq-gen/codegen/tests.BaseStruct"
}
