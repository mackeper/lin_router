package lexer

import (
	"errors"
	"fmt"
)

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\n' || ch == '\t' || ch == '\r'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

func readNextToken(data string, pos int) (Token, int, error) {
	for pos < len(data) {
		switch {
		case isWhitespace(data[pos]):
			pos++
			continue
		case data[pos] == '(':
			return Token{Type: OPEN_PAREN, Value: "("}, pos + 1, nil
		case data[pos] == ')':
			return Token{Type: CLOSE_PAREN, Value: ")"}, pos + 1, nil
		case data[pos] == '"':
			value := ""
			pos++
			for data[pos] != '"' {
				value += string(data[pos])
				pos++
			}
			return Token{Type: STRING, Value: value}, pos + 1, nil
		case isDigit(data[pos]) || (data[pos] == '-' && isDigit(data[pos+1])):
			value := ""
			for pos < len(data) && !isWhitespace(data[pos]) && data[pos] != ')' {
				value += string(data[pos])
				pos++
			}
			return Token{Type: NUMBER, Value: value}, pos, nil
		case isLetter(data[pos]):
			value := ""
			for pos < len(data) && (isLetter(data[pos]) || isDigit(data[pos]) || data[pos] == '-') {
				value += string(data[pos])
				pos++
			}
			return Token{Type: IDENTIFIER, Value: value}, pos, nil
		default:
			pos++
		}
	}
	return Token{}, pos + 1, errors.ErrUnsupported
}

func Tokenize(data string) ([]Token, error) {
	tokens := []Token{}
	pos := 0
	for pos < len(data) {
		token, new_pos, err := readNextToken(data, pos)
		if err != nil {
			return nil, fmt.Errorf("failed to read token at position %d: %w", pos, err)
		}
		fmt.Printf("Token: %s, Value: %s\n", token.Type.String(), token.Value)
		pos = new_pos
		tokens = append(tokens, token)
	}

	return tokens, nil
}
