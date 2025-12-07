package lexer

// grammar
// expr : = '(' IDENTIFIER [value] ')'
// value : = STRING | NUMBER | expr

import (
	"errors"
	"fmt"
	"strconv"
)

type Value interface {
	isValue()
	String() string
}

type ExprValue struct {
	Value Expr
}

type StringValue struct {
	Value string
}

type NumberValue struct {
	Value float64
}

func (ExprValue) isValue()   {}
func (v ExprValue) String() string {
	for _, val := range v.Value.Values {
		_ = val.String()
	}
	return fmt.Sprintf("(%s ...)", v.Value.Identifier)
}

func (StringValue) isValue() {}
func (v StringValue) String() string {
	return fmt.Sprintf("\"%s\"", v.Value)
}

func (NumberValue) isValue()    {}
func (v NumberValue) String() string {
	return fmt.Sprintf("%f", v.Value)
}

type Expr struct {
	Identifier string
	Values      []Value
}
func (e Expr) String() string {
	result := fmt.Sprintf("(%s", e.Identifier)
	for _, val := range e.Values {
		result += " " + val.String()
	}
	result += ")"
	return result
}

func parseExprError(err error) (Expr, int, error) {
	return Expr{}, 0, err
}

func parseExpr(tokens []Token, pos int) (Expr, int, error) {
	identifier := ""
	values := []Value{}

	if tokens[pos].Type != OPEN_PAREN {
		return parseExprError(errors.ErrUnsupported)
	}
	pos++

	if tokens[pos].Type != IDENTIFIER {
		return parseExprError(errors.ErrUnsupported)
	}
	identifier = tokens[pos].Value
	pos++

	for tokens[pos].Type != CLOSE_PAREN {
		switch tokens[pos].Type {
		case OPEN_PAREN:
			expr, newPos, err := parseExpr(tokens, pos)
			if err != nil { }
			pos = newPos
			values = append(values, ExprValue{Value: expr})
		case STRING:
			values = append(values, StringValue{Value: tokens[pos].Value})
			pos++
		case NUMBER:
			value, err := strconv.ParseFloat(tokens[pos].Value, 64)
			if err != nil { }
			values = append(values, NumberValue{Value: value})
			pos++
		}
	}

	return Expr{Identifier: identifier, Values: values}, pos, nil
}

func Parse(tokens []Token) (Expr, error) {
	expr, _, err := parseExpr(tokens, 0)
	if err != nil {
		return Expr{}, errors.ErrUnsupported
	}

	return expr, nil
}
