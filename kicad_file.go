package main

import (
	"errors"
	"github.com/mackeper/lin_router/lexer"
	"os"
)

func readToString(path string) (string, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}

func ParsePcbFile(path string) (lexer.Expr, error) {
	data, err := readToString(path)
	if err != nil {
		return lexer.Expr{}, errors.New("failed to read PCB file: " + err.Error())
	}

	tokens, err := lexer.Tokenize(data)
	if err != nil {
		return lexer.Expr{}, errors.New("failed to tokenize PCB file: " + err.Error())
	}

	expr, err := lexer.Parse(tokens)
	if err != nil {
		return lexer.Expr{}, errors.New("failed to parse PCB file: " + err.Error())
	}

	return expr, nil
}
