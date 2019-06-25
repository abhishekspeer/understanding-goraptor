package main

import (
	"fmt"
	"os"

	"github.com/deltamobile/goraptor"
	"github.com/spdx/tools-golang/v0/spdx"
)

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

	ch := parserfile.rdfparser.ParseFile(input, "") // takes in input and baseuri
	for {
		statement, ok := <-ch
		// fmt.Println(statement)
		if !ok {
			break
		}
		err := parserfile.processTriple(statement)
		if err != nil {
			fmt.Println("ERROR HERE!!!!!!!!!!!!!!!")
		}
	}

	parserinstance := parserfile
	fmt.Printf("%#v", parserinstance.rdfparser)

}

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

func (p *Parser) Parse() (*spdx.Document2_1, error) {

	return p.doc, nil
}

func (p *Parser) Free() {
	p.rdfparser.Free()
	p.doc = nil
}

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
