package rdf2v1

import (
	"fmt"

	"github.com/deltamobile/goraptor"
)

type Document struct {
	CreationInfo *CreationInfo
	// Packages      []*Package
	// OtherLicenses []*OtherLicense
	// Relationships []*Relationship
	// Annotations []*Annotation
}

func (p *Parser) MapDocument(doc *Document) *builder {
	fmt.Println("\n\n///MAPDOCUMENT\n")
	builder := &builder{t: typeDocument, ptr: doc}
	fmt.Printf("\nBUILDER BEFORE UPDATE: \n")
	fmt.Println(builder)

	builder.updaters = map[string]updater{
		"CreationInfo": func(obj goraptor.Term) error {
			ci, err := p.requestCreationInfo(obj)

			fmt.Println(ci)
			fmt.Println(err)
			if err != nil {
				return err
			}
			doc.CreationInfo = ci
			return nil
		},
	}
	// fmt.Println("\n\nLLLLLLLLLLLLLLLLLLLLLLL\n\n")
	fmt.Printf("\nBUILDER UPDATED: \n")
	fmt.Printf("%#v", builder)
	fmt.Println("\n///MAPDOCUMENT DONE\n\n")
	// fmt.Println(doc.CreationInfo)

	// fmt.Println("\n\nLLLLLLLLLLLLLLLLLLLLLLL\n\n")

	return builder
}
