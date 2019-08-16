package main

import (
	"fmt"
	"os"
	"spdx/tools-golang/v0/spdx"
	"strconv"
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
	// doc2v1 := TransferDocument(spdxdoc, sp)
	fmt.Printf("%v%T\n\n\n", spdxdoc.Relationship[1].Package[0].PackageLicenseDeclared, sp.SnippetFromFile.FileDependency)

	// invert := CollectDocument(doc2v1)
	// WRITER
	output := os.Stdout
	err = rdf2v1.Write(output, spdxdoc, sp)

	// fmt.Printf("%v\n\n\n", spdxdoc.Relationship[0])
	// fmt.Printf("%#v\n\n\n", invert.Relationship[0].File[0])
	// fmt.Printf("%#v\n\n\n", invert.Relationship[1].File[0])

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

func CollectDocument(doc2v1 *spdx.Document2_1) *rdf2v1.Document {

	stdDoc := rdf2v1.Document{
		SPDXVersion:            rdf2v1.Str(doc2v1.CreationInfo.SPDXVersion),
		DataLicense:            rdf2v1.Str(doc2v1.CreationInfo.DataLicense),
		Review:                 collectReview(doc2v1),
		DocumentName:           rdf2v1.Str(doc2v1.CreationInfo.DocumentName),
		DocumentNamespace:      rdf2v1.Str(doc2v1.CreationInfo.DocumentNamespace),
		SPDXID:                 rdf2v1.Str(doc2v1.CreationInfo.SPDXIdentifier),
		DocumentComment:        rdf2v1.Str(doc2v1.CreationInfo.DocumentComment),
		ExtractedLicensingInfo: collectExtractedLicInfo(doc2v1),
		Relationship:           collectRelationships(doc2v1),
		// License:
		CreationInfo:        collectCreationInfo(doc2v1),
		Annotation:          collectDocAnnotation(doc2v1),
		ExternalDocumentRef: collectExternalDocumentRef(doc2v1),
	}
	return &stdDoc
}

func transferCreationInfo(spdxdoc *rdf2v1.Document) *spdx.CreationInfo2_1 {

	var listExtDocRef []string
	listExtDocRef = append(listExtDocRef, spdxdoc.ExternalDocumentRef.ExternalDocumentId.Val)
	listExtDocRef = append(listExtDocRef, spdxdoc.ExternalDocumentRef.SPDXDocument.Val)
	listExtDocRef = append(listExtDocRef, rdf2v1.ExtractChecksumAlgo(spdxdoc.ExternalDocumentRef.Checksum.Algorithm.Val))
	listExtDocRef = append(listExtDocRef, spdxdoc.ExternalDocumentRef.Checksum.ChecksumValue.Val)

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

func collectExternalDocumentRef(doc2v1 *spdx.Document2_1) *rdf2v1.ExternalDocumentRef {
	stdEdr := rdf2v1.ExternalDocumentRef{

		ExternalDocumentId: rdf2v1.Str(doc2v1.CreationInfo.ExternalDocumentReferences[0]),
		SPDXDocument:       rdf2v1.Str(doc2v1.CreationInfo.ExternalDocumentReferences[1]),
		Checksum:           collectDocChecksum(doc2v1),
	}
	return &stdEdr
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

func collectVerificationCode(pkg2_1 *spdx.Package2_1) *rdf2v1.PackageVerificationCode {

	stdVc := rdf2v1.PackageVerificationCode{

		PackageVerificationCode:             rdf2v1.Str(pkg2_1.PackageVerificationCode),
		PackageVerificationCodeExcludedFile: rdf2v1.Str(pkg2_1.PackageVerificationCodeExcludedFile),
	}
	return &stdVc
}

func collectPackageChecksum(pkg2_1 *spdx.Package2_1) *rdf2v1.Checksum {

	stdPc := rdf2v1.Checksum{

		Algorithm:     rdf2v1.Str(PkgChecksumAlgo(pkg2_1)),
		ChecksumValue: rdf2v1.Str(PkgChecksumValue(pkg2_1)),
	}
	return &stdPc
}

func collectFileChecksum(File2_1 *spdx.File2_1) *rdf2v1.Checksum {

	stdFc := rdf2v1.Checksum{

		Algorithm:     rdf2v1.Str(FileChecksumAlgo(File2_1)),
		ChecksumValue: rdf2v1.Str(FileChecksumValue(File2_1)),
	}
	return &stdFc
}

func collectDocChecksum(doc2v1 *spdx.Document2_1) *rdf2v1.Checksum {

	stdFc := rdf2v1.Checksum{

		Algorithm:     rdf2v1.Str(doc2v1.CreationInfo.ExternalDocumentReferences[2]),
		ChecksumValue: rdf2v1.Str(doc2v1.CreationInfo.ExternalDocumentReferences[3]),
	}
	return &stdFc
}

func PkgChecksumAlgo(pkg2_1 *spdx.Package2_1) string {
	if pkg2_1.PackageChecksumSHA1 != "" {
		return rdf2v1.InsertChecksumAlgo("SHA1")
	}
	if pkg2_1.PackageChecksumSHA256 != "" {
		return rdf2v1.InsertChecksumAlgo("SHA256")
	}
	if pkg2_1.PackageChecksumMD5 != "" {
		return rdf2v1.InsertChecksumAlgo("MD5")
	}
	return ""
}

func PkgChecksumValue(pkg2_1 *spdx.Package2_1) string {
	if pkg2_1.PackageChecksumSHA1 != "" {
		return pkg2_1.PackageChecksumSHA1
	}
	if pkg2_1.PackageChecksumSHA256 != "" {
		return pkg2_1.PackageChecksumSHA256
	}
	if pkg2_1.PackageChecksumMD5 != "" {
		return pkg2_1.PackageChecksumMD5
	}
	return ""
}

func FileChecksumAlgo(File2_1 *spdx.File2_1) string {
	if File2_1.FileChecksumSHA1 != "" {
		return rdf2v1.InsertChecksumAlgo("SHA1")
	}
	if File2_1.FileChecksumSHA256 != "" {
		return rdf2v1.InsertChecksumAlgo("SHA256")
	}
	if File2_1.FileChecksumMD5 != "" {
		return rdf2v1.InsertChecksumAlgo("MD5")
	}
	return ""
}

func FileChecksumValue(File2_1 *spdx.File2_1) string {
	if File2_1.FileChecksumSHA1 != "" {
		return File2_1.FileChecksumSHA1
	}
	if File2_1.FileChecksumSHA256 != "" {
		return File2_1.FileChecksumSHA256
	}
	if File2_1.FileChecksumMD5 != "" {
		return File2_1.FileChecksumMD5
	}
	return ""
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
				RefA: spdxdoc.SPDXID.Val,
				// RefB:
				// Relationship:        rdf2v1.ExtractRelType(a.RelationshipType.Val),
				// RelationshipComment: a.RelationshipComment.Val,
			}
			pointer := &stdRel
			arrRel = append(arrRel, pointer)
		}
	}

	return arrRel
}

