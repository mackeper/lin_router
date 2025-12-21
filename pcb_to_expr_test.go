package main

import (
	"testing"

	"github.com/mackeper/lin_router/lexer"
	"github.com/mackeper/lin_router/pcb"
)

// TODO: Assert the contents of the segments added to the expression
// do this by adding a segment and checking that its start, end, width, and layer

func TestAddSegmentsToExpr_NoSegments(t *testing.T) {
	// Arrange
	expr := lexer.Expr{
		Type:   lexer.ExprUnknown,
		Values: []lexer.Value{},
	}

	// Act
	board, err := AddSegmentsToExpr(pcb.NewBoard(), &expr)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(board.Values) != 0 {
		t.Fatalf("Expected 0 segments in expression, got %d", len(board.Values))
	}
}

func TestAddSegmentsToExpr_OneSegment(t *testing.T) {
	// Arrange
	expr := lexer.Expr{
		Type:   lexer.ExprUnknown,
		Values: []lexer.Value{},
	}

	board := pcb.NewBoard()
	seg := pcb.Segment{
		Start: pcb.Position{X: 0, Y: 0},
		End:   pcb.Position{X: 10, Y: 0},
		Width: 1,
		Layer: "F.Cu",
	}
	board.Segments = append(board.Segments, seg)

	// Act
	resultExpr, err := AddSegmentsToExpr(board, &expr)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(resultExpr.Values) != 1 {
		t.Fatalf("Expected 1 segments in expression, got %d", len(resultExpr.Values))
	}
	if resultExpr.Values[0].(lexer.ExprValue).Value.Type != lexer.ExprSegment {
		t.Fatalf("Expected segment expression type, got %v", resultExpr.Values[0].(lexer.ExprValue).Value.Type)
	}

	validateSegmentExpr(t, seg, resultExpr.Values[0].(lexer.ExprValue).Value)
}

func TestAddSegmentsToExpr_MultipleSegments(t *testing.T) {
	// Arrange
	expr := lexer.Expr{
		Type:   lexer.ExprUnknown,
		Values: []lexer.Value{},
	}

	board := pcb.NewBoard()
	board.Segments = append(board.Segments,
		pcb.Segment{
			Start: pcb.Position{X: 10, Y: 0},
			End:   pcb.Position{X: 10, Y: 10},
			Width: 1,
		},
		pcb.Segment{
			Start: pcb.Position{X: 10, Y: 10},
			End:   pcb.Position{X: 0, Y: 10},
			Width: 1,
		},
		pcb.Segment{
			Start: pcb.Position{X: 10, Y: 10},
			End:   pcb.Position{X: 5, Y: 10},
			Width: 1,
		},
	)

	// Act
	resultExpr, err := AddSegmentsToExpr(board, &expr)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(resultExpr.Values) != 3 {
		t.Fatalf("Expected 3 segments in expression, got %d", len(resultExpr.Values))
	}
}

func findExprValueByType(values []lexer.Value, exprType lexer.ExprType) (lexer.ExprValue, bool) {
	for _, v := range values {
		ev, ok := v.(lexer.ExprValue)
		if !ok {
			continue
		}
		if ev.Value.Type == exprType {
			return ev, true
		}
	}
	return lexer.ExprValue{}, false
}

func validateSegmentExpr(t *testing.T, expected pcb.Segment, actual lexer.Expr) {
	startValue, found := findExprValueByType(actual.Values, lexer.ExprStart)
	if !found {
		t.Fatalf("Expected start position value in segment expression")
	}
	endValue, found := findExprValueByType(actual.Values, lexer.ExprEnd)
	if !found {
		t.Fatalf("Expected end position value in segment expression")
	}
	widthValue, found := findExprValueByType(actual.Values, lexer.ExprWidth)
	if !found {
		t.Fatalf("Expected width value in segment expression")
	}

	startX := startValue.Value.Values[0].(lexer.NumberValue).Value
	startY := startValue.Value.Values[1].(lexer.NumberValue).Value
	if startX != expected.Start.X || startY != expected.Start.Y {
		t.Fatalf("Expected start position (%v,%v), got (%v,%v)",
			expected.Start.X, expected.Start.Y,
			startX, startY)
	}

	endX := endValue.Value.Values[0].(lexer.NumberValue).Value
	endY := endValue.Value.Values[1].(lexer.NumberValue).Value
	if endX != expected.End.X || endY != expected.End.Y {
		t.Fatalf("Expected end position (%v,%v), got (%v,%v)",
			expected.End.X, expected.End.Y,
			endX, endY)
	}

	width := widthValue.Value.Values[0].(lexer.NumberValue).Value
	if width != expected.Width {
		t.Fatalf("Expected width %v, got %v",
			expected.Width,
			width)
	}
}
