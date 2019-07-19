package rdf2v1

import (
	"fmt"

	"github.com/deltamobile/goraptor"
)

type CreationInfo struct {
	SPDXIdentifier     ValueStr
	LicenseListVersion ValueStr
	Creator            []ValueStr
	Create             ValueDate
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
func (p *Parser) MapCreationInfo(ci *CreationInfo) *builder {
	fmt.Println("\n\n///MAPCREATIONINFO\n")
	builder := &builder{t: typeCreationInfo, ptr: ci}
	builder.updaters = map[string]updater{
		"licenseListVersion": update(&ci.LicenseListVersion),
		"creator":            updateList(&ci.Creator),
		"created":            updateDate(&ci.Create),
	}
	return builder
}
