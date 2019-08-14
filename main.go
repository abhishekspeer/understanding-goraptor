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
	fmt.Printf("FFFFFFFFFFFFFFFFFF%#v\n\n\n", spdxdoc.SPDXID)

	if err != nil {
		fmt.Println("Parsing Error")
		return
	}
	doc2v1 := TransferDocument(spdxdoc, sp)
	fmt.Printf("%T\n\n\n", doc2v1)

	// WRITER
	output := os.Stdout
	err = rdf2v1.Write(output, spdxdoc, sp)

}

func Parse(input string) (*rdf2v1.Document, *rdf2v1.Snippet, error) {
	parser := rdf2v1.NewParser(input)
	defer fmt.Printf("RDF Document parsed successfully.\n")
	defer parser.Free()
	return parser.Parse()
}
func TransferDocument(spdxdoc *rdf2v1.Document, sp *rdf2v1.Snippet) *spdx.Document2_1 {

	stdDoc := spdx.Document2_1{

		CreationInfo:  transferCreationInfo(spdxdoc),
		Packages:      transferPackages(spdxdoc, sp),
		OtherLicenses: transferOtherLicenses(spdxdoc),
		Relationships: transferRelationships(spdxdoc),
		Annotations:   transferAnnotation(spdxdoc),
		Reviews:       transferReview(spdxdoc),
	}
	return &stdDoc
}

// func CollectDocument(doc2v1 *spdx.Document2_1, sp []*rdf2v1.Snippet) *rdf2v1.Document {

// 	stdDoc := rdf2v1.Document{

// 		CreationInfo:  collectCreationInfo(doc2v1),
// 		Packages:      transferPackages(spdxdoc, sp),
// 		OtherLicenses: transferOtherLicenses(spdxdoc),
// 		Relationships: transferRelationships(spdxdoc),
// 		Annotations:   transferAnnotation(spdxdoc),
// 		Review:        collectReview(doc2v1),
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
		SPDXIdentifier:             spdxdoc.SPDXID.Val,
		DocumentName:               spdxdoc.DocumentName.Val,
		DocumentNamespace:          spdxdoc.DocumentNamespace.Val,
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

