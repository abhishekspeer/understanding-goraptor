package rdf2v1

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/deltamobile/goraptor"
)

func Write(output *os.File, doc *Document, sn *Snippet) error {
	f := NewFormatter(output, "rdfxml-abbrev")
	_, docerr := f.Document(doc)
	if docerr != nil {
		return nil
	}
	_, snippet := f.Snippet(sn)
	if snippet != nil {
		return nil
	}
	f.Close()
	return nil
}

// Formatter struct to write the output
type Formatter struct {
	serializer *goraptor.Serializer
	nodeIds    map[string]int
	fileIds    map[string]goraptor.Term
}

// NewFormatter initialses a new Formatter Interface
func NewFormatter(output *os.File, format string) *Formatter {

	// a new goraptor.NewSerializer
	s := goraptor.NewSerializer(format)

	s.StartStream(output, baseUri)

	// goraptor.NamespaceHandler:
	// handler function to be called when the parser encounters a namespace.
	s.SetNamespace("rdf", "http://www.w3.org/1999/02/22-rdf-syntax-ns#")
	s.SetNamespace("spdx", "http://spdx.org/rdf/terms#")
	s.SetNamespace("rdfs", "http://www.w3.org/2000/01/rdf-schema#")
	s.SetNamespace("doap:", "http://usefulinc.com/ns/doap#")
	s.SetNamespace("j.0:", "http://www.w3.org/2009/pointers#")

	return &Formatter{
		serializer: s,
		nodeIds:    make(map[string]int),
		fileIds:    make(map[string]goraptor.Term),
	}
}

// NodeId method to set new node ID for a particular prefix
func (f *Formatter) NodeId(prefix string) *goraptor.Blank {

	f.nodeIds[prefix]++

	id := goraptor.Blank(prefix + strconv.Itoa(f.nodeIds[prefix]))
	return &id
}

// Sets node type
func (f *Formatter) setNodeType(node, t goraptor.Term) error {
	return f.add(node, prefix("ns:type"), t)
}

// Add 'keys' to 'values' for subject 'to'
func (f *Formatter) add(to, key, value goraptor.Term) error {
	return f.serializer.Add(&goraptor.Statement{
		Subject:   to,
		Predicate: key,
		Object:    value,
	})
}
func (f *Formatter) addTerm(to goraptor.Term, key string, value goraptor.Term) error {
	return f.add(to, prefix(key), value)
}

func (f *Formatter) addLiteral(to goraptor.Term, key, value string) error {
	if value == "" {
		return nil
	}
	return f.add(to, prefix(key), &goraptor.Literal{Value: value})
}
func (f *Formatter) addPairs(to goraptor.Term, pairs ...pair) error {
	for _, p := range pairs {
		if err := f.addLiteral(to, p.key, p.val); err != nil {
			return err
		}
	}
	return nil
}

func (f *Formatter) Document(doc *Document) (docId goraptor.Term, err error) {

	// _docId := goraptor.Blank("doc")
	if doc == nil {
		return nil, errors.New("Nil document.")
	}

	// docId = &_docId
	docId = blank("doc")

	if err = f.setNodeType(docId, typeDocument); err != nil {
		return
	}

	if err = f.addLiteral(docId, "specVersion", doc.SPDXVersion.Val); err != nil {
		return
	}

	if doc.DataLicense.Val != "" {
		if err = f.addTerm(docId, "dataLicense", uri(licenseUri+doc.DataLicense.Val)); err != nil {
			return
		}
	}
	if doc.DocumentName.Val != "" {
		if err = f.addTerm(docId, "name", uri(licenseUri+doc.DataLicense.Val)); err != nil {
			return
		}
	}
	if err = f.addLiteral(docId, "rdfs:comment", doc.DocumentComment.Val); err != nil {
		return
	}

	if id, err := f.CreationInfo(doc.CreationInfo); err == nil {
		if err = f.addTerm(docId, "creationInfo", id); err != nil {
			return docId, err
		}
	} else {
		return docId, err
	}

	if id, err := f.License(doc.License); err == nil {
		if err = f.addTerm(docId, "dataLicense", id); err != nil {
			return docId, err
		}
	} else {
		return docId, err
	}

	if id, err := f.ExternalDocumentRef(doc.ExternalDocumentRef); err == nil {
		if err = f.addTerm(docId, "externalDocumentRef", id); err != nil {
			return docId, err
		}
	} else {
		return docId, err
	}
	if err = f.Relationships(docId, "relationship", doc.Relationship); err != nil {
		return
	}

	if err = f.Reviews(docId, "reviewed", doc.Review); err != nil {
		return
	}
	if err = f.Annotations(docId, "annotation", doc.Annotation); err != nil {
		return
	}
	if err = f.ExtractedLicInfos(docId, "hasExtractedLicensingInfo", doc.ExtractedLicensingInfo); err != nil {
		return
	}
	return docId, nil
}

