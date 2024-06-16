/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/ashish10alex/dj/cmd"
	_ "github.com/ashish10alex/dj/cmd/table_ops"
	_ "github.com/ashish10alex/dj/cmd/tag_ops"
)

func main() {
	cmd.Execute()
}
