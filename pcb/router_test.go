package pcb

import (
	"testing"
)

func TestRouteBoard_TwoPadsSameNetCloseEnough(t *testing.T) {
	board := NewBoard()
	board.AddPad(Pad{
		Position: Position{0, 0},
		Net:      Net{Number: 1, Name: "VCC"},
		Layer:    "F.Cu",
	})
	board.AddPad(Pad{
		Position: Position{2, 0},
		Net:      Net{Number: 1, Name: "VCC"},
		Layer:    "F.Cu",
	})

	RouteBoard(board)

	if len(board.Segments) != 1 {
		t.Errorf("Expected 1 segment, got %d", len(board.Segments))
	}
	if board.Segments[0].Net != 1 {
		t.Errorf("Expected segment on net 1, got %d", board.Segments[0].Net)
	}
	if board.Segments[0].Layer != "F.Cu" {
		t.Errorf("Expected segment on F.Cu, got %s", board.Segments[0].Layer)
	}
}

func TestRouteBoard_TwoPadsTooFarApart(t *testing.T) {
	board := NewBoard()
	board.AddPad(Pad{
		Position: Position{0, 0},
		Net:      Net{Number: 1, Name: "VCC"},
		Layer:    "F.Cu",
	})
	board.AddPad(Pad{
		Position: Position{5, 0},
		Net:      Net{Number: 1, Name: "VCC"},
		Layer:    "F.Cu",
	})

	RouteBoard(board)

	if len(board.Segments) != 0 {
		t.Errorf("Expected 0 segments, got %d", len(board.Segments))
	}
}

func TestRouteBoard_TwoPadsDifferentLayers(t *testing.T) {
	board := NewBoard()
	board.AddPad(Pad{
		Position: Position{0, 0},
		Net:      Net{Number: 1, Name: "VCC"},
		Layer:    "F.Cu",
	})
	board.AddPad(Pad{
		Position: Position{1, 0},
		Net:      Net{Number: 1, Name: "VCC"},
		Layer:    "B.Cu",
	})

	RouteBoard(board)

	if len(board.Segments) != 0 {
		t.Errorf("Expected 0 segments for different layers, got %d", len(board.Segments))
	}
}

func TestRouteBoard_ThreePadsFormingTriangle(t *testing.T) {
	board := NewBoard()
	board.AddPad(Pad{
		Position: Position{0, 0},
		Net:      Net{Number: 1, Name: "VCC"},
		Layer:    "F.Cu",
	})
	board.AddPad(Pad{
		Position: Position{2, 0},
		Net:      Net{Number: 1, Name: "VCC"},
		Layer:    "F.Cu",
	})
	board.AddPad(Pad{
		Position: Position{1, 1},
		Net:      Net{Number: 1, Name: "VCC"},
		Layer:    "F.Cu",
	})

	RouteBoard(board)

	if len(board.Segments) != 3 {
		t.Errorf("Expected 3 segments for triangle, got %d", len(board.Segments))
	}
}

func TestRouteBoard_DifferentNets(t *testing.T) {
	board := NewBoard()
	board.AddPad(Pad{
		Position: Position{0, 0},
		Net:      Net{Number: 1, Name: "VCC"},
		Layer:    "F.Cu",
	})
	board.AddPad(Pad{
		Position: Position{1, 0},
		Net:      Net{Number: 2, Name: "GND"},
		Layer:    "F.Cu",
	})

	RouteBoard(board)

	if len(board.Segments) != 0 {
		t.Errorf("Expected 0 segments for different nets, got %d", len(board.Segments))
	}
}

func TestRouteBoard_EmptyBoard(t *testing.T) {
	board := NewBoard()

	RouteBoard(board)

	if len(board.Segments) != 0 {
		t.Errorf("Expected 0 segments for empty board, got %d", len(board.Segments))
	}
}

func TestRouteBoard_SinglePad(t *testing.T) {
	board := NewBoard()
	board.AddPad(Pad{
		Position: Position{0, 0},
		Net:      Net{Number: 1, Name: "VCC"},
		Layer:    "F.Cu",
	})

	RouteBoard(board)

	if len(board.Segments) != 0 {
		t.Errorf("Expected 0 segments for single pad, got %d", len(board.Segments))
	}
}

func TestRouteBoard_MultipleNets(t *testing.T) {
	board := NewBoard()
	board.AddPad(Pad{Position: Position{0, 0}, Net: Net{Number: 1, Name: "VCC"}, Layer: "F.Cu"})
	board.AddPad(Pad{Position: Position{1, 0}, Net: Net{Number: 1, Name: "VCC"}, Layer: "F.Cu"})
	board.AddPad(Pad{Position: Position{10, 10}, Net: Net{Number: 2, Name: "GND"}, Layer: "F.Cu"})
	board.AddPad(Pad{Position: Position{11, 10}, Net: Net{Number: 2, Name: "GND"}, Layer: "F.Cu"})

	RouteBoard(board)

	if len(board.Segments) != 2 {
		t.Errorf("Expected 2 segments (one per net), got %d", len(board.Segments))
	}

	net1Count := 0
	net2Count := 0
	for _, seg := range board.Segments {
		if seg.Net == 1 {
			net1Count++
		}
		if seg.Net == 2 {
			net2Count++
		}
	}

	if net1Count != 1 || net2Count != 1 {
		t.Errorf("Expected 1 segment per net, got %d for net 1 and %d for net 2", net1Count, net2Count)
	}
}
