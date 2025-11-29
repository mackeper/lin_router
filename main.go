package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, Lin Router!")
	_, err := ParsePcbFile("main.kicad_pcb")
	if err != nil {
		fmt.Println("Error parsing PCB file:", err)
	} else {
		// fmt.Println("Parsed PCB:", PCB)
	}
}
