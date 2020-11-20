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

func save(filePath string, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		log.Panicf("Failed to serialize config: %v", err)
	}
	err = ioutil.WriteFile(filePath, b, 0644)
	if err != nil {
		log.Panicf("Failed to write config: %v", err)
	}
}
