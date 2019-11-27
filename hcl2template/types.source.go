package hcl2template

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/packer/packer"
)

// A source field in an HCL file will load into the Source type.
//
type Source struct {
	// Type of source; ex: virtualbox-iso
	Type string
	// Given name; if any
	Name string

	startBuilder func() (packer.Builder, hcl.Diagnostics)

	HCL2Ref HCL2Ref
}

func (p *Parser) decodeSource(block *hcl.Block) (*Source, hcl.Diagnostics) {
	source := &Source{
		Type: block.Labels[0],
		Name: block.Labels[1],
	}
	source.HCL2Ref.DeclRange = block.DefRange

	starter := func() (packer.Builder, hcl.Diagnostics) {
		var diags hcl.Diagnostics

		// calling BuilderSchemas will start a new builder plugin to ask about
		// the schema of the builder; but we do not know yet if the builder is
		// actually going to be used. This also allows to call the same builder
		// more than once.
		builder, err := p.BuilderSchemas(source.Type)
		if err != nil {
			diags = append(diags, &hcl.Diagnostic{
				Summary: "Failed to load " + sourceLabel + " type",
				Detail:  err.Error(),
				Subject: &block.LabelRanges[0],
			})
			return builder, diags
		}

		decoded, moreDiags := decodeHCL2Spec(block, nil, builder)
		diags = append(diags, moreDiags...)
		warning, err := builder.Prepare(decoded)
		moreDiags = warningErrorsToDiags(block, warning, err)
		diags = append(diags, moreDiags...)
		return builder, diags
	}

	source.startBuilder = starter

	return source, nil
}

func (source *Source) Ref() SourceRef {
	return SourceRef{
		Type: source.Type,
		Name: source.Name,
	}
}

type SourceRef struct {
	Type string
	Name string
}

// NoSource is the zero value of sourceRef, representing the absense of an
// source.
var NoSource SourceRef

func sourceRefFromAbsTraversal(t hcl.Traversal) (SourceRef, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	if len(t) != 3 {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid " + sourceLabel + " reference",
			Detail:   "A " + sourceLabel + " reference must have three parts separated by periods: the keyword \"" + sourceLabel + "\", the builder type name, and the source name.",
			Subject:  t.SourceRange().Ptr(),
		})
		return NoSource, diags
	}

	if t.RootName() != sourceLabel {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid " + sourceLabel + " reference",
			Detail:   "The first part of an source reference must be the keyword \"" + sourceLabel + "\".",
			Subject:  t[0].SourceRange().Ptr(),
		})
		return NoSource, diags
	}
	btStep, ok := t[1].(hcl.TraverseAttr)
	if !ok {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid " + sourceLabel + " reference",
			Detail:   "The second part of an " + sourceLabel + " reference must be an identifier giving the builder type of the " + sourceLabel + ".",
			Subject:  t[1].SourceRange().Ptr(),
		})
		return NoSource, diags
	}
	nameStep, ok := t[2].(hcl.TraverseAttr)
	if !ok {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid " + sourceLabel + " reference",
			Detail:   "The third part of an " + sourceLabel + " reference must be an identifier giving the name of the " + sourceLabel + ".",
			Subject:  t[2].SourceRange().Ptr(),
		})
		return NoSource, diags
	}

	return SourceRef{
		Type: btStep.Name,
		Name: nameStep.Name,
	}, diags
}

func (r SourceRef) String() string {
	return fmt.Sprintf("%s.%s", r.Type, r.Name)
}
