package rdf2v1

import (
	"github.com/deltamobile/goraptor"
)

type File struct {
	FileName           ValueStr
	FileChecksum       *FileChecksum
	LicenseInfoInFile  []ValueStr
	FileCopyrightText  ValueStr
	FileComment        ValueStr
	FileSPDXIdentifier ValueStr
	FileType           ValueStr
	FileChecksumSHA1   ValueStr
	FileChecksumSHA256 ValueStr
	FileChecksumMD5    ValueStr
	LicenseConcluded   ValueStr
	FileNotice         ValueStr
	FileContributor    []ValueStr
	//Snippets 			[]*Snippet

}
type FileChecksum struct {
	Algorithm     ValueStr
	ChecksumValue ValueStr
}
type DisjunctiveLicenseSet struct {
	Member []ValueStr
}

func (p *Parser) requestFile(node goraptor.Term) (*File, error) {
	obj, err := p.requestElementType(node, typeFile)
	if err != nil {
		return nil, err
	}
	return obj.(*File), err
}
func (p *Parser) requestFileChecksum(node goraptor.Term) (*FileChecksum, error) {
	obj, err := p.requestElementType(node, typeChecksum)
	if err != nil {
		return nil, err
	}
	return obj.(*FileChecksum), err
}
func (p *Parser) requestDisjunctiveLicenseSet(node goraptor.Term) (*DisjunctiveLicenseSet, error) {
	obj, err := p.requestElementType(node, typeDisjunctiveLicenseSet)
	if err != nil {
		return nil, err
	}
	return obj.(*DisjunctiveLicenseSet), err
}

// Returns a builder for file.
func (p *Parser) MapFile(file *File) *builder {
	builder := &builder{t: typeFile, ptr: file}
	builder.updaters = map[string]updater{
		"fileName": update(&file.FileName),
		"checksum": func(obj goraptor.Term) error {
			filecksum, err := p.requestFileChecksum(obj)
			file.FileChecksum = filecksum
			return err
		},
		// "FileType":          updateTrimPrefix("http://spdx.org/rdf/terms#", &file.FileType),
		"LicenseConcluded":  update(&file.LicenseConcluded),
		"LicenseInfoInFile": updateList(&file.LicenseInfoInFile),
		"FileCopyrightText": update(&file.FileCopyrightText), //
		// "LicenseComments":   update(&file.LicenseComments),
		// "rdfs:comment":      update(&file.FileComment),
		// "FilenNoticeText":   update(&file.FileNotice),
		// "FileContributor":   updateList(&file.FileContributor),
		//snippet
	}
	return builder
}

func (p *Parser) MapDisjunctiveLicenseSet(dls *DisjunctiveLicenseSet) *builder {
	builder := &builder{t: typeDisjunctiveLicenseSet, ptr: dls}
	builder.updaters = map[string]updater{
		"member": updateList(&dls.Member),
	}
	return builder
}
