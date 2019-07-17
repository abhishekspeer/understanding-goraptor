package main

import (
	"fmt"
	"os"
	"ug/understanding-goraptor/rdf2v1"
)

func main() {
	// check that we've received the right number of arguments

	args := os.Args
	if len(args) != 2 {
		fmt.Printf("Usage: %v <spdx-file-in>\n", args[0])
		fmt.Printf("  Load SPDX 2.1 RDF file <spdx-file-in>, and\n")
		fmt.Printf("  print its contents.\n")
		return
	}
	var spdxdoc *rdf2v1.Document
	var err error

	input := args[1]
	spdxdoc, err = Parse(input)

	if err != nil {
		fmt.Println("Parsing Error")
		return
	}

	fmt.Printf("%#v\n\n", spdxdoc.CreationInfo)
}

func Parse(input string) (*rdf2v1.Document, error) {
	parser := rdf2v1.NewParser(input)
	defer fmt.Println("RDF Doc PARSED")
	defer parser.Free()
	return parser.Parse()
}
