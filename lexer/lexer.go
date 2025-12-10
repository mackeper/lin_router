package lexer

import (
	"fmt"
	"regexp"
)

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\n' || ch == '\t' || ch == '\r'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isPositiveDigit(ch byte) bool {
	return '1' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

var numberRegex = regexp.MustCompile(`^-?[0-9]+(\.[0-9]+)?$`)

func isValidNumber(s string) bool {
	return numberRegex.MatchString(s)
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
		case isPositiveDigit(data[pos]) || // positive number, e.g. 42
			(data[pos] == '-' && pos+1 < len(data) && isDigit(data[pos+1])) || // negative number, e.g. -42
			(data[pos] == '0' && pos+1 < len(data) && data[pos+1] != 'x'): // E.g. 0.5 or just 0
			value := ""
			for pos < len(data) && !isWhitespace(data[pos]) && data[pos] != ')' {
				value += string(data[pos])
				pos++
			}
			if isValidNumber(value) {
				return Token{Type: NUMBER, Value: value}, pos, nil
			}
			return Token{Type: IDENTIFIER, Value: value}, pos, nil
		case isLetter(data[pos]) ||
			(data[pos] == '0' && pos+1 < len(data) && data[pos+1] == 'x'): // handle hex identifiers like 0x1A2B
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
	return Token{}, pos + 1, fmt.Errorf("no more tokens at position %d", pos)
}

func Tokenize(data string) ([]Token, error) {
	tokens := []Token{}
	pos := 0
	for pos < len(data) {
		token, new_pos, err := readNextToken(data, pos)
		if err != nil {
			break
		}
		pos = new_pos
		tokens = append(tokens, token)
	}

	return tokens, nil
}
