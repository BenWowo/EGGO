package main

import (
	"eggo/scanner"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: go run main.go <inputfile>\n")
	}

	filePath := os.Args[1]

	scanner := scanner.New(filePath)

	scanner.ScanFile()
}
