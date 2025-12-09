package lexer

import (
	"testing"
)

func (e ExprValue) Equal(other ExprValue) bool {
	if e.Value.Identifier != other.Value.Identifier {
		return false
	}
	if len(e.Value.Values) != len(other.Value.Values) {
		return false
	}
	for i := range e.Value.Values {
		switch v := e.Value.Values[i].(type) {
		case ExprValue:
			ov, ok := other.Value.Values[i].(ExprValue)
			if !ok || !v.Equal(ov) {
				return false
			}
		case StringValue:
			ov, ok := other.Value.Values[i].(StringValue)
			if !ok || v.Value != ov.Value {
				return false
			}
		case NumberValue:
			ov, ok := other.Value.Values[i].(NumberValue)
			if !ok || v.Value != ov.Value {
				return false
			}
		}
	}
	return true
}

func TestParseExpr(t *testing.T) {
	tests := []struct {
		name     string
		input    []Token
		expected Expr
	}{
		{
			"Simple expr",
			[]Token{
				{Type: OPEN_PAREN, Value: "("},
				{Type: IDENTIFIER, Value: "hej"},
				{Type: CLOSE_PAREN, Value: ")"},
			},
			Expr{Identifier: "hej", Values: []Value{}},
		},
		{
			"expr with number",
			[]Token{
				{Type: OPEN_PAREN, Value: "("},
				{Type: IDENTIFIER, Value: "hej"},
				{Type: NUMBER, Value: "123.45"},
				{Type: CLOSE_PAREN, Value: ")"},
			},
			Expr{Identifier: "hej", Values: []Value{
				NumberValue{Value: 123.45},
			}},
		},
		{
			"expr with number",
			[]Token{
				{Type: OPEN_PAREN, Value: "("},
				{Type: IDENTIFIER, Value: "hej"},
				{Type: STRING, Value: "tjenis1337"},
				{Type: CLOSE_PAREN, Value: ")"},
			},
			Expr{Identifier: "hej", Values: []Value{
				StringValue{Value: "tjenis1337"},
			}},
		},
		{
			"complex expr",
			[]Token{
				{Type: OPEN_PAREN, Value: "("},
				{Type: IDENTIFIER, Value: "hej"},
				{Type: STRING, Value: "tjenis1337"},
				{Type: OPEN_PAREN, Value: "("},
				{Type: IDENTIFIER, Value: "nested"},
				{Type: NUMBER, Value: "42"},
				{Type: CLOSE_PAREN, Value: ")"},
				{Type: CLOSE_PAREN, Value: ")"},
			},
			Expr{Identifier: "hej", Values: []Value{
				StringValue{Value: "tjenis1337"},
				ExprValue{Value: Expr{
					Identifier: "nested",
					Values: []Value{
						NumberValue{Value: 42},
					},
				}},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := Parse(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			if expr.Identifier != tt.expected.Identifier {
				t.Errorf("Expected identifier: %s, Actual identifier: %s", expr.Identifier, tt.expected.Identifier)
			}
			if len(expr.Values) != len(tt.expected.Values) {
				t.Errorf("Expected values length: %d, Actual values length: %d", len(tt.expected.Values), len(expr.Values))
			}
			for i, v := range expr.Values {
				if ev, ok := v.(ExprValue); ok {
					ov, ok := tt.expected.Values[i].(ExprValue)
					if !ok || !ev.Equal(ov) {
						t.Errorf("Expected ExprValue: %v, Actual ExprValue: %v", ov, ev)
					}
					continue
				}
				if v != tt.expected.Values[i] {
					t.Errorf("Expected value: %v, Actual value: %v", tt.expected.Values[i], v)
				}
			}
		})
	}
}

func TestValueString(t *testing.T) {
	tests := []struct {
		name     string
		value    Value
		expected string
	}{
		{
			"StringValue",
			StringValue{Value: "hello"},
			"\"hello\"",
		},
		{
			"NumberValue",
			NumberValue{Value: 3.14},
			"3.140000",
		},
		{
			"ExprValue",
			ExprValue{Value: Expr{
				Identifier: "greet",
				Values: []Value{
					StringValue{Value: "world"},
					NumberValue{Value: 42},
				},
			}},
			"(greet \"world\" 42.000000)",
		},
		{
			"Nested ExprValue",
			ExprValue{Value: Expr{
				Identifier: "outer",
				Values: []Value{
					ExprValue{Value: Expr{
						Identifier: "inner",
						Values: []Value{
							StringValue{Value: "nested"},
						},
					}},
					NumberValue{Value: 123},
				},
			}},
			"(outer (inner \"nested\") 123.000000)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.value.String()
			if result != tt.expected {
				t.Errorf("Expected: %s, Actual: %s", tt.expected, result)
			}
		})
	}
}

func TestTokenizeCount(t *testing.T) {
	input := "(add 1 2.5)"
	tokens, err := Tokenize(input)
	if err != nil {
		t.Fatalf("Tokenize failed: %v", err)
	}
	t.Logf("Token count: %d", len(tokens))
	for i, tok := range tokens {
		t.Logf("%d: %s = '%s'", i, tok.Type.String(), tok.Value)
	}
}

func TestRoundtrip(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			"Simple expression",
			"(hello \"world\")",
		},
		{
			"Expression with number",
			"(add 1 2.5)",
		},
		{
			"Nested expression",
			"(outer (inner \"value\" 42) \"text\")",
		},
		{
			"KiCad-like structure",
			"(kicad_pcb (version 20240108) (general (thickness 1.6)))",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the original input
			tokens1, err := Tokenize(tt.input)
			if err != nil {
				t.Fatalf("Tokenize failed: %v", err)
			}
			expr1, err := Parse(tokens1)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			// Convert back to string
			result := expr1.String()

			// Parse the result
			tokens2, err := Tokenize(result)
			if err != nil {
				t.Fatalf("Tokenize result failed: %v", err)
			}
			expr2, err := Parse(tokens2)
			if err != nil {
				t.Fatalf("Parse result failed: %v", err)
			}

			// Compare the two Exprs
			if expr1.Identifier != expr2.Identifier {
				t.Errorf("Identifier mismatch: %s != %s", expr1.Identifier, expr2.Identifier)
			}
			if len(expr1.Values) != len(expr2.Values) {
				t.Errorf("Values length mismatch: %d != %d", len(expr1.Values), len(expr2.Values))
			}

			// Deep comparison of values
			for i := range expr1.Values {
				if ev1, ok := expr1.Values[i].(ExprValue); ok {
					ev2, ok := expr2.Values[i].(ExprValue)
					if !ok || !ev1.Equal(ev2) {
						t.Errorf("ExprValue mismatch at index %d", i)
					}
				} else if expr1.Values[i] != expr2.Values[i] {
					t.Errorf("Value mismatch at index %d: %v != %v", i, expr1.Values[i], expr2.Values[i])
				}
			}
		})
	}
}
