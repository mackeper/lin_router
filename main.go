package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/mackeper/lin_router/pcb"
)

func main() {
	inputPath := flag.String("i", "", "Path to the KiCad PCB file to process (required)")
	verbose := flag.Bool("v", false, "Enable verbose output")
	maxDistance := flag.Float64("max-distance", 3.0, "Maximum routing distance in mm")
	flag.Parse()

	// Setup logging
	level := slog.LevelError
	if *verbose {
		level = slog.LevelDebug
	}
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	if *inputPath == "" {
		slog.Error("-i flag is required")
		flag.Usage()
		os.Exit(1)
	}

	slog.Debug("Parsing PCB file", "path", *inputPath)
	expr, err := ParsePcbFile(*inputPath)
	if err != nil {
		slog.Error("Error parsing PCB file", "error", err)
		os.Exit(1)
	}

	slog.Debug("Converting expression to PCB structure")
	board, err := ExprToPCB(expr)
	if err != nil {
		slog.Error("Error converting expression to PCB", "error", err)
		os.Exit(1)
	}

	slog.Debug("Adding trivial segments to PCB", "max_distance", *maxDistance)
	pcb.AddTrivialSegments(board, *maxDistance)

	slog.Debug("Converting PCB structure back to expression")
	expr, err = AddSegmentsToExpr(board, &expr)
	if err != nil {
		slog.Error("Error converting PCB back to expression", "error", err)
		os.Exit(1)
	}

	fmt.Println(expr.String())
}
