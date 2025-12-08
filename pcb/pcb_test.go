package pcb

import (
	"math"
	"testing"
)

func TestPositionDistance(t *testing.T) {
	tests := []struct {
		name     string
		p1       Position
		p2       Position
		expected float64
	}{
		{
			"same position",
			Position{0, 0},
			Position{0, 0},
			0,
		},
		{
			"horizontal distance",
			Position{0, 0},
			Position{3, 0},
			3,
		},
		{
			"vertical distance",
			Position{0, 0},
			Position{0, 4},
			4,
		},
		{
			"diagonal distance (3-4-5 triangle)",
			Position{0, 0},
			Position{3, 4},
			5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.p1.Distance(tt.p2)
			if math.Abs(result-tt.expected) > 0.0001 {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestPadDistance(t *testing.T) {
	pad1 := Pad{
		Position: Position{0, 0},
		Net:      Net{Number: 1, Name: "GND"},
		Number:   "1",
	}
	pad2 := Pad{
		Position: Position{3, 4},
		Net:      Net{Number: 1, Name: "GND"},
		Number:   "2",
	}

	dist := pad1.Distance(pad2)
	if math.Abs(dist-5.0) > 0.0001 {
		t.Errorf("Expected 5.0, got %f", dist)
	}
}

func TestBoardAddPad(t *testing.T) {
	board := NewBoard()
	pad := Pad{
		Position: Position{10, 20},
		Net:      Net{Number: 1, Name: "VCC"},
		Number:   "1",
	}

	board.AddPad(pad)

	if len(board.Pads) != 1 {
		t.Errorf("Expected 1 pad, got %d", len(board.Pads))
	}
	if board.Pads[0].Position.X != 10 {
		t.Errorf("Expected X=10, got %f", board.Pads[0].Position.X)
	}
}

func TestBoardAddSegment(t *testing.T) {
	board := NewBoard()
	seg := Segment{
		Start: Position{0, 0},
		End:   Position{10, 10},
		Width: 0.25,
		Layer: "F.Cu",
		Net:   1,
	}

	board.AddSegment(seg)

	if len(board.Segments) != 1 {
		t.Errorf("Expected 1 segment, got %d", len(board.Segments))
	}
}

func TestBoardGetPadsByNet(t *testing.T) {
	board := NewBoard()
	board.AddPad(Pad{Position: Position{0, 0}, Net: Net{Number: 1, Name: "VCC"}})
	board.AddPad(Pad{Position: Position{10, 0}, Net: Net{Number: 1, Name: "VCC"}})
	board.AddPad(Pad{Position: Position{20, 0}, Net: Net{Number: 2, Name: "GND"}})

	net1Pads := board.GetPadsByNet(1)
	if len(net1Pads) != 2 {
		t.Errorf("Expected 2 pads on net 1, got %d", len(net1Pads))
	}

	net2Pads := board.GetPadsByNet(2)
	if len(net2Pads) != 1 {
		t.Errorf("Expected 1 pad on net 2, got %d", len(net2Pads))
	}

	net3Pads := board.GetPadsByNet(3)
	if len(net3Pads) != 0 {
		t.Errorf("Expected 0 pads on net 3, got %d", len(net3Pads))
	}
}

func TestSegmentLength(t *testing.T) {
	seg := Segment{
		Start: Position{0, 0},
		End:   Position{3, 4},
	}

	length := seg.Length()
	if math.Abs(length-5.0) > 0.0001 {
		t.Errorf("Expected 5.0, got %f", length)
	}
}
