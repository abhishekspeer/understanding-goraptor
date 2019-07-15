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

	//   PARSE FILE method - Takes the file location as an input
	input := args[1]
	parserfile := rdf2v1.NewParser(input)
	defer parserfile.Free()

	// Parse the file using goraptor's ParseFile method and return a Statement.
	ch := parserfile.Rdfparser.ParseFile(input, "") // takes in input and baseuri
	for {
		statement, ok := <-ch
		fmt.Println(statement)
		if !ok {
			break
		}
		err := parserfile.ProcessTriple(statement)
		if err != nil {
			fmt.Println("Processing Failed")
		}
	}

	parserinstance := parserfile
	// fmt.Println(parserinstance.buffer)
	fmt.Println(len(parserinstance.Buffer))
}
