package provisioning

import (
	"encoding/json"
	"fmt"

	"github.com/OpenNMS/onmsctl/rest"
	"github.com/urfave/cli"
)

// GetRequisitionNames gets the requisition list
func GetRequisitionNames() (RequisitionsList, error) {
	jsonRequisitions, err := rest.Instance.Get("/rest/requisitionNames")
	if err != nil {
		return RequisitionsList{}, fmt.Errorf("Cannot retrieve requisition names: %s", err)
	}
	requisitions := RequisitionsList{}
	json.Unmarshal(jsonRequisitions, &requisitions)
	return requisitions, nil
}

// RequisitionExists verifies if a given requisition exists
func RequisitionExists(foreignSource string) bool {
	requisitions, err := GetRequisitionNames()
	if err != nil {
		return false
	}
	var found = false
	for _, fs := range requisitions.ForeignSources {
		if fs == foreignSource {
			found = true
			break
		}
	}
	return found
}

// GetNode gets a node from ReST using CLI context
func GetNode(c *cli.Context) (Node, error) {
	if !c.Args().Present() {
		return Node{}, fmt.Errorf("Requisition name and foreign ID required")
	}
	foreignSource := c.Args().Get(0)
	if !RequisitionExists(foreignSource) {
		return Node{}, fmt.Errorf("Requisition %s doesn't exist", foreignSource)
	}
	foreignID := c.Args().Get(1)
	if foreignID == "" {
		return Node{}, fmt.Errorf("Foreign ID required")
	}
	jsonBytes, err := rest.Instance.Get("/rest/requisitions/"+foreignSource+"/nodes/"+foreignID)
	if err != nil {
		return Node{}, fmt.Errorf("Cannot retrieve node %s from requisition %s", foreignID, foreignSource)
	}
	node := Node{}
	json.Unmarshal(jsonBytes, &node)
	return node, nil
}
