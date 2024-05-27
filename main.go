/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"dataform_json_parser/cmd"
	_ "dataform_json_parser/cmd/table_ops"
	_ "dataform_json_parser/cmd/tag_ops"
)

func main() {
	cmd.Execute()
}
