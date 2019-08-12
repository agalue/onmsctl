package provisioning

import (
	"encoding/json"
	"fmt"

	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/urfave/cli"
)

// Formats the available file formats for requisitions and foreign source definitions
var Formats = []string{"xml", "json", "yaml"}

// GetRequisitionNames gets the requisition list
func GetRequisitionNames() (model.RequisitionsList, error) {
	jsonRequisitions, err := rest.Instance.Get("/rest/requisitionNames")
	if err != nil {
		return model.RequisitionsList{}, fmt.Errorf("Cannot retrieve requisition names: %s", err)
	}
	requisitions := model.RequisitionsList{}
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
func GetNode(c *cli.Context) (model.RequisitionNode, error) {
	node := model.RequisitionNode{}
	if !c.Args().Present() {
		return node, fmt.Errorf("Requisition name and foreign ID required")
	}
	foreignSource := c.Args().Get(0)
	if !RequisitionExists(foreignSource) {
		return node, fmt.Errorf("Requisition %s doesn't exist", foreignSource)
	}
	foreignID := c.Args().Get(1)
	if foreignID == "" {
		return node, fmt.Errorf("Foreign ID required")
	}
	jsonBytes, err := rest.Instance.Get("/rest/requisitions/" + foreignSource + "/nodes/" + foreignID)
	if err != nil {
		return node, fmt.Errorf("Cannot retrieve node %s from requisition %s", foreignID, foreignSource)
	}
	json.Unmarshal(jsonBytes, &node)
	return node, nil
}

// GetForeignSourceDef gets a foreign source definition from ReST using CLI context
func GetForeignSourceDef(c *cli.Context) (model.ForeignSourceDef, error) {
	fsDef := model.ForeignSourceDef{}
	if !c.Args().Present() {
		return fsDef, fmt.Errorf("Foreign source name required")
	}
	foreignSource := c.Args().Get(0)
	if foreignSource != "default" && !RequisitionExists(foreignSource) {
		return fsDef, fmt.Errorf("Foreign source %s doesn't exist", foreignSource)
	}
	jsonBytes, err := rest.Instance.Get("/rest/foreignSources/" + foreignSource)
	if err != nil {
		return fsDef, fmt.Errorf("Cannot retrieve foreign source definition %s", foreignSource)
	}
	json.Unmarshal(jsonBytes, &fsDef)
	return fsDef, nil
}
