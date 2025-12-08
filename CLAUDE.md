# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

lin_router is a PCB router for KiCad files written in Go. It parses KiCad PCB files (`.kicad_pcb` format) and implements simple linear routing algorithms to connect pads on the same net.

## Build and Test Commands

```bash
# Build the project
make build
# Output: bin/lin_router

# Run tests
make test

# Run a single test
go test -v -run TestName ./path/to/package

# Run all tests in a specific package
go test -v ./lexer
go test -v ./pcb

# Format code
make fmt

# Vet code
make vet

# Clean build artifacts
make clean

# Run without building
make run

# Build and run
make build-run
```

## Architecture

### Three-Stage Pipeline

1. **Lexer** (`lexer/lexer.go`) - Tokenizes KiCad's S-expression format into tokens (parentheses, identifiers, numbers, strings)
2. **Parser** (`lexer/parser.go`) - Transforms tokens into a hierarchical expression tree (Expr)
3. **PCB Domain** (`pcb/`) - Converts expressions into PCB objects and performs routing

### Key Components

**Lexer Package** (`lexer/`)
- `lexer.go` - Tokenizer for KiCad S-expression syntax
- `parser.go` - Recursive descent parser that builds expression trees
- `token.go` - Token type definitions (OPEN_PAREN, CLOSE_PAREN, IDENTIFIER, NUMBER, STRING)
- `expr_type.go` - Maps KiCad identifiers (footprint, pad, segment, net, etc.) to typed enums
- Grammar: `expr := '(' IDENTIFIER [value]* ')'` where `value := STRING | NUMBER | expr`

**PCB Package** (`pcb/`)
- `pcb.go` - Board model with Pads and Segments collections
- `router.go` - Simple linear routing algorithm that connects pads within MaxRoutingDistance (3.0mm) on the same layer
- `pad.go`, `segment.go`, `net.go` - Domain models for PCB elements
- Routing constants: `MaxRoutingDistance = 3.0`, `DefaultTraceWidth = 0.2`

**Main Entry** (`main.go`, `kicad_file.go`)
- `main.go` - Entry point, calls ParsePcbFile
- `kicad_file.go` - Coordinates the lexer/parser to read KiCad files

### Data Flow

```
.kicad_pcb file → Tokenize() → []Token → Parse() → Expr tree → (future: convert to PCB model) → AddTrivialSegments() → output
```

Currently the parser outputs the expression tree. The PCB domain models exist but aren't yet connected to the parser output.

## Important Patterns

### Expression Value Types
The parser uses an interface-based value system where `Value` can be `StringValue`, `NumberValue`, or `ExprValue` (nested expression). All implement `isValue()` and `String()`.

### Routing Algorithm
`AddTrivialSegments()` in `pcb/router.go` implements a simple all-pairs algorithm:
- For each net, get all pads on that net
- For each pair of pads, if they're on the same layer and within MaxRoutingDistance, create a segment connecting them
- This is a naive approach suitable for simple linear traces

## Testing Philosophy

Tests are located alongside source files with `_test.go` suffix. Key test files:
- `lexer/lexer_test.go` - Tokenization tests
- `lexer/parser_test.go` - Parser tests for expression tree building
- `lexer/expr_type_test.go` - Expression type mapping tests
- `pcb/pcb_test.go` - Board model tests
- `pcb/router_test.go` - Routing algorithm tests
