package rdf2v1

import (
	"os"

	"github.com/deltamobile/goraptor"
)

type Formatter struct {
	serializer *goraptor.Serializer
	nodeIds    map[string]int
	fileIds    map[string]goraptor.Term
}

func NewFormatter(output *os.File, format string) *Formatter {

	// Initialses a new goraptor.NewSerializer
	s := goraptor.NewSerializer(format)

	s.StartStream(output, baseUri)

	// goraptor.NamespaceHandler:
	// A handler function to be called when the parser encounters a namespace.
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
