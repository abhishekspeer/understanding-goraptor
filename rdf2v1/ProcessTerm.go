package rdf2v1

import (
	"strings"

	"github.com/deltamobile/goraptor"
)

const (
	baseUri    = "http://spdx.org/rdf/terms#"
	licenseUri = "http://spdx.org/licenses/"
)

var rdfPrefixes = map[string]string{
	"rdf:":  "http://www.w3.org/1999/02/22-rdf-syntax-ns#",
	"doap:": "http://usefulinc.com/ns/doap#",
	"rdfs:": "http://www.w3.org/2000/01/rdf-schema#",
	"j.0:":  "http://www.w3.org/2009/pointers#",
	"":      baseUri,
}

// Converts typeX to its full URI accorinding to rdfPrefixes,
// if no : is found in the string it'll assume it as "spdx:" and expand to baseUri
func prefix(k string) *goraptor.Uri {
	var pref string = ""
	rest := k
	if i := strings.Index(k, ":"); i >= 0 {
		pref = k[:i+1]
		rest = k[i+1:]
	}
	if long, ok := rdfPrefixes[pref]; ok {
		pref = long
	}
	uri := goraptor.Uri(pref + rest)
	return &uri
}

// Change the RDF prefixes to their short forms.
func shortPrefix(t goraptor.Term) string {
	str := termStr(t)
	for short, long := range rdfPrefixes {
		if strings.HasPrefix(str, long) {
			str = strings.Replace(str, long, short, 1)
			return strings.Replace(str, long, short, 1)
		}
	}
	return str
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
