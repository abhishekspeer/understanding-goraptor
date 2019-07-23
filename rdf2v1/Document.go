package rdf2v1

import (
	"fmt"

	"github.com/deltamobile/goraptor"
)

type Document struct {
	SPDXVersion            ValueStr
	DataLicense            ValueStr
	CreationInfo           *CreationInfo
	Review                 []*Review
	DocumentName           ValueStr
	DocumentComment        ValueStr
	ExtractedLicensingInfo []*ExtractedLicensingInfo
	Relationship           []*Relationship
	License                *License
	Annotation             *Annotation
	ExternalDocumentRef    *ExternalDocumentRef
}
type Review struct {
	ReviewComment ValueStr
	ReviewDate    ValueDate
	Reviewer      ValueStr
}

type ExternalDocumentRef struct {
	ExternalDocumentId ValueStr
	Checksum           *Checksum
	SPDXDocument       ValueStr
}

func (p *Parser) requestDocument(node goraptor.Term) (*Document, error) {
	obj, err := p.requestElementType(node, typeDocument)
	if err != nil {
		return nil, err
	}
	return obj.(*Document), err
}
func (p *Parser) requestReview(node goraptor.Term) (*Review, error) {
	obj, err := p.requestElementType(node, typeReview)
	if err != nil {
		return nil, err
	}
	return obj.(*Review), err
}

func (p *Parser) requestExternalDocumentRef(node goraptor.Term) (*ExternalDocumentRef, error) {
	obj, err := p.requestElementType(node, typeExternalDocumentRef)
	if err != nil {
		return nil, err
	}
	return obj.(*ExternalDocumentRef), err
}

func (p *Parser) MapExternalDocumentRef(edr *ExternalDocumentRef) *builder {
	builder := &builder{t: typeExternalDocumentRef, ptr: edr}
	builder.updaters = map[string]updater{
		"externalDocumentId": update(&edr.ExternalDocumentId),
		"checksum": func(obj goraptor.Term) error {
			cksum, err := p.requestChecksum(obj)

			edr.Checksum = cksum
			return err
		},
		"spdxDocument": update(&edr.SPDXDocument),
	}
	return builder

}
func (p *Parser) MapReview(rev *Review) *builder {
	builder := &builder{t: typeReview, ptr: rev}
	builder.updaters = map[string]updater{
		"rdfs:comment": update(&rev.ReviewComment),
		"reviewDate":   updateDate(&rev.ReviewDate),
		"reviewer":     update(&rev.Reviewer),
	}
	return builder

}

func (p *Parser) MapDocument(doc *Document) *builder {
	fmt.Println("\n\n///MAPDOCUMENT\n")
	builder := &builder{t: typeDocument, ptr: doc}
	// fmt.Printf("\nBUILDER BEFORE UPDATE: \n")
	// fmt.Println(builder)

	builder.updaters = map[string]updater{
		"specVersion": update(&doc.SPDXVersion),
		// Example: gets CC0-1.0 from "http://spdx.org/licenses/CC0-1.0"
		"dataLicense": func(obj goraptor.Term) error {
			lic, err := p.requestLicense(obj)
			doc.License = lic
			return err
		},
		"creationInfo": func(obj goraptor.Term) error {
			ci, err := p.requestCreationInfo(obj)
			fmt.Println(ci)
			fmt.Println(err)
			doc.CreationInfo = ci
			return err
		},
		"reviewed": func(obj goraptor.Term) error {
			rev, err := p.requestReview(obj)
			if err != nil {
				return err
			}
			doc.Review = append(doc.Review, rev)
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
			if err != nil {
				return err
			}
			doc.Relationship = append(doc.Relationship, rel)
			return nil
		},
		"annotation": func(obj goraptor.Term) error {
			an, err := p.requestAnnotation(obj)
			doc.Annotation = an
			return err
		},
		"externalDocumentRef": func(obj goraptor.Term) error {
			edr, err := p.requestExternalDocumentRef(obj)
			doc.ExternalDocumentRef = edr
			return err
		},
		// "Package": func(obj goraptor.Term) error {
		// 	pkg, err := p.requestPackage(obj)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	doc.Package = append(doc.Package, pkg)
		// 	return nil
		// },
	}
	// fmt.Println("\n\nLLLLLLLLLLLLLLLLLLLLLLL\n\n")
	// fmt.Printf("\nBUILDER UPDATED: \n")
	// fmt.Printf("%#v", builder)
	// fmt.Println("\n///MAPDOCUMENT DONE\n\n")
	// fmt.Println(doc.CreationInfo)

	// fmt.Println("\n\nLLLLLLLLLLLLLLLLLLLLLLL\n\n")

	return builder
}
