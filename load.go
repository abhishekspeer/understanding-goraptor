package main

import (
	"fmt"
	"io"

	"github.com/deltamobile/goraptor"
	"spdx/tools-golang/v0/spdx"
)
// type Document2_1 struct {
// 	CreationInfo  *CreationInfo2_1
// 	Packages      []*Package2_1
// 	OtherLicenses []*OtherLicense2_1
// 	Relationships []*Relationship2_1
// 	Annotations   []*Annotation2_1

// 	// DEPRECATED in version 2.0 of spec
// 	Reviews []*Review2_1

type Parser struct {
	rdfparser *goraptor.Parser
	input     io.Reader
	index     map[string]*builder
	buffer    map[string][]bufferEntry
	doc       *spdx.Document2_1
}

func Parse(input io.Reader, format string) (*spdx.Document2_1, error) {
	parser := NewParser(input, format)
	defer parser.Free()
	return parser.Parse()
}

func NewParser(input io.Reader, format string) *Parser {
	if format == "rdf" {
		format = "guess"
	}

	return &Parser{
		rdfparser: goraptor.NewParser(format),
		input:     input,
		index:     make(map[string]*builder),
		buffer:    make(map[string][]bufferEntry),
	}
}

func (p *Parser) Parse() (*spdx.Document2_1, error) {
	ch := p.rdfparser.Parse(p.input, baseUri)
	locCh := p.rdfparser.LocatorChan()
	var err error
	for statement := range ch {
		locator := <-locCh
		meta := spdx.NewMetaL(locator.Line)
		if err = p.processTruple(statement, meta); err != nil {
			break
		}
	}
	// Consume input channel in case of error. Otherwise goraptor will keep the goroutine busy.
	for _ = range ch {
		<-locCh
	}
	return p.doc, err
}

func (p *Parser) Free() {
	p.rdfparser.Free()
	p.doc = nil
}

type builder struct {
	t        goraptor.Term // type of element this builder represents
	ptr      interface{}   // the spdx element that this builder builds
	updaters map[string]updater
}

func (b *builder) apply(pred, obj goraptor.Term, meta *spdx.Meta) error {
	property := shortPrefix(pred)
	f, ok := b.updaters[property]
	if !ok {
		return spdx.NewParseError(fmt.Sprintf(msgPropertyNotSupported, property, b.t), meta)
	}
	return f(obj, meta)
}

func (b *builder) has(pred string) bool {
	_, ok := b.updaters[pred]
	return ok
}

type updater func(goraptor.Term, *spdx.Meta) error

type bufferEntry struct {
	*goraptor.Statement
	*spdx.Meta
}

