package main

import (
	"eggo/gen"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("usage: go run main.go <inputfilePath> <outputfilePath>\n")
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	gen.GenerateLLVM(inputPath, outputPath)
}
