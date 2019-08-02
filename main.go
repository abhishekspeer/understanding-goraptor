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
	var sp []*rdf2v1.Snippet
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
	fmt.Printf("%T\n\n", sp)
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

	fmt.Println("\n\nFINAL TRANSLATED DOCUMENT")
	// fmt.Printf("%v\n", spdxdoc.DocumentNamespace)
	// var doc2v1 *spdx.Document2_1
	doc2v1 := TransferDocument(spdxdoc)
	fmt.Printf("%#v\n", doc2v1)

}

func Parse(input string) (*rdf2v1.Document, []*rdf2v1.Snippet, error) {
	parser := rdf2v1.NewParser(input)
	defer fmt.Printf("RDF Document parsed successfully.\n")
	defer parser.Free()
	return parser.Parse()
}
func TransferDocument(spdxdoc *rdf2v1.Document) *spdx.Document2_1 {

	stdDoc := spdx.Document2_1{

		CreationInfo:  transferCreationInfo(spdxdoc),
		Packages:      transferPackages(spdxdoc),
		OtherLicenses: transferOtherLicenses(spdxdoc),
		Relationships: transferRelationships(spdxdoc),
		Annotations:   transferAnnotation(spdxdoc),
		Reviews:       transferReview(spdxdoc),
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
		LicenseListVersion:         spdxdoc.CreationInfo.LicenseListVersion.Val,
		CreatorPersons:             rdf2v1.ExtractCreator(spdxdoc.CreationInfo, "Person"),
		CreatorOrganizations:       rdf2v1.ExtractCreator(spdxdoc.CreationInfo, "Organization"),
		CreatorTools:               rdf2v1.ExtractCreator(spdxdoc.CreationInfo, "Tool"),
		Created:                    spdxdoc.CreationInfo.Create.Val,
		CreatorComment:             spdxdoc.CreationInfo.Comment.Val,
		DocumentComment:            spdxdoc.DocumentComment.Val,
	}
	return &stdCi
}

func transferAnnotation(spdxdoc *rdf2v1.Document) []*spdx.Annotation2_1 {
	var arrAnn []*spdx.Annotation2_1
	for _, an := range spdxdoc.Annotation {
		stdAnn := spdx.Annotation2_1{
			Annotator:                rdf2v1.ExtractKeyValue(an.Annotator.Val, "subvalue"),
			AnnotatorType:            rdf2v1.ExtractKeyValue(an.Annotator.Val, "subkey"),
			AnnotationType:           an.AnnotationType.Val,
			AnnotationDate:           an.AnnotationDate.Val,
			AnnotationComment:        an.AnnotationComment.Val,
			AnnotationSPDXIdentifier: "",
		}
		pointer := &stdAnn
		arrAnn = append(arrAnn, pointer)
	}

	for _, a := range spdxdoc.Relationship {
		if a != nil {
			if a.Package != nil {
				for _, b := range a.Package {
					if b != nil {
						for _, an := range b.Annotation {

							stdAnn := spdx.Annotation2_1{
								Annotator:                rdf2v1.ExtractKeyValue(an.Annotator.Val, "subvalue"),
								AnnotatorType:            rdf2v1.ExtractKeyValue(an.Annotator.Val, "subkey"),
								AnnotationType:           an.AnnotationType.Val,
								AnnotationDate:           an.AnnotationDate.Val,
								AnnotationComment:        an.AnnotationComment.Val,
								AnnotationSPDXIdentifier: "",
							}
							pointer := &stdAnn
							arrAnn = append(arrAnn, pointer)
						}
					}
				}
			}
		}

	}
	for _, a := range spdxdoc.Relationship {
		if a != nil {
			if a.File != nil {
				for _, b := range a.File {
					if b != nil {
						for _, an := range b.Annotation {

							stdAnn := spdx.Annotation2_1{
								Annotator:                rdf2v1.ExtractKeyValue(an.Annotator.Val, "subvalue"),
								AnnotatorType:            rdf2v1.ExtractKeyValue(an.Annotator.Val, "subkey"),
								AnnotationType:           an.AnnotationType.Val,
								AnnotationDate:           an.AnnotationDate.Val,
								AnnotationComment:        an.AnnotationComment.Val,
								AnnotationSPDXIdentifier: "",
							}
							pointer := &stdAnn
							arrAnn = append(arrAnn, pointer)
						}
					}
				}
			}
		}

	}
	return arrAnn
}

