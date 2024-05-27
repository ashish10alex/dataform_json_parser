package cmd

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"
	"sync"
	"time"
)

func (j JsonStruct) GetGitMetadata(compact bool) GitMetadata {
	return GitMetadata{
		GitRepositoryId: j.GetGitRepository(),
		GitBranch:       j.GetGitBranch(compact),
	}
}

func (j JsonStruct) GetProjectConfig() ProjectConfig {
	return j.ProjectConfig
}

// convert `/path/to/file_name.sqlx` to `file_name`
func getShortFileNameFromBaseFileName(fileName string) string {
	shortFileName := strings.Split(fileName, "/")[len(strings.Split(fileName, "/"))-1]
	return strings.Split(shortFileName, ".")[0]
}

// Neovim front end will pass the file name, we need to infer the target tables from the file name
func (j JsonStruct) GetTargetFromFileName(fileName string) Targets {

	targets := Targets{}

	for _, table := range j.Tables {

		tableShortFileName := getShortFileNameFromBaseFileName(table.FileName)

		if tableShortFileName == fileName {
			tableTarget := TargetMetadata{table.Target.Schema, table.Target.Name, table.Target.Database, table.Disabled}
			targets.TableTargets = append(targets.TableTargets, tableTarget)
		}
	}

	for _, assertion := range j.Assertions {

		assertionShortFileName := getShortFileNameFromBaseFileName(assertion.FileName)
		if assertionShortFileName == fileName {
			assertionTarget := TargetMetadata{assertion.Target.Schema, assertion.Target.Name, assertion.Target.Database, false}
			targets.AssertionTargets = append(targets.AssertionTargets, assertionTarget)
		}
	}

    for _, operation := range j.Operations {
        operationFileName := getShortFileNameFromBaseFileName(operation.FileName)
        if operationFileName == fileName {
            operationTarget := TargetMetadata{operation.Target.Schema, operation.Target.Name, operation.Target.Database, false}
            targets.OperationTargets = append(targets.OperationTargets, operationTarget)
        }
    }

	return targets
}

// This is O(N) can be optimized to O(1) by using a map ?
func (j JsonStruct) GetQueryForTable(tableName *string, includeAssertions bool) (*OutputQuery, error) {

	outputQuery := OutputQuery{}
	incrementalQuery := ""
	nonIncrementalQuery := ""
	incrementalPreOpsQuery := ""
	assertionQuery := ""
    operationsQuery := ""

	for _, table := range j.Tables {
		if table.Target.Name == *tableName {
			if table.Type == "incremental" {
				nonIncrementalQuery = "-- non IncrementalQuery \n" + table.Query + "; \n"

				incrementalPreOpsQuery = "--incrementalPreOpsQuery \n"
				for i := 0; i < len(table.IncrementalPreOps); i++ {
					incrementalPreOpsQuery += table.IncrementalPreOps[i]
				}
				incrementalPreOpsQuery += "; \n"

				incrementalQuery = "--IncrementalQuery \n" + table.IncrementalQuery + "; \n"
				outputQuery = OutputQuery{nonIncrementalQuery, incrementalPreOpsQuery, incrementalQuery, assertionQuery, operationsQuery}

			} else {
				regularQuery := "--Query \n" + table.Query + "; \n"
				outputQuery = OutputQuery{regularQuery, incrementalPreOpsQuery, incrementalQuery, assertionQuery, operationsQuery}
			}
		}
	}
	// TODO: If table is not found in the tables, then we should check in the assertions ?
	if includeAssertions {
		for _, assertion := range j.Assertions {
			if assertion.Target.Name == *tableName {
				assertionQuery = "-- Assertion \n" + assertion.Query + "; \n"
				outputQuery = OutputQuery{nonIncrementalQuery, incrementalPreOpsQuery, incrementalQuery, assertionQuery, operationsQuery}
			}
		}
	}

    for _, operation := range j.Operations{
        if operation.Target.Name == *tableName {
            operationsQuery = "-- Operation \n" + operation.Queries[0] + "; \n"
            outputQuery = OutputQuery{nonIncrementalQuery, incrementalPreOpsQuery, incrementalQuery, assertionQuery, operationsQuery}
        }

    }

	if outputQuery.Query == "" && outputQuery.IncrementalQuery == "" && outputQuery.Assertion == "" && outputQuery.OperationsQuery == ""{
		return &outputQuery, ErrorTableNotFound
	}
	return &outputQuery, nil
}

