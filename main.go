package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	inputPath := flag.String("i", "", "Path to the KiCad PCB file to process (required)")
	verbose := flag.Bool("v", false, "Enable verbose output")
	flag.Parse()

	if *inputPath == "" {
		fmt.Println("Error: -i flag is required")
		flag.Usage()
		os.Exit(1)
	}

	if *verbose {
		fmt.Println("Hello, Lin Router!")
	}

	_, err := ParsePcbFile(*inputPath)
	if err != nil {
		fmt.Println("Error parsing PCB file:", err)
		os.Exit(1)
	}
}
