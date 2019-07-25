package main

import (
	"spdx/tools-golang/v0/spdx"
	"ug/understanding-goraptor/rdf2v1"
)

// func transferDocument(spdxdoc *rdf2v1.Document) *spdx.Document2_1 {

// 	stdDoc := spdx.Document2_1{

// 		CreationInfo:  transferCreationInfo(spdxdoc),
// 		Packages:      transferPackages(spdxdoc),
// 		OtherLicenses: transferOtherLicenses(spdxdoc),
// 		Relationships: transferRelationships(spdxdoc),
// 		Annotations:   transferAnnotation(spdxdoc),
// 		Reviews:       transferReview(spdxdoc),
// 	}
// 	return &stdDoc
// }

func transferCreationInfo(spdxdoc *rdf2v1.Document) *spdx.CreationInfo2_1 {

	var listExtDocRef []string
	listExtDocRef = append(listExtDocRef, spdxdoc.ExternalDocumentRef.ExternalDocumentId.V())
	listExtDocRef = append(listExtDocRef, spdxdoc.ExternalDocumentRef.SPDXDocument.V())
	listExtDocRef = append(listExtDocRef, spdxdoc.ExternalDocumentRef.Checksum.Algorithm.V())

	stdCi := spdx.CreationInfo2_1{

		SPDXVersion:                spdxdoc.SPDXVersion.Val,
		DataLicense:                spdxdoc.DataLicense.Val,
		SPDXIdentifier:             "",
		DocumentName:               spdxdoc.DocumentName.Val,
		DocumentNamespace:          "",
		ExternalDocumentReferences: listExtDocRef,
		LicenseListVersion:         "",
		Created:                    spdxdoc.CreationInfo.Create.Val,
		CreatorComment:             spdxdoc.CreationInfo.Comment.Val,
		DocumentComment:            spdxdoc.DocumentComment.Val,
	}
	return &stdCi
}

func transferAnnotation(spdxdoc *rdf2v1.Document) *spdx.Annotation2_1 {

	stdAnn := spdx.Annotation2_1{
		Annotator:                spdxdoc.Annotation.Annotator.Val,
		AnnotationType:           spdxdoc.Annotation.AnnotationType.Val,
		AnnotationDate:           spdxdoc.Annotation.AnnotationDate.Val,
		AnnotationComment:        spdxdoc.Annotation.AnnotationComment.Val,
		AnnotationSPDXIdentifier: "",
		AnnotatorType:            "",
	}

	return &stdAnn
}

// Discuss location
// func transferReview(spdxdoc *rdf2v1.Document) *([]spdx.Review2_1) {

// 	stdRev := spdx.Review2_1{
// 		Reviewer:      spdxdoc.Review.Reviewer.Val,
// 		ReviewerType:  "",
// 		ReviewDate:    spdxdoc.Review.ReviewDate.Val,
// 		ReviewComment: spdxdoc.Review.ReviewComment.Val,
// 	}

// 	return &stdRev
// }

func transferFile(spdxdoc *rdf2v1.Document) []*spdx.File2_1 {
	var arrFile []spdx.File2_1
	for i, _ := range spdxdoc.Relationship {
		k := spdxdoc.Relationship[i]
		for _, b := range k.File {
			file := b

			stdFile := spdx.File2_1{

				FileName:           file.FileName.Val,
				FileSPDXIdentifier: "",
				FileType:           rdf2v1.ValueList(file.FileType),
				FileChecksumSHA1:   "",
				FileChecksumSHA256: "",
				FileChecksumMD5:    "",
				// LicenseConcluded:   "",
				LicenseInfoInFile: rdf2v1.ValueList(file.LicenseInfoInFile),
				LicenseComments:   file.FileLicenseComments.V(),
				FileCopyrightText: file.FileCopyrightText.V(),
				// ArtifactOfProjects: file.Annotation.Annotator.Val,
				FileComment:     file.FileComment.Val,
				FileNotice:      file.FileNoticeText.Val,
				FileContributor: rdf2v1.ValueList(file.FileContributor),
				// FileDependencies:   "",//DISCUSS
				// Snippets:           "",//DISCUSS
			}
			arrFile = append(arrFile, stdFile)
		}
	}
	return arrFile
}

func transferPackages(spdxdoc *rdf2v1.Document) *[]*spdx.Package2_1 {
	var arrPkg []spdx.Package2_1
	for i, _ := range spdxdoc.Relationship {
		k := spdxdoc.Relationship[i]
		for _, b := range k.Package {
			pkg := b

			stdPkg := spdx.Package2_1{
				// IsUnpackaged: "",
				PackageName:           pkg.PackageName.Val,
				PackageSPDXIdentifier: "",
				PackageVersion:        pkg.PackageVersionInfo.Val,
				PackageFileName:       pkg.PackageFileName.Val,

				PackageSupplierPerson:       "",
				PackageSupplierOrganization: "",
				// PackageSupplierNOASSERTION:  "",

				PackageOriginatorPerson:       "",
				PackageOriginatorOrganization: "",
				// PackageOriginatorNOASSERTION:  "",

				PackageDownloadLocation: pkg.PackageDownloadLocation.Val,
				// FilesAnalyzed:                      "",
				// IsFilesAnalyzedTagPresent:          "",
				PackageVerificationCode:             pkg.Annotation.AnnotationComment.Val,
				PackageVerificationCodeExcludedFile: "", //DISCUSS
				PackageChecksumSHA1:                 pkg.PackageName.Val,
				PackageChecksumSHA256:               "",
				PackageChecksumMD5:                  "",
				PackageHomePage:                     "",
				PackageSourceInfo:                   "",
				PackageLicenseConcluded:             "",
				// PackageLicenseInfoFromFiles: "",
				PackageLicenseDeclared: pkg.PackageLicenseDeclared.Val,
				PackageLicenseComments: pkg.PackageLicenseComments.Val,
				PackageCopyrightText:   pkg.PackageCopyrightText.Val,
				PackageSummary:         pkg.PackageSummary.Val,
				PackageDescription:     pkg.PackageDescription.Val,
				PackageComment:         "", //not in rdf
				// PackageExternalReferences:   "",
				Files: transferFile(spdxdoc),
			}
			arrPkg = append(arrPkg, stdPkg)
		}
	}
	return &arrPkg
}

func transferOtherLicenses(spdxdoc *rdf2v1.Document) []*spdx.OtherLicense2_1 {

	stdOl := spdx.Review2_1{
		LicenseIdentifier:      spdxdoc.Review.Reviewer.Val,
		ExtractedText:          "",
		LicenseName:            spdxdoc.Review.ReviewDate.Val,
		LicenseCrossReferences: spdxdoc.Review.ReviewComment.Val,
		LicenseComment:         "",
	}

	return &stdOl
}

func transferArtifactOfProject(spdxdoc *rdf2v1.Document) []*spdx.ArtifactOfProject2_1 {
	var arrAop []*spdx.ArtifactOfProject2_1
	for i, _ := range spdxdoc.Relationship {
		k := spdxdoc.Relationship[i]
		for _, b := range k.File {
			file := b
			for _, a := range file.Project {

				stdAop := spdx.ArtifactOfProject2_1{
					Name:     a.Name.Val,
					HomePage: a.Homepage.Val,
					URI:      "",
				}

				pointer := &stdAop
				arrAop = append(arrAop, pointer)
			}
		}
	}
	return arrAop
}
