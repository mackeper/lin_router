package main

import (
	"github.com/mackeper/lin_router/lexer"
	"testing"
)

func TestExprToPCB_SimplePad(t *testing.T) {
	expr := lexer.Expr{
		Type: lexer.ExprUnknown,
		Values: []lexer.Value{
			lexer.ExprValue{
				Value: lexer.Expr{
					Type:       lexer.ExprPad,
					Identifier: "pad",
					Values: []lexer.Value{
						lexer.ExprValue{
							Value: lexer.Expr{
								Type:       lexer.ExprAt,
								Identifier: "at",
								Values: []lexer.Value{
									lexer.NumberValue{Value: 10.0},
									lexer.NumberValue{Value: 20.0},
								},
							},
						},
						lexer.ExprValue{
							Value: lexer.Expr{
								Type:       lexer.ExprNet,
								Identifier: "net",
								Values: []lexer.Value{
									lexer.NumberValue{Value: 1},
									lexer.StringValue{Value: "GND"},
								},
							},
						},
						lexer.ExprValue{
							Value: lexer.Expr{
								Type:       lexer.ExprLayer,
								Identifier: "layer",
								Values: []lexer.Value{
									lexer.StringValue{Value: "F.Cu"},
								},
							},
						},
					},
				},
			},
		},
	}

	board, err := ExprToPCB(expr)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(board.Pads) != 1 {
		t.Fatalf("Expected 1 pad, got %d", len(board.Pads))
	}

	pad := board.Pads[0]
	if pad.Position.X != 10.0 || pad.Position.Y != 20.0 {
		t.Errorf("Expected pad position (10.0, 20.0), got (%f, %f)", pad.Position.X, pad.Position.Y)
	}
	if pad.Net.Number != 1 || pad.Net.Name != "GND" {
		t.Errorf("Expected pad net (1, GND), got (%d, %s)", pad.Net.Number, pad.Net.Name)
	}
	if pad.Layer != "F.Cu" {
		t.Errorf("Expected pad layer F.Cu, got %s", pad.Layer)
	}
}
