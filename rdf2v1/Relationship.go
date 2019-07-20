package rdf2v1

import "github.com/deltamobile/goraptor"

type Relationship struct {
	Relationship ValueStr
	Package      []*Package
	// RelationshipComment string
}

func (p *Parser) requestRelationship(node goraptor.Term) (*Relationship, error) {
	obj, err := p.requestElementType(node, typeRelationship)
	if err != nil {
		return nil, err
	}
	return obj.(*Relationship), err
}

func (p *Parser) MapRelationship(rel *Relationship) *builder {
	builder := &builder{t: typeRelationship, ptr: rel}
	builder.updaters = map[string]updater{
		"relationshipType": updateTrimPrefix(baseUri, &rel.Relationship),
		"relatedSpdxElement": func(obj goraptor.Term) error {
			pkg, err := p.requestPackage(obj)
			if err != nil {
				return err
			}
			rel.Package = append(rel.Package, pkg)
			return nil
		},
	}
	return builder
}
