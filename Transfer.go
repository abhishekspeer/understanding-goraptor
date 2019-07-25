package main

import (
	"spdx/tools-golang/v0/spdx"
	"ug/understanding-goraptor/rdf2v1"
)

func transferCreationInfo(spdxdoc *rdf2v1.Document) spdx.CreationInfo2_1 {

	var list []string
	list = append(list, spdxdoc.ExternalDocumentRef.ExternalDocumentId.V())
	list = append(list, spdxdoc.ExternalDocumentRef.SPDXDocument.V())
	list = append(list, spdxdoc.ExternalDocumentRef.Checksum.Algorithm.V())

	stdCi := spdx.CreationInfo2_1{

		SPDXVersion:                spdxdoc.SPDXVersion.V(),
		DataLicense:                spdxdoc.CreationInfo.V(),
		SPDXIdentifier:             "",
		DocumentName:               spdxdoc.DocumentName.V(),
		DocumentNamespace:          "",
		ExternalDocumentReferences: list,
		LicenseListVersion:         "",
		Created:                    spdxdoc.CreationInfo.Create.Val(),
		CreatorComment:             spdxdoc.CreationInfo.Comment.V(),
		DocumentComment:            spdxdoc.DocumentComment.V(),
	}
	return stdCi
}

func transferAnnotation(spdxdoc *rdf2v1.Document) spdx.Annotation2_1 {

	var list []string
	list = append(list, spdxdoc.ExternalDocumentRef.ExternalDocumentId.V())
	list = append(list, spdxdoc.ExternalDocumentRef.SPDXDocument.V())
	list = append(list, spdxdoc.ExternalDocumentRef.Checksum.Algorithm.V())

	stdAnn := spdx.Annotation2_1{
		Annotator:                spdxdoc.Annotation.Annotator.Val,
		AnnotationType:           spdxdoc.Annotation.AnnotationType.Val,
		AnnotationDate:           spdxdoc.Annotation.AnnotationDate.Val(),
		AnnotationComment:        spdxdoc.Annotation.AnnotationComment.Val,
		AnnotationSPDXIdentifier: "",
		AnnotatorType:            "",
	}

	return stdAnn
}
