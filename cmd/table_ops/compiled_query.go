package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var tableOpsQueryCmd = &cobra.Command{
	Use:   "query",
	Short: "Get compiled query for a specific table",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		jsonFile := cmd.Flag("json-file").Value.String()

		if cmd.Flag("file").Value.String() != "" {

			fileName := cmd.Flag("file").Value.String()
			jsonData, err := readJson(jsonFile)
            if err != nil {
                fmt.Println(err.Error())
            }
			targets := jsonData.GetTargetFromFileName(fileName)

            allTargets := append(targets.TableTargets, targets.AssertionTargets...)
            allTargets = append(allTargets, targets.OperationTargets...)

			includeAssertions := true
			queryForFileName := ""

			for _, target := range allTargets {
				outputQuery, err := jsonData.GetQueryForTable(&target.Name, includeAssertions)
				if err != nil {
					cmd.Println(err)
					return
				}
				queryForFileName += outputQuery.Query + outputQuery.IncrementalPreOpsQuery + outputQuery.IncrementalQuery + outputQuery.Assertion + outputQuery.OperationsQuery
			}

			outputFile := cmd.Flag("out-file").Value.String()
			if outputFile != "" {
				os.WriteFile(outputFile, []byte(*&queryForFileName), 0644)
			} else {
				fmt.Println(*&queryForFileName)
			}
		}
	},
}

func init() {
	tableOpsQueryCmd.Flags().StringP("file", "f", "", "Get compiled query for a specific file")
	tableOpsQueryCmd.Flags().StringP("out-file", "o", "", "SQL file to write the compiled query to")
	tableopsCmd.AddCommand(tableOpsQueryCmd)

}
