package hcl2template

import (
	"testing"
)

func TestParse_communicator(t *testing.T) {
	defaultParser := getBasicParser()

	tests := []parseTest{
		{"no marshalling",
			defaultParser,
			parseTestArgs{"testdata/communicator/basic.pkr.hcl"},
			&PackerConfig{
				Communicators: map[CommunicatorRef]*Communicator{
					CommunicatorRef{
						Type: "ssh",
						Name: "vagrant",
					}: &Communicator{
						Type:               "ssh",
						Name:               "vagrant",
						CommunicatorConfig: basicMockCommunicator,
					},
				},
			},
			false,
		},
		{"marshalling",
			defaultParser,
			parseTestArgs{"testdata/communicator/basic.pkr.hcl"},
			&PackerConfig{
				Communicators: map[CommunicatorRef]*Communicator{
					CommunicatorRef{
						Type: "ssh",
						Name: "vagrant",
					}: &Communicator{
						Type:               "ssh",
						Name:               "vagrant",
						CommunicatorConfig: basicMockCommunicator,
					},
				},
			},
			false,
		},
	}
	testParse(t, tests)
}
