package main

import (
	"fmt"
	//"io"
	"os"

	//"spdx/tools-golang/v0/spdx"

	//"github.com/deltamobile/goraptor"
)

func main() {
	// defer fmt.Println("DONE!")
	// parser := goraptor.NewParser("guess")
	// defer parser.Free()

	// ch := parser.ParseUri("http://spdx.org/licenses/CC0-1.0", "")
	// for {
	// 	statement, ok := <-ch
	// 	if !ok {
	// 		break
	// 	}
	// 	fmt.Println(parser.Parse())
	// 	// fmt.Printf("%T",statement) // returning a goraptor statement (custom type)
	// 	// fmt.Println(statement)

	//	}




	// Checking right number of arguments
	args := os.Args

	if len(args) != 2 {
		fmt.Printf("Usage: %v <spdx-rdf-in>\n", args[0])
		return
	}

	// assigning filename 
	filename := args[1]
	
	// opening local file
	r, err := os.Open(filename)
	
	defer r.Close()
	defer fmt.Printf("\nDONE\n")

	if err != nil {
		fmt.Printf("Error while opening %v for reading: %v", filename, err)
		return
	}
	fmt.Println("File opened")
	// fmt.Printf("%s and Type: %T", r, r)
	
	data := make([]byte, 10)

	count, err := r.Read(data)

	if err != nil {
		fmt.Printf("Error while writing %v: %v", filename, err)
		return
	}
	fmt.Println("File Read") 
	fmt.Printf("Read %d bytes\nContents of File: \n%q\n%T", count, string(data[:]),data)


//	Parse(r, "kjbk")

}

// type Parser struct {
// 	rdfparser *goraptor.Parser
// 	input     io.Reader
// 	index     map[string]*builder
// 	buffer    map[string][]bufferEntry
// 	doc       *spdx.Document2_1
// }

// func Parse(input io.Reader, format string) (*spdx.Document2_1, error) {
// 	parser := NewParser(input, format)
// 	defer parser.Free()
// 	return parser.Parse()
// }

// func NewParser(input io.Reader, format string) *Parser {
// 	if format == "rdf" {
// 		format = "guess"
// 	}

// 	return &Parser{
// 		rdfparser: goraptor.NewParser(format),
// 		input:     input,
// 		index:     make(map[string]*builder),
// 		buffer:    make(map[string][]bufferEntry),
// 	}
// }

// func (p *Parser) Parse() (*spdx.Document2_1, error) {
// 	ch := p.rdfparser.Parse(p.input, baseUri)
// 	locCh := p.rdfparser.LocatorChan()
// 	var err error
// 	for statement := range ch {
// 		locator := <-locCh
// 		meta := spdx.NewMetaL(locator.Line)
// 		if err = p.processTruple(statement, meta); err != nil {
// 			break
// 		}
// 	}
// 	// Consume input channel in case of error. Otherwise goraptor will keep the goroutine busy.
// 	for _ = range ch {
// 		<-locCh
// 	}
// 	return p.doc, err
// }

// func (p *Parser) Free() {
// 	p.rdfparser.Free()
// 	p.doc = nil
// }

// type builder struct {
// 	t        goraptor.Term // type of element this builder represents
// 	ptr      interface{}   // the spdx element that this builder builds
// 	updaters map[string]updater
// }

// func (b *builder) apply(pred, obj goraptor.Term, meta *spdx.Meta) error {
// 	property := shortPrefix(pred)
// 	f, ok := b.updaters[property]
// 	if !ok {
// 		return spdx.NewParseError(fmt.Sprintf(msgPropertyNotSupported, property, b.t), meta)
// 	}
// 	return f(obj, meta)
// }

// func (b *builder) has(pred string) bool {
// 	_, ok := b.updaters[pred]
// 	return ok
// }

// type updater func(goraptor.Term, *spdx.Meta) error

// type bufferEntry struct {
// 	*goraptor.Statement
// 	*spdx.Meta
// }
