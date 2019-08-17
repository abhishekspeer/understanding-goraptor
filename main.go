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
	var sp *rdf2v1.Snippet
	input := args[1]
	spdxdoc, sp, err := Parse2_1(input)

	if err != nil {
		fmt.Println("Parsing Error")
		return
	}
	fmt.Printf("ORIGINAL:\n %v\n", spdxdoc.Relationship[2].Package[0])

	_ = sp
	doc2v1 := rdf2v1.TransferDocument(spdxdoc, sp)
	fmt.Printf("\n\n%#v\n\n", doc2v1.OtherLicenses)
	new2v1 := rdf2v1.CollectDocument(doc2v1)
	fmt.Printf("NEW:\n %v\n", new2v1.Relationship[0].Package[0])

	// WRITER
	// output := os.Stdout
	// errdoc := rdf2v1.WriteDocument(output, spdxdoc, sp)
	// if errdoc != nil {
	// 	fmt.Println("Cannot Write Document")
	// 	return
	// }
}
func Parse2_1(input string) (*rdf2v1.Document, *rdf2v1.Snippet, error) {
	parser := rdf2v1.NewParser(input)
	defer fmt.Printf("Successfully loaded %s\n\n", input)
	defer parser.Free()
	return parser.Parse()
}
