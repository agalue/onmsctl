package model

import (
	"testing"

	"gotest.tools/assert"
)

func TestRequisitionObject(t *testing.T) {
	req := Requisition{
		Name: "Test",
		Nodes: []Node{
			{
				ForeignID: "opennms.com",
				Interfaces: []Interface{
					{
						IPAddress: "www.opennms.com",
					},
				},
			},
		},
	}
	var err error

	err = req.IsValid()
	assert.NilError(t, err)
	assert.Equal(t, req.Nodes[0].NodeLabel, req.Nodes[0].ForeignID)
	assert.Equal(t, req.Nodes[0].Interfaces[0].IPAddress, "34.194.50.139")
}
