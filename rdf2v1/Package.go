package rdf2v1

import (
	"github.com/deltamobile/goraptor"
)

type Package struct {
	PackageName                 ValueStr
	PackageFileName             ValueStr
	PackageDownloadLocation     ValueStr
	PackageVerificationCode     *PackageVerificationCode
	PackageChecksum             *Checksum
	PackageLicenseComments      ValueStr
	PackageLicenseConcluded     []ValueStr
	PackageLicenseInfoFromFiles ValueStr
	PackageLicenseDeclared      ValueStr
	PackageCopyrightText        ValueStr
	Files                       []*File
}
type PackageVerificationCode struct {
	PackageVerificationCode ValueStr
}

func (p *Parser) requestPackage(node goraptor.Term) (*Package, error) {
	obj, err := p.requestElementType(node, typePackage)
	if err != nil {
		return nil, err
	}
	return obj.(*Package), err
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
		"packageFileName":  update(&pkg.PackageFileName),
		"downloadLocation": update(&pkg.PackageDownloadLocation),
		"packageVerificationCode": func(obj goraptor.Term) error {
			pkgvc, err := p.requestPackageVerificationCode(obj)
			pkg.PackageVerificationCode = pkgvc
			return err
		},
		"checksum": func(obj goraptor.Term) error {
			pkgcksum, err := p.requestChecksum(obj)
			pkg.PackageChecksum = pkgcksum
			return err
		},
		"licenseComments":      update(&pkg.PackageLicenseComments),
		"licenseConcluded":     updateList(&pkg.PackageLicenseConcluded),
		"licenseDeclared":      update(&pkg.PackageLicenseDeclared),
		"licenseInfoFromFiles": update(&pkg.PackageLicenseInfoFromFiles),
		"copyrightText":        update(&pkg.PackageCopyrightText),
		"hasFile": func(obj goraptor.Term) error {
			file, err := p.requestFile(obj)
			if err != nil {
				return err
			}
			pkg.Files = append(pkg.Files, file)
			return nil
		},
	}
	return builder
}

func (p *Parser) MapPackageVerificationCode(pkgvc *PackageVerificationCode) *builder {
	builder := &builder{t: typePackageVerificationCode, ptr: pkgvc}
	builder.updaters = map[string]updater{
		"packageVerificationCodeValue": update(&pkgvc.PackageVerificationCode),
	}
	return builder
}
