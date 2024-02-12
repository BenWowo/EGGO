package main

import (
	"eggo/parser"
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

	ast := parser.ParseBinaryExpression()

	fmt.Printf("%+v", ast)
}
