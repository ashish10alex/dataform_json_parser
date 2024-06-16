package cmd

import (
	"github.com/ashish10alex/dj/cmd"
	"fmt"

	"github.com/spf13/cobra"
)

var readJsonFile = cmd.ReadJson
var JsonStruct cmd.JsonStruct
var jsonPrintTagReponse = cmd.JsonPrintTagReponse
var JsonPrint = cmd.JsonPrint

// tagopsCmd represents the tagops command
var tagopsCmd = &cobra.Command{
	Use:   "tag-ops",
	Short: "All tag related operations",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		jsonFile := cmd.Flag("json-file").Value.String()
		if cmd.Flag("unique").Value.String() == "true" {
			jsonData, err := readJsonFile(jsonFile)
			if err != nil {
				fmt.Println(err.Error())
			}
			tags := jsonData.GetUniqueTags()
			JsonPrint(tags)
		}
	},
}

var tagOpsCostCmd = &cobra.Command{
	Use:   "cost",
	Short: "Get cost for a specific tag or all tags",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		jsonFile := cmd.Flag("json-file").Value.String()
		location := cmd.Flag("location").Value.String()
		includeAssertions, _ := cmd.Flags().GetBool("include-assertions")
		compact, _ := cmd.Flags().GetBool("compact")
		jsonData, err := readJsonFile(jsonFile)
		if err != nil {
			fmt.Println(err.Error())
		}

		getAllTags, _ := cmd.Flags().GetBool("all")
		tag := cmd.Flag("tag").Value.String()
		projectId := jsonData.GetTargetGcpProjectId()
		keyFile := ""
		gitMetadata := jsonData.GetGitMetadata(compact)

		if getAllTags {
			tag := ""
			tagReponse := jsonData.GetDataProcessedForTags(&tag, true, keyFile, projectId, gitMetadata, location, includeAssertions)
			jsonPrintTagReponse(tagReponse)
		} else if tag != "" {
			tagReponse := jsonData.GetDataProcessedForTags(&tag, false, keyFile, projectId, gitMetadata, location, includeAssertions)
			jsonPrintTagReponse(tagReponse)
		}
	},
}

func init() {

	tagopsCmd.Flags().BoolP("unique", "u", false, "Get unique tags from compiled Dataform json")
	tagOpsCostCmd.Flags().BoolP("compact", "c", true, "compact af")
	tagOpsCostCmd.Flags().BoolP("include-assertions", "i", false, "Include assertions in the cost")
	tagOpsCostCmd.Flags().BoolP("all", "a", false, "Get cost for all tags")
	tagOpsCostCmd.Flags().StringP("tag", "t", "", "Get cost for a specific tag")
	tagopsCmd.AddCommand(tagOpsCostCmd)
	cmd.RootCmd.AddCommand(tagopsCmd)

}
