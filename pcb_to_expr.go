package main

import (
	"log/slog"

	"github.com/mackeper/lin_router/lexer"
	"github.com/mackeper/lin_router/pcb"
)

func AddSegmentsToExpr(board *pcb.Board, expr *lexer.Expr) (lexer.Expr, error) {
	segmentExprs := []lexer.Expr{}
	for _, seg := range board.Segments {
		slog.Debug("Add segment",
			"start_x", seg.Start.X, "start_y", seg.Start.Y,
			"end_x", seg.End.X, "end_y", seg.End.Y,
			"width", seg.Width, "layer", seg.Layer)
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
				lexer.ExprValue{Value: lexer.Expr{
					Type: lexer.ExprNet,
					Identifier: "net",
					Values: []lexer.Value{
						lexer.NumberValue{Value: float64(seg.Net)},
					},
				}},
				lexer.ExprValue{Value: lexer.Expr{
					Type: lexer.ExprUUID,
					Identifier: "uuid",
					Values: []lexer.Value{
						lexer.StringValue{Value: seg.UUID},
					},
				}},
			},
		}
		segmentExprs = append(segmentExprs, segExpr)
		slog.Debug("Created segment expression", "expr", segExpr.String())
	}

	for _, segExpr := range segmentExprs {
		expr.Values = append(expr.Values, lexer.ExprValue{Value: segExpr})
	}

	slog.Debug("Added segments to expression", "count", len(segmentExprs))
	return *expr, nil
}
