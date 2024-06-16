/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/ashish10alex/dataform_json_parser/cmd"
	_ "github.com/ashish10alex/dataform_json_parser/cmd/table_ops"
	_ "github.com/ashish10alex/dataform_json_parser/cmd/tag_ops"
)

func main() {
	cmd.Execute()
}
