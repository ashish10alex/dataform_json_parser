package cmd

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ashish10alex/dj/cmd"
)


func TestJsonStruct_GetUniqueTags(t *testing.T) {

    jsonFile := "../../internal/test/compiles.json"
    jsonData, err := readJsonFile(jsonFile)
    if err != nil {
        fmt.Println(err.Error())
    }

	tests := []struct {
		name     string
		input    cmd.JsonStruct
		expected cmd.Tags
	}{
		{
			name: "Unique tags example one",
			input: jsonData,
			expected: cmd.Tags{Tags: []string{"TAG_1", "FOO", "TAG_2", "BAR"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.GetUniqueTags()
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetUniqueTags() = %v, want %v", result, tt.expected)
			}
		})
	}
}
