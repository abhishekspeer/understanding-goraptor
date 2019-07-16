package rdf2v1

import (
	"github.com/deltamobile/goraptor"
)

type CreationInfo struct {
	SPDXVersion                ValueStr
	DataLicense                ValueStr
	SPDXIdentifier             ValueStr
	DocumentName               ValueStr
	DocumentNamespace          ValueStr
	ExternalDocumentReferences []ValueStr
	LicenseListVersion         ValueStr
	CreatorPersons             []ValueStr
	CreatorOrganizations       []ValueStr
	CreatorTools               []ValueStr
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
		// "specVersion":                update(&ci.SPDXVersion),
		"DataLicense":                update(&ci.DataLicense),
		"SPDXIdentifier":             update(&ci.SPDXIdentifier),
		"DocumentName":               update(&ci.DocumentName),
		"DocumentNamespace":          update(&ci.DocumentNamespace),
		"ExternalDocumentReferences": updateList(&ci.ExternalDocumentReferences),
		"LicenseListVersion":         update(&ci.LicenseListVersion),
		"CreatorPersons":             updateList(&ci.CreatorPersons),
		"CreatorOrganizations":       updateList(&ci.CreatorOrganizations),
		"CreatorTools":               updateList(&ci.CreatorTools),
		"Created":                    updateDate(&ci.Create),
		"rdfs:CreatorComment":        update(&ci.CreatorComment),
		"rdfs:DocumentComment":       update(&ci.DocumentComment),
	}
	return builder
}
