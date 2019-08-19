package rdf2v1

import (
	"strconv"
	"tools-golang/v0/spdx"
)

// Done
func CollectDocument(doc2v1 *spdx.Document2_1) *Document {

	stdDoc := Document{
		SPDXVersion:            Str(doc2v1.CreationInfo.SPDXVersion),
		DataLicense:            Str(doc2v1.CreationInfo.DataLicense),
		Review:                 collectReview(doc2v1),
		DocumentName:           Str(doc2v1.CreationInfo.DocumentName),
		DocumentNamespace:      Str(doc2v1.CreationInfo.DocumentNamespace),
		SPDXID:                 Str(doc2v1.CreationInfo.SPDXIdentifier),
		DocumentComment:        Str(doc2v1.CreationInfo.DocumentComment),
		ExtractedLicensingInfo: collectExtractedLicInfo(doc2v1),
		Relationship:           collectRelationships(doc2v1),
		CreationInfo:           collectCreationInfo(doc2v1),
		Annotation:             collectDocAnnotation(doc2v1),
		ExternalDocumentRef:    collectExternalDocumentRef(doc2v1),
		License:                collectLicense(doc2v1),
	}
	return &stdDoc
}

// Done
func collectCreationInfo(doc2v1 *spdx.Document2_1) *CreationInfo {

	stdCi := CreationInfo{

		SPDXIdentifier:     Str(doc2v1.CreationInfo.SPDXIdentifier),
		LicenseListVersion: Str(doc2v1.CreationInfo.LicenseListVersion),
		Creator:            InsertCreator(doc2v1.CreationInfo),
		Create:             Str(doc2v1.CreationInfo.Created),
		Comment:            Str(doc2v1.CreationInfo.CreatorComment),
	}
	return &stdCi
}

// Done
func collectExternalDocumentRef(doc2v1 *spdx.Document2_1) *ExternalDocumentRef {
	stdEdr := ExternalDocumentRef{

		ExternalDocumentId: Str(doc2v1.CreationInfo.ExternalDocumentReferences[0]),
		SPDXDocument:       Str(doc2v1.CreationInfo.ExternalDocumentReferences[1]),
		Checksum:           collectDocChecksum(doc2v1),
	}
	return &stdEdr
}

func collectLicense(doc2v1 *spdx.Document2_1) *License {
	stdEdr := License{
		LicenseId: Str(LicenseUri + doc2v1.CreationInfo.DataLicense),
	}
	return &stdEdr
}

// Done
func collectVerificationCode(pkg2_1 *spdx.Package2_1) *PackageVerificationCode {

	stdVc := PackageVerificationCode{

		PackageVerificationCode:             Str(pkg2_1.PackageVerificationCode),
		PackageVerificationCodeExcludedFile: Str(pkg2_1.PackageVerificationCodeExcludedFile),
	}
	return &stdVc
}

// Done
func collectPackageChecksum(pkg2_1 *spdx.Package2_1) *Checksum {

	stdPc := Checksum{

		Algorithm:     Str(PkgChecksumAlgo(pkg2_1)),
		ChecksumValue: Str(PkgChecksumValue(pkg2_1)),
	}
	return &stdPc
}

// Done
func collectFileChecksum(File2_1 *spdx.File2_1) *Checksum {

	stdFc := Checksum{

		Algorithm:     Str(FileChecksumAlgo(File2_1)),
		ChecksumValue: Str(FileChecksumValue(File2_1)),
	}
	return &stdFc
}

// Done
func collectDocChecksum(doc2v1 *spdx.Document2_1) *Checksum {

	stdFc := Checksum{

		Algorithm:     Str(doc2v1.CreationInfo.ExternalDocumentReferences[2]),
		ChecksumValue: Str(doc2v1.CreationInfo.ExternalDocumentReferences[3]),
	}
	return &stdFc
}

// Done
func collectDocAnnotation(doc2v1 *spdx.Document2_1) []*Annotation {
	var arrAnn []*Annotation
	for _, an := range doc2v1.Annotations {
		if an.AnnotationSPDXIdentifier == doc2v1.CreationInfo.SPDXIdentifier {
			stdAnn := Annotation{
				Annotator:                Str(an.AnnotatorType + an.Annotator),
				AnnotationType:           Str(an.AnnotationType),
				AnnotationDate:           Str(an.AnnotationDate),
				AnnotationComment:        Str(an.AnnotationComment),
				AnnotationSPDXIdentifier: Str(an.AnnotationSPDXIdentifier),
			}
			pointer := &stdAnn
			arrAnn = append(arrAnn, pointer)
		}
	}
	return arrAnn
}

