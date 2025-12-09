package lexer

import (
	"testing"
)

func TestReadNextToken(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		pos      int
		wantType TokenType
		wantVal  string
	}{
		// Pos 0
		{"open paren", "(", 0, OPEN_PAREN, "("},
		{"identifier", "hello ", 0, IDENTIFIER, "hello"},
		{"number", "123 ", 0, NUMBER, "123"},
		{"negative number", "-45.67 ", 0, NUMBER, "-45.67"},
		{"string", `"test string"`, 0, STRING, "test string"},
		{"close paren", ")", 0, CLOSE_PAREN, ")"},
		// Pos > 0
		{"open paren,pos 2", "(h (", 2, OPEN_PAREN, "("},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, _, err := readNextToken(tt.input, tt.pos)
			if err != nil {
				t.Fatal(err)
			}
			if token.Type != tt.wantType {
				t.Errorf("got type %v, want %v", token.Type, tt.wantType)
			}
			if token.Value != tt.wantVal {
				t.Errorf("got value %q, want %q", token.Value, tt.wantVal)
			}
		})
	}
}
