package main

import (
	"fmt"
	"github.com/cloudquery/cq-gen/codegen"
)

func main() {
	err := codegen.Generate("config.hcl")
	if err != nil {
		fmt.Println(err)
		return
	}
}
