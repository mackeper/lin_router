package main

import (
	"github.com/mackeper/lin_router/lexer"
	"github.com/mackeper/lin_router/pcb"
	"testing"
)

func TestExprToPCB_SimplePad(t *testing.T) {
	// Arrange
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

	// Act
	board, err := ExprToPCB(expr)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	validateBoard(t, board, 1)
	validatePad(t, board.Pads[0], 10.0, 20.0, 1, "GND", "F.Cu")
}

func TestExprToPCB_MultiplePads(t *testing.T) {
	// Arrange
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
			lexer.ExprValue{
				Value: lexer.Expr{
					Type:       lexer.ExprUnknown,
					Identifier: "some_other_thing",
					Values:     []lexer.Value{},
				},
			},
			lexer.ExprValue{
				Value: lexer.Expr{
					Type:       lexer.ExprUnknown,
					Identifier: "another_thing",
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
												lexer.NumberValue{Value: 11.0},
												lexer.NumberValue{Value: 21.1},
											},
										},
									},
									lexer.ExprValue{
										Value: lexer.Expr{
											Type:       lexer.ExprNet,
											Identifier: "net",
											Values: []lexer.Value{
												lexer.NumberValue{Value: 2},
												lexer.StringValue{Value: "RAW"},
											},
										},
									},
									lexer.ExprValue{
										Value: lexer.Expr{
											Type:       lexer.ExprLayer,
											Identifier: "layer",
											Values: []lexer.Value{
												lexer.StringValue{Value: "*.Cu *.Mask"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// Act
	board, err := ExprToPCB(expr)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	validateBoard(t, board, 2)
	validatePad(t, board.Pads[0], 11.0, 21.1, 2, "RAW", "*.Cu *.Mask")
	validatePad(t, board.Pads[1], 10.0, 20.0, 1, "GND", "F.Cu")
}

func validateBoard(t *testing.T, board *pcb.Board, expectedPadCount int) {
	if len(board.Pads) != expectedPadCount {
		t.Errorf("Expected %d pads, got %d", expectedPadCount, len(board.Pads))
	}
}

func validatePad(t *testing.T, pad pcb.Pad, expectedX, expectedY float64, expectedNetNumber int, expectedNetName, expectedLayers string) {
	if pad.Position.X != expectedX || pad.Position.Y != expectedY {
		t.Errorf("Expected pad position (%f, %f), got (%f, %f)", expectedX, expectedY, pad.Position.X, pad.Position.Y)
	}
	if pad.Net.Number != expectedNetNumber || pad.Net.Name != expectedNetName {
		t.Errorf("Expected pad net (%d, %s), got (%d, %s)", expectedNetNumber, expectedNetName, pad.Net.Number, pad.Net.Name)
	}
	if len(pad.Layers) == 0 || pad.Layers[0] != expectedLayers {
		t.Errorf("Expected pad layers %s, got %v", expectedLayers, pad.Layers)
	}
}

func TestExprToPCB_FootprintWithPads(t *testing.T) {
	// Arrange
	expr := lexer.Expr{
		Type: lexer.ExprKicadPcb,
		Values: []lexer.Value{
			lexer.ExprValue{
				Value: lexer.Expr{
					Type:       lexer.ExprFootprint,
					Identifier: "footprint",
					Values: []lexer.Value{
						lexer.StringValue{Value: "TestLib:TestFootprint"},
						lexer.ExprValue{
							Value: lexer.Expr{
								Type:       lexer.ExprAt,
								Identifier: "at",
								Values: []lexer.Value{
									lexer.NumberValue{Value: 100.0},
									lexer.NumberValue{Value: 200.0},
								},
							},
						},
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
												lexer.NumberValue{Value: 5.0},
												lexer.NumberValue{Value: 10.0},
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
				},
			},
		},
	}

	// Act
	board, err := ExprToPCB(expr)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	validateBoard(t, board, 1)
	// Pad at relative (5, 10) + footprint at (100, 200) = absolute (105, 210)
	validatePad(t, board.Pads[0], 105.0, 210.0, 1, "GND", "F.Cu")
}

func TestExprToPCB_FootprintMissingPosition(t *testing.T) {
	// Arrange
	expr := lexer.Expr{
		Type: lexer.ExprKicadPcb,
		Values: []lexer.Value{
			lexer.ExprValue{
				Value: lexer.Expr{
					Type:       lexer.ExprFootprint,
					Identifier: "footprint",
					Values: []lexer.Value{
						lexer.StringValue{Value: "TestLib:TestFootprint"},
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
												lexer.NumberValue{Value: 5.0},
												lexer.NumberValue{Value: 10.0},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// Act
	board, err := ExprToPCB(expr)

	// Assert
	if err == nil {
		t.Fatalf("Expected error for footprint missing position, got nil")
	}
	if board != nil {
		t.Errorf("Expected nil board on error, got %v", board)
	}
}
