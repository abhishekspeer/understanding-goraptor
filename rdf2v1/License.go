package rdf2v1

import "github.com/deltamobile/goraptor"

type License struct {
	LicenseComment                ValueStr
	LicenseName                   ValueStr
	LicenseText                   ValueStr
	StandardLicenseHeader         ValueStr
	LicenseSeeAlso                []ValueStr
	LicenseIsFsLibre              ValueStr
	StandardLicenseTemplate       ValueStr
	StandardLicenseHeaderTemplate ValueStr
	LicenseId                     ValueStr
	LicenseisOsiApproved          ValueStr
}

func (p *Parser) requestLicense(node goraptor.Term) (*License, error) {
	obj, err := p.requestElementType(node, typeLicense)
	if err != nil {
		return nil, err
	}
	return obj.(*License), err
}

func (p *Parser) MapLicense(lic *License) *builder {
	builder := &builder{t: typeLicense, ptr: lic}
	builder.updaters = map[string]updater{
		"rdfs:comment":                  update(&lic.LicenseComment),
		"name":                          update(&lic.LicenseName),
		"licenseText":                   update(&lic.LicenseText),
		"licenseId":                     update(&lic.LicenseId),
		"rdfs:seeAlso":                  updateList(&lic.LicenseSeeAlso),
		"isFsfLibre":                    update(&lic.LicenseIsFsLibre),
		"isOsiApproved":                 update(&lic.LicenseisOsiApproved),
		"standardLicenseHeader":         update(&lic.StandardLicenseHeader),
		"standardLicenseTemplate":       update(&lic.StandardLicenseTemplate),
		"standardLicenseHeaderTemplate": update(&lic.StandardLicenseTemplate),
	}
	return builder
}
