package services

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/test"
	"gotest.tools/assert"
)

var mockRequisition = model.Requisition{
	Name: "Test1",
	Nodes: []model.RequisitionNode{
		{
			ForeignID: "n1",
			NodeLabel: "n1",
			Interfaces: []model.RequisitionInterface{
				{
					IPAddress:   "10.0.0.1",
					SnmpPrimary: "P",
					Services: []model.RequisitionMonitoredService{
						{Name: "HTTP"},
					},
				},
			},
			Categories: []model.RequisitionCategory{
				{Name: "Server"},
			},
			Assets: []model.RequisitionAsset{
				{Name: "city", Value: "Durham"},
			},
			MetaData: []model.RequisitionMetaData{
				{Key: "owner", Value: "agalue"},
			},
		},
	},
}

type mockRequisitionsRest struct {
	t *testing.T
}

func (api mockRequisitionsRest) Get(path string) ([]byte, error) {
	switch path {
	case "/rest/requisitionNames":
		return []byte(`{"count":2,"foreign-source":["Test1","Test2"]}`), nil
	case "/rest/foreignSourcesConfig/assets":
		return []byte(`{"count":3,"element":["address1","city","zip"]}`), nil
	case "/rest/foreignSourcesConfig/policies":
		return []byte(test.PoliciesJSON), nil
	case "/rest/foreignSourcesConfig/detectors":
		return []byte(test.DetectorsJSON), nil
	case "/rest/requisitions/deployed/stats":
		return []byte(`{"count":1,"foreign-source":[{"count":1,"last-imported":1567526173532,"name":"Test"}]}`), nil
	case "/rest/requisitions/Test1":
		return json.Marshal(mockRequisition)
	case "/rest/requisitions/Test1/nodes/n1":
		return json.Marshal(mockRequisition.Nodes[0])
	case "/rest/requisitions/Test1/nodes/n1/interfaces/10.0.0.1":
		return json.Marshal(mockRequisition.Nodes[0].Interfaces[0])
	}
	return nil, fmt.Errorf("GET: sould not be called with path %s", path)
}

func (api mockRequisitionsRest) Post(path string, jsonBytes []byte) error {
	switch path {
	case "/rest/requisitions":
		fs := &model.Requisition{}
		if err := json.Unmarshal(jsonBytes, fs); err != nil {
			return err
		}
		assert.Assert(api.t, fs.Name != "")
		return nil
	case "/rest/requisitions/Test1/nodes":
		node := &model.RequisitionNode{}
		if err := json.Unmarshal(jsonBytes, node); err != nil {
			return err
		}
		assert.Assert(api.t, node.ForeignID != "")
		return nil
	case "/rest/requisitions/Test1/nodes/n1/interfaces":
		intf := &model.RequisitionInterface{}
		if err := json.Unmarshal(jsonBytes, intf); err != nil {
			return err
		}
		assert.Assert(api.t, intf.IPAddress != "")
		return nil
	case "/rest/requisitions/Test1/nodes/n1/interfaces/10.0.0.1/services":
		svc := &model.RequisitionMonitoredService{}
		if err := json.Unmarshal(jsonBytes, svc); err != nil {
			return err
		}
		assert.Assert(api.t, svc.Name != "")
		return nil
	case "/rest/requisitions/Test1/nodes/n1/categories":
		cat := &model.RequisitionCategory{}
		if err := json.Unmarshal(jsonBytes, cat); err != nil {
			return err
		}
		assert.Assert(api.t, cat.Name != "")
		return nil
	case "/rest/requisitions/Test1/nodes/n1/assets":
		asset := &model.RequisitionAsset{}
		if err := json.Unmarshal(jsonBytes, asset); err != nil {
			return err
		}
		assert.Assert(api.t, asset.Name != "")
		return nil
	}
	return fmt.Errorf("POST: sould not be called with %s", path)
}

func (api mockRequisitionsRest) Delete(path string) error {
	switch path {
	case "/rest/requisitions/deployed/Test1":
		return nil
	case "/rest/requisitions/Test1":
		return nil
	case "/rest/foreignSources/deployed/Test1":
		return nil
	case "/rest/foreignSources/Test1":
		return nil
	case "/rest/requisitions/Test1/nodes/n1":
		return nil
	case "/rest/requisitions/Test1/nodes/n1/interfaces/10.0.0.1":
		return nil
	case "/rest/requisitions/Test1/nodes/n1/interfaces/10.0.0.1/services/HTTP":
		return nil
	case "/rest/requisitions/Test1/nodes/n1/categories/Routers":
		return nil
	case "/rest/requisitions/Test1/nodes/n1/assets/city":
		return nil
	}
	return fmt.Errorf("DELETE: sould not be called with %s", path)
}

func (api mockRequisitionsRest) Put(path string, jsonBytes []byte, contentType string) error {
	switch path {
	case "/rest/requisitions/Test1/import?rescanExisting=false":
		return nil
	}
	return fmt.Errorf("PUT: sould not be called with %s", path)
}

func TestGetRequisitionsStats(t *testing.T) {
	api := GetRequisitionsAPI(&mockRequisitionsRest{t})
	stats, err := api.GetRequisitionsStats()
	assert.NilError(t, err)
	assert.Equal(t, 1, stats.Count)
	assert.Equal(t, "Test", stats.ForeignSources[0].Name)
}