func (f *Formatter) Snippet(snip *Snippet) (snipId goraptor.Term, err error) {

	if snip == nil {
		return nil, errors.New("Nil Snippet.")
	}

	// docId = &_docId
	snipId = blank("snip")

	if err = f.setNodeType(snipId, typeDocument); err != nil {
		return
	}
	if err = f.addLiteral(snipId, "name", snip.SnippetName.Val); err != nil {
		return
	}
	if err = f.addLiteral(snipId, "copyrightText", snip.SnippetCopyrightText.Val); err != nil {
		return
	}
	if err = f.addLiteral(snipId, "licenseComments", snip.SnippetLicenseComments.Val); err != nil {
		return
	}
	if err = f.addLiteral(snipId, "rdfs:comment", snip.SnippetComment.Val); err != nil {
		return
	}
	for _, lis := range snip.LicenseInfoInSnippet {
		if err = f.addLiteral(snipId, "creator", lis.Val); err != nil {
			return
		}
	}
	if snip.SnippetLicenseConcluded.Val != "" {
		if err = f.addTerm(snipId, "licenseConcluded", prefix(snip.SnippetLicenseConcluded.Val)); err != nil {
			return
		}
	}
	if snip.SnippetFromFile != nil {
		sfId, err := f.File(snip.SnippetFromFile)
		if err != nil {
			return snipId, err
		}
		if err = f.addTerm(snipId, "snippetFromFile", sfId); err != nil {
			return snipId, err
		}
	}
	return snipId, nil
}
func (f *Formatter) ExternalDocumentRef(edr *ExternalDocumentRef) (id goraptor.Term, err error) {
	id = f.NodeId("edr")

	if err = f.setNodeType(id, typeExternalDocumentRef); err != nil {
		return
	}

	err = f.addPairs(id,
		pair{"externalDocumentId", edr.ExternalDocumentId.Val},
		pair{"spdxDocument", edr.SPDXDocument.Val},
	)

	if err != nil {
		return
	}

	if edr.Checksum != nil {
		cksumId, err := f.Checksum(edr.Checksum)
		if err != nil {
			return id, err
		}
		if err = f.addTerm(id, "checksum", cksumId); err != nil {
			return id, err
		}
	}

	return id, nil
}

func (f *Formatter) CreationInfo(ci *CreationInfo) (id goraptor.Term, err error) {
	id = f.NodeId("cri")

	if err = f.setNodeType(id, typeCreationInfo); err != nil {
		return
	}

	err = f.addPairs(id,
		pair{"created", ci.Create.Val},
		pair{"rdfs:comment", ci.Comment.Val},
		pair{"licenseListVersion", ci.LicenseListVersion.Val},
	)

	if err != nil {
		return
	}

	for _, creator := range ci.Creator {
		if err = f.addLiteral(id, "creator", creator.Val); err != nil {
			return
		}
	}

	return id, nil
}

