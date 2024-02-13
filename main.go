package main

import (
	"eggo/parser"
	"eggo/repl"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: go run main.go <inputfile>\n")
	}

	filePath := os.Args[1]

	parser := parser.New(filePath)

	ast := parser.ParseBinaryOperation(0)

	fmt.Printf("%f\n", repl.InterpretAST(ast))

	// fmt.Printf("%v\n", ast)
}
