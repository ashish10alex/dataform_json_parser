package cmd

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// Gets the line and column number where the error message is using regex by looking for [lineNumber:columnNumber] pattern in the string
// If error position cannot be inferred 0, 0 are returned
func getLineAndColumnNumberFromDryRunError(stderrMessage string) (int, int) {
	pattern := regexp.MustCompile(`(\d+):(\d+)`)
	match := pattern.FindStringSubmatch(stderrMessage)

	if len(match) >= 3 {
		lineNumber, _ := strconv.Atoi(match[1])
		columnNumber, _ := strconv.Atoi(match[2])
		// fmt.Printf("Error found at line %d, column %d\n", lineNumber, columnNumber)
		return lineNumber, columnNumber
	} else {
		return 0, 0
	}
}

func jsonMarshallStdOutErr(tableMetadata TableMetadata,  gitMetadata *GitMetadata, isError bool, lineNumber int,
	columnNumber int, errorMessage string, bytesProcessed float32,
	gbProcessed float32, compact bool, runDateTime time.Time) {
	if compact {
        tableMetadata.Query = ""
	}
	responseDataStruct := DryRunReponse{
		FileName: tableMetadata.TargetMetadata.Name,
		Schema:   tableMetadata.TargetMetadata.Schema,
		Database: tableMetadata.TargetMetadata.Database,
		Query:    tableMetadata.Query,
		Compact:  compact,
		Error: ErrorDetails{
			IsError:      isError,
			LineNumber:   lineNumber,
			ColumnNumber: columnNumber,
			ErrorMsg:     errorMessage,
			Disabled:     tableMetadata.TargetMetadata.Disabled,
		},
		Cost:        (gbProcessed / 1000.) * COST_IN_POUNDS_FOR_TERRABYTE,
		GitMetadata: *gitMetadata,
		GBProcessed: gbProcessed,
		RunDateTime: runDateTime,
	}

	responseDataJson, _ := json.Marshal(responseDataStruct)
	fmt.Printf("%s \n", responseDataJson)

}

func ProcessErrorAndOutput(tableMetadata TableMetadata,  gitMetadata *GitMetadata, bytesProcessed float32, err error, compact bool, runDateTime time.Time) {
	if err != nil {
		error_message := err.Error()
		lineNumber, columnNumber := getLineAndColumnNumberFromDryRunError(err.Error())
		jsonMarshallStdOutErr(tableMetadata,  gitMetadata, true, lineNumber, columnNumber, error_message, 0., 0., compact, runDateTime)
	} else {
		gbProcessed := bytesProcessed / 1e9
		jsonMarshallStdOutErr(tableMetadata,  gitMetadata, false, 0, 0, "", bytesProcessed, gbProcessed, compact, runDateTime)
	}
}