// Write Review
func (f *Formatter) Review(r *Review) (id goraptor.Term, err error) {
	id = f.NodeId("rev")

	if err = f.setNodeType(id, typeReview); err != nil {
		return
	}

	err = f.addPairs(id,
		pair{"reviewer", r.Reviewer.Val},
		pair{"reviewDate", r.ReviewDate.Val},
		pair{"rdfs:comment", r.ReviewComment.Val},
	)

	return id, err
}
func (f *Formatter) Project(pro *Project) (id goraptor.Term, err error) {
	id = f.NodeId("pro")

	if err = f.setNodeType(id, typeProject); err != nil {
		return
	}

	err = f.addPairs(id,
		pair{"homepage", pro.Homepage.Val},
		pair{"name", pro.Name.Val},
	)

	return id, err
}
func (f *Formatter) PackageVerificationCode(pvc *PackageVerificationCode) (id goraptor.Term, err error) {
	id = f.NodeId("pvc")

	if err = f.setNodeType(id, typePackageVerificationCode); err != nil {
		return
	}

	err = f.addPairs(id,
		pair{"packageVerificationCodeValue", pvc.PackageVerificationCode.Val},
		pair{"packageVerificationCodeExcludedFile", pvc.PackageVerificationCodeExcludedFile.Val},
	)

	return id, err
}

// Write Annotation
func (f *Formatter) Annotation(an *Annotation) (id goraptor.Term, err error) {
	id = f.NodeId("an")

	if err = f.setNodeType(id, typeAnnotation); err != nil {
		return
	}
	// type 1: add pairs
	err = f.addPairs(id,
		pair{"annotationDate", an.AnnotationDate.Val},
		pair{"rdfs:comment", an.AnnotationComment.Val},
		pair{"annotator", an.Annotator.Val},
	)
	if err != nil {
		return
	}
	// type 2: add pairs with prefix
	if an.AnnotationType.Val != "" {
		if err = f.addTerm(id, "annotationType", prefix(an.AnnotationType.Val)); err != nil {
			return
		}
	}
	return id, err
}

func (f *Formatter) Reviews(parent goraptor.Term, element string, rs []*Review) error {

	if len(rs) == 0 {
		return nil
	}

	for _, r := range rs {
		revId, err := f.Review(r)
		if err != nil {
			return err
		}
		if revId == nil {
			continue
		}
		if err = f.addTerm(parent, element, revId); err != nil {
			return err
		}
	}
	return nil
}

func (f *Formatter) Annotations(parent goraptor.Term, element string, ans []*Annotation) error {

	if len(ans) == 0 {
		return nil
	}

	for _, an := range ans {
		annId, err := f.Annotation(an)
		if err != nil {
			return err
		}
		if annId == nil {
			continue
		}
		if err = f.addTerm(parent, element, annId); err != nil {
			return err
		}
	}
	return nil
}
func (f *Formatter) Projects(parent goraptor.Term, element string, pros []*Project) error {

	if len(pros) == 0 {
		return nil
	}

	for _, pro := range pros {
		proId, err := f.Project(pro)
		if err != nil {
			return err
		}
		if proId == nil {
			continue
		}
		if err = f.addTerm(parent, element, proId); err != nil {
			return err
		}
	}
	return nil
}
func (f *Formatter) Checksum(cksum *Checksum) (id goraptor.Term, err error) {
	id = f.NodeId("cksum")

	if err = f.setNodeType(id, typeChecksum); err != nil {
		return
	}

	err = f.addLiteral(id, "checksumValue", cksum.ChecksumValue.Val)
	if err != nil {
		return
	}

	algo := strings.ToLower(cksum.Algorithm.Val)
	if algo == "sha1" {
		err = f.addTerm(id, "algorithm", prefix("checksumAlgorithm_sha1"))

	} else if algo == "md5" {
		err = f.addTerm(id, "algorithm", prefix("checksumAlgorithm_md5"))
	} else if algo == "sha256" {
		err = f.addTerm(id, "algorithm", prefix("checksumAlgorithm_sha256"))
	} else {
		err = f.addLiteral(id, "algorithm", algo)
	}

	return id, err
}