func collectCreationInfo(doc2v1 *spdx.Document2_1) *rdf2v1.CreationInfo {

	stdCi := rdf2v1.CreationInfo{

		SPDXIdentifier:     (rdf2v1.InsertId(doc2v1.CreationInfo.SPDXIdentifier)),
		LicenseListVersion: rdf2v1.Str(doc2v1.CreationInfo.LicenseListVersion),
		Creator:            rdf2v1.InsertCreator(doc2v1.CreationInfo),
		Create:             rdf2v1.Str(doc2v1.CreationInfo.Created),
		Comment:            rdf2v1.Str(doc2v1.CreationInfo.CreatorComment),
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
			AnnotationSPDXIdentifier: spdxdoc.SPDXID.Val,
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

func collectDocAnnotation(doc2v1 *spdx.Document2_1) []*rdf2v1.Annotation {
	var arrAnn []*rdf2v1.Annotation
	for _, an := range doc2v1.Annotations {
		if an.AnnotationSPDXIdentifier == doc2v1.CreationInfo.SPDXIdentifier {
			stdAnn := rdf2v1.Annotation{
				Annotator:                rdf2v1.Str(an.AnnotatorType + an.Annotator),
				AnnotationType:           rdf2v1.Str(an.AnnotationType),
				AnnotationDate:           rdf2v1.Str(an.AnnotationDate),
				AnnotationComment:        rdf2v1.Str(an.AnnotationComment),
				AnnotationSPDXIdentifier: rdf2v1.Str(an.AnnotationSPDXIdentifier),
			}
			pointer := &stdAnn
			arrAnn = append(arrAnn, pointer)
		}
	}
	return arrAnn
}

func collectFileAnnotation(doc2v1 *spdx.Document2_1) []*rdf2v1.Annotation {
	var arrAnn []*rdf2v1.Annotation
	for _, pkg := range doc2v1.Packages {
		for _, file := range pkg.Files {
			for _, an := range doc2v1.Annotations {
				if an.AnnotationSPDXIdentifier == file.FileSPDXIdentifier {
					stdAnn := rdf2v1.Annotation{
						Annotator:                rdf2v1.Str(an.AnnotatorType + an.Annotator),
						AnnotationType:           rdf2v1.Str(an.AnnotationType),
						AnnotationDate:           rdf2v1.Str(an.AnnotationDate),
						AnnotationComment:        rdf2v1.Str(an.AnnotationComment),
						AnnotationSPDXIdentifier: rdf2v1.Str(an.AnnotationSPDXIdentifier),
					}
					pointer := &stdAnn
					arrAnn = append(arrAnn, pointer)
				}
			}
		}
	}
	return arrAnn
}

func collectPackageAnnotation(doc2v1 *spdx.Document2_1) []*rdf2v1.Annotation {
	var arrAnn []*rdf2v1.Annotation
	for _, pkg := range doc2v1.Packages {
		for _, an := range doc2v1.Annotations {
			if an.AnnotationSPDXIdentifier == pkg.PackageSPDXIdentifier {
				stdAnn := rdf2v1.Annotation{
					Annotator:                rdf2v1.Str(an.AnnotatorType + an.Annotator),
					AnnotationType:           rdf2v1.Str(an.AnnotationType),
					AnnotationDate:           rdf2v1.Str(an.AnnotationDate),
					AnnotationComment:        rdf2v1.Str(an.AnnotationComment),
					AnnotationSPDXIdentifier: rdf2v1.Str(an.AnnotationSPDXIdentifier),
				}
				pointer := &stdAnn
				arrAnn = append(arrAnn, pointer)
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

func collectReview(doc2v1 *spdx.Document2_1) []*rdf2v1.Review {
	var arrRev []*rdf2v1.Review
	for _, a := range doc2v1.Reviews {
		if a != nil {
			stdRev := rdf2v1.Review{
				Reviewer:      rdf2v1.Str(a.Reviewer),
				ReviewDate:    rdf2v1.Str(a.ReviewDate),
				ReviewComment: rdf2v1.Str(a.ReviewComment),
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

// Complete it later
// func collectRelationships(doc2v1 *spdx.Document2_1) []*rdf2v1.Relationship {
// 	var arrRel []*rdf2v1.Relationship
// 	for _, a := range doc2v1.Relationships {
// 		if a != nil {
// 			stdRel := rdf2v1.Relationship{
// 				Relationship:        a.RelationshipType.Val,
// 				RelationshipComment: a.RelationshipComment.Val,
// 			}
// 			pointer := &stdRel
// 			arrRel = append(arrRel, pointer)
// 		}
// 	}

// 	return arrRel
// }
func transferFile(spdxdoc *rdf2v1.Document, sp *rdf2v1.Snippet) []*spdx.File2_1 {
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
							FileChecksumSHA1:   rdf2v1.AlgoIdentifier(b.FileChecksum, "SHA1"),
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
							Snippets: transferSnippets(sp),
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

// func collectFile(doc2v1 *spdx.Document2_1) []*rdf2v1.File {
// 	var arrFile []*rdf2v1.File
// 	for _, a := range spdxdoc.Relationship {
// 		if a != nil {
// 			if a.File != nil {
// 				for _, b := range a.File {
// 					if b != nil {
// 						stdFile := spdx.File2_1{

// 							FileName:           b.FileName.Val,
// 							FileSPDXIdentifier: "",
// 							FileType:           rdf2v1.ValueList(b.FileType),
// 							FileChecksumSHA1:   rdf2v1.AlgoIdentifier(b.FileChecksum, "SHA1"),
// 							FileChecksumSHA256: rdf2v1.AlgoIdentifier(b.FileChecksum, "SHA256"),
// 							FileChecksumMD5:    rdf2v1.AlgoIdentifier(b.FileChecksum, "MD5"),
// 							// LicenseConcluded:   "", //DISCUSS
// 							LicenseInfoInFile:  rdf2v1.ValueList(b.LicenseInfoInFile),
// 							LicenseComments:    b.FileLicenseComments.V(),
// 							FileCopyrightText:  b.FileCopyrightText.V(),
// 							ArtifactOfProjects: transferArtifactOfProject(spdxdoc),
// 							FileComment:        b.FileComment.Val,
// 							FileNotice:         b.FileNoticeText.Val,
// 							FileContributor:    rdf2v1.ValueList(b.FileContributor),
// 							// FileDependencies:   "",//DISCUSS
// 							Snippets: transferSnippets(sp),
// 						}
// 						pointer := &stdFile
// 						arrFile = append(arrFile, pointer)
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return arrFile
// }

func transferPackages(spdxdoc *rdf2v1.Document, sp *rdf2v1.Snippet) []*spdx.Package2_1 {
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
							PackageChecksumSHA1:                 rdf2v1.AlgoIdentifier(b.PackageChecksum, "SHA1"),
							PackageChecksumSHA256:               rdf2v1.AlgoIdentifier(b.PackageChecksum, "SHA256"),
							PackageChecksumMD5:                  rdf2v1.AlgoIdentifier(b.PackageChecksum, "MD5"),
							PackageHomePage:                     b.PackageHomepage.Val,
							PackageSourceInfo:                   b.PackageSourceInfo.Val,
							PackageLicenseConcluded:             "",
							// PackageLicenseInfoFromFiles: //DISCUSS
							PackageLicenseDeclared: b.PackageLicenseDeclared.Val,
							PackageLicenseComments: b.PackageLicenseComments.Val,
							PackageCopyrightText:   b.PackageCopyrightText.Val,
							PackageSummary:         b.PackageSummary.Val,
							PackageDescription:     b.PackageDescription.Val,
							PackageComment:         b.PackageComment.Val,
							// PackageExternalReferences:   "", //DISCUSS
							Files: transferFile(spdxdoc, sp),
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
func collectOtherLicenses(doc2v1 *spdx.Document2_1) []*rdf2v1.ExtractedLicensingInfo {
	var arrEl []*rdf2v1.ExtractedLicensingInfo
	for _, a := range doc2v1.OtherLicenses {
		if a != nil {
			stdEl := rdf2v1.ExtractedLicensingInfo{
				LicenseIdentifier: rdf2v1.Str(a.LicenseIdentifier),
				ExtractedText:     rdf2v1.Str(a.ExtractedText),
				// LicenseName:            a.LicenseName,// DISCUSS: more than one fields in RDF but string in standard struct
				// LicenseCrossReferences: rdf2v1.ValueList(a),//DISCUSS: Not in RDF file
				LicenseComment: rdf2v1.Str(a.LicenseComment),
			}
			pointer := &stdEl
			arrEl = append(arrEl, pointer)
		}
	}
	return arrEl
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

// func collectArtifactOfProject(spdxdoc *spdx.Document2_1) []*rdf2v1.ArtifactOfProject2_1 {
// 	var arrAop []*spdx.ArtifactOfProject2_1
// 	for _, a := range spdxdoc.Relationship {
// 		if a != nil {
// 			if a.File != nil {
// 				for _, b := range a.File {
// 					if b != nil {
// 						for _, c := range b.Project {
// 							stdAop := spdx.ArtifactOfProject2_1{
// 								Name:     c.Name.Val,
// 								HomePage: c.Homepage.Val,
// 								URI:      "",
// 							}

// 							pointer := &stdAop
// 							arrAop = append(arrAop, pointer)
// 						}
// 					}

// 				}
// 			}
// 		}
// 	}

// 	return arrAop
// }
func transferSnippets(sp *rdf2v1.Snippet) []*spdx.Snippet2_1 {
	var arrSn []*spdx.Snippet2_1
	// for _, a := range sp {
	if sp != nil {
		stdSn := spdx.Snippet2_1{
			SnippetLicenseComments:  sp.SnippetLicenseComments.Val,
			SnippetCopyrightText:    sp.SnippetCopyrightText.Val,
			SnippetLicenseConcluded: sp.SnippetLicenseConcluded.Val, //DISCUSS: Not in RDF file
			SnippetComment:          sp.SnippetComment.Val,
			// LicenseInfoInSnippet:    a.LicenseInfoInSnippet, // DISCUSS: more than one fields in RDF but string in standard struct
		}
		pointer := &stdSn
		arrSn = append(arrSn, pointer)
	}
	// }
	return arrSn
}

func collectSnippets(sp *spdx.Snippet2_1) []*rdf2v1.Snippet {
	var arrSn []*rdf2v1.Snippet
	// for _, a := range sp {
	if sp != nil {
		stdSn := rdf2v1.Snippet{
			SnippetLicenseComments:  rdf2v1.Str(sp.SnippetLicenseComments),
			SnippetCopyrightText:    rdf2v1.Str(sp.SnippetCopyrightText),
			SnippetLicenseConcluded: rdf2v1.Str(sp.SnippetLicenseConcluded), //DISCUSS: Not in RDF file
			SnippetComment:          rdf2v1.Str(sp.SnippetComment),
			// LicenseInfoInSnippet:    a.LicenseInfoInSnippet, // DISCUSS: more than one fields in RDF but string in standard struct
		}
		pointer := &stdSn
		arrSn = append(arrSn, pointer)
	}
	// }
	return arrSn
}
