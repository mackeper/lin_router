package main

import (
	"log/slog"

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

		slog.Debug("Processing expr", "type", current.Type)
		if current.Type == lexer.ExprPad {
			slog.Debug("Found pad expression")
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
		slog.Debug("Parsing pad sub-expression", "val", val)
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
				layer := subExpr.Values[0].(lexer.StringValue).Value
				if layer == "*.Cu" {
					pad.Layers = []string{"F.Cu", "B.Cu"}
				} else {
					pad.Layers = []string{layer}
				}
			case lexer.ExprLayers:
				for _, layerVal := range subExpr.Values {
					if strVal, ok := layerVal.(lexer.StringValue); ok {
						if strVal.Value == "*.Cu" {
							pad.Layers = append(pad.Layers, "F.Cu", "B.Cu")
						} else {
							pad.Layers = append(pad.Layers, strVal.Value)
						}
					}
				}
			}
		}
	}
	return pad, nil
}
