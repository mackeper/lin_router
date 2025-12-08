package lexer

import "testing"

func TestIdentifierToExprType(t *testing.T) {
	tests := []struct {
		identifier string
		expected   ExprType
	}{
		{"kicad_pcb", ExprKicadPcb},
		{"footprint", ExprFootprint},
		{"pad", ExprPad},
		{"segment", ExprSegment},
		{"via", ExprVia},
		{"net", ExprNet},
		{"layer", ExprLayer},
		{"gr_line", ExprGrLine},
		{"gr_arc", ExprGrArc},
		{"at", ExprAt},
		{"unknown_type", ExprUnknown},
		{"", ExprUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.identifier, func(t *testing.T) {
			result := IdentifierToExprType(tt.identifier)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestExprTypeString(t *testing.T) {
	tests := []struct {
		exprType ExprType
		expected string
	}{
		{ExprKicadPcb, "kicad_pcb"},
		{ExprFootprint, "footprint"},
		{ExprPad, "pad"},
		{ExprSegment, "segment"},
		{ExprVia, "via"},
		{ExprNet, "net"},
		{ExprLayer, "layer"},
		{ExprGrLine, "gr_line"},
		{ExprGrArc, "gr_arc"},
		{ExprAt, "at"},
		{ExprUnknown, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.exprType.String()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestParseWithExprType(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedType ExprType
		expectedId   string
	}{
		{
			"pad expression",
			"(pad 1 2)",
			ExprPad,
			"pad",
		},
		{
			"via expression",
			"(via)",
			ExprVia,
			"via",
		},
		{
			"footprint expression",
			"(footprint \"R_0805\")",
			ExprFootprint,
			"footprint",
		},
		{
			"unknown expression",
			"(custom_thing 123)",
			ExprUnknown,
			"custom_thing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := Tokenize(tt.input)
			if err != nil {
				t.Fatalf("Tokenize failed: %v", err)
			}
			expr, err := Parse(tokens)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if expr.Type != tt.expectedType {
				t.Errorf("Expected type %v, got %v", tt.expectedType, expr.Type)
			}
			if expr.Identifier != tt.expectedId {
				t.Errorf("Expected identifier %s, got %s", tt.expectedId, expr.Identifier)
			}
		})
	}
}
