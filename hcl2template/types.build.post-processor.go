package hcl2template

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/packer/packer"
)

// PostProcessor represents a parsed PostProcessor
type PostProcessor struct {
	PType string
	// Cfg is a parsed config
	PostProcessor packer.PostProcessor
}

type PostProcessorGroup struct {
	CommunicatorRef CommunicatorRef

	PostProcessors []PostProcessor
	HCL2Ref        HCL2Ref
}

// PostProcessorGroups is a slice of provision blocks; which contains
// provisioners
type PostProcessorGroups []*PostProcessorGroup

func (p *Parser) decodePostProcessorGroup(block *hcl.Block, provisionerSpecs packer.PostProcessorStore) (*PostProcessorGroup, hcl.Diagnostics) {
	var b struct {
		Communicator string   `hcl:"communicator,optional"`
		Remain       hcl.Body `hcl:",remain"`
	}

	diags := gohcl.DecodeBody(block.Body, nil, &b)

	pg := &PostProcessorGroup{}
	pg.CommunicatorRef = communicatorRefFromString(b.Communicator)
	pg.HCL2Ref.DeclRange = block.DefRange

	buildSchema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{},
	}
	for _, k := range provisionerSpecs.List() {
		buildSchema.Blocks = append(buildSchema.Blocks, hcl.BlockHeaderSchema{
			Type: k,
		})
	}

	content, moreDiags := b.Remain.Content(buildSchema)
	diags = append(diags, moreDiags...)
	for _, block := range content.Blocks {
		postProcessor, err := provisionerSpecs.Get(block.Type)
		if err != nil {
			diags = append(diags, &hcl.Diagnostic{
				Summary: "Failed loading " + block.Type,
				Subject: block.LabelRanges[0].Ptr(),
				Detail:  err.Error(),
			})
			continue
		}
		flatProvisinerCfg, moreDiags := decodeHCL2Spec(block, nil, postProcessor)
		diags = append(diags, moreDiags...)
		err = postProcessor.Configure(flatProvisinerCfg)
		if err != nil {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Failed preparing " + block.Type,
				Detail:   err.Error(),
				Subject:  block.DefRange.Ptr(),
			})
		}
		pg.PostProcessors = append(pg.PostProcessors, PostProcessor{
			PType:         block.Type,
			PostProcessor: postProcessor,
		})
	}

	return pg, diags
}

func (pgs PostProcessorGroups) FirstCommunicatorRef() CommunicatorRef {
	if len(pgs) == 0 {
		return NoCommunicator
	}
	return pgs[0].CommunicatorRef
}
