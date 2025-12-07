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
