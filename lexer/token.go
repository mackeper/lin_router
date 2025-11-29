package lexer

type TokenType int

const (
	OPEN_PAREN = iota
	CLOSE_PAREN
	IDENTIFIER
	NUMBER
	STRING
)

func (t TokenType) String() string {
	return [...]string{
		"OPEN_PAREN",
		"CLOSE_PAREN",
		"IDENTIFIER",
		"NUMBER",
		"STRING",
	}[t]
}

type Token struct {
	Type  TokenType
	Value string
}

