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
	// var sp *rdf2v1.Snippet
	var err error

	input := args[1]
	spdxdoc, sp, err := Parse(input)

	if err != nil {
		fmt.Println("Parsing Error")
		return
	}

	fmt.Println("===================================================\n")
	fmt.Println("Some Information Printed from the Document Returned\n")
	fmt.Println("===================================================\n")
	fmt.Println("%#v\n\n", sp)
	// fmt.Printf("Relationship: %v\n\n", spdxdoc.Relationship[0].Package[0])
	fmt.Printf("\nRelationship: %v\n\n", spdxdoc.Relationship[3].File)
	// fmt.Printf("Relationship: %v\n\n", spdxdoc.Relationship[2])
	// fmt.Printf("Relationship: %v\n\n", spdxdoc.Relationship[3])
	// fmt.Printf("SpecVersion: %v\n\n", spdxdoc.SPDXVersion.Val)
	fmt.Printf("\n\nCreationInfo Creator: %v\n\n", spdxdoc.CreationInfo)
	fmt.Printf("CreationInfo Create:%v\n\n", spdxdoc.CreationInfo.Create)
	fmt.Printf("DocumentName: %v\n\n", spdxdoc.DocumentName.Val)
	fmt.Printf("DocumentComment: %v\n\n", spdxdoc.DocumentComment)

	// iniEdr := spdxdoc.ExternalDocumentRef
	// intRel := spdxdoc.Relationship
	stdRel := make([]spdx.Relationship2_1, 4)
	fmt.Printf("Docummment: %v\n\n", stdRel[0])

	// for i, s := range intRel {
	// 	stdRel[i] = s.RelationshipType.
	// 	fmt.Println("%#v",stdRel[i])

	// }
	// intRel := spdxdoc.Relationship

	// stdRel := spdx.Relationship2_1{
	// 	Relationship: intRel.RelationshipType.Val,
	// 	// RelationshipComment:
	// }

	// FINAL TRANSLATED DOCUMENT
	doc2v1 := transferDocument(spdxdoc)
	fmt.Printf("%T", doc2v1)

}

func Parse(input string) (*rdf2v1.Document, *rdf2v1.Snippet, error) {
	parser := rdf2v1.NewParser(input)
	defer fmt.Println("RDF Document parsed successfully.\n")
	defer parser.Free()
	return parser.Parse()
}

func transferDocument(spdxdoc *rdf2v1.Document) *spdx.Document2_1 {

	stdDoc := spdx.Document2_1{

		// CreationInfo: transferCreationInfo(spdxdoc),
		// Packages:      transferPackages(spdxdoc),
		// OtherLicenses: transferOtherLicenses(spdxdoc),
		// Relationships: transferRelationships(spdxdoc),
		// Annotations:   transferAnnotation(spdxdoc),
		// Reviews:       transferReview(spdxdoc),
	}
	return &stdDoc
}
