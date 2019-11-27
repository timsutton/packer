package hcl2template

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

const (
	buildFromLabel = "from"

	buildProvisionnersLabel = "provision"

	buildPostProcessLabel = "post-process"
)

var buildSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{Type: buildFromLabel, LabelNames: []string{"src"}},
		{Type: buildProvisionnersLabel},
		{Type: buildPostProcessLabel},
	},
}

type Build struct {
	// Ordered list of provisioner groups
	ProvisionerGroups ProvisionerGroups

	// Ordered list of post-provisioner groups
	PostProvisionerGroups PostProcessorGroups

	// Ordered list of output stanzas
	Froms []SourceRef

	HCL2Ref HCL2Ref
}

type Builds []*Build

func (p *Parser) decodeBuildConfig(block *hcl.Block) (*Build, hcl.Diagnostics) {
	build := &Build{}

	var b struct {
		FromSources []string `hcl:"from_sources"`
		Config      hcl.Body `hcl:",remain"`
	}
	diags := gohcl.DecodeBody(block.Body, nil, &b)

	for _, buildFrom := range b.FromSources {
		ref := sourceRefFromString(buildFrom)

		if ref == NoSource {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid " + sourceLabel + " reference",
				Detail: "A " + sourceLabel + " type must start with a letter and " +
					"may contain only letters, digits, underscores, and dashes." +
					"A valid source reference looks like: `src.type.name`",
				Subject: &block.LabelRanges[0],
			})
		}
		if !hclsyntax.ValidIdentifier(ref.Type) ||
			!hclsyntax.ValidIdentifier(ref.Name) {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid " + sourceLabel + " reference",
				Detail: "A " + sourceLabel + " type must start with a letter and " +
					"may contain only letters, digits, underscores, and dashes." +
					"A valid source reference looks like: `src.type.name`",
				Subject: &block.LabelRanges[0],
			})
		}

		build.Froms = append(build.Froms, ref)
	}

	content, diags := b.Config.Content(buildSchema)
	for _, block := range content.Blocks {
		switch block.Type {
		case buildProvisionnersLabel:
			pg, moreDiags := p.decodeProvisionerGroup(block, p.ProvisionersSchemas)
			diags = append(diags, moreDiags...)
			build.ProvisionerGroups = append(build.ProvisionerGroups, pg)
		case buildPostProcessLabel:
			pg, moreDiags := p.decodePostProcessorGroup(block, p.PostProcessorsSchemas)
			diags = append(diags, moreDiags...)
			build.PostProvisionerGroups = append(build.PostProvisionerGroups, pg)
		}
	}

	return build, diags
}
