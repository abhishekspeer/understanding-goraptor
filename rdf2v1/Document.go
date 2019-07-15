package rdf2v1

import "github.com/deltamobile/goraptor"

type Document struct {
	CreationInfo *CreationInfo
	// Packages      []*Package
	// OtherLicenses []*OtherLicense
	// Relationships []*Relationship
	// Annotations []*Annotation
}

func (p *Parser) MapDocument(doc *Document) *builder {
	builder := &builder{t: typeDocument, ptr: doc}
	builder.updaters = map[string]updater{
		"creationInfo": func(obj goraptor.Term) error {
			ci, err := p.requestCreationInfo(obj)
			if err != nil {
				return err
			}
			doc.CreationInfo = ci
			return nil
		},
	}

	return builder
}
