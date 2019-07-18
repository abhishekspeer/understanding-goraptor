package rdf2v1

import (
	"fmt"

	"github.com/deltamobile/goraptor"
)

type Document struct {
	SPDXVersion            ValueStr
	DataLicense            ValueStr
	CreationInfo           *CreationInfo
	DocumentName           ValueStr
	DocumentComment        ValueStr
	ExtractedLicensingInfo []*ExtractedLicensingInfo
	Relationship           *Relationship

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
		"specVersion": update(&doc.SPDXVersion),
		// Example: gets CC0-1.0 from "http://spdx.org/licenses/CC0-1.0"
		"dataLicense": updateTrimPrefix(licenseUri, &doc.DataLicense),
		"creationInfo": func(obj goraptor.Term) error {
			ci, err := p.requestCreationInfo(obj)
			fmt.Println(ci)
			fmt.Println(err)
			doc.CreationInfo = ci
			return err
		},
		"name":         update(&doc.DocumentName),
		"rdfs:comment": update(&doc.DocumentComment),
		"hasExtractedLicensingInfo": func(obj goraptor.Term) error {
			eli, err := p.requestExtractedLicensingInfo(obj)
			if err != nil {
				return err
			}
			doc.ExtractedLicensingInfo = append(doc.ExtractedLicensingInfo, eli)
			return nil
		},
		"relationship": func(obj goraptor.Term) error {
			rel, err := p.requestRelationship(obj)
			doc.Relationship = rel
			return err
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
