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
	// var spdxdoc *rdf2v1.Document
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
	// fmt.Printf("Relationship: %#v\n\n", spdxdoc.Relationship[3].File[0])
	// fmt.Printf("Relationship: %v\n\n", spdxdoc.Relationship[2])
	// fmt.Printf("Relationship: %v\n\n", spdxdoc.Relationship[3])
	// fmt.Printf("SpecVersion: %v\n\n", spdxdoc.SPDXVersion.Val)
	fmt.Printf("\n\nCreationInfo Creator: %v\n\n", spdxdoc.CreationInfo)
	fmt.Printf("CreationInfo Create:%v\n\n", spdxdoc.CreationInfo.Create)
	fmt.Printf("DocumentName: %v\n\n", spdxdoc.DocumentName.Val)
	fmt.Printf("DocumentComment: %v\n\n", spdxdoc.DocumentComment)

	iniSpec := spdxdoc.SPDXVersion
	intDi := spdxdoc.DataLicense
	intCi := spdxdoc.CreationInfo
	intAnn := spdxdoc.Annotation
	iniEdr := spdxdoc.ExternalDocumentRef
	// intRel := spdxdoc.Relationship
	stdRel := make([]spdx.Relationship2_1, 4)
	fmt.Printf("Docummment: %v\n\n", stdRel[0])

	// for i, s := range intRel {
	// 	stdRel[i] = s.RelationshipType.
	// 	fmt.Println("%#v",stdRel[i])

	// }
	// intRel := spdxdoc.Relationship

	stdCi := spdx.CreationInfo2_1{

		SPDXVersion:                iniSpec.Val,
		DataLicense:                intDi.Val,
		DocumentName:               spdxdoc.DocumentName.Val,
		ExternalDocumentReferences: iniEdr.ExternalDocumentId.L(),
		Created:                    intCi.Create.Val(),
		CreatorComment:             intCi.Comment.V(),
		DocumentComment:            spdxdoc.DocumentComment.V(),
	}
	fmt.Println("===================================================")
	fmt.Println("CreationInfo2_1\n")
	fmt.Println("===================================================")
	fmt.Printf("%#v\n\n", stdCi.ExternalDocumentReferences)

	stdAnn := spdx.Annotation2_1{
		Annotator:         intAnn.Annotator.Val,
		AnnotatorType:     intAnn.AnnotationType.Val,
		AnnotationDate:    intAnn.AnnotationDate.Val(),
		AnnotationComment: intAnn.AnnotationComment.Val,
	}
	fmt.Println("===================================================")
	fmt.Println("Annotation2_1\n")
	fmt.Println("===================================================")
	fmt.Printf("%#v\n\n", stdAnn)

	// stdRel := spdx.Relationship2_1{
	// 	Relationship: intRel.RelationshipType.Val,
	// 	// RelationshipComment:
	// }
}

func Parse(input string) (*rdf2v1.Document, *rdf2v1.Snippet, error) {
	parser := rdf2v1.NewParser(input)
	defer fmt.Println("RDF Document parsed successfully.\n")
	defer parser.Free()
	return parser.Parse()
}