// Get cost of creating all tables in a Dataform project
func (j JsonStruct) GetDryRunForAllTables(keyfile string, projectId string, location string, includeAssertions bool, compact bool) {
	gitMetadata := j.GetGitMetadata(compact)
	runDateTime := time.Now()
	var wg sync.WaitGroup
	for i := 1; i <= len(j.Tables); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			table := j.Tables[i-1]
			cost, _, err := j.DryRun(&table.Target.Name, "", j.GetTargetGcpProjectId(), location, includeAssertions)

			targetMetadata := TargetMetadata{
				table.Target.Schema,
				table.Target.Name,
				table.Target.Database,
				table.Disabled,
			}

			tableMetadata := TableMetadata{
				targetMetadata,
				table.Query,
			}

			ProcessErrorAndOutput(tableMetadata, &gitMetadata, cost, err, true, runDateTime)
		}(i)
	}
	wg.Wait()
}

// Get cost of running a query associated with a table created using Dataform
func (j JsonStruct) DryRun(tableName *string, keyfile string, projectId string, location string, includeAssertions bool) (float32, *string, error) {
	outputQuery, err := j.GetQueryForTable(tableName, includeAssertions)
	dryRunQuery := ""
	if outputQuery.IncrementalQuery != "" {
		dryRunQuery = outputQuery.IncrementalPreOpsQuery + outputQuery.IncrementalQuery + outputQuery.Assertion
	} else {
		dryRunQuery = outputQuery.Query + outputQuery.Assertion + outputQuery.OperationsQuery
	}
	bytesProcessed, err := queryDryRun(os.Stdout, &dryRunQuery, projectId, keyfile, location)
	if err != nil {
		return 0.0, &dryRunQuery, err
	}
	return bytesProcessed, &dryRunQuery, nil
}

func (j JsonStruct) GetTargetGcpProjectId() string {
	return j.ProjectConfig.DefaultDatabase
}

func (j JsonStruct) GetGitRepository() string {
	return j.ProjectConfig.DefaultSchema
}

func (j JsonStruct) getGitBranchValue() string {
	var out bytes.Buffer
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(out.String())
}

func (j JsonStruct) GetGitBranch(compact bool) string {
	// NOTE: Would this error out if the git branch is not available ?
	switch compact {
	case true:
		return "main"
	default:
		return j.getGitBranchValue()
	}
}

func (j JsonStruct) GetUniqueTags() Tags {
	tags := []string{}

	for _, table := range j.Tables {
		for _, tag := range table.Tags {
			if slices.Contains(tags, tag) == false {
				tags = append(tags, tag)
			}
		}
	}

	for _, assertion := range j.Assertions {
		for _, tag := range assertion.Tags {
			if slices.Contains(tags, tag) == false {
				tags = append(tags, tag)
			}
		}
	}

	response := Tags{
		Tags: tags,
	}
	return response
}

func (j JsonStruct) ListTables() {
	tables := []string{}

	for _, table := range j.Tables {
		targetTable := table.Target.Name
		if slices.Contains(tables, targetTable) == false {
			tables = append(tables, targetTable)
		}
	}

	for _, assertion := range j.Assertions {
		targetTable := assertion.Target.Name
		if slices.Contains(tables, targetTable) == false {
			tables = append(tables, targetTable)
		}
	}

	JsonPrint(tables)
}

func (j JsonStruct) IsTableDisabled(tableName *string) (bool, error) {

	for _, table := range j.Tables {
		if table.Target.Name == *tableName {
			return table.Disabled, nil
		}
	}
	return false, ErrorTableNotFound
}
