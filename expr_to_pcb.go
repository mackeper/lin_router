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
			pad, _ := parsePadExpr(current)
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
		switch v := val.(type) {
		case lexer.ExprValue:
			subExpr := v.Value
			switch subExpr.Type {
			case lexer.ExprAt:
				pad.Position = pcb.Position{
					X: subExpr.Values[0].(lexer.NumberValue).Value,
					Y: subExpr.Values[1].(lexer.NumberValue).Value,
				}
			case lexer.ExprNet:
				pad.Net = pcb.Net{
					Number: int(subExpr.Values[0].(lexer.NumberValue).Value),
					Name:   subExpr.Values[1].(lexer.StringValue).Value,
				}
			case lexer.ExprLayer:
				pad.Layer = subExpr.Values[0].(lexer.StringValue).Value
			}
		}
	}
	return pad, nil
}
		
