package cmd

import (
	"slices"
	"sync"
	"time"
)

// Get a list of all tables in a tag
func (j JsonStruct) getTablesInTag(tag string) []string {
	tables := []string{}
	for _, table := range j.Tables {
		if slices.Contains(table.Tags, tag) {
			tables = append(tables, table.Target.Name)
		}
	}
	return tables
}

// Gets the cost of running one or more tags in a Dataform project
//
// TODO: This has to made faster!!!
func (j JsonStruct) GetDataProcessedForTags(tag *string, all bool, keyfile string, projectId string, gitMetadata GitMetadata, location string,
	includeAssertions bool) *[]TagReponse {

	filteredTables := []string{}
	tableDataProcessed := make(map[string]float32)
	tableErrorDetails := make(map[string]ErrorDetails)
	mapOfTags := make(map[string]TagReponse)
	finalResult := []TagReponse{}
	runDateTime := time.Now()
	someMapMutex := sync.RWMutex{} // Might need to use this when we are writing to the map in parallel ?

	if !all {
		filteredTables = j.getTablesInTag(*tag)
	} else {
		for _, table := range j.Tables {
			table_name := table.Target.Name
			filteredTables = append(filteredTables, table_name)
		}
	}

	var wg sync.WaitGroup
	for i := 1; i <= len(filteredTables); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			table := filteredTables[i-1]
			cost, _, err := j.DryRun(&table, keyfile, projectId, location, includeAssertions)
			disabled, err := j.IsTableDisabled(&table)

			someMapMutex.Lock()
			if err != nil {
				tableDataProcessed[table] = 0.0
				tableErrorDetails[table] = ErrorDetails{
					IsError:  true,
					ErrorMsg: err.Error(),
					Disabled: disabled,
				}
			} else {
				tableDataProcessed[table] = cost / 1e9
			}
			someMapMutex.Unlock()
		}(i)
	}
	wg.Wait()

	// When a specific tag is requested
	if *tag != "" {
		for table, dataProcessed := range tableDataProcessed {

			tablesDryRunMetadata := []TablesDryRunMetadata{}
			errorMessage := tableErrorDetails[table].ErrorMsg
			lineNumber, columnNumber := getLineAndColumnNumberFromDryRunError(errorMessage)

			error := TablesDryRunMetadata{
				Table: table,
				Error: ErrorDetails{
					IsError:      tableErrorDetails[table].IsError,
					LineNumber:   lineNumber,
					ColumnNumber: columnNumber,
					ErrorMsg:     tableErrorDetails[table].ErrorMsg,
					Disabled:     tableErrorDetails[table].Disabled,
				},
			}
			tablesDryRunMetadata = append(tablesDryRunMetadata, error)

			if _, ok := mapOfTags[*tag]; !ok {

				mapOfTags[*tag] = TagReponse{
					Tag:                  *tag,
					GBProcessed:          dataProcessed,
					Cost:                 (dataProcessed / 1000.) * COST_IN_POUNDS_FOR_TERRABYTE,
					RunDateTime:          runDateTime,
					TablesDryRunMetadata: tablesDryRunMetadata,
					GitMetadata:          gitMetadata,
					HasError:             tableErrorDetails[table].IsError,
				}
			} else {

				if entry, ok := mapOfTags[*tag]; ok {

					entry.TablesDryRunMetadata = append(entry.TablesDryRunMetadata, tablesDryRunMetadata...)
					entry.GBProcessed += dataProcessed
					entry.Cost = (entry.GBProcessed / 1000.) * COST_IN_POUNDS_FOR_TERRABYTE
					entry.HasError = entry.HasError || tableErrorDetails[table].IsError
					mapOfTags[*tag] = entry

				}
			}
		}
		finalResult = append(finalResult, mapOfTags[*tag])
		return &finalResult
	}

	// When all tags are requested
	for _, table := range j.Tables {
		for _, tag := range table.Tags {
			tableName := table.Target.Name
			tables := []string{table.Target.Name}
			dataProcessed := tableDataProcessed[tables[0]]

			tablesDryRunMetadata := []TablesDryRunMetadata{}

			errorMessage := tableErrorDetails[tables[0]].ErrorMsg
			lineNumber, columnNumber := getLineAndColumnNumberFromDryRunError(errorMessage)

			error := TablesDryRunMetadata{
				Table: tableName,
				Error: ErrorDetails{
					IsError:      tableErrorDetails[tables[0]].IsError,
					LineNumber:   lineNumber,
					ColumnNumber: columnNumber,
					ErrorMsg:     tableErrorDetails[tables[0]].ErrorMsg,
					Disabled:     tableErrorDetails[tables[0]].Disabled,
				},
			}
			tablesDryRunMetadata = append(tablesDryRunMetadata, error)

			if _, ok := mapOfTags[tag]; !ok {

				mapOfTags[tag] = TagReponse{
					Tag:                  tag,
					GBProcessed:          dataProcessed,
					Cost:                 (dataProcessed / 1000.) * COST_IN_POUNDS_FOR_TERRABYTE,
					RunDateTime:          runDateTime,
					TablesDryRunMetadata: tablesDryRunMetadata,
					GitMetadata:          gitMetadata,
					HasError:             tableErrorDetails[tables[0]].IsError,
				}
			} else {

				if entry, ok := mapOfTags[tag]; ok {

					entry.TablesDryRunMetadata = append(entry.TablesDryRunMetadata, tablesDryRunMetadata...)
					entry.GBProcessed += dataProcessed
					entry.Cost = (entry.GBProcessed / 1000.) * COST_IN_POUNDS_FOR_TERRABYTE
					entry.HasError = entry.HasError || tableErrorDetails[tables[0]].IsError
					mapOfTags[tag] = entry

				}
			}
		}
	}

	for _, tag := range mapOfTags {
		finalResult = append(finalResult, tag)
	}

	return &finalResult
}
