package main

import (
	"errors"
	"fmt"
	"github.com/mackeper/lin_router/lexer"
	"os"
)

// PCB represents a parsed KiCad PCB file.
type PCB struct {
	// Add relevant fields here
}

func readToString(path string) (string, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}

func ParsePcbFile(path string) (string, error) {
	data, err := readToString(path)
	if err != nil {
		return "", errors.New("failed to read PCB file: " + err.Error())
	}

	tokens, err := lexer.Tokenize(data)
	if err != nil {
		return "", errors.New("failed to tokenize PCB file: " + err.Error())
	}
	for _, token := range tokens {
		fmt.Println(token.Type.String(), ":", token.Value)
	}

	// Simulate successful parsing
	fmt.Println("Successfully read PCB file:", path)
	return data, nil
}
