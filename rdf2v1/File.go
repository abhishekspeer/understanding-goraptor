package rdf2v1

import (
	"fmt"

	"github.com/deltamobile/goraptor"
)

type File struct {
	FileName              ValueStr
	FileChecksum          *Checksum
	LicenseInfoInFile     []ValueStr
	FileCopyrightText     ValueStr
	DisjunctiveLicenseSet *DisjunctiveLicenseSet
	ConjunctiveLicenseSet *ConjunctiveLicenseSet
	FileContributor       []ValueStr
	FileComment           ValueStr
	FileLicenseComments   ValueStr
	FileType              ValueStr
	FileNoticeText        ValueStr
	Annotation            *Annotation
	Project               *Project
	SnippetLicense        *License
	// //Snippets 			[]*Snippet

}
type Project struct {
	Homepage ValueStr
	Name     ValueStr
}

type DisjunctiveLicenseSet struct {
	Member []ValueStr
}
type ConjunctiveLicenseSet struct {
	License                *License
	ExtractedLicensingInfo *ExtractedLicensingInfo
}

// ERROR
func (p *Parser) requestFile(node goraptor.Term) (*File, error) {
	obj, err := p.requestElementType(node, typeFile)
	if err != nil {
		fmt.Println("TTTTTTTTT")
		return nil, err
	}
	return obj.(*File), err
}
func (p *Parser) requestFileChecksum(node goraptor.Term) (*Checksum, error) {
	obj, err := p.requestElementType(node, typeChecksum)
	if err != nil {
		return nil, err
	}
	return obj.(*Checksum), err
}
func (p *Parser) requestDisjunctiveLicenseSet(node goraptor.Term) (*DisjunctiveLicenseSet, error) {
	obj, err := p.requestElementType(node, typeDisjunctiveLicenseSet)
	if err != nil {
		return nil, err
	}
	return obj.(*DisjunctiveLicenseSet), err
}
func (p *Parser) requestConjunctiveLicenseSet(node goraptor.Term) (*ConjunctiveLicenseSet, error) {
	obj, err := p.requestElementType(node, typeConjunctiveLicenseSet)
	if err != nil {
		return nil, err
	}
	return obj.(*ConjunctiveLicenseSet), err
}
func (p *Parser) requestProject(node goraptor.Term) (*Project, error) {
	obj, err := p.requestElementType(node, typeProject)
	if err != nil {
		return nil, err
	}
	return obj.(*Project), err
}

// Returns a builder for file.
func (p *Parser) MapFile(file *File) *builder {
	builder := &builder{t: typeFile, ptr: file}
	builder.updaters = map[string]updater{
		"fileName": update(&file.FileName),
		"checksum": func(obj goraptor.Term) error {
			cksum, err := p.requestChecksum(obj)
			file.FileChecksum = cksum
			return err
		},
		"fileType": updateTrimPrefix(baseUri, &file.FileType),
		"licenseConcluded": func(obj goraptor.Term) error {
			switch {
			case obj == typeLicense:
				lic, err := p.requestLicense(obj)
				if err != nil {
					return err
				}
				file.SnippetLicense = lic
			case obj == typeDisjunctiveLicenseSet:
				dls, err := p.requestDisjunctiveLicenseSet(obj)
				if err != nil {
					return err
				}
				file.DisjunctiveLicenseSet = dls
			}
			return nil
		},
		"licenseInfoInFile": updateList(&file.LicenseInfoInFile),
		"copyrightText":     update(&file.FileCopyrightText), //
		"licenseComments":   update(&file.FileLicenseComments),
		"rdfs:comment":      update(&file.FileComment),
		"noticeText":        update(&file.FileNoticeText),
		"fileContributor":   updateList(&file.FileContributor),
		"annotation": func(obj goraptor.Term) error {
			an, err := p.requestAnnotation(obj)
			file.Annotation = an
			return err
		},
		"artifactOf": func(obj goraptor.Term) error {
			pro, err := p.requestProject(obj)
			file.Project = pro
			return err
		},
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

func (p *Parser) MapProject(pro *Project) *builder {
	builder := &builder{t: typeProject, ptr: pro}
	builder.updaters = map[string]updater{
		"doap:homepage": update(&pro.Homepage),
		"doap:name":     update(&pro.Name),
	}
	return builder
}
func (p *Parser) MapConjunctiveLicenseSet(cls *ConjunctiveLicenseSet) *builder {
	builder := &builder{t: typeConjunctiveLicenseSet, ptr: cls}
	builder.updaters = map[string]updater{
		"member": func(obj goraptor.Term) error {
			switch {
			case obj == typeLicense:
				lic, err := p.requestLicense(obj)
				if err != nil {
					return err
				}
				cls.License = lic
			case obj == typeExtractedLicensingInfo:
				eli, err := p.requestExtractedLicensingInfo(obj)
				if err != nil {
					return err
				}
				cls.ExtractedLicensingInfo = eli
			}
			return nil
		},
	}
	return builder
}
