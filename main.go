package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mackeper/lin_router/pcb"
)

func verbosePrint(verbose bool, message string) {
	if verbose {
		fmt.Println(message)
	}
}

func main() {
	inputPath := flag.String("i", "", "Path to the KiCad PCB file to process (required)")
	verbose := flag.Bool("v", false, "Enable verbose output")
	flag.Parse()

	if *inputPath == "" {
		fmt.Println("Error: -i flag is required")
		flag.Usage()
		os.Exit(1)
	}

	verbosePrint(*verbose, "Parsing PCB file: "+*inputPath)
	expr, err := ParsePcbFile(*inputPath)
	if err != nil {
		fmt.Println("Error parsing PCB file:", err)
		os.Exit(1)
	}

	verbosePrint(*verbose, "Converting expression to PCB structure")
	board, err := ExprToPCB(expr)
	if err != nil {
		fmt.Println("Error converting expression to PCB:", err)
		os.Exit(1)
	}

	verbosePrint(*verbose, "Adding trivial segments to PCB")
	pcb.AddTrivialSegments(board)

	verbosePrint(*verbose, "Converting PCB structure back to expression")
	AddSegmentsToExpr(board, &expr)

	fmt.Printf("Writing modified PCB to stdout\n")
	fmt.Println(expr.String())
}
