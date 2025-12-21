package main

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/mackeper/lin_router/lexer"
	"github.com/mackeper/lin_router/pcb"
)

type exprWithOffset struct {
	expr     lexer.Expr
	offset   pcb.Position
	rotation float64
}

func ExprToPCB(expr lexer.Expr) (*pcb.Board, error) {
	board := pcb.NewBoard()

	stack := []exprWithOffset{{expr: expr, offset: pcb.Position{X: 0, Y: 0}, rotation: 0}}
	pads := []pcb.Pad{}
	vias := []pcb.Via{}
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		slog.Debug("Processing expr", "type", current.expr.Type)
		if current.expr.Type == lexer.ExprPad {
			slog.Debug("Found pad expression")
			pad, err := parsePadExpr(current.expr, current.offset, current.rotation)
			if err != nil {
				return nil, fmt.Errorf("failed to parse pad: %w", err)
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
			// Check if this is a footprint and extract its position and rotation
			offset := current.offset
			rotation := current.rotation
			if current.expr.Type == lexer.ExprFootprint {
				footprintPos, footprintRot, err := extractAtPositionAndRotation(current.expr)
				if err != nil {
					return nil, fmt.Errorf("footprint missing position: %w", err)
				}
				offset = footprintPos
				rotation = footprintRot
				slog.Debug("Found footprint", "offset_x", offset.X, "offset_y", offset.Y, "rotation", rotation)
			}

			for _, val := range current.expr.Values {
				if v, ok := val.(lexer.ExprValue); ok {
					stack = append(stack, exprWithOffset{expr: v.Value, offset: offset, rotation: rotation})
				}
			}
		}
	}

	board.Pads = pads
	board.Vias = vias
	return board, nil
}

func extractAtPositionAndRotation(expr lexer.Expr) (pcb.Position, float64, error) {
	for _, val := range expr.Values {
		if v, ok := val.(lexer.ExprValue); ok {
			if v.Value.Type == lexer.ExprAt {
				if len(v.Value.Values) < 2 {
					return pcb.Position{}, 0, fmt.Errorf("at expression requires at least 2 values")
				}
				xVal, ok := v.Value.Values[0].(lexer.NumberValue)
				if !ok {
					return pcb.Position{}, 0, fmt.Errorf("expected NumberValue for X coordinate")
				}
				yVal, ok := v.Value.Values[1].(lexer.NumberValue)
				if !ok {
					return pcb.Position{}, 0, fmt.Errorf("expected NumberValue for Y coordinate")
				}

				rotation := 0.0
				if len(v.Value.Values) >= 3 {
					if rotVal, ok := v.Value.Values[2].(lexer.NumberValue); ok {
						rotation = rotVal.Value
					}
				}

				return pcb.Position{
					X: xVal.Value,
					Y: yVal.Value,
				}, rotation, nil
			}
		}
	}
	return pcb.Position{}, 0, fmt.Errorf("no at position found")
}

func parsePadExpr(expr lexer.Expr, offset pcb.Position, rotation float64) (pcb.Pad, error) {
	pad := pcb.Pad{}

	for _, val := range expr.Values {
		slog.Debug("Parsing pad sub-expression", "val", val)
		if v, ok := val.(lexer.ExprValue); ok {
			subExpr := v.Value
			slog.Debug("Pad sub-expr type", "type", subExpr.Type)
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

				// Apply rotation transformation
				rotatedX, rotatedY := rotatePoint(relX, relY, rotation)

				pad.Position = pcb.Position{
					X: rotatedX + offset.X,
					Y: rotatedY + offset.Y,
				}
				slog.Debug("Pad position", "rel_x", relX, "rel_y", relY, "rotation", rotation, "abs_x", pad.Position.X, "abs_y", pad.Position.Y)
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
					var layerStr string
					switch v := layerVal.(type) {
					case lexer.IdentifierValue:
						layerStr = v.Value
					case lexer.StringValue:
						layerStr = v.Value
					default:
						continue
					}

					if layerStr == "*.Cu" {
						pad.Layers = append(pad.Layers, "F.Cu", "B.Cu")
					} else {
						pad.Layers = append(pad.Layers, layerStr)
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

func rotatePoint(x, y, degrees float64) (float64, float64) {
	radians := -degrees * math.Pi / 180.0
	cos := math.Cos(radians)
	sin := math.Sin(radians)
	return x*cos - y*sin, x*sin + y*cos
}
