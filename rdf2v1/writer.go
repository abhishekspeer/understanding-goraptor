package rdf2v1

import (
	"os"
	"strconv"

	"github.com/deltamobile/goraptor"
)

func Write(output *os.File, doc *Document) error {
	f := NewFormatter(output, "rdfxml-abbrev")
	_, err := f.Document(doc)
	f.Close()
	return err
}

// Formatter struct to write the output
type Formatter struct {
	serializer *goraptor.Serializer
	nodeIds    map[string]int
	fileIds    map[string]goraptor.Term
}

// NewFormatter initialses a new Formatter Interface
func NewFormatter(output *os.File, format string) *Formatter {

	// a new goraptor.NewSerializer
	s := goraptor.NewSerializer(format)

	s.StartStream(output, baseUri)

	// goraptor.NamespaceHandler:
	// handler function to be called when the parser encounters a namespace.
	s.SetNamespace("rdf", "http://www.w3.org/1999/02/22-rdf-syntax-ns#")
	s.SetNamespace("", "http://spdx.org/rdf/terms#")
	s.SetNamespace("rdfs", "http://www.w3.org/2000/01/rdf-schema#")
	s.SetNamespace("doap:", "http://usefulinc.com/ns/doap#")
	s.SetNamespace("j.0:", "http://www.w3.org/2009/pointers#")

	return &Formatter{
		serializer: s,
		nodeIds:    make(map[string]int),
		fileIds:    make(map[string]goraptor.Term),
	}
}

// NodeId method to set new node ID for a particular prefix
func (f *Formatter) NodeId(prefix string) *goraptor.Blank {

	f.nodeIds[prefix]++

	id := goraptor.Blank(prefix + strconv.Itoa(f.nodeIds[prefix]))
	return &id
}

// Sets node type
func (f *Formatter) setNodeType(node goraptor.Term, t string) error {
	return f.add(node, prefix("ns:type"), prefix(t))
}

// Add 'keys' to 'values' for subject 'to'
func (f *Formatter) add(to, key, value goraptor.Term) error {
	return f.serializer.Add(&goraptor.Statement{
		Subject:   to,
		Predicate: key,
		Object:    value,
	})
}
func (f *Formatter) addTerm(to goraptor.Term, key string, value goraptor.Term) error {
	return f.add(to, prefix(key), value)
}

func (f *Formatter) addPairs(to goraptor.Term, pairs ...pair) error {
	for _, p := range pairs {
		if err := f.addLiteral(to, p.key, p.val); err != nil {
			return err
		}
	}
	return nil
}

func (f *Formatter) addLiteral(to goraptor.Term, key, value string) error {
	if value == "" {
		return nil
	}
	return f.add(to, prefix(key), &goraptor.Literal{Value: value})
}

func (f *Formatter) Document(doc *Document) (docId goraptor.Term, err error) {

	_docId := goraptor.Blank("doc")
	docId = &_docId

	if err = f.setNodeType(docId, "SpdxDocument"); err != nil {
		return
	}

	if err = f.addLiteral(docId, "specVersion", doc.SPDXVersion.Val); err != nil {
		return
	}

	if err = f.addTerm(docId, "dataLicense", uri(licenseUri+doc.DataLicense.Val)); err != nil {
		return
	}

	// if id, err := f.CreationInfo(doc.CreationInfo); err == nil {
	// 	if err = f.addTerm(docId, "creationInfo", id); err != nil {
	// 		return docId, err
	// 	}
	// } else {
	// 	return docId, err
	// }

	if err = f.addLiteral(docId, "rdfs:comment", doc.DocumentComment.Val); err != nil {
		return
	}

	/*
		if err = f.Reviews(docId, "reviewed", doc.Reviews); err != nil {
			return
		}
	*/
	// if err = f.Packages(docId, "describesPackage", doc.Relationship.Package); err != nil {
	// 	return
	// }

	/*
			   if err = f.Files(docId, "referencesFile", doc.Files); err != nil {
			       return
			   }
		 	   if err = f.ExtrLicInfos(docId, "hasExtractedLicensingInfo", doc.ExtractedLicenceInfo); err != nil {
			       return
			   }
	*/

	return docId, nil
}

// Close to free the serializer
func (f *Formatter) Close() {
	f.serializer.EndStream()
	f.serializer.Free()
}
