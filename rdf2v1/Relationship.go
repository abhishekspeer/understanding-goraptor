package rdf2v1

import "github.com/deltamobile/goraptor"

type Relationship struct {
	RelationshipType ValueStr
	SpdxElement      *SpdxElement
	File             *File
	Package          []*Package
	// RelationshipComment string
}
type SpdxElement struct {
}

func (p *Parser) requestRelationship(node goraptor.Term) (*Relationship, error) {
	obj, err := p.requestElementType(node, typeRelationship)
	if err != nil {
		return nil, err
	}
	return obj.(*Relationship), err
}
func (p *Parser) requestSpdxElement(node goraptor.Term) (*SpdxElement, error) {
	obj, err := p.requestElementType(node, typeSpdxElement)
	if err != nil {
		return nil, err
	}
	return obj.(*SpdxElement), err
}

func (p *Parser) MapRelationship(rel *Relationship) *builder {
	builder := &builder{t: typeRelationship, ptr: rel}
	builder.updaters = map[string]updater{
		"relationshipType": update(&rel.RelationshipType),
		// "relatedSpdxElement": func(obj goraptor.Term) error {
		// 	pkg, err := p.requestPackage(obj)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	rel.Package = append(rel.Package, pkg)
		// 	return nil
		// },
		"relatedSpdxElement": func(obj goraptor.Term) error {
			file, _ := p.requestFile(obj)
			rel.File = file
			return nil
		},
	}
	return builder
}

func (p *Parser) MapSpdxElement(se *SpdxElement) *builder {
	builder := &builder{t: typeSpdxElement, ptr: se}
	builder.updaters = map[string]updater{}
	return builder
}
