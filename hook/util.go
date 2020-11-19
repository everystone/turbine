package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func load(filePath string, out interface{}) {
	var data []byte

	var err error
	data, err = ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file %s %v", filePath, err)
	}

	if err := json.Unmarshal(data, &out); err != nil {
		log.Fatalf("Failed to parse %s: %v", filePath, err)
	}
}
