package main

import (
	"fmt"
	"github.com/mackeper/lin_router/lexer"
	"github.com/mackeper/lin_router/pcb"
)

func AddSegmentsToExpr(board *pcb.Board, expr *lexer.Expr) (lexer.Expr, error) {
	segmentExprs := []lexer.Expr{}
	for _, seg := range board.Segments {
		fmt.Printf("Add segment from (%f, %f) to (%f, %f) width %f layer %s\n",
			seg.Start.X, seg.Start.Y, seg.End.X, seg.End.Y, seg.Width, seg.Layer)
		segExpr := lexer.Expr{
			Type: lexer.ExprSegment,
			Identifier: "segment",
			Values: []lexer.Value{
				lexer.ExprValue{Value: lexer.Expr{
					Type: lexer.ExprStart,
					Identifier: "start",
					Values: []lexer.Value{
						lexer.NumberValue{Value: seg.Start.X},
						lexer.NumberValue{Value: seg.Start.Y},
					},
				}},
				lexer.ExprValue{Value: lexer.Expr{
					Type: lexer.ExprEnd,
					Identifier: "end",
					Values: []lexer.Value{
						lexer.NumberValue{Value: seg.End.X},
						lexer.NumberValue{Value: seg.End.Y},
					},
				}},
				lexer.ExprValue{Value: lexer.Expr{
					Type: lexer.ExprWidth,
					Identifier: "width",
					Values: []lexer.Value{
						lexer.NumberValue{Value: seg.Width},
					},
				}},
				lexer.ExprValue{Value: lexer.Expr{
					Type: lexer.ExprLayer,
					Identifier: "layer",
					Values: []lexer.Value{
						lexer.StringValue{Value: seg.Layer},
					},
				}},
			},
		}
		segmentExprs = append(segmentExprs, segExpr)
		fmt.Printf("Created segment expression: %s\n", segExpr.String())
	}

	for _, segExpr := range segmentExprs {
		expr.Values = append(expr.Values, lexer.ExprValue{Value: segExpr})
	}

	fmt.Printf("Added %d segments to expression\n", len(segmentExprs))
	return *expr, nil
}
