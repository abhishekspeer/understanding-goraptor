package rdf2v1

import (
	"github.com/deltamobile/goraptor"
)

type CreationInfo struct {
	SPDXIdentifier     ValueStr
	LicenseListVersion ValueStr
	Creator            []ValueStr
	Create             ValueDate
	Comment            ValueStr
}

func (p *Parser) requestCreationInfo(node goraptor.Term) (*CreationInfo, error) {

	obj, err := p.requestElementType(node, typeCreationInfo)
	if err != nil {
		return nil, err
	}
	return obj.(*CreationInfo), err
}

func (p *Parser) MapCreationInfo(ci *CreationInfo) *builder {
	builder := &builder{t: typeCreationInfo, ptr: ci}
	builder.updaters = map[string]updater{
		"licenseListVersion": update(&ci.LicenseListVersion),
		"creator":            updateList(&ci.Creator),
		"created":            updateDate(&ci.Create),
		"rdfs:comment":       update(&ci.Comment),
	}
	return builder
}
