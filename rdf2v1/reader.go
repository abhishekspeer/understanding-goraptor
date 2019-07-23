package rdf2v1

import (
	"fmt"

	"github.com/deltamobile/goraptor"
)

var (
	URInsType = uri("http://www.w3.org/1999/02/22-rdf-syntax-ns#type")

	typeDocument                = prefix("SpdxDocument")
	typeCreationInfo            = prefix("CreationInfo")
	typeExtractedLicensingInfo  = prefix("ExtractedLicensingInfo")
	typeRelationship            = prefix("Relationship")
	typePackage                 = prefix("Package")
	typePackageVerificationCode = prefix("PackageVerificationCode")
	typeChecksum                = prefix("Checksum")
	typeDisjunctiveLicenseSet   = prefix("DisjunctiveLicenseSet")
	typeFile                    = prefix("File")
	typeSpdxElement             = prefix("SpdxElement")
	typeSnippet                 = prefix("Snippet")
	typeLicenseConcluded        = prefix("licenseConcluded")
	typeReview                  = prefix("Review")
	typeAnnotation              = prefix("Annotation")
	typeLicense                 = prefix("License")
	typeExternalDocumentRef     = prefix("ExternalDocumentRef")
	typeProject                 = prefix("doap:Project")
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

func (p *Parser) Parse() (*Document, error) {
	// PARSE FILE method - Takes the file location as an input
	ch := p.Rdfparser.ParseFile(p.Input, "")
	var err error
	fmt.Println("PARSEFILE APPLIED")

	for statement := range ch {
		if err = p.ProcessTriple(statement); err != nil {
			fmt.Println(err)
			break
		}
	}

	return p.Doc, err
}

// Free the goraptor parser.
func (p *Parser) Free() {
	p.Rdfparser.Free()
	p.Doc = nil
}

func (p *Parser) ProcessTriple(stm *goraptor.Statement) error {
	// defer fmt.Println("Works")
	node := termStr(stm.Subject)

	// fmt.Println("\n///PARSEFILE")
	// fmt.Println("\nNODE:" + node)
	// fmt.Println("\nPREDICATE:")
	// fmt.Println(stm.Predicate)
	// fmt.Println("\nOBJECT:")
	// fmt.Println(stm.Object)
	// fmt.Println("\nURINS:")
	// fmt.Println(URInsType)
	// defer fmt.Println("***********************************")

	fmt.Println("\nstm.Predicate.Equals(URInsType)?", stm.Predicate.Equals(URInsType))
	if stm.Predicate.Equals(URInsType) {
		_, err := p.setNodeType(stm.Subject, stm.Object)

		// fmt.Println("\na:\n")
		fmt.Println(err)
		return err
	}

	// apply function if it's a builder
	builder, ok := p.Index[node]
	// fmt.Printf("BUILDER: %#v\n", builder)
	// fmt.Printf("BUILDER Creationinfo: %v\n", builder.updaters["creationInfo"])
	if ok {
		// PRINT BUILDER EACH TIME IT IS RETURNED
		// defer fmt.Printf("APPLIED: %v\n", builder.ptr)
		// defer fmt.Println("============================")
		// defer fmt.Printf("%v:\n", shortPrefix(builder.t))
		// defer fmt.Println("============================")
		return builder.apply(stm.Predicate, stm.Object)

	}

	// buffer statement
	if _, ok := p.Buffer[node]; !ok {
		// fmt.Println("\nbuffer0000000000000000000000\n")
		p.Buffer[node] = make([]*goraptor.Statement, 0)
	}

	p.Buffer[node] = append(p.Buffer[node], stm)
	// fmt.Println("\n\n\n\n99999999999999999999999\n\n\n\n")
	// fmt.Println(p.Buffer)
	return nil
}

func (p *Parser) setNodeType(node, t goraptor.Term) (interface{}, error) {
	// fmt.Printf("\n///SetNodeType\n")
	nodeStr := termStr(node)
	// fmt.Printf("\nNODESTR:" + nodeStr + "\n")
	builder, ok := p.Index[nodeStr] ////
	// fmt.Println(ok)

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
	fmt.Printf("\n//BUILDER\n")
	switch {
	// t is goraptor Object
	case t.Equals(typeDocument):
		p.Doc = new(Document)
		builder = p.MapDocument(p.Doc)
		// fmt.Printf("\n///BUILDER BACK FROM NEW BUILDER \n", builder)
		// // fmt.Printf("\n%#v\n", builder)

	case t.Equals(typeCreationInfo):
		builder = p.MapCreationInfo(new(CreationInfo))
		// fmt.Printf("%\n#v98765432\n", builder)

	case t.Equals(typeExtractedLicensingInfo):
		builder = p.MapExtractedLicensingInfo(new(ExtractedLicensingInfo))

	case t.Equals(typeRelationship):
		builder = p.MapRelationship(new(Relationship))

	case t.Equals(typePackage):
		builder = p.MapPackage(new(Package))

	case t.Equals(typePackageVerificationCode):
		builder = p.MapPackageVerificationCode(new(PackageVerificationCode))

	case t.Equals(typeChecksum):
		builder = p.MapChecksum(new(Checksum))

	case t.Equals(typeDisjunctiveLicenseSet):
		builder = p.MapDisjunctiveLicenseSet(new(DisjunctiveLicenseSet))

	case t.Equals(typeFile):
		builder = p.MapFile(new(File))

	case t.Equals(typeReview):
		builder = p.MapReview(new(Review))

	case t.Equals(typeLicense):
		builder = p.MapLicense(new(License))

	case t.Equals(typeAnnotation):
		builder = p.MapAnnotation(new(Annotation))

	case t.Equals(typeExternalDocumentRef):
		builder = p.MapExternalDocumentRef(new(ExternalDocumentRef))

	case t.Equals(typeProject):
		builder = p.MapProject(new(Project))

	case t.Equals(typeSnippet):
		builder = p.MapSnippet(new(Snippet))

	case t.Equals(typeSpdxElement):
		builder = p.MapSpdxElement(new(SpdxElement))

	default:
		fmt.Println(t)
		return nil, fmt.Errorf("ErrorTypeMatch")
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
	fmt.Println(builder)
	if ok {
		if !checkCompatibleTypes(builder.t, t) {
			return nil, fmt.Errorf("Incompatible Type")
		}
		return builder.ptr, nil
	}
	return p.setNodeType(node, t)
}

// Builder Struct and associated methods
type builder struct {
	t        goraptor.Term // type of element this builder represents
	ptr      interface{}   // the spdx element that this builder builds
	updaters map[string]updater
}

func (b *builder) apply(pred, obj goraptor.Term) error {
	property := shortPrefix(pred)
	fmt.Printf("Property: %#v\n", property)
	f, ok := b.updaters[property]
	// fmt.Printf("\nF: %#v", f)

	// fmt.Printf("\nOK: %#v", ok)
	if !ok {
		return fmt.Errorf("Property %s is not supported for %s.", property, b.t)
	}
	fmt.Println("\napplyfunctionadone")
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
