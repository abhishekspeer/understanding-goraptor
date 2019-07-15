package rdf2v1

import (
	"os"

	"github.com/deltamobile/goraptor"
)

func WriteInput(input *os.File, output *os.File) error {

	inputFormat := "guess"
	outputFormat := "rdfxml-abbrev"

	parser := goraptor.NewParser(inputFormat)
	defer parser.Free()

	serializer := goraptor.NewSerializer(outputFormat)
	defer serializer.Free()

	parser.SetNamespaceHandler(func(pfx, uri string) { serializer.SetNamespace(pfx, uri) })

	statements := parser.Parse(input, baseUri)
	err := serializer.SetFile(output, baseUri)
	if err != nil {
		return err
	}
	return serializer.AddN(statements)
}

type Formatter struct {
	serializer *goraptor.Serializer
	nodeIds    map[string]int
	fileIds    map[string]goraptor.Term
}

func NewFormatter(output *os.File, format string) *Formatter {
	s := goraptor.NewSerializer(format)
	s.StartStream(output, baseUri)

	s.SetNamespace("rdf", "http://www.w3.org/1999/02/22-rdf-syntax-ns#")
	s.SetNamespace("", "http://spdx.org/rdf/terms#")
	s.SetNamespace("rdfs", "http://www.w3.org/2000/01/rdf-schema#")

	return &Formatter{
		serializer: s,
		nodeIds:    make(map[string]int),
		fileIds:    make(map[string]goraptor.Term),
	}
}