// Done
func collectFileAnnotation(doc2v1 *spdx.Document2_1) []*Annotation {
	var arrAnn []*Annotation
	for _, pkg := range doc2v1.Packages {
		for _, file := range pkg.Files {
			for _, an := range doc2v1.Annotations {
				if an.AnnotationSPDXIdentifier == file.FileSPDXIdentifier {
					stdAnn := Annotation{
						Annotator:                Str(an.AnnotatorType + an.Annotator),
						AnnotationType:           Str(an.AnnotationType),
						AnnotationDate:           Str(an.AnnotationDate),
						AnnotationComment:        Str(an.AnnotationComment),
						AnnotationSPDXIdentifier: Str(an.AnnotationSPDXIdentifier),
					}
					pointer := &stdAnn
					arrAnn = append(arrAnn, pointer)
				}
			}
		}
	}
	return arrAnn
}

// Done
func collectPackageAnnotation(doc2v1 *spdx.Document2_1) []*Annotation {
	var arrAnn []*Annotation
	for _, pkg := range doc2v1.Packages {
		for _, an := range doc2v1.Annotations {
			if an.AnnotationSPDXIdentifier == pkg.PackageSPDXIdentifier {
				stdAnn := Annotation{
					Annotator:                Str(an.AnnotatorType + an.Annotator),
					AnnotationType:           Str(an.AnnotationType),
					AnnotationDate:           Str(an.AnnotationDate),
					AnnotationComment:        Str(an.AnnotationComment),
					AnnotationSPDXIdentifier: Str(an.AnnotationSPDXIdentifier),
				}
				pointer := &stdAnn
				arrAnn = append(arrAnn, pointer)
			}
		}
	}
	return arrAnn
}

// Done
func collectReview(doc2v1 *spdx.Document2_1) []*Review {
	var arrRev []*Review
	for _, a := range doc2v1.Reviews {
		if a != nil {
			stdRev := Review{
				Reviewer:      Str(a.Reviewer),
				ReviewDate:    Str(a.ReviewDate),
				ReviewComment: Str(a.ReviewComment),
			}
			pointer := &stdRev
			arrRev = append(arrRev, pointer)
		}
	}

	return arrRev
}

// Done
func collectRelationships(doc2v1 *spdx.Document2_1) []*Relationship {
	var arrRel []*Relationship
	for _, a := range doc2v1.Relationships {
		if a != nil {
			stdRel := Relationship{
				RelationshipType:    Str(a.Relationship),
				RelationshipComment: Str(a.RelationshipComment),
				Package:             collectPackages(doc2v1),
			}
			pointer := &stdRel
			arrRel = append(arrRel, pointer)
		}
	}

	return arrRel
}

func collectFilesfromPackages(doc2v1 *spdx.Document2_1) []*File {
	var arrFile []*File
	for _, a := range doc2v1.Packages {
		if a != nil {
			if a.Files != nil {
				for _, b := range a.Files {
					if b != nil {
						stdFile := File{

							FileName:           Str(b.FileName),
							FileSPDXIdentifier: Str(b.FileSPDXIdentifier),
							FileType:           ValueStrList(b.FileType),
							FileChecksum:       collectFileChecksum(b),
							// LicenseConcluded:   "", //DISCUSS
							LicenseInfoInFile:   ValueStrList(b.LicenseInfoInFile),
							FileLicenseComments: Str(b.LicenseComments),
							FileCopyrightText:   Str(b.FileCopyrightText),
							Project:             collectArtifactOfProject(doc2v1),
							FileComment:         Str(b.FileComment),
							FileNoticeText:      Str(b.FileNotice),
							FileContributor:     ValueStrList(b.FileContributor),
							// FileDependencies:   ,//DISCUSS
							Annotation: collectFileAnnotation(doc2v1),
							// ExtractedLicensingInfo
							// FileRelationship: coll
							// SnippetLicense,
						}
						pointer := &stdFile
						arrFile = append(arrFile, pointer)
					}
				}
			}
		}
	}
	return arrFile
}

