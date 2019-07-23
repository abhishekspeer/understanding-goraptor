package rdf2v1

import (
	"fmt"

	"github.com/deltamobile/goraptor"
)

type Snippet struct {
	SnippetName            ValueStr
	SnippetCopyrightText   ValueStr
	SnippetLicenseComments ValueStr
	SnippetFromFile        *File
	// 	SnippetSPDXIdentifier string
	// 	SnippetByteRangeStart int
	// 	SnippetByteRangeEnd   int
	// 	SnippetLineRangeStart int
	// 	SnippetLineRangeEnd   int
	// 	SnippetLicenseConcluded string
	// 	SnippetComment string
	// 	LicenseInfoInSnippet []string
}

// type SnippetFromFile struct {
// 	SnippetFileCopyrightText ValueStr
// 	SnippetFileContributor   []ValueStr
// 	SnippetFileName          ValueStr
// 	SnippetFileChecksum      *Checksum
// 	SnippetLicense           *License
// }

func (p *Parser) requestSnippet(node goraptor.Term) (*Snippet, error) {
	obj, err := p.requestElementType(node, typeSnippet)
	if err != nil {
		return nil, err
	}
	return obj.(*Snippet), err
}

// Returns a builder for cri.
func (p *Parser) MapSnippet(s *Snippet) *builder {
	fmt.Println("\n\n///Snippet\n")
	builder := &builder{t: typeSnippet, ptr: s}
	builder.updaters = map[string]updater{
		"name":            update(&s.SnippetName),
		"copyrightText":   update(&s.SnippetCopyrightText),
		"licenseComments": update(&s.SnippetLicenseComments),
		"snippetFromFile": func(obj goraptor.Term) error {
			file, _ := p.requestFile(obj)
			s.SnippetFromFile = file
			return nil
		},
	}
	return builder
}
