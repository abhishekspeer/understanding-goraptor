package rdf2v1

import "github.com/deltamobile/goraptor"

type Package struct {
	IsUnpackaged                        ValueBool
	PackageName                         ValueStr
	PackageSPDXIdentifier               ValueStr
	PackageVersion                      ValueStr
	PackageFileName                     ValueStr
	PackageSupplierPerson               ValueStr
	PackageSupplierOrganization         ValueStr
	PackageSupplierNOASSERTION          ValueBool
	PackageOriginatorPerson             ValueStr
	PackageOriginatorOrganization       ValueStr
	PackageOriginatorNOASSERTION        ValueBool
	PackageDownloadLocation             ValueStr
	FilesAnalyzed                       ValueBool
	IsFilesAnalyzedTagPresent           ValueBool
	PackageVerificationCode             ValueStr
	PackageVerificationCodeExcludedFile ValueStr
	PackageChecksumSHA1                 ValueStr
	PackageChecksumSHA256               ValueStr
	PackageChecksumMD5                  ValueStr
	PackageHomePage                     ValueStr
	PackageSourceInfo                   ValueStr
	PackageLicenseConcluded             ValueStr
	PackageLicenseInfoFromFiles         []ValueStr
	PackageLicenseDeclared              ValueStr
	PackageLicenseComments              ValueStr
	PackageCopyrightText                ValueStr
	PackageSummary                      ValueStr
	PackageDescription                  ValueStr
	PackageComment                      ValueStr
	PackageExternalReferences           []*PackageExternalReference
	Files                               []*File
}
type PackageExternalReference struct {
	Category           ValueStr
	RefType            ValueStr
	Locator            ValueStr
	ExternalRefComment ValueStr
}

func (p *Parser) requestPackage(node goraptor.Term) (*Package, error) {
	obj, err := p.requestElementType(node, typePackage)
	if err != nil {
		return nil, err
	}
	return obj.(*Package), err
}
