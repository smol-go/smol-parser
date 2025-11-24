package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		printUsage()
		return
	}
}

func printUsage() {
	fmt.Println("smol-parser - A JSON lexer and parser")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  smol-parser parse <file.json>     - Parse and validate JSON")
	fmt.Println("  smol-parser lex <file.json>       - Tokenize JSON (show tokens)")
	fmt.Println("  smol-parser pretty <file.json>    - Pretty print JSON")
	fmt.Println("  smol-parser compact <file.json>   - Compact print JSON")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  smol-parser parse example.json")
	fmt.Println("  smol-parser pretty example.json")
}
