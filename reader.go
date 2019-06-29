package main

import (
	"fmt"
	"os"

	"github.com/deltamobile/goraptor"
	"github.com/spdx/tools-golang/v0/spdx"
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

// Parser Struct and associated methods
type Parser struct {
	rdfparser *goraptor.Parser
	input     string
	index     map[string]*builder
	buffer    map[string][]*goraptor.Statement
	doc       *spdx.Document2_1
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

// Process the goraptor statement and Apply buffer and builder operations
func (p *Parser) processTriple(stm *goraptor.Statement) error {
	node := termStr(stm.Subject)
	////
	// apply function if it's a builder
	bldr, ok := p.index[node]
	if ok {
		return bldr.apply(stm.Predicate, stm.Object)
	}

	// buffer statement
	if _, ok := p.buffer[node]; !ok {
		p.buffer[node] = make([]*goraptor.Statement, 0)
	}
	p.buffer[node] = append(p.buffer[node], stm)

	return nil
}

// returns the SPDX document
func (p *Parser) Parse() (*spdx.Document2_1, error) {

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

type updater func(goraptor.Term) error

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
