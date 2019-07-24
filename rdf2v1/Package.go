package rdf2v1

import (
	"github.com/deltamobile/goraptor"
)

type Package struct {
	PackageName                 ValueStr
	PackageVersionInfo          ValueStr
	PackageFileName             ValueStr
	PackageDownloadLocation     ValueStr
	PackageVerificationCode     *PackageVerificationCode
	PackageChecksum             []*Checksum
	PackageLicenseComments      ValueStr
	DisjunctiveLicenseSet       *DisjunctiveLicenseSet
	PackageLicenseInfoFromFiles []ValueStr
	PackageLicenseDeclared      ValueStr
	PackageCopyrightText        ValueStr
	File                        []*File
	PackageRelationship         *PackageRelationship
	PackageHomepage             ValueStr
	PackageSupplier             ValueStr
}
type PackageVerificationCode struct {
	PackageVerificationCode             ValueStr
	PackageVerificationCodeExcludedFile ValueStr
}
type PackageRelationship struct {
	Relationshiptype   ValueStr
	relatedSpdxElement ValueStr
}

func (p *Parser) requestPackage(node goraptor.Term) (*Package, error) {
	obj, err := p.requestElementType(node, typePackage)
	if err != nil {
		return nil, err
	}
	return obj.(*Package), err
}
func (p *Parser) requestPackageRelationship(node goraptor.Term) (*PackageRelationship, error) {
	obj, err := p.requestElementType(node, typeRelationship)
	if err != nil {
		return nil, err
	}
	return obj.(*PackageRelationship), err
}

func (p *Parser) requestPackageVerificationCode(node goraptor.Term) (*PackageVerificationCode, error) {
	obj, err := p.requestElementType(node, typePackageVerificationCode)
	if err != nil {
		return nil, err
	}
	return obj.(*PackageVerificationCode), err
}

func (p *Parser) MapPackage(pkg *Package) *builder {
	builder := &builder{t: typePackage, ptr: pkg}
	builder.updaters = map[string]updater{
		"name":             update(&pkg.PackageName),
		"versionInfo":      update(&pkg.PackageVersionInfo),
		"packageFileName":  update(&pkg.PackageFileName),
		"downloadLocation": update(&pkg.PackageDownloadLocation),
		"packageVerificationCode": func(obj goraptor.Term) error {
			pkgvc, err := p.requestPackageVerificationCode(obj)
			pkg.PackageVerificationCode = pkgvc
			return err
		},
		"checksum": func(obj goraptor.Term) error {
			pkgcksum, err := p.requestChecksum(obj)
			if err != nil {
				return err
			}
			pkg.PackageChecksum = append(pkg.PackageChecksum, pkgcksum)
			return err
		},
		"licenseComments": update(&pkg.PackageLicenseComments),
		"licenseConcluded": func(obj goraptor.Term) error {
			pkgdls, err := p.requestDisjunctiveLicenseSet(obj)
			pkg.DisjunctiveLicenseSet = pkgdls
			return err
		},
		"licenseDeclared":      update(&pkg.PackageLicenseDeclared),
		"licenseInfoFromFiles": updateList(&pkg.PackageLicenseInfoFromFiles),
		"copyrightText":        update(&pkg.PackageCopyrightText),
		"hasFile": func(obj goraptor.Term) error {
			file, err := p.requestFile(obj)
			if err != nil {
				return err
			}
			pkg.File = append(pkg.File, file)
			return nil
		},
		"relationship": func(obj goraptor.Term) error {
			rel, err := p.requestPackageRelationship(obj)
			pkg.PackageRelationship = rel
			return err
		},
		"doap:homepage": update(&pkg.PackageHomepage),
		"supplier":      update(&pkg.PackageSupplier),
	}
	return builder
}

func (p *Parser) MapPackageVerificationCode(pkgvc *PackageVerificationCode) *builder {
	builder := &builder{t: typePackageVerificationCode, ptr: pkgvc}
	builder.updaters = map[string]updater{
		"packageVerificationCodeValue":        update(&pkgvc.PackageVerificationCode),
		"packageVerificationCodeExcludedFile": update(&pkgvc.PackageVerificationCodeExcludedFile),
	}
	return builder
}
func (p *Parser) MapPackageRelationship(pkgrel *PackageRelationship) *builder {
	builder := &builder{t: typePackageVerificationCode, ptr: pkgrel}
	builder.updaters = map[string]updater{
		"relationshipType":   update(&pkgrel.Relationshiptype),
		"relatedSpdxElement": update(&pkgrel.Relationshiptype),
	}
	return builder
}