func collectRelationships(doc2v1 *spdx.Document2_1) []*rdf2v1.Relationship {
	var arrRel []*rdf2v1.Relationship
	for _, a := range doc2v1.Relationships {
		if a != nil {
			stdRel := rdf2v1.Relationship{
				RelationshipType:    rdf2v1.Str(a.Relationship),
				RelationshipComment: rdf2v1.Str(a.RelationshipComment),
				Package:             collectPackages(doc2v1),
				File:                collectFile(doc2v1),
			}
			pointer := &stdRel
			arrRel = append(arrRel, pointer)
		}
	}

	return arrRel
}

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
							FileChecksumSHA1:   rdf2v1.AlgoValue(b.FileChecksum, "SHA1"),
							FileChecksumSHA256: rdf2v1.AlgoValue(b.FileChecksum, "SHA256"),
							FileChecksumMD5:    rdf2v1.AlgoValue(b.FileChecksum, "MD5"),
							// LicenseConcluded:   "", //DISCUSS
							LicenseInfoInFile:  rdf2v1.ValueList(b.LicenseInfoInFile),
							LicenseComments:    b.FileLicenseComments.Val,
							FileCopyrightText:  b.FileCopyrightText.Val,
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

func collectFile(doc2v1 *spdx.Document2_1) []*rdf2v1.File {
	var arrFile []*rdf2v1.File
	for _, a := range doc2v1.Packages {
		if a != nil {
			if a.Files != nil {
				for _, b := range a.Files {
					if b != nil {
						stdFile := rdf2v1.File{

							FileName: rdf2v1.Str(b.FileName),
							// FileSPDXIdentifier: "",
							FileType:     rdf2v1.ValueStrList(b.FileType),
							FileChecksum: collectFileChecksum(b),
							// LicenseConcluded:   "", //DISCUSS
							LicenseInfoInFile:   rdf2v1.ValueStrList(b.LicenseInfoInFile),
							FileLicenseComments: rdf2v1.Str(b.LicenseComments),
							FileCopyrightText:   rdf2v1.Str(b.FileCopyrightText),
							Project:             collectArtifactOfProject(doc2v1),
							FileComment:         rdf2v1.Str(b.FileComment),
							FileNoticeText:      rdf2v1.Str(b.FileNotice),
							FileContributor:     rdf2v1.ValueStrList(b.FileContributor),
							// FileDependencies:   "",//DISCUSS
							Annotation: collectFileAnnotation(doc2v1),
							// ExtractedLicensingInfo
							// DisjunctiveLicenseSet
							// ConjunctiveLicenseSet
							// FileRelationship:
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
							PackageChecksumSHA1:                 rdf2v1.AlgoValue(b.PackageChecksum, "SHA1"),
							PackageChecksumSHA256:               rdf2v1.AlgoValue(b.PackageChecksum, "SHA256"),
							PackageChecksumMD5:                  rdf2v1.AlgoValue(b.PackageChecksum, "MD5"),
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

func collectPackages(doc2v1 *spdx.Document2_1) []*rdf2v1.Package {
	var arrPkg []*rdf2v1.Package
	for _, a := range doc2v1.Packages {
		if a != nil {
			stdPkg := rdf2v1.Package{
				PackageName:             rdf2v1.Str(a.PackageName),
				PackageVersionInfo:      rdf2v1.Str(a.PackageVersion),
				PackageFileName:         rdf2v1.Str(a.PackageFileName),
				PackageDownloadLocation: rdf2v1.Str(a.PackageDownloadLocation),
				PackageVerificationCode: collectVerificationCode(a), //passing specific package
				PackageComment:          rdf2v1.Str(a.PackageComment),
				PackageChecksum:         collectPackageChecksum(a),
				// // PackageLicense            :  *License
				PackageLicenseComments: rdf2v1.Str(a.PackageLicenseComments),

				// CHECK SPEC
				// DisjunctiveLicenseSet       :*DisjunctiveLicenseSet
				// ConjunctiveLicenseSet       :*ConjunctiveLicenseSet
				// PackageLicenseInfoFromFiles: rdf2v1.ValueStrList()

				PackageLicenseDeclared: rdf2v1.Str(a.PackageLicenseDeclared),
				PackageCopyrightText:   rdf2v1.Str(a.PackageCopyrightText),
				// File                 :       []*File
				// PackageRelationship   :      *Relationship
				PackageHomepage: rdf2v1.Str(a.PackageHomePage),
				PackageSupplier: rdf2v1.Str(InsertSupplier(a)),
				// PackageExternalRef       :   *ExternalRef
				PackageOriginator:  rdf2v1.Str(InsertOriginator(a)),
				PackageSourceInfo:  rdf2v1.Str(a.PackageSummary),
				FilesAnalyzed:      rdf2v1.Str(strconv.FormatBool(a.FilesAnalyzed)),
				PackageSummary:     rdf2v1.Str(a.PackageSummary),
				PackageDescription: rdf2v1.Str(a.PackageDescription),
				Annotation:         collectPackageAnnotation(doc2v1),
			}

			pointer := &stdPkg
			arrPkg = append(arrPkg, pointer)
		}
	}
	return arrPkg
}

func InsertSupplier(a *spdx.Package2_1) string {
	if a.PackageSupplierPerson != "" {
		return ("Person: " + a.PackageSupplierPerson)
	}
	if a.PackageSupplierOrganization != "" {
		return ("Organization: " + a.PackageSupplierPerson)
	}
	return ""
}

func InsertOriginator(a *spdx.Package2_1) string {
	if a.PackageOriginatorPerson != "" {
		return ("Person: " + a.PackageOriginatorPerson)
	}
	if a.PackageOriginatorOrganization != "" {
		return ("Organization: " + a.PackageOriginatorPerson)
	}
	return ""
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

func collectExtractedLicInfo(doc2v1 *spdx.Document2_1) []*rdf2v1.ExtractedLicensingInfo {
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
								HomePage: c.HomePage.Val,
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

func collectArtifactOfProject(doc2v1 *spdx.Document2_1) []*rdf2v1.Project {
	var arrp []*rdf2v1.Project
	for _, a := range doc2v1.Packages {
		if a != nil {
			if a.Files != nil {
				for _, b := range a.Files {
					if b != nil {
						for _, c := range b.ArtifactOfProjects {
							stdp := rdf2v1.Project{
								Name:     rdf2v1.Str(c.Name),
								HomePage: rdf2v1.Str(c.HomePage),
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
