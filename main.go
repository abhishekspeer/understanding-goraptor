package main

import (
	"fmt"
	"os"
	"spdx/tools-golang/v0/spdx"
	"ug/understanding-goraptor/rdf2v1"
)

func main() {
	// check that we've received the right number of arguments

	args := os.Args
	if len(args) != 2 {
		fmt.Printf("Usage: %v <spdx-file-in>\n", args[0])
		fmt.Printf("  Load SPDX 2.1 RDF file <spdx-file-in>, and\n")
		fmt.Printf("  print its contents.\n")
		return
	}
	var spdxdoc *rdf2v1.Document
	var sp *rdf2v1.Snippet
	var err error

	input := args[1]
	spdxdoc, sp, err = Parse(input)

	if err != nil {
		fmt.Println("Parsing Error")
		return
	}

	// fmt.Println("===================================================\n")
	// fmt.Println("Some Information Printed from the Document Returned\n")
	// fmt.Println("===================================================\n")
	fmt.Println("%T\n\n", sp)
	// // fmt.Printf("Relationship: %v\n\n", spdxdoc.Relationship[0].Package[0])
	// fmt.Printf("\nRelationship: %v\n\n", spdxdoc.Relationship[3].File)
	// // fmt.Printf("Relationship: %v\n\n", spdxdoc.Relationship[2])
	// // fmt.Printf("Relationship: %v\n\n", spdxdoc.Relationship[3])
	// // fmt.Printf("SpecVersion: %v\n\n", spdxdoc.SPDXVersion.Val)
	// fmt.Printf("\n\nCreationInfo Creator: %v\n\n", spdxdoc.CreationInfo)
	// fmt.Printf("CreationInfo Create:%v\n\n", spdxdoc.CreationInfo.Create)
	// fmt.Printf("DocumentName: %v\n\n", spdxdoc.DocumentName.Val)
	// fmt.Printf("DocumentComment: %v\n\n", spdxdoc.DocumentComment)

	// iniEdr := spdxdoc.ExternalDocumentRef
	// intRel := spdxdoc.Relationship
	// stdRel := make([]spdx.Relationship2_1, 4)
	// fmt.Printf("Docummment: %v\n\n", stdRel[0])

	fmt.Println("FINAL TRANSLATED DOCUMENT")
	// var doc2v1 *spdx.Document2_1
	doc2v1 := TransferDocument(spdxdoc)
	fmt.Printf("%T", doc2v1.Packages[0].Files[0])

}

func Parse(input string) (*rdf2v1.Document, *rdf2v1.Snippet, error) {
	parser := rdf2v1.NewParser(input)
	defer fmt.Printf("RDF Document parsed successfully.\n")
	defer parser.Free()
	return parser.Parse()
}

// func transferDocument(spdxdoc *rdf2v1.Document) *spdx.Document2_1 {

// 	stdDoc := spdx.Document2_1{

// 		CreationInfo:  transferCreationInfo(spdxdoc),
// 		Packages:      transferPackages(spdxdoc),
// 		OtherLicenses: transferOtherLicenses(spdxdoc),
// 		// Relationships: transferRelationships(spdxdoc),
// 		// Annotations:   transferAnnotation(spdxdoc),
// 		// Reviews:       transferReview(spdxdoc),
// 	}
// 	return &stdDoc
// }
func TransferDocument(spdxdoc *rdf2v1.Document) *spdx.Document2_1 {

	stdDoc := spdx.Document2_1{

		CreationInfo:  transferCreationInfo(spdxdoc),
		Packages:      transferPackages(spdxdoc),
		OtherLicenses: transferOtherLicenses(spdxdoc),
		// Relationships: transferRelationships(spdxdoc),
		// Annotations:   transferAnnotation(spdxdoc),
		// Reviews:       transferReview(spdxdoc),
	}
	return &stdDoc
}

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
	var arrFile []*spdx.File2_1
	for i, _ := range spdxdoc.Relationship {
		k := spdxdoc.Relationship[i]
		for _, b := range k.File {
			file := b

			stdFile := spdx.File2_1{

				FileName:           file.FileName.Val,
				FileSPDXIdentifier: "",
				FileType:           rdf2v1.ValueList(file.FileType),
				FileChecksumSHA1:   "", // DISCUSS
				FileChecksumSHA256: "",
				FileChecksumMD5:    "",
				// LicenseConcluded:   "", //DISCUSS
				LicenseInfoInFile:  rdf2v1.ValueList(file.LicenseInfoInFile),
				LicenseComments:    file.FileLicenseComments.V(),
				FileCopyrightText:  file.FileCopyrightText.V(),
				ArtifactOfProjects: transferArtifactOfProject(spdxdoc),
				FileComment:        file.FileComment.Val,
				FileNotice:         file.FileNoticeText.Val,
				FileContributor:    rdf2v1.ValueList(file.FileContributor),
				// FileDependencies:   "",//DISCUSS
				// Snippets:           "",//DISCUSS
			}
			pointer := &stdFile
			arrFile = append(arrFile, pointer)
		}
	}
	return arrFile
}

func transferPackages(spdxdoc *rdf2v1.Document) []*spdx.Package2_1 {
	var arrPkg []*spdx.Package2_1
	for _, a := range spdxdoc.Relationship {
		rel := a
		for _, b := range rel.Package {
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
			pointer := &stdPkg
			arrPkg = append(arrPkg, pointer)
		}
	}
	return arrPkg
}

func transferOtherLicenses(spdxdoc *rdf2v1.Document) []*spdx.OtherLicense2_1 {
	var arrOl []*spdx.OtherLicense2_1
	for _, a := range spdxdoc.ExtractedLicensingInfo {
		stdOl := spdx.OtherLicense2_1{
			LicenseIdentifier: a.LicenseIdentifier.Val,
			ExtractedText:     a.ExtractedText.Val,
			// LicenseName:            a.LicenseName,// DISCUSS: more than one fields in RDF but string in standard struct
			// LicenseCrossReferences: rdf2v1.ValueList(a),//DISCUSS: Not in RDF file
			LicenseComment: a.LicenseComment.Val,
		}
		pointer := &stdOl
		arrOl = append(arrOl, pointer)
	}
	return arrOl
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
