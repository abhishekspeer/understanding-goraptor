package rdf2v1

import (
	"fmt"

	"github.com/deltamobile/goraptor"
)

var (
	URInsType = uri("http://www.w3.org/1999/02/22-rdf-syntax-ns#type")

	typeDocument     = prefix("SpdxDocument")
	typeCreationInfo = prefix("CreationInfo")
	typePackage      = prefix("Package")
	typeFile         = prefix("File")
)

// Parser Struct and associated methods
type Parser struct {
	Rdfparser *goraptor.Parser
	Input     string
	Index     map[string]*builder
	Buffer    map[string][]*goraptor.Statement
	Doc       *Document
}

// NewParser uses goraptor.NewParser to initialse a new parser interface
func NewParser(input string) *Parser {

	return &Parser{
		Rdfparser: goraptor.NewParser("guess"),
		Input:     input,
		Index:     make(map[string]*builder),
		Buffer:    make(map[string][]*goraptor.Statement),
	}
}
func (p *Parser) ProcessTriple(stm *goraptor.Statement) error {
	node := termStr(stm.Subject)
	if stm.Predicate.Equals(URInsType) {
		_, err := p.setNodeType(stm.Subject, stm.Object)
		return err
	}

	// apply function if it's a builder
	builder, ok := p.Index[node]
	if ok {
		return builder.apply(stm.Predicate, stm.Object)
	}

	// buffer statement
	if _, ok := p.Buffer[node]; !ok {
		p.Buffer[node] = make([]*goraptor.Statement, 0)
	}
	p.Buffer[node] = append(p.Buffer[node], stm)

	return nil
}

func (p *Parser) setNodeType(node, t goraptor.Term) (interface{}, error) {
	nodeStr := termStr(node)
	builder, ok := p.Index[nodeStr]
	if ok {
		if !checkRaptorTypes(builder.t, t) && builder.checkPredicate("ns:type") {
			//apply the type change
			if err := builder.apply(uri("ns:type"), t); err != nil {
				return nil, err
			}
			return builder.ptr, nil
		}
		if !checkCompatibleTypes(builder.t, t) {
			return nil, fmt.Errorf("Incompatible Type")
		}
		return builder.ptr, nil
	}
	// new builder by type
	switch {
	case t.Equals(typeDocument):
		builder = p.documentMap(new(Document))
	case t.Equals(typeCreationInfo):
		builder = p.mapCreationInfo(new(CreationInfo))
	}

	p.Index[nodeStr] = builder

	// run buffer
	buf := p.Buffer[nodeStr]
	for _, stm := range buf {
		if err := builder.apply(stm.Predicate, stm.Object); err != nil {
			return nil, err
		}
	}
	delete(p.Buffer, nodeStr)

	return builder.ptr, nil
}

func checkRaptorTypes(found goraptor.Term, need ...goraptor.Term) bool {
	for _, b := range need {
		if found == b || found.Equals(b) {
			return true
		}
	}
	return false
}

func checkCompatibleTypes(input, required goraptor.Term) bool {
	if checkRaptorTypes(input, required) {
		return true
	}
	return false
}

func (p *Parser) requestElementType(node, t goraptor.Term) (interface{}, error) {
	builder, ok := p.Index[termStr(node)]
	if ok {
		if !checkCompatibleTypes(builder.t, t) {
			return nil, fmt.Errorf("Incompatible Type")
		}
		return builder.ptr, nil
	}
	return p.setNodeType(node, t)
}

// // return builders for document

func (p *Parser) documentMap(doc *Document) *builder {
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

// returns the SPDX document
func (p *Parser) Parse() (*Document, error) {

	return p.Doc, nil
}

// Free the Parser
func (p *Parser) Free() {
	p.Rdfparser.Free()
	p.Doc = nil
}

// Builder Struct and associated methods
type builder struct {
	t        goraptor.Term // type of element this builder represents
	ptr      interface{}   // the spdx element that this builder builds
	updaters map[string]updater
}

func (b *builder) apply(pred, obj goraptor.Term) error {
	property := termStr(pred)
	f, ok := b.updaters[termStr(pred)]
	if !ok {
		return fmt.Errorf("Property %s is not supported for %s.", property, b.t)
	}
	return f(obj)
}

// Converts goraptor.Term (Subject, Predicate and Object) to string.
func termStr(term goraptor.Term) string {
	switch t := term.(type) {
	case *goraptor.Uri:
		return string(*t)
	case *goraptor.Blank:
		return string(*t)
	case *goraptor.Literal:
		return t.Value
	default:
		return ""
	}
}

// to check if builder contains a predicate
func (b *builder) checkPredicate(pred string) bool {
	_, ok := b.updaters[pred]
	return ok
}

// Uri, Literal and Blank are goraptors named types
// Return *goraptor.Uri
func uri(uri string) *goraptor.Uri {
	return (*goraptor.Uri)(&uri)
}

// Return *goraptor.Literal
func literal(lit string) *goraptor.Literal {
	return &goraptor.Literal{Value: lit}
}

// Return *goraptor.Blank from string
func blank(b string) *goraptor.Blank {
	return (*goraptor.Blank)(&b)
}