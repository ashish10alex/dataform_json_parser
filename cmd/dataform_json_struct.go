package cmd

import "time"

type JsonStruct struct {
	Tables        []Tables       `json:"tables"`
	Assertions    []Assertions   `json:"assertions"`
	Declarations  []Declarations `json:"declarations"`
	Operations    []Operations   `json:"operations"`
	ProjectConfig ProjectConfig  `json:"projectConfig"`
	GitMetadata   GitMetadata
	Targets       []Target `json:"targets"`
}

type ProjectConfig struct {
	Warehouse       string `json:"warehouse"`
	DefaultSchema   string `json:"defaultSchema"`
	AssertionSchema string `json:"assertionSchema"`
	DefaultDatabase string `json:"defaultDatabase"`
	DefaultLocation string `json:"defaultLocation"`
	TablePrefix     string `json:"tablePrefix"`
}

type Operations struct {
	FileName string   `json:"fileName"`
	Queries  []string `json:"queries"`
	Target   Target   `json:"target"`
	Tags     []string `json:"tags"`
}

type Tables struct {
	Query             string   `json:"query"`
	IncrementalQuery  string   `json:"incrementalQuery"`
	IncrementalPreOps []string `json:"incrementalPreOps"`
	FileName          string   `json:"fileName"`
	Type              string   `json:"type"`
	Target            Target   `json:"target"`
	Tags              []string `json:"tags"`
	Disabled          bool     `json:"disabled"`
}

type TableMetadata struct {
	TargetMetadata TargetMetadata
	Query          string
}

type Assertions struct {
	Query    string   `json:"query"`
	FileName string   `json:"fileName"`
	Target   Target   `json:"target"`
	Tags     []string `json:"tags"`
}

type Declarations struct {
	FileName string `json:"fileName"`
	Target   Target `json:"target"`
}

type OutputSources struct {
	Declarations []string
	Targets      []string
}

type OutputQuery struct {
	Query                  string
	IncrementalPreOpsQuery string
	IncrementalQuery       string
	Assertion              string
	OperationsQuery        string
}

type Target struct {
	Schema   string `json:"schema"`
	Name     string `json:"name"`
	Database string `json:"database"`
}

type TargetMetadata struct {
	Schema   string
	Name     string
	Database string
	Disabled bool
}

type Targets struct {
	TableTargets     []TargetMetadata
	AssertionTargets []TargetMetadata
	OperationTargets []TargetMetadata
}

type parentAction struct {
	Schema   string `json:"schema"`
	Name     string `json:"name"`
	Database string `json:"database"`
}

type DataformJson interface {
	getTargetGcpProjectId() string
	getGitRepository() string
	dryRun(*string) (float32, error) // dryRun(tableName *string) (float32, error)
	getUniqueTags() []string
	getQueryForTable(string) string
	getCostForEachTag() map[string]float32
	getTablesInTag(string) []string
}

type ErrorDetails struct {
	IsError      bool
	LineNumber   int
	ColumnNumber int
	ErrorMsg     string
	Disabled     bool
}

type TablesDryRunMetadata struct {
	Table string
	Error ErrorDetails
}

type Tags struct {
	Tags []string `json:"tags"`
}

type TagReponse struct {
	Tag                  string
	GBProcessed          float32
	Cost                 float32
	HasError             bool
	RunDateTime          time.Time
	GitMetadata          GitMetadata
	TablesDryRunMetadata []TablesDryRunMetadata
}

type GitMetadata struct {
	GitRepositoryId string
	GitBranch       string
}

type DryRunReponse struct {
	FileName    string
	Schema      string
	Database    string
	Cost        float32
	GBProcessed float32
	Error       ErrorDetails
	GitMetadata GitMetadata
	RunDateTime time.Time
	Compact     bool
	Query       string
}
