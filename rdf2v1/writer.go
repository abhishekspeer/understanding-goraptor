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
	// _, docerr := f.Document(doc)
	// if docerr != nil {
	// 	return nil
	// }
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

	s.StartStream(output, BaseUri)

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

// Sets node Type
func (f *Formatter) setNodeType(node, t goraptor.Term) error {
	return f.add(node, Prefix("ns:Type"), t)
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
	return f.add(to, Prefix(key), value)
}

func (f *Formatter) addLiteral(to goraptor.Term, key, value string) error {
	if value == "" {
		return nil
	}
	return f.add(to, Prefix(key), &goraptor.Literal{Value: value})
}
func (f *Formatter) addPairs(to goraptor.Term, Pairs ...Pair) error {
	for _, p := range Pairs {
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
	docId = Blank("doc")

	if err = f.setNodeType(docId, TypeDocument); err != nil {
		return
	}

	if err = f.addLiteral(docId, "specVersion", doc.SPDXVersion.Val); err != nil {
		return
	}

	if doc.DataLicense.Val != "" {
		if err = f.addTerm(docId, "dataLicense", Uri(LicenseUri+doc.DataLicense.Val)); err != nil {
			return
		}
	}
	if doc.DocumentName.Val != "" {
		if err = f.addTerm(docId, "name", Uri(LicenseUri+doc.DataLicense.Val)); err != nil {
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
	snipId = Blank("snip")

	if err = f.setNodeType(snipId, TypeSnippet); err != nil {
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

	if snip.SnippetLicenseConcluded.Val != "" {
		if err = f.addTerm(snipId, "licenseConcluded", Prefix(snip.SnippetLicenseConcluded.Val)); err != nil {
			return
		}
	}
	for _, li := range snip.LicenseInfoInSnippet {
		if err = f.addLiteral(snipId, "licenseInfoInSnippet", li.Val); err != nil {
			return
		}
	}
	if err = f.SnippetStartEndPointers(snipId, "range", snip.SnippetStartEndPointer); err != nil {

		return

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

	if err = f.setNodeType(id, TypeExternalDocumentRef); err != nil {
		return
	}

	err = f.addPairs(id,
		Pair{"externalDocumentId", edr.ExternalDocumentId.Val},
		Pair{"spdxDocument", edr.SPDXDocument.Val},
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

	if err = f.setNodeType(id, TypeCreationInfo); err != nil {
		return
	}

	err = f.addPairs(id,
		Pair{"created", ci.Create.Val},
		Pair{"rdfs:comment", ci.Comment.Val},
		Pair{"licenseListVersion", ci.LicenseListVersion.Val},
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

	if err = f.setNodeType(id, TypeReview); err != nil {
		return
	}

	err = f.addPairs(id,
		Pair{"reviewer", r.Reviewer.Val},
		Pair{"reviewDate", r.ReviewDate.Val},
		Pair{"rdfs:comment", r.ReviewComment.Val},
	)

	return id, err
}
func (f *Formatter) Project(pro *Project) (id goraptor.Term, err error) {
	id = f.NodeId("pro")

	if err = f.setNodeType(id, TypeProject); err != nil {
		return
	}

	err = f.addPairs(id,
		Pair{"homepage", pro.HomePage.Val},
		Pair{"name", pro.Name.Val},
	)

	return id, err
}
func (f *Formatter) PackageVerificationCode(pvc *PackageVerificationCode) (id goraptor.Term, err error) {
	id = f.NodeId("pvc")

	if err = f.setNodeType(id, TypePackageVerificationCode); err != nil {
		return
	}

	err = f.addPairs(id,
		Pair{"packageVerificationCodeValue", pvc.PackageVerificationCode.Val},
		Pair{"packageVerificationCodeExcludedFile", pvc.PackageVerificationCodeExcludedFile.Val},
	)

	return id, err
}

// Write Annotation
func (f *Formatter) Annotation(an *Annotation) (id goraptor.Term, err error) {
	id = f.NodeId("an")

	if err = f.setNodeType(id, TypeAnnotation); err != nil {
		return
	}
	// Type 1: add Pairs
	err = f.addPairs(id,
		Pair{"annotationDate", an.AnnotationDate.Val},
		Pair{"rdfs:comment", an.AnnotationComment.Val},
		Pair{"annotator", an.Annotator.Val},
	)
	if err != nil {
		return
	}
	// Type 2: add Pairs with Prefix
	if an.AnnotationType.Val != "" {
		if err = f.addTerm(id, "annotationType", Prefix(an.AnnotationType.Val)); err != nil {
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

	if err = f.setNodeType(id, TypeChecksum); err != nil {
		return
	}

	err = f.addLiteral(id, "checksumValue", cksum.ChecksumValue.Val)
	if err != nil {
		return
	}

	algo := strings.ToLower(cksum.Algorithm.Val)
	if algo == "sha1" {
		err = f.addTerm(id, "algorithm", Prefix("checksumAlgorithm_sha1"))

	} else if algo == "md5" {
		err = f.addTerm(id, "algorithm", Prefix("checksumAlgorithm_md5"))
	} else if algo == "sha256" {
		err = f.addTerm(id, "algorithm", Prefix("checksumAlgorithm_sha256"))
	} else {
		err = f.addLiteral(id, "algorithm", algo)
	}

	return id, err
}

func (f *Formatter) ExtractedLicInfo(lic *ExtractedLicensingInfo) (id goraptor.Term, err error) {
	id = f.NodeId("lic")

	if err = f.setNodeType(id, TypeExtractedLicensingInfo); err != nil {
		return
	}

	err = f.addPairs(id,
		Pair{"licenseId", lic.LicenseIdentifier.Val},
		Pair{"extractedText", lic.ExtractedText.Val},
		Pair{"rdfs:comment", lic.LicenseComment.Val},
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

	if err = f.setNodeType(id, TypeFile); err != nil {
		return
	}

	err = f.addPairs(id,
		Pair{"fileName", file.FileName.Val},
		Pair{"licenseComments", file.FileLicenseComments.Val},
		Pair{"copyrightText", file.FileCopyrightText.Val},
		Pair{"rdfs:comment", file.FileComment.Val},
		Pair{"noticeText", file.FileNoticeText.Val},
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
		if err = f.addTerm(id, "licenseInfoInFile", Prefix(lif.Val)); err != nil {
			return
		}
	}
	for _, ft := range file.FileType {
		if err = f.addTerm(id, "fileType", Prefix(ft.Val)); err != nil {
			return
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

	if err = f.setNodeType(id, TypeRelationship); err != nil {
		return
	}

	err = f.addPairs(id,
		Pair{"rdfs:comment", rel.RelationshipComment.Val},
	)
	if err != nil {
		return
	}

	if rel.RelationshipType.Val != "" {
		if err = f.addTerm(id, "relationshipType", Prefix(rel.RelationshipType.Val)); err != nil {
			return
		}
	}
	if rel.RelatedSpdxElement.Val != "" {
		if err = f.addTerm(id, "relatedSpdxElement", Prefix(rel.RelatedSpdxElement.Val)); err != nil {
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

	if err = f.Files(id, "relatedSpdxElement", rel.File); err != nil {
		if err = f.Packages(id, "relatedSpdxElement", rel.Package); err != nil {
			return
		}
	}

	return id, err
}
func (f *Formatter) SpdxElement(se *SpdxElement) (id goraptor.Term, err error) {
	id = f.NodeId("se")

	if err = f.setNodeType(id, TypeSpdxElement); err != nil {
		return
	}
	if se.SpdxElement.Val != "" {
		if err = f.addTerm(id, "SpdxElement", Prefix(se.SpdxElement.Val)); err != nil {
			return
		}
	}
	return id, err
}
func (f *Formatter) License(lic *License) (id goraptor.Term, err error) {
	id = f.NodeId("lic")

	if err = f.setNodeType(id, TypeLicense); err != nil {
		return
	}

	err = f.addPairs(id,
		Pair{"rdfs:comment", lic.LicenseComment.Val},
		Pair{"name", lic.LicenseName.Val},
		Pair{"licenseText", lic.LicenseText.Val},
		Pair{"standardLicenseHeader", lic.StandardLicenseHeader.V()},
		Pair{"standardLicenseTemplate", lic.StandardLicenseTemplate.V()},
		Pair{"standardLicenseHeaderTemplate", lic.StandardLicenseHeaderTemplate.Val},
		Pair{"isFsfLibre", lic.LicenseIsFsLibre.Val},
		Pair{"licenseId", lic.LicenseId.Val},
		Pair{"licenseOsiApproved", lic.LicenseisOsiApproved.Val},
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

	if err = f.setNodeType(id, TypeConjunctiveLicenseSet); err != nil {
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

	if err = f.setNodeType(id, TypeDisjunctiveLicenseSet); err != nil {
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

	if err = f.setNodeType(id, TypePackage); err != nil {
		return
	}

	err = f.addPairs(id,
		Pair{"name", pkg.PackageName.Val},
		Pair{"versionInfo", pkg.PackageVersionInfo.Val},
		Pair{"packageFileName", pkg.PackageFileName.Val},
		Pair{"downloadLocation", pkg.PackageDownloadLocation.Val},
		Pair{"rdfs:comment", pkg.PackageComment.Val},
		Pair{"licenseComments", pkg.PackageLicenseComments.Val},
		Pair{"copyrightText", pkg.PackageCopyrightText.Val},
		Pair{"doap:homepage", pkg.PackageHomepage.Val},
		Pair{"supplier", pkg.PackageSupplier.Val},
		Pair{"originator", pkg.PackageOriginator.V()},
		Pair{"sourceInfo", pkg.PackageSourceInfo.Val},
		Pair{"filesAnalyzed", pkg.FilesAnalyzed.Val},
		Pair{"summary", pkg.PackageSummary.Val},
		Pair{"description", pkg.PackageDescription.Val},
	)
	if err != nil {
		return
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
		if err = f.addTerm(id, "licenseInfoFromFiles", Prefix(lif.Val)); err != nil {
			return
		}
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

	if pkg.PackageLicenseDeclared.Val != "" {
		if err = f.addTerm(id, "licenseDeclared", Prefix(pkg.PackageLicenseDeclared.Val)); err != nil {
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
			if err = f.addTerm(id, "licenseDeclared", pkglicId); err != nil {
				return id, err
			}
		}

	}
	return id, err
}

func (f *Formatter) ExternalRef(er *ExternalRef) (id goraptor.Term, err error) {
	id = f.NodeId("er")

	if err = f.setNodeType(id, TypeExternalRef); err != nil {
		return
	}

	err = f.addPairs(id,
		Pair{"referenceLocator", er.ReferenceLocator.Val},
		Pair{"rdfs:comment", er.ReferenceComment.Val},
	)
	if id, err := f.ReferenceType(er.ReferenceType); err == nil {
		if err = f.addTerm(id, "referenceType", id); err != nil {
			return id, err
		}
	}
	if er.ReferenceCategory.Val != "" {
		if err = f.addTerm(id, "referenceCategory", Prefix(er.ReferenceCategory.Val)); err != nil {
			return
		}
	}
	return id, err
}

func (f *Formatter) ReferenceType(rt *ReferenceType) (id goraptor.Term, err error) {
	id = f.NodeId("rt")

	if err = f.setNodeType(id, TypeReferenceType); err != nil {
		return
	}

	if rt.ReferenceType.Val != "" {
		if err = f.addTerm(id, "referenceType", Prefix(rt.ReferenceType.Val)); err != nil {
			return
		}
	}
	return id, err
}

func (f *Formatter) SnippetStartEndPointer(se *SnippetStartEndPointer) (id goraptor.Term, err error) {
	id = f.NodeId("ssep")

	if err = f.setNodeType(id, TypeSnippetStartEndPointer); err != nil {
		return
	}

	if err = f.ByteOffsetPointers(id, "j.0:endPointer", se.ByteOffsetPointer); err != nil {
		if err = f.LineCharPointers(id, "j.0:endPointer", se.LineCharPointer); err != nil {
			return
		}
	}
	if err = f.LineCharPointers(id, "j.0:startPointer", se.LineCharPointer); err != nil {
		if err = f.ByteOffsetPointers(id, "j.0:startPointer", se.ByteOffsetPointer); err != nil {
			return
		}
	}

	return id, nil
}

func (f *Formatter) LineCharPointer(lcp *LineCharPointer) (id goraptor.Term, err error) {
	id = f.NodeId("lc")

	if err = f.setNodeType(id, TypeLineCharPointer); err != nil {
		return
	}

	err = f.addPairs(id,
		Pair{"j.0:reference", lcp.Reference.Val},
		Pair{"j.0:lineNumber", lcp.LineNumber.Val},
	)

	if err != nil {
		return
	}

	return id, nil
}
func (f *Formatter) ByteOffsetPointer(bop *ByteOffsetPointer) (id goraptor.Term, err error) {
	id = f.NodeId("bo")

	if err = f.setNodeType(id, TypeByteOffsetPointer); err != nil {
		return
	}

	err = f.addPairs(id,
		Pair{"j.0:reference", bop.Reference.Val},
		Pair{"j.0:offset", bop.Offset.Val},
	)

	if err != nil {
		return
	}

	return id, nil
}

func (f *Formatter) SnippetStartEndPointers(parent goraptor.Term, element string, ses []*SnippetStartEndPointer) error {

	if len(ses) == 0 {
		return nil
	}

	for _, se := range ses {
		if se != nil {
		}
		sepId, err := f.SnippetStartEndPointer(se)

		if err != nil {
			return err
		}
		if sepId == nil {
			continue
		}
		if err = f.addTerm(parent, element, sepId); err != nil {
			return err
		}
	}

	return nil
}

func (f *Formatter) ByteOffsetPointers(parent goraptor.Term, element string, bos []*ByteOffsetPointer) error {

	if len(bos) == 0 {
		return nil
	}

	for _, bo := range bos {
		if bo != nil {
			bopId, err := f.ByteOffsetPointer(bo)
			if err != nil {
				return err
			}
			if bopId == nil {
				continue
			}
			if err = f.addTerm(parent, element, bopId); err != nil {
				return err
			}
		}
	}
	return nil
}
func (f *Formatter) LineCharPointers(parent goraptor.Term, element string, lcs []*LineCharPointer) error {

	if len(lcs) == 0 {
		return nil
	}

	for _, lc := range lcs {
		if lc != nil {
			lcId, err := f.LineCharPointer(lc)
			if err != nil {
				return err
			}
			if lcId == nil {
				continue
			}
			if err = f.addTerm(parent, element, lcId); err != nil {
				return err
			}
		}
	}
	return nil
}

// Close to free the serializer
func (f *Formatter) Close() {
	f.serializer.EndStream()
	f.serializer.Free()
}
