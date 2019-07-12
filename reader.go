package main

import (
	"fmt"
	"os"

	"github.com/deltamobile/goraptor"
)

func main() {

	// check that we've received the right number of arguments

	args := os.Args
	if len(args) != 2 {
		fmt.Printf("Usage: %v <spdx-file-in>\n", args[0])
		fmt.Printf("  Load SPDX 2.1 RDF file <spdx-file-in>, and\n")
		fmt.Printf("  print its contents.\n")
		return
	}

	//   PARSE FILE method - Takes the file location as an input
	input := args[1]
	parserfile := NewParser(input)
	defer parserfile.Free()

	// Parse the file using goraptor's ParseFile method and return a Statement.
	ch := parserfile.rdfparser.ParseFile(input, "") // takes in input and baseuri
	for {
		statement, ok := <-ch
		// fmt.Println(statement)
		if !ok {
			break
		}
		err := parserfile.processTriple(statement)
		if err != nil {
			fmt.Println("Processing Failed")
		}
	}

	parserinstance := parserfile
	// fmt.Println(parserinstance.buffer)
	fmt.Println(len(parserinstance.buffer))
}

var (
	URInsType = uri("http://www.w3.org/1999/02/22-rdf-syntax-ns#type")

	typeDocument     = prefix("SpdxDocument")
	typeCreationInfo = prefix("CreationInfo")
	typePackage      = prefix("Package")
	typeFile         = prefix("File")
)

// Parser Struct and associated methods
type Parser struct {
	rdfparser *goraptor.Parser
	input     string
	index     map[string]*builder
	buffer    map[string][]*goraptor.Statement
	doc       *Document
}

// NewParser uses goraptor.NewParser to initialse a new parser interface
func NewParser(input string) *Parser {

	return &Parser{
		rdfparser: goraptor.NewParser("guess"),
		input:     input,
		index:     make(map[string]*builder),
		buffer:    make(map[string][]*goraptor.Statement),
	}
}

func (p *Parser) setNodeType(node, t goraptor.Term) (interface{}, error) {
	nodeStr := termStr(node)
	builder, ok := p.index[nodeStr]
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

	p.index[nodeStr] = builder

	// run buffer
	buf := p.buffer[nodeStr]
	for _, stm := range buf {
		if err := builder.apply(stm.Predicate, stm.Object); err != nil {
			return nil, err
		}
	}
	delete(p.buffer, nodeStr)

	return builder.ptr, nil
}

func (p *Parser) processTriple(stm *goraptor.Statement) error {
	node := termStr(stm.Subject)
	if stm.Predicate.Equals(URInsType) {
		_, err := p.setNodeType(stm.Subject, stm.Object)
		return err
	}

	// apply function if it's a builder
	builder, ok := p.index[node]
	if ok {
		return builder.apply(stm.Predicate, stm.Object)
	}

	// buffer statement
	if _, ok := p.buffer[node]; !ok {
		p.buffer[node] = make([]*goraptor.Statement, 0)
	}
	p.buffer[node] = append(p.buffer[node], stm)

	return nil
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
	builder, ok := p.index[termStr(node)]
	if ok {
		if !checkCompatibleTypes(builder.t, t) {
			return nil, fmt.Errorf("Incompatible Type")
		}
		return builder.ptr, nil
	}
	return p.setNodeType(node, t)
}

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

// possiblity of a mistake (node type)
func (p *Parser) reqSomething(node, t goraptor.Term) (interface{}, error) {
	builder, ok := p.index[termStr(node)]
	if ok {
		if !builder.t.Equals(t) {
			return nil, fmt.Errorf("Incompatible Type", node, builder.t, t)
		}
		return builder.ptr, nil
	}
	return p.setNodeType(node, t)
}

// returns the SPDX document
func (p *Parser) Parse() (*Document, error) {

	return p.doc, nil
}

// Free the Parser
func (p *Parser) Free() {
	p.rdfparser.Free()
	p.doc = nil
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

type updater func(goraptor.Term) error
