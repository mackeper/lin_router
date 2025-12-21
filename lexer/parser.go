package lexer

// grammar
// expr : = '(' IDENTIFIER [value] ')'
// value : = STRING | NUMBER | IDENTIFIER | expr

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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

type IdentifierValue struct {
	Value string
}

func (ExprValue) isValue() {}
func (v ExprValue) String() string {
	return v.Value.String()
}
func (v ExprValue) stringWithIndent(indent int) string {
	return v.Value.stringWithIndent(indent)
}

func (StringValue) isValue() {}
func (v StringValue) String() string {
	return fmt.Sprintf("%q", v.Value)
}

func (NumberValue) isValue() {}
func (v NumberValue) String() string {
	// Format integers without decimals
	if v.Value == float64(int(v.Value)) {
		return fmt.Sprintf("%d", int(v.Value))
	}
	return fmt.Sprintf("%f", v.Value)
}

func (IdentifierValue) isValue() {}
func (v IdentifierValue) String() string {
	return v.Value
}

type Expr struct {
	Type       ExprType
	Identifier string
	Values     []Value
}

func (e Expr) String() string {
	return e.stringWithIndent(0)
}

func (e Expr) stringWithIndent(indent int) string {
	// Check if this is kicad_pcb root or has many nested children
	isRootOrComplex := e.Identifier == "kicad_pcb" || e.Identifier == "module" || e.Identifier == "footprint"

	if !isRootOrComplex {
		// Simple expression - keep on one line
		result := fmt.Sprintf("(%s", e.Identifier)
		for _, val := range e.Values {
			result += " " + val.String()
		}
		result += ")"
		return result
	}

	// Complex expression - format with newlines
	result := fmt.Sprintf("(%s", e.Identifier)
	for _, val := range e.Values {
		if exprVal, ok := val.(ExprValue); ok {
			result += "\n" + strings.Repeat("  ", indent+1) + exprVal.stringWithIndent(indent+1)
		} else {
			result += " " + val.String()
		}
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
		return parseExprError(fmt.Errorf("Expr expected '(', got %s", tokens[pos].Type.String()))
	}
	pos++

	if tokens[pos].Type != IDENTIFIER {
		// return parseExprError(fmt.Errorf("Expr expected IDENTIFIER, got %s. Context: %v", tokens[pos].Type.String(), tokens[pos]))
		// Layers can look like this (34 "B.Paste" user)
		tokens[pos].Type = IDENTIFIER
	}
	identifier = tokens[pos].Value
	pos++

	for pos < len(tokens) && tokens[pos].Type != CLOSE_PAREN {
		switch tokens[pos].Type {
		case OPEN_PAREN:
			expr, newPos, err := parseExpr(tokens, pos)
			if err != nil {
				return parseExprError(fmt.Errorf("failed to parse nested expression: %w", err))
			}
			pos = newPos
			values = append(values, ExprValue{Value: expr})
		case STRING:
			values = append(values, StringValue{Value: tokens[pos].Value})
			pos++
		case NUMBER:
			value, err := strconv.ParseFloat(tokens[pos].Value, 64)
			if err != nil {
				return parseExprError(fmt.Errorf("failed to parse number: %w", err))
			}
			values = append(values, NumberValue{Value: value})
			pos++
		case IDENTIFIER:
			values = append(values, IdentifierValue{Value: tokens[pos].Value})
			pos++
		default:
			return parseExprError(fmt.Errorf("unexpected token %s", tokens[pos].Type.String()))
		}
	}

	if pos >= len(tokens) {
		return parseExprError(errors.New("unexpected end of tokens"))
	}

	return Expr{
		Type:       IdentifierToExprType(identifier),
		Identifier: identifier,
		Values:     values,
	}, pos + 1, nil
}

func Parse(tokens []Token) (Expr, error) {
	expr, _, err := parseExpr(tokens, 0)
	if err != nil {
		return Expr{}, err
	}

	return expr, nil
}