func transferReview(spdxdoc *rdf2v1.Document) []*spdx.Review2_1 {
	var arrRev []*spdx.Review2_1
	for _, a := range spdxdoc.Review {
		if a != nil {
			stdRev := spdx.Review2_1{
				Reviewer:      a.Reviewer.Val,
				ReviewerType:  rdf2v1.ExtractKey(a.Reviewer.Val),
				ReviewDate:    a.ReviewDate.Val,
				ReviewComment: a.ReviewComment.Val,
			}
			pointer := &stdRev
			arrRev = append(arrRev, pointer)
		}
	}

	return arrRev
}

func transferRelationships(spdxdoc *rdf2v1.Document) []*spdx.Relationship2_1 {
	var arrRel []*spdx.Relationship2_1
	for _, a := range spdxdoc.Relationship {
		if a != nil {
			stdRel := spdx.Relationship2_1{
				Relationship:        a.RelationshipType.Val,
				RelationshipComment: a.RelationshipComment.Val,
			}
			pointer := &stdRel
			arrRel = append(arrRel, pointer)
		}
	}

	return arrRel
}
func transferFile(spdxdoc *rdf2v1.Document) []*spdx.File2_1 {
	var arrFile []*spdx.File2_1
	for _, a := range spdxdoc.Relationship {
		if a != nil {
			if a.File != nil {
				for _, b := range a.File {
					if b != nil {
						stdFile := spdx.File2_1{

							FileName:           b.FileName.Val,
							FileSPDXIdentifier: "",
							FileType:           rdf2v1.ValueList(b.FileType),
							FileChecksumSHA1:   rdf2v1.AlgoIdentifier(b.FileChecksum, "SHA1"), // DISCUSS
							FileChecksumSHA256: rdf2v1.AlgoIdentifier(b.FileChecksum, "SHA256"),
							FileChecksumMD5:    rdf2v1.AlgoIdentifier(b.FileChecksum, "MD5"),
							// LicenseConcluded:   "", //DISCUSS
							LicenseInfoInFile:  rdf2v1.ValueList(b.LicenseInfoInFile),
							LicenseComments:    b.FileLicenseComments.V(),
							FileCopyrightText:  b.FileCopyrightText.V(),
							ArtifactOfProjects: transferArtifactOfProject(spdxdoc),
							FileComment:        b.FileComment.Val,
							FileNotice:         b.FileNoticeText.Val,
							FileContributor:    rdf2v1.ValueList(b.FileContributor),
							// FileDependencies:   "",//DISCUSS
							// Snippets:           "",//DISCUSS
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

func transferPackages(spdxdoc *rdf2v1.Document) []*spdx.Package2_1 {
	var arrPkg []*spdx.Package2_1
	for _, a := range spdxdoc.Relationship {
		if a != nil {
			if a.Package != nil {
				for _, b := range a.Package {

					if b != nil {
						stdPkg := spdx.Package2_1{
							IsUnpackaged:          b.PackageName.Val == "",
							PackageName:           b.PackageName.Val,
							PackageSPDXIdentifier: "",
							PackageVersion:        b.PackageVersionInfo.Val,
							PackageFileName:       b.PackageFileName.Val,

							PackageSupplierPerson:       rdf2v1.ExtractValueType(b.PackageSupplier.Val, "Person"),
							PackageSupplierOrganization: rdf2v1.ExtractValueType(b.PackageSupplier.Val, "Organization"),
							PackageSupplierNOASSERTION:  b.PackageSupplier.Val == "NOASSERTION",

							PackageOriginatorPerson:       rdf2v1.ExtractValueType(b.PackageOriginator.Val, "Person"),
							PackageOriginatorOrganization: rdf2v1.ExtractValueType(b.PackageOriginator.Val, "Organization"),
							PackageOriginatorNOASSERTION:  b.PackageOriginator.Val == "NOASSERTION",

							PackageDownloadLocation:             b.PackageDownloadLocation.Val,
							FilesAnalyzed:                       !(b.PackageName.Val == ""),
							IsFilesAnalyzedTagPresent:           b.PackageName.Val == "",
							PackageVerificationCode:             b.PackageVerificationCode.PackageVerificationCode.Val,
							PackageVerificationCodeExcludedFile: b.PackageVerificationCode.PackageVerificationCodeExcludedFile.Val,
							// PackageChecksumSHA1:                 rdf2v1.AlgoIdentifier(b.PackageChecksum, "SHA1"),
							// PackageChecksumSHA256:               rdf2v1.AlgoIdentifier(b.PackageChecksumMD5, "SHA256"),
							// PackageChecksumMD5:                  rdf2v1.AlgoIdentifier(b.PackageChecksum, "MD5"),
							PackageHomePage:         b.PackageHomepage.Val,
							PackageSourceInfo:       b.PackageSourceInfo.Val,
							PackageLicenseConcluded: "",
							// PackageLicenseInfoFromFiles: //DISCUSS
							PackageLicenseDeclared: b.PackageLicenseDeclared.Val,
							PackageLicenseComments: b.PackageLicenseComments.Val,
							PackageCopyrightText:   b.PackageCopyrightText.Val,
							PackageSummary:         b.PackageSummary.Val,
							PackageDescription:     b.PackageDescription.Val,
							PackageComment:         b.PackageComment.Val,
							// PackageExternalReferences:   "", //DISCUSS
							Files: transferFile(spdxdoc),
						}

						pointer := &stdPkg
						arrPkg = append(arrPkg, pointer)
					}
				}
			}
		}
	}
	return arrPkg
}

func transferOtherLicenses(spdxdoc *rdf2v1.Document) []*spdx.OtherLicense2_1 {
	var arrOl []*spdx.OtherLicense2_1
	for _, a := range spdxdoc.ExtractedLicensingInfo {
		if a != nil {
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
	}
	return arrOl
}

func transferArtifactOfProject(spdxdoc *rdf2v1.Document) []*spdx.ArtifactOfProject2_1 {
	var arrAop []*spdx.ArtifactOfProject2_1
	for _, a := range spdxdoc.Relationship {
		if a != nil {
			if a.File != nil {
				for _, b := range a.File {
					if b != nil {
						for _, c := range b.Project {
							stdAop := spdx.ArtifactOfProject2_1{
								Name:     c.Name.Val,
								HomePage: c.Homepage.Val,
								URI:      "",
							}

							pointer := &stdAop
							arrAop = append(arrAop, pointer)
						}
					}

				}
			}
		}
	}

	return arrAop
}

// func transferSnippets(sp []*rdf2v1.Snippet) []*spdx.Snippet2_1 {
// 	var arrSn []*spdx.Snippet2_1
// 	for _, a := range sp {
// 		if a != nil {
// 			stdSn := spdx.Snippet2_1{
// 				SnippetLicenseComments:  a.SnippetLicenseComments.Val,
// 				SnippetCopyrightText:    a.SnippetCopyrightText.Val,
// 				SnippetLicenseConcluded: a.SnippetLicenseConcluded.Val, //DISCUSS: Not in RDF file
// 				SnippetComment:          a.SnippetComment.Val,
// 				LicenseInfoInSnippet:    a.LicenseName,                 // DISCUSS: more than one fields in RDF but string in standard struct
// 			}
// 			pointer := &stdSn
// 			arrSn = append(arrSn, pointer)
// 		}
// 	}
// 	return arrSn
// }