func (f *Formatter) ExtractedLicInfo(lic *ExtractedLicensingInfo) (id goraptor.Term, err error) {
	id = f.NodeId("lic")

	if err = f.setNodeType(id, typeExtractedLicensingInfo); err != nil {
		return
	}

	err = f.addPairs(id,
		pair{"licenseId", lic.LicenseIdentifier.Val},
		pair{"extractedText", lic.ExtractedText.Val},
		pair{"rdfs:comment", lic.LicenseComment.Val},
	)

	if err != nil {
		return
	}

	for _, name := range lic.LicenseName {
		if err = f.addLiteral(id, "name", name.Val); err != nil {
			return
		}
	}

	for _, seealso := range lic.LicenseSeeAlso {
		if err = f.addLiteral(id, "rdfs:seeAlso", seealso.Val); err != nil {
			return
		}
	}

	return id, err
}

func (f *Formatter) ExtractedLicInfos(parent goraptor.Term, element string, lics []*ExtractedLicensingInfo) error {

	if len(lics) == 0 {
		return nil
	}

	for _, lic := range lics {
		licId, err := f.ExtractedLicInfo(lic)
		if err != nil {
			return err
		}
		if licId == nil {
			continue
		}
		if err = f.addTerm(parent, element, licId); err != nil {
			return err
		}
	}
	return nil
}

// Write a slice of files.
func (f *Formatter) Files(parent goraptor.Term, element string, files []*File) error {
	if len(files) == 0 {
		return nil
	}
	for _, file := range files {
		fId, err := f.File(file)
		if err != nil {
			return err
		}
		if fId == nil {
			continue
		}
		if err = f.addTerm(parent, element, fId); err != nil {
			return err
		}
	}
	return nil
}

// Write a file.
func (f *Formatter) File(file *File) (id goraptor.Term, err error) {
	id, ok := f.fileIds[file.FileName.Val]
	if ok {
		return
	}

	id = f.NodeId("file")
	f.fileIds[file.FileName.Val] = id

	if err = f.setNodeType(id, typeFile); err != nil {
		return
	}

	err = f.addPairs(id,
		pair{"fileName", file.FileName.Val},
		pair{"licenseComments", file.FileLicenseComments.Val},
		pair{"copyrightText", file.FileCopyrightText.Val},
		pair{"rdfs:comment", file.FileComment.Val},
		pair{"noticeText", file.FileNoticeText.Val},
	)

	if err != nil {
		return
	}
	if file.FileChecksum != nil {
		cksumId, err := f.Checksum(file.FileChecksum)
		if err != nil {
			return id, err
		}
		if err = f.addTerm(id, "checksum", cksumId); err != nil {
			return id, err
		}
	}

	if file.ExtractedLicensingInfo != nil {
		exlicId, err := f.ExtractedLicInfo(file.ExtractedLicensingInfo)
		if err != nil {
			return id, err
		}
		if err = f.addTerm(id, "hasExtractedLicensingInfo", exlicId); err != nil {
			return id, err
		}
	}
	if file.DisjunctiveLicenseSet != nil {
		dlsId, err := f.DisjunctiveLicenseSet(file.DisjunctiveLicenseSet)
		if err != nil {
			return id, err
		}
		if err = f.addTerm(id, "member", dlsId); err != nil {
			return id, err
		}
	}
	if file.ConjunctiveLicenseSet != nil {
		clsId, err := f.ConjunctiveLicenseSet(file.ConjunctiveLicenseSet)
		if err != nil {
			return id, err
		}
		if err = f.addTerm(id, "member", clsId); err != nil {
			return id, err
		}
	}

	for _, fc := range file.FileContributor {
		if err = f.addLiteral(id, "fileContributor", fc.Val); err != nil {
			return
		}
	}
	for _, lif := range file.LicenseInfoInFile {
		if err = f.addTerm(id, "licenseInfoInFile", prefix(lif.Val)); err != nil {
			return
		}
	}
	for _, ft := range file.FileType {
		if err = f.addTerm(id, "fileType", prefix(ft.Val)); err != nil {
			return
		}
	}

	if file.FileChecksum != nil {
		cksumId, err := f.Checksum(file.FileChecksum)
		if err != nil {
			return id, err
		}
		if err = f.addTerm(id, "checksum", cksumId); err != nil {
			return id, err
		}
	}
	//checkaftersnippets
	if file.SnippetLicense != nil {
		filelicId, err := f.License(file.SnippetLicense)
		if err != nil {
			filelicId, err = f.DisjunctiveLicenseSet(file.DisjunctiveLicenseSet)
			if err != nil {
				filelicId, err = f.ExtractedLicInfo(file.ExtractedLicensingInfo)
				if err != nil {
					return id, err
				}
			}
		}
		if err = f.addTerm(id, "licenseConcluded", filelicId); err != nil {
			return id, err

		}
	}
	if file.FileDependency != nil {
		fdId, err := f.File(file.FileDependency)
		if err != nil {
			return id, err
		}
		if err = f.addTerm(id, "fileDependency", fdId); err != nil {
			return id, err
		}
	}

	if file.FileRelationship != nil {
		frId, err := f.Relationship(file.FileRelationship)
		if err != nil {
			return id, err
		}
		if err = f.addTerm(id, "relationship", frId); err != nil {
			return id, err
		}
	}

	if err = f.Annotations(id, "annotation", file.Annotation); err != nil {
		return
	}
	if err = f.Projects(id, "artifactOf", file.Project); err != nil {
		return
	}
	return id, err

}