func TestCreateRequisition(t *testing.T) {
	api := GetRequisitionsAPI(&mockRequisitionsRest{t})
	err := api.CreateRequisition("Example")
	assert.NilError(t, err)

	err = api.CreateRequisition(mockRequisition.Name)
	assert.ErrorContains(t, err, "already exist")
}

func TestGetRequisition(t *testing.T) {
	api := GetRequisitionsAPI(&mockRequisitionsRest{t})
	req, err := api.GetRequisition(mockRequisition.Name)
	assert.NilError(t, err)
	assert.Equal(t, mockRequisition.Name, req.Name)
}

func TestSetRequisition(t *testing.T) {
	api := GetRequisitionsAPI(&mockRequisitionsRest{t})
	err := api.SetRequisition(mockRequisition)
	assert.NilError(t, err)
}

func TestDeleteRequisition(t *testing.T) {
	api := GetRequisitionsAPI(&mockRequisitionsRest{t})
	err := api.DeleteRequisition(mockRequisition.Name)
	assert.NilError(t, err)
}

func TestImportRequisition(t *testing.T) {
	api := GetRequisitionsAPI(&mockRequisitionsRest{t})
	err := api.ImportRequisition(mockRequisition.Name, "false")
	assert.NilError(t, err)
}

func TestGetNode(t *testing.T) {
	api := GetRequisitionsAPI(&mockRequisitionsRest{t})
	node, err := api.GetNode(mockRequisition.Name, mockRequisition.Nodes[0].ForeignID)
	assert.NilError(t, err)
	assert.Equal(t, mockRequisition.Nodes[0].ForeignID, node.NodeLabel)
}

func TestSetNode(t *testing.T) {
	api := GetRequisitionsAPI(&mockRequisitionsRest{t})
	err := api.SetNode(mockRequisition.Name, mockRequisition.Nodes[0])
	assert.NilError(t, err)
}

func TestDeleteNode(t *testing.T) {
	api := GetRequisitionsAPI(&mockRequisitionsRest{t})
	err := api.DeleteNode(mockRequisition.Name, mockRequisition.Nodes[0].ForeignID)
	assert.NilError(t, err)
}

func TestGetInterface(t *testing.T) {
	api := GetRequisitionsAPI(&mockRequisitionsRest{t})
	n := mockRequisition.Nodes[0]
	intf, err := api.GetInterface(mockRequisition.Name, n.ForeignID, n.Interfaces[0].IPAddress)
	assert.NilError(t, err)
	assert.Equal(t, n.Interfaces[0].IPAddress, intf.IPAddress)
}

func TestSetInterface(t *testing.T) {
	api := GetRequisitionsAPI(&mockRequisitionsRest{t})
	n := mockRequisition.Nodes[0]
	err := api.SetInterface(mockRequisition.Name, n.ForeignID, n.Interfaces[0])
	assert.NilError(t, err)
}

func TestDeleteInterface(t *testing.T) {
	api := GetRequisitionsAPI(&mockRequisitionsRest{t})
	n := mockRequisition.Nodes[0]
	err := api.DeleteInterface(mockRequisition.Name, n.ForeignID, n.Interfaces[0].IPAddress)
	assert.NilError(t, err)
}

func TestSetService(t *testing.T) {
	api := GetRequisitionsAPI(&mockRequisitionsRest{t})
	n := mockRequisition.Nodes[0]
	s := model.RequisitionMonitoredService{Name: "HTTP"}
	err := api.SetService(mockRequisition.Name, n.ForeignID, n.Interfaces[0].IPAddress, s)
	assert.NilError(t, err)
}

func TestDeleteService(t *testing.T) {
	api := GetRequisitionsAPI(&mockRequisitionsRest{t})
	n := mockRequisition.Nodes[0]
	err := api.DeleteInterface(mockRequisition.Name, n.ForeignID, n.Interfaces[0].IPAddress)
	assert.NilError(t, err)
}

func TestSetCategory(t *testing.T) {
	api := GetRequisitionsAPI(&mockRequisitionsRest{t})
	c := model.RequisitionCategory{Name: "Routers"}
	err := api.SetCategory(mockRequisition.Name, mockRequisition.Nodes[0].ForeignID, c)
	assert.NilError(t, err)
}

func TestDeleteCategory(t *testing.T) {
	api := GetRequisitionsAPI(&mockRequisitionsRest{t})
	err := api.DeleteCategory(mockRequisition.Name, mockRequisition.Nodes[0].ForeignID, "Routers")
	assert.NilError(t, err)
}

func TestSetAseet(t *testing.T) {
	api := GetRequisitionsAPI(&mockRequisitionsRest{t})
	a := model.RequisitionAsset{Name: "city", Value: "Apex"}
	err := api.SetAsset(mockRequisition.Name, mockRequisition.Nodes[0].ForeignID, a)
	assert.NilError(t, err)
}

func TestDeleteAsset(t *testing.T) {
	api := GetRequisitionsAPI(&mockRequisitionsRest{t})
	err := api.DeleteAsset(mockRequisition.Name, mockRequisition.Nodes[0].ForeignID, "city")
	assert.NilError(t, err)
}
