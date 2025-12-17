package main

import (
	"fmt"

	"github.com/mackeper/lin_router/lexer"
	"github.com/mackeper/lin_router/pcb"
)

func ExprToPCB(expr lexer.Expr) (*pcb.Board, error) {
	board := pcb.NewBoard()

	stack := []lexer.Expr{expr}
	pads := []pcb.Pad{}
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		fmt.Printf("Processing expr of type: %v\n", current.Type)
		if current.Type == lexer.ExprPad {
			pad, err := parsePadExpr(current)
			if err != nil {
				return nil, fmt.Errorf("failed to parse pad: %w", err)
			}
			pads = append(pads, pad)
		} else {
			for _, val := range current.Values {
				if v, ok := val.(lexer.ExprValue); ok {
					stack = append(stack, v.Value)
				}
			}
		}
	}

	board.Pads = pads
	return board, nil
}

func parsePadExpr(expr lexer.Expr) (pcb.Pad, error) {
	pad := pcb.Pad{}

	for _, val := range expr.Values {
		if v, ok := val.(lexer.ExprValue); ok {
			subExpr := v.Value
			switch subExpr.Type {
			case lexer.ExprAt:
				if len(subExpr.Values) < 2 {
					return pad, fmt.Errorf("at expression requires 2 values")
				}
				xVal, ok := subExpr.Values[0].(lexer.NumberValue)
				if !ok {
					return pad, fmt.Errorf("expected NumberValue for X coordinate")
				}
				yVal, ok := subExpr.Values[1].(lexer.NumberValue)
				if !ok {
					return pad, fmt.Errorf("expected NumberValue for Y coordinate")
				}
				pad.Position = pcb.Position{
					X: xVal.Value,
					Y: yVal.Value,
				}
			case lexer.ExprNet:
				if len(subExpr.Values) < 2 {
					return pad, fmt.Errorf("net expression requires 2 values")
				}
				numVal, ok := subExpr.Values[0].(lexer.NumberValue)
				if !ok {
					return pad, fmt.Errorf("expected NumberValue for net number")
				}
				nameVal, ok := subExpr.Values[1].(lexer.StringValue)
				if !ok {
					return pad, fmt.Errorf("expected StringValue for net name")
				}
				pad.Net = pcb.Net{
					Number: int(numVal.Value),
					Name:   nameVal.Value,
				}
			case lexer.ExprLayer:
				if len(subExpr.Values) < 1 {
					return pad, fmt.Errorf("layer expression requires 1 value")
				}
				layerVal, ok := subExpr.Values[0].(lexer.StringValue)
				if !ok {
					return pad, fmt.Errorf("expected StringValue for layer")
				}
				pad.Layer = layerVal.Value
			}
		}
	}
	return pad, nil
}
