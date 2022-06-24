package aws

func (p *Provider) buildTypeMapping() map[string]mapper {
	mapping := map[string]mapper{}

	p.registerGeneratedTypes(mapping)
	p.register_s3(mapping)
	p.register_eks(mapping)
	p.register_cloudfront(mapping)
	p.register_iam_manual(mapping)

	return mapping
}
