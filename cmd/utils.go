package cmd

import (
	"encoding/json"
	"io"
	"os"
)

func getJsonIoReaderForFile(fileName string) io.Reader {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	return file
}

func ReadJson(fileName string) (JsonStruct, error) {
	var fileReader io.Reader

	if fileName == "" {
		fileReader = os.Stdin
	} else {
		fileReader = getJsonIoReaderForFile(fileName)
	}

	var jsonData JsonStruct
	err := json.NewDecoder(fileReader).Decode(&jsonData)

	if err != nil {
		return JsonStruct{}, err
	}

	return jsonData, nil
}

func JsonPrint(j interface{}) {
	jsonBytes, _ := json.Marshal(j)
	os.Stdout.Write(jsonBytes)
}

func JsonPrintTagReponse(j interface{}) {
	for _, tagReponse := range *j.(*[]TagReponse) {
		tagReponsesOutput, _ := json.Marshal(tagReponse)
		os.Stdout.Write(tagReponsesOutput)
		os.Stdout.Write([]byte("\n"))
	}
}