func (f *Formatter) Relationships(parent goraptor.Term, element string, rels []*Relationship) error {
	if len(rels) == 0 {
		return nil
	}
	for _, rel := range rels {
		relId, err := f.Relationship(rel)
		if err != nil {
			return err
		}
		if relId == nil {
			continue
		}
		if err = f.addTerm(parent, element, relId); err != nil {
			return err
		}
	}
	return nil
}
func (f *Formatter) Relationship(rel *Relationship) (id goraptor.Term, err error) {
	id = f.NodeId("rel")

	if err = f.setNodeType(id, typeRelationship); err != nil {
		return
	}

	err = f.addPairs(id,
		pair{"rdfs:comment", rel.RelationshipComment.Val},
	)
	if err != nil {
		return
	}

	if rel.RelationshipType.Val != "" {
		if err = f.addTerm(id, "relationshipType", prefix(rel.RelationshipType.Val)); err != nil {
			return
		}
	}
	if rel.relatedSpdxElement.Val != "" {
		if err = f.addTerm(id, "relatedSpdxElement", prefix(rel.relatedSpdxElement.Val)); err != nil {
			return
		}
	}
	if rel.SpdxElement != nil {
		seId, err := f.SpdxElement(rel.SpdxElement)
		if err != nil {
			return id, err
		}
		if err = f.addTerm(id, "relatedSpdxElement", seId); err != nil {
			return id, err
		}
	}
	// if err = f.Files(id, "referencesFile", rel.File); err != nil {
	// 	return
	// }
	/*
		if rel.relatedSpdxElement.Val != "" {
			rseid, err := f.Package(rel.Package)
			if err != nil {
				rseid, err := f.VerificationCode(pkg.VerificationCode)
				if err != nil {
					pkgid, err := f.VerificationCode(pkg.VerificationCode)
					if err != nil {
						return id, err
					}
				}

			}
			if err = f.addTerm(id, "relatedSpdxElement", prefix(rel.relatedSpdxElement.Val)); err != nil {
				return id ,err
			}}*/

	return id, err
}
func (f *Formatter) SpdxElement(se *SpdxElement) (id goraptor.Term, err error) {
	id = f.NodeId("se")

	if err = f.setNodeType(id, typeSpdxElement); err != nil {
		return
	}
	if se.SpdxElement.Val != "" {
		if err = f.addTerm(id, "SpdxElement", prefix(se.SpdxElement.Val)); err != nil {
			return
		}
	}
	return id, err
}
func (f *Formatter) License(lic *License) (id goraptor.Term, err error) {
	id = f.NodeId("lic")

	if err = f.setNodeType(id, typeLicense); err != nil {
		return
	}

	err = f.addPairs(id,
		pair{"rdfs:comment", lic.LicenseComment.Val},
		pair{"name", lic.LicenseName.Val},
		pair{"licenseText", lic.LicenseText.Val},
		pair{"standardLicenseHeader", lic.StandardLicenseHeader.V()},
		pair{"standardLicenseTemplate", lic.StandardLicenseTemplate.V()},
		pair{"standardLicenseHeaderTemplate", lic.StandardLicenseHeaderTemplate.Val},
		pair{"isFsfLibre", lic.LicenseIsFsLibre.Val},
		pair{"licenseId", lic.LicenseId.Val},
		pair{"licenseOsiApproved", lic.LicenseisOsiApproved.Val},
	)
	for _, sa := range lic.LicenseSeeAlso {
		if err = f.addLiteral(id, "rdfs:seeAlso", sa.Val); err != nil {
			return
		}
	}

	return id, err
}
func (f *Formatter) ConjunctiveLicenseSet(cls *ConjunctiveLicenseSet) (id goraptor.Term, err error) {
	id = f.NodeId("cls")

	if err = f.setNodeType(id, typeConjunctiveLicenseSet); err != nil {
		return
	}

	if id, err := f.License(cls.License); err == nil {
		if err = f.addTerm(id, "member", id); err != nil {
			return id, err
		}
	} else if id, err := f.ExtractedLicInfo(cls.ExtractedLicensingInfo); err == nil {
		if err = f.addTerm(id, "member", id); err != nil {
			return id, err
		}
	} else {
		return id, err
	}

	return id, err
}

