package cmd

import (
	"dataform_json_parser/cmd"
	"encoding/json"
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

type OutputSources = cmd.OutputSources

var tableOpsSourcesCmd = &cobra.Command{
	Use:   "declarations-and-targets",
	Short: "Get all the base name of declarations and targets. This can be used to fill in completion menu by the editor",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		jsonFile := cmd.Flag("json-file").Value.String()
        jsonData, err := readJson(jsonFile)
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
        outputSources := OutputSources{}

        for _, declaration := range jsonData.Declarations {
            outputSources.Declarations = append(outputSources.Declarations, declaration.Target.Name)
        }

        for _, target  := range jsonData.Targets{
            outputSources.Targets = append(outputSources.Targets, target.Name)
        }

        sourcesJson, err := json.Marshal(outputSources)
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
        fmt.Println(string(sourcesJson))
	},
}

func init() {
	tableopsCmd.AddCommand(tableOpsSourcesCmd)
}
