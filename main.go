package main

import (
	"fmt"
	"os"
	"spdx/tools-golang/v0/spdx"
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

	fmt.Println("===================================================\n")
	fmt.Println("Some Information Printed from the Document Returned\n")
	fmt.Println("===================================================\n")
	fmt.Printf("Relationship: %v\n\n", *(spdxdoc.Relationship[0]))
	fmt.Printf("Relationship[2]: %v\n\n", spdxdoc.Relationship[2])
	fmt.Printf("SpecVersion: %v\n\n", spdxdoc.SPDXVersion.Val)
	fmt.Printf("CreationInfo Creator: %v\n\n", spdxdoc.CreationInfo.Creator[0])
	fmt.Printf("CreationInfo Create:%v\n\n", spdxdoc.CreationInfo.Create)
	fmt.Printf("DocumentName: %v\n\n", spdxdoc.DocumentName.Val)
	fmt.Printf("DocumentComment: %v\n\n", spdxdoc.DocumentComment)

	spec := spdxdoc.SPDXVersion
	dl := spdxdoc.DataLicense
	ci := spdxdoc.CreationInfo
	ci2v1 := spdx.CreationInfo2_1{

		CreatorComment:  ci.Comment.Val,
		Created:         ci.Create.V(),
		SPDXVersion:     spec.Val,
		DataLicense:     dl.Val,
		DocumentName:    spdxdoc.DocumentName.Val,
		DocumentComment: spdxdoc.DocumentComment.Val,
	}
	fmt.Printf("Creationinfo2v1: %v\n\n", ci2v1)

}
func Parse(input string) (*rdf2v1.Document, error) {
	parser := rdf2v1.NewParser(input)
	defer fmt.Println("RDF Doc PARSED")
	defer parser.Free()
	return parser.Parse()
}