func (f *Formatter) DisjunctiveLicenseSet(dls *DisjunctiveLicenseSet) (id goraptor.Term, err error) {
	id = f.NodeId("dls")

	if err = f.setNodeType(id, typeDisjunctiveLicenseSet); err != nil {
		return
	}

	for _, mem := range dls.Member {
		if err = f.addLiteral(id, "member", mem.Val); err != nil {
			fmt.Println("\n\n\nERR\n\n\n") //check
			return
		}
	}

	return id, err
}

func (f *Formatter) Packages(parent goraptor.Term, element string, pkgs []*Package) error {
	if len(pkgs) == 0 {
		return nil
	}
	for _, pkg := range pkgs {
		pkgid, err := f.Package(pkg)
		if err != nil {
			return err
		}
		if err = f.addTerm(parent, element, pkgid); err != nil {
			return err
		}
	}
	return nil
}

func (f *Formatter) Package(pkg *Package) (id goraptor.Term, err error) {
	id = f.NodeId("pkg")

	if err = f.setNodeType(id, typePackage); err != nil {
		return
	}

	err = f.addPairs(id,
		pair{"name", pkg.PackageName.Val},
		pair{"versionInfo", pkg.PackageVersionInfo.Val},
		pair{"packageFileName", pkg.PackageFileName.Val},
		pair{"downloadLocation", pkg.PackageDownloadLocation.Val},
		pair{"rdfs:comment", pkg.PackageComment.Val},
		pair{"licenseComments", pkg.PackageLicenseComments.Val},
		pair{"copyrightText", pkg.PackageCopyrightText.Val},
		pair{"doap:homepage", pkg.PackageHomepage.Val},
		pair{"supplier", pkg.PackageSupplier.Val},
		pair{"originator", pkg.PackageOriginator.V()},
		pair{"sourceInfo", pkg.PackageSourceInfo.Val},
		pair{"filesAnalyzed", pkg.FilesAnalyzed.Val},
		pair{"summary", pkg.PackageSummary.Val},
		pair{"description", pkg.PackageDescription.Val},
	)
	if err != nil {
		return
	}
	if pkg.PackageLicense != nil {
		pkglicId, err := f.License(pkg.PackageLicense)
		if err != nil {
			pkglicId, err = f.DisjunctiveLicenseSet(pkg.DisjunctiveLicenseSet)
			if err != nil {
				pkglicId, err = f.ConjunctiveLicenseSet(pkg.ConjunctiveLicenseSet)
				if err != nil {
					return id, err
				}
			}
		}
		if err = f.addTerm(id, "licenseConcluded", pkglicId); err != nil {
			return id, err

		}
	}
	if pkg.PackageVerificationCode != nil {
		pkgid, err := f.PackageVerificationCode(pkg.PackageVerificationCode)
		if err != nil {
			return id, err
		}
		if err = f.addTerm(id, "packageVerificationCode", pkgid); err != nil {
			return id, err
		}
	}

	if pkg.PackageChecksum != nil {
		cksumId, err := f.Checksum(pkg.PackageChecksum)
		if err != nil {
			return id, err
		}
		if err = f.addTerm(id, "checksum", cksumId); err != nil {
			return id, err
		}
	}

	if pkg.PackageLicenseDeclared.Val != "" {
		if err = f.addTerm(id, "licenseDeclared", prefix(pkg.PackageLicenseDeclared.Val)); err != nil {
			return
		}
	}
	if err = f.Annotations(id, "annotation", pkg.Annotation); err != nil {
		return
	}
	if err = f.Files(id, "hasFile", pkg.File); err != nil {
		return
	}

	if pkg.PackageExternalRef != nil {
		pkgErId, err := f.ExternalRef(pkg.PackageExternalRef)
		if err != nil {
			return id, err
		}
		if err = f.addTerm(id, "externalRef", pkgErId); err != nil {
			return id, err
		}
	}

	if pkg.PackageRelationship != nil {
		pkgRelId, err := f.Relationship(pkg.PackageRelationship)
		if err != nil {
			return id, err
		}
		if err = f.addTerm(id, "relationship", pkgRelId); err != nil {
			return id, err
		}
	}
	for _, lif := range pkg.PackageLicenseInfoFromFiles {
		if err = f.addTerm(id, "licenseInfoFromFiles", prefix(lif.Val)); err != nil {
			return
		}
	}
	// if pkg.PackageLicenseDeclared != "" {
	// 	licId, err := f.License(pkg.PackageLicenseDeclared)
	// 	if err != nil {
	// 		return id, err
	// 	}
	// 	if err = f.addTerm(id, "licenseDeclared", licId); err != nil {
	// 		return id, err
	// 	}
	// }

	return id, err
}

// Write Review
func (f *Formatter) ExternalRef(er *ExternalRef) (id goraptor.Term, err error) {
	id = f.NodeId("er")

	if err = f.setNodeType(id, typeExternalRef); err != nil {
		return
	}

	err = f.addPairs(id,
		pair{"referenceLocator", er.ReferenceLocator.Val},
		pair{"rdfs:comment", er.ReferenceComment.Val},
	)
	if er.ReferenceCategory.Val != "" {
		if err = f.addTerm(id, "annotationType", prefix(er.ReferenceCategory.Val)); err != nil {
			return
		}
	}
	return id, err
}

func (f *Formatter) ReferenceType(rt *ReferenceType) (id goraptor.Term, err error) {
	id = f.NodeId("rt")

	if err = f.setNodeType(id, typeReferenceType); err != nil {
		return
	}

	if rt.ReferenceType.Val != "" {
		if err = f.addTerm(id, "referenceType", prefix(rt.ReferenceType.Val)); err != nil {
			return
		}
	}
	return id, err
}

// Close to free the serializer
func (f *Formatter) Close() {
	f.serializer.EndStream()
	f.serializer.Free()
}
