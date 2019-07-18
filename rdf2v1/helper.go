package rdf2v1

import (
	"strings"

	"github.com/deltamobile/goraptor"
)

func prefix(k string) *goraptor.Uri {
	var pref string
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
			return strings.Replace(str, long, short, 1)
		}
	}
	return str
}

const (
	baseUri    = "http://spdx.org/rdf/terms#"
	licenseUri = "http://spdx.org/licenses/"
)

var rdfPrefixes = map[string]string{
	"rdf:":  "http://www.w3.org/1999/02/22-rdf-syntax-ns#",
	"rdfs:": "http://www.w3.org/2000/01/rdf-schema#",
	"spdx:": baseUri,
}