func collectPackages(doc2v1 *spdx.Document2_1) []*Package {
	var arrPkg []*Package
	for _, a := range doc2v1.Packages {
		if a != nil {
			stdPkg := Package{
				PackageName:                 Str(a.PackageName),
				PackageVersionInfo:          Str(a.PackageVersion),
				PackageFileName:             Str(a.PackageFileName),
				PackageSPDXIdentifier:       Str(a.PackageSPDXIdentifier),
				PackageDownloadLocation:     Str(a.PackageDownloadLocation),
				PackageVerificationCode:     collectVerificationCode(a), //passing specific package
				PackageComment:              Str(a.PackageComment),
				PackageChecksum:             collectPackageChecksum(a),
				PackageLicenseComments:      Str(a.PackageLicenseComments),
				PackageLicenseInfoFromFiles: ValueStrList(a.PackageLicenseInfoFromFiles),
				PackageLicenseDeclared:      Str(a.PackageLicenseDeclared),
				PackageCopyrightText:        Str(a.PackageCopyrightText),
				PackageHomepage:             Str(a.PackageHomePage),
				PackageSupplier:             Str(InsertSupplier(a)),
				PackageExternalRef:          collectPkgExternalRef(a),
				PackageOriginator:           Str(InsertOriginator(a)),
				PackageSourceInfo:           Str(a.PackageSummary),
				FilesAnalyzed:               Str(strconv.FormatBool(a.FilesAnalyzed)),
				PackageSummary:              Str(a.PackageSummary),
				PackageDescription:          Str(a.PackageDescription),
				Annotation:                  collectPackageAnnotation(doc2v1),
				File:                        collectFilesfromPackages(doc2v1),
			}

			pointer := &stdPkg
			arrPkg = append(arrPkg, pointer)
		}
	}
	return arrPkg
}

// Done
func collectExtractedLicInfo(doc2v1 *spdx.Document2_1) []*ExtractedLicensingInfo {
	var arrEl []*ExtractedLicensingInfo
	for _, a := range doc2v1.OtherLicenses {
		if a != nil {
			stdEl := ExtractedLicensingInfo{
				LicenseIdentifier: Str(a.LicenseIdentifier),
				LicenseName:       Str(a.LicenseName),
				ExtractedText:     Str(a.ExtractedText),
				LicenseComment:    Str(a.LicenseComment),
				LicenseSeeAlso:    ValueStrList(a.LicenseCrossReferences),
			}
			pointer := &stdEl
			arrEl = append(arrEl, pointer)
		}
	}
	return arrEl
}

// Done
func collectArtifactOfProject(doc2v1 *spdx.Document2_1) []*Project {
	var arrp []*Project
	for _, a := range doc2v1.Packages {
		if a != nil {
			if a.Files != nil {
				for _, b := range a.Files {
					if b != nil {
						for _, c := range b.ArtifactOfProjects {
							stdp := Project{
								Name:     Str(c.Name),
								HomePage: Str(c.HomePage),
							}

							pointer := &stdp
							arrp = append(arrp, pointer)
						}
					}

				}
			}
		}
	}

	return arrp
}

func CollectSnippets(doc2v1 *spdx.Document2_1) *Snippet {
	for _, pkg := range doc2v1.Packages {
		if pkg != nil {
			for _, file := range pkg.Files {
				if file != nil {
					for _, sp := range file.Snippets {
						if sp != nil {
							stdSn := Snippet{
								SnippetName:             Str(sp.SnippetName),
								SnippetSPDXIdentifier:   Str(sp.SnippetSPDXIdentifier),
								SnippetLicenseComments:  Str(sp.SnippetLicenseComments),
								SnippetCopyrightText:    Str(sp.SnippetCopyrightText),
								SnippetLicenseConcluded: Str(sp.SnippetLicenseConcluded),
								SnippetComment:          Str(sp.SnippetComment),
								LicenseInfoInSnippet:    ValueStrList(sp.LicenseInfoInSnippet),
							}

							return &stdSn
						}
					}
				}
			}
		}
	}
	return nil
}

func collectPkgExternalRef(pkg *spdx.Package2_1) []*ExternalRef {
	var arrPer []*ExternalRef
	for _, a := range pkg.PackageExternalReferences {
		if a != nil {
			stdEl := ExternalRef{
				ReferenceLocator:  Str(a.Locator),
				ReferenceType:     collectReferenceType(a),
				ReferenceCategory: Str(a.Category),
				ReferenceComment:  Str(a.ExternalRefComment),
			}
			pointer := &stdEl
			arrPer = append(arrPer, pointer)
		}
	}
	return arrPer
}

func collectReferenceType(pkger *spdx.PackageExternalReference2_1) *ReferenceType {

	stdRt := ReferenceType{
		ReferenceType: Str(pkger.RefType),
	}
	return &stdRt
}

// func FindPackagefromFile(file *File) ValueStr {
// 	for key, value := range PackagetoFile {

// 		for _, f := range value {
// 			if f == file {
// 				return key
// 			}
// 		}

// 	}
// 	return Str("")
// }
