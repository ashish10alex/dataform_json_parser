/*
Copyright © 2024 Ashish Alex
*/
package cmd

import (
	"github.com/ashish10alex/dj/cmd"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
)

var JsonFile string
var readJson = cmd.ReadJson
var ProcessErrorAndOutput = cmd.ProcessErrorAndOutput

type Target = cmd.Target
type TableMetadata = cmd.TableMetadata

// tableopsCmd represents the tableops command
var tableopsCmd = &cobra.Command{
	Use:   "table-ops",
	Short: "All table related operations",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		listTables, _ := cmd.Flags().GetBool("list-tables")
		if listTables {
			jsonFile := cmd.Flag("json-file").Value.String()
			jsonData, err := readJson(jsonFile)
			if err != nil {
				fmt.Println(err.Error())
			}
			jsonData.ListTables()
		} else {
			cmd.Help()
		}
	},
}

var tableOpsCostCmd = &cobra.Command{
	Use:   "cost",
	Short: "Get cost for a specific file or all tables",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		jsonFile := cmd.Flag("json-file").Value.String()
		keyFile := cmd.Flag("key-file").Value.String()
		location := cmd.Flag("location").Value.String()
		includeAssertions, _ := cmd.Flags().GetBool("include-assertions")
		compact, _ := cmd.Flags().GetBool("compact")

		if cmd.Flag("all").Value.String() == "true" {

			jsonData, err := readJson(jsonFile)
			if err != nil {
				fmt.Println(err.Error())
			}
			projectId := jsonData.GetTargetGcpProjectId()
			jsonData.GetDryRunForAllTables(projectId, keyFile, location, includeAssertions, compact)

		} else if cmd.Flag("file").Value.String() != "" {
			fileName := cmd.Flag("file").Value.String()
			jsonData, err := readJson(jsonFile)
			if err != nil {
                log.Fatal("Error parsing json: ", err.Error())
			}

			targets := jsonData.GetTargetFromFileName(fileName)
			if len(targets.TableTargets) == 0 && len(targets.AssertionTargets) == 0 && len(targets.OperationTargets) == 0 {
				log.Fatal("Table name not for file name: ", fileName)
			}

			projectId := jsonData.GetTargetGcpProjectId()
			gitMetadata := jsonData.GetGitMetadata(compact)

			var dryRunError error
			runDateTime := time.Now()

			allTargets := append(targets.TableTargets, targets.AssertionTargets...)
            allTargets = append(allTargets, targets.OperationTargets...)

			for _, target := range allTargets {
				bytesProcessed, query, err := jsonData.DryRun(&target.Name, keyFile, projectId, location, includeAssertions)
				if err != nil {
					dryRunError = err
				}
				tableMetadata := TableMetadata{TargetMetadata: target, Query: *query}
				ProcessErrorAndOutput(tableMetadata, &gitMetadata, bytesProcessed, dryRunError, compact, runDateTime)
			}

		}
	},
}

func init() {
	cmd.RootCmd.AddCommand(tableopsCmd)

	tableOpsCostCmd.Flags().BoolP("all", "a", false, "Cost for all tables")
	tableOpsCostCmd.Flags().BoolP("compact", "c", true, "Compact JSON response. Excludes compiled query")
	tableOpsCostCmd.Flags().BoolP("include-assertions", "i", false, "Include assertions in the cost")
	tableOpsCostCmd.Flags().StringP("file", "f", "", "Get cost for a specific file, file name should be passed without the extension")
	tableopsCmd.Flags().BoolP("list-tables", "u", false, "List all tables in the project")
	tableopsCmd.AddCommand(tableOpsCostCmd)
}
