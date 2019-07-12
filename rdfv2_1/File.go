package main

import (
	"github.com/deltamobile/goraptor"
)

type File struct {
	FileName           ValueStr
	FileSPDXIdentifier ValueStr
	FileType           ValueStr
	FileChecksumSHA1   ValueStr
	FileChecksumSHA256 ValueStr
	FileChecksumMD5    ValueStr
	LicenseConcluded   ValueStr
	LicenseInfoInFile  []ValueStr
	LicenseComments    ValueStr
	FileCopyrightText  ValueStr
	FileComment        ValueStr
	FileNotice         ValueStr
	FileContributor    []ValueStr
	//Snippets 			[]*Snippet

}

func (p *Parser) requestFile(node goraptor.Term) (*File, error) {
	obj, err := p.requestElementType(node, typeFile)
	if err != nil {
		return nil, err
	}
	return obj.(*File), err
}

// Returns a builder for file.
func (p *Parser) mapFile(file *File) *builder {
	builder := &builder{t: typeFile, ptr: file}
	builder.updaters = map[string]updater{
		"FileName":           update(&file.FileName),
		"FileSPDXIdentifier": update(&file.FileSPDXIdentifier),
		"FileType":           updateTrimPrefix("http://spdx.org/rdf/terms#", &file.FileType),
		"FileChecksumSHA1":   update(&file.FileChecksumSHA1),
		"FileChecksumSHA256": update(&file.FileChecksumSHA256),
		"FileChecksumMD5":    update(&file.FileChecksumMD5),
		"LicenseConcluded":   update(&file.LicenseConcluded),
		"LicenseInfoInFile":  updateList(&file.LicenseInfoInFile),
		"LicenseComments":    update(&file.LicenseComments),
		"rdfs:comment":       update(&file.FileComment),
		"FileCopyrightText":  update(&file.FileCopyrightText),
		"FilenNoticeText":    update(&file.FileNotice),
		"FileContributor":    updateList(&file.FileContributor),
		//snippet
	}
	return builder
}
