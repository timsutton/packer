package hcl2template

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/packer/packer"
)

// Provisioner represents a parsed provisioner
type Provisioner struct {
	PType string
	// Cfg is a parsed config
	Provisioner packer.Provisioner
}

type ProvisionerGroup struct {
	CommunicatorRef CommunicatorRef

	Provisioners []Provisioner
	HCL2Ref      HCL2Ref
}

// ProvisionerGroups is a slice of provision blocks; which contains
// provisioners
type ProvisionerGroups []*ProvisionerGroup

func (p *Parser) decodeProvisionerGroup(block *hcl.Block, provisionerSpecs packer.ProvisionerStore) (*ProvisionerGroup, hcl.Diagnostics) {
	var b struct {
		Communicator string   `hcl:"communicator,optional"`
		Remain       hcl.Body `hcl:",remain"`
	}

	diags := gohcl.DecodeBody(block.Body, nil, &b)

	pg := &ProvisionerGroup{}
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
		provisioner, err := provisionerSpecs.Get(block.Type)
		if err != nil {
			diags = append(diags, &hcl.Diagnostic{
				Summary: "Failed loading " + block.Type,
				Subject: block.LabelRanges[0].Ptr(),
				Detail:  err.Error(),
			})
			continue
		}
		flatProvisinerCfg, moreDiags := decodeHCL2Spec(block, nil, provisioner)
		diags = append(diags, moreDiags...)
		err = provisioner.Prepare(flatProvisinerCfg)
		if err != nil {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Failed preparing " + block.Type,
				Detail:   err.Error(),
				Subject:  block.DefRange.Ptr(),
			})
		}
		pg.Provisioners = append(pg.Provisioners, Provisioner{
			PType:       block.Type,
			Provisioner: provisioner,
		})
	}

	return pg, diags
}

func (pgs ProvisionerGroups) FirstCommunicatorRef() CommunicatorRef {
	if len(pgs) == 0 {
		return NoCommunicator
	}
	return pgs[0].CommunicatorRef
}
