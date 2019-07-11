package main

import (
	"github.com/deltamobile/goraptor"
)

type CreationInfo struct {
	SPDXVersion                ValueStr
	DataLicense                ValueStr
	SPDXIdentifier             ValueStr
	DocumentName               ValueStr
	DocumentNamespace          ValueStr
	ExternalDocumentReferences []ValueCreator
	LicenseListVersion         ValueStr
	CreatorPersons             []ValueCreator
	CreatorOrganizations       []ValueCreator
	CreatorTools               []ValueCreator
	Create                     ValueDate
	CreatorComment             ValueStr
	DocumentComment            ValueStr
}

// Requests type, returns a pointer.
func (p *Parser) requestCreationInfo(node goraptor.Term) (*CreationInfo, error) {
	obj, err := p.requestElementType(node, typeCreationInfo)
	if err != nil {
		return nil, err
	}
	return obj.(*CreationInfo), err
}

// Returns a builder for cri.
func (p *Parser) mapCreationInfo(ci *CreationInfo) *builder {
	builder := &builder{t: typeCreationInfo, ptr: ci}
	builder.updaters = map[string]updater{
		"SPDXVersion":                update(&ci.SPDXVersion),
		"DataLicense":                update(&ci.DataLicense),
		"SPDXIdentifier":             update(&ci.SPDXIdentifier),
		"DocumentName":               update(&ci.DocumentName),
		"DocumentNamespace":          update(&ci.DocumentNamespace),
		"ExternalDocumentReferences": updateListCreator(&ci.ExternalDocumentReferences),
		"LicenseListVersion":         update(&ci.LicenseListVersion),
		"CreatorPersons":             updateListCreator(&ci.CreatorPersons),
		"CreatorOrganizations":       updateListCreator(&ci.CreatorOrganizations),
		"CreatorTools":               updateListCreator(&ci.CreatorTools),
		"Created":                    updateDate(&ci.Create),
		"rdfs:CreatorComment":        update(&ci.CreatorComment),
		"rdfs:DocumentComment":       update(&ci.DocumentComment),
	}
	return builder
}
