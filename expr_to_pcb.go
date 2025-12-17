package main

import (
	"fmt"
	"log/slog"

	"github.com/mackeper/lin_router/lexer"
	"github.com/mackeper/lin_router/pcb"
)

type exprWithOffset struct {
	expr   lexer.Expr
	offset pcb.Position
}

func ExprToPCB(expr lexer.Expr) (*pcb.Board, error) {
	board := pcb.NewBoard()

	stack := []exprWithOffset{{expr: expr, offset: pcb.Position{X: 0, Y: 0}}}
	pads := []pcb.Pad{}
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		slog.Debug("Processing expr", "type", current.expr.Type)
		if current.expr.Type == lexer.ExprPad {
			slog.Debug("Found pad expression")
			pad, err := parsePadExpr(current.expr, current.offset)
			if err != nil {
				return nil, fmt.Errorf("failed to parse pad: %w", err)
			}
			pads = append(pads, pad)
		} else {
			// Check if this is a footprint and extract its position
			offset := current.offset
			if current.expr.Type == lexer.ExprFootprint {
				footprintPos, err := extractAtPosition(current.expr)
				if err != nil {
					return nil, fmt.Errorf("footprint missing position: %w", err)
				}
				offset = footprintPos
				slog.Debug("Found footprint", "offset_x", offset.X, "offset_y", offset.Y)
			}

			for _, val := range current.expr.Values {
				if v, ok := val.(lexer.ExprValue); ok {
					stack = append(stack, exprWithOffset{expr: v.Value, offset: offset})
				}
			}
		}
	}

	board.Pads = pads
	return board, nil
}

func extractAtPosition(expr lexer.Expr) (pcb.Position, error) {
	for _, val := range expr.Values {
		if v, ok := val.(lexer.ExprValue); ok {
			if v.Value.Type == lexer.ExprAt {
				if len(v.Value.Values) < 2 {
					return pcb.Position{}, fmt.Errorf("at expression requires 2 values")
				}
				xVal, ok := v.Value.Values[0].(lexer.NumberValue)
				if !ok {
					return pcb.Position{}, fmt.Errorf("expected NumberValue for X coordinate")
				}
				yVal, ok := v.Value.Values[1].(lexer.NumberValue)
				if !ok {
					return pcb.Position{}, fmt.Errorf("expected NumberValue for Y coordinate")
				}
				return pcb.Position{
					X: xVal.Value,
					Y: yVal.Value,
				}, nil
			}
		}
	}
	return pcb.Position{}, fmt.Errorf("no at position found")
}

func parsePadExpr(expr lexer.Expr, offset pcb.Position) (pcb.Pad, error) {
	pad := pcb.Pad{}

	for _, val := range expr.Values {
		slog.Debug("Parsing pad sub-expression", "val", val)
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
				relX := xVal.Value
				relY := yVal.Value
				pad.Position = pcb.Position{
					X: relX + offset.X,
					Y: relY + offset.Y,
				}
				slog.Debug("Pad position", "rel_x", relX, "rel_y", relY, "abs_x", pad.Position.X, "abs_y", pad.Position.Y)
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
				layer := layerVal.Value
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
