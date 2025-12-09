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
	vias := []pcb.Via{}
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		slog.Debug("Processing expr", "type", current.expr.Type)
		if current.expr.Type == lexer.ExprPad {
			slog.Debug("Found pad expression")
			pad, err := parsePadExpr(current.expr, current.offset)
			if err != nil {
				return nil, err
			}
			pads = append(pads, pad)
		} else if current.expr.Type == lexer.ExprVia {
			slog.Debug("Found via expression")
			via, err := parseViaExpr(current.expr)
			if err != nil {
				return nil, err
			}
			vias = append(vias, via)
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
	board.Vias = vias
	return board, nil
}

func extractAtPosition(expr lexer.Expr) (pcb.Position, error) {
	for _, val := range expr.Values {
		if v, ok := val.(lexer.ExprValue); ok {
			if v.Value.Type == lexer.ExprAt {
				return pcb.Position{
					X: v.Value.Values[0].(lexer.NumberValue).Value,
					Y: v.Value.Values[1].(lexer.NumberValue).Value,
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
		switch v := val.(type) {
		case lexer.ExprValue:
			subExpr := v.Value
			switch subExpr.Type {
			case lexer.ExprAt:
				relX := subExpr.Values[0].(lexer.NumberValue).Value
				relY := subExpr.Values[1].(lexer.NumberValue).Value
				pad.Position = pcb.Position{
					X: relX + offset.X,
					Y: relY + offset.Y,
				}
				slog.Debug("Pad position", "rel_x", relX, "rel_y", relY, "abs_x", pad.Position.X, "abs_y", pad.Position.Y)
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

func parseViaExpr(expr lexer.Expr) (pcb.Via, error) {
	via := pcb.Via{}

	for _, val := range expr.Values {
		slog.Debug("Parsing via sub-expression", "val", val)
		switch v := val.(type) {
		case lexer.ExprValue:
			subExpr := v.Value
			switch subExpr.Type {
			case lexer.ExprAt:
				via.Position = pcb.Position{
					X: subExpr.Values[0].(lexer.NumberValue).Value,
					Y: subExpr.Values[1].(lexer.NumberValue).Value,
				}
			case lexer.ExprNet:
				via.Net = int(subExpr.Values[0].(lexer.NumberValue).Value)
			case lexer.ExprLayers:
				for _, layerVal := range subExpr.Values {
					if strVal, ok := layerVal.(lexer.StringValue); ok {
						via.Layers = append(via.Layers, strVal.Value)
					}
				}
			case lexer.ExprUUID:
				via.UUID = subExpr.Values[0].(lexer.StringValue).Value
			}
		}
	}
	return via, nil
}
