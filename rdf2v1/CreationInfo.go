package rdf2v1

import (
	"fmt"

	"github.com/deltamobile/goraptor"
)

type CreationInfo struct {
	SPDXIdentifier     ValueStr
	LicenseListVersion ValueStr
	Creator            []ValueCreator
	Create             ValueDate
}

// Requests type, returns a pointer.
func (p *Parser) requestCreationInfo(node goraptor.Term) (*CreationInfo, error) {
	fmt.Println("GGGGGGGGGGGGGGGGGGGG")

	obj, err := p.requestElementType(node, typeCreationInfo)
	if err != nil {
		return nil, err
	}
	return obj.(*CreationInfo), err
}

// Returns a builder for cri.
func (p *Parser) MapCreationInfo(ci *CreationInfo) *builder {
	builder := &builder{t: typeCreationInfo, ptr: ci}
	builder.updaters = map[string]updater{
		"licenseListVersion": update(&ci.LicenseListVersion),
		"creator":            updateListCreator(&ci.Creator),
		"created":            updateDate(&ci.Create),
	}
	return builder
}
