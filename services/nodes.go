package services

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/OpenNMS/onmsctl/api"
	"github.com/OpenNMS/onmsctl/model"
)

const defaultLimit int = 10

type nodesAPI struct {
	rest api.RestAPI
}

// GetNodesAPI Obtain an implementation of the Events API
func GetNodesAPI(rest api.RestAPI) api.NodesAPI {
	return &nodesAPI{rest}
}

func (api nodesAPI) GetNodes() (*model.OnmsNodeList, error) {
	bytes, err := api.rest.Get(fmt.Sprintf("/api/v2/nodes?limit=%d&offset=0&orderBy=label", defaultLimit))
	if err != nil {
		return nil, err
	}
	list := &model.OnmsNodeList{}
	if bytes != nil && len(bytes) > 0 {
		if err = json.Unmarshal(bytes, list); err != nil {
			return nil, err
		}
	}
	if list.TotalCount > list.Count {
		pages := list.TotalCount / defaultLimit
		if list.TotalCount%defaultLimit > 0 {
			pages++
		}
		m := sync.Mutex{}
		wg := &sync.WaitGroup{}
		for i := 1; i < pages; i++ {
			wg.Add(1)
			go func(page int, wg *sync.WaitGroup) {
				defer wg.Done()
				url := fmt.Sprintf("/api/v2/nodes?limit=%d&offset=%d&orderBy=label", defaultLimit, defaultLimit*page)
				if bytes, err = api.rest.Get(url); err != nil {
					return
				}
				temp := &model.OnmsNodeList{}
				if bytes != nil && len(bytes) > 0 {
					if err = json.Unmarshal(bytes, temp); err != nil {
						return
					}
				}
				m.Lock()
				list.Nodes = append(list.Nodes, temp.Nodes...)
				m.Unlock()
			}(i, wg)
		}
		wg.Wait()
		list.Count = len(list.Nodes)
	}
	return list, nil
}

func (api nodesAPI) GetNode(nodeCriteria string) (*model.OnmsNode, error) {
	if err := api.isCriteriaValid(nodeCriteria); err != nil {
		return nil, err
	}
	bytes, err := api.rest.Get("/api/v2/nodes/" + nodeCriteria)
	if err != nil {
		return nil, err
	}
	node := &model.OnmsNode{}
	if err = json.Unmarshal(bytes, node); err != nil {
		return nil, err
	}
	return node, nil
}

func (api nodesAPI) AddNode(node *model.OnmsNode) error {
	if err := node.Validate(); err != nil {
		return err
	}
	log.Printf("Adding node %s", node.Label)
	// Search for existing nodes
	filter := fmt.Sprintf("(node.label==*%s*)", node.Label)
	if node.ForeignSource != "" || node.ForeignID != "" {
		filter += fmt.Sprintf(",(node.foreignSource==%s;node.foreignId==%s)", node.ForeignSource, node.ForeignID)
	}
	list, err := api.searchNodes(filter)
	if err != nil {
		return fmt.Errorf("Cannot search for existing nodes: %v", err)
	}
	if len(list.Nodes) > 0 {
		return fmt.Errorf("Cannot add node because found one with ID %s that matches either the Label or ForeignSource/ForeignID combination", list.Nodes[0].ID)
	}
	// Create node
	jsonBytes, err := json.Marshal(node.ExtractBasic())
	if err != nil {
		return err
	}
	response, err := api.rest.PostRaw("/api/v2/nodes", jsonBytes, "application/json")
	if err != nil {
		return nil
	}
	if err = api.rest.IsValid(response); err != nil {
		return err
	}
	// Extract nodeID from location header
	re := regexp.MustCompile(`\/(\d+)$`)
	match := re.FindStringSubmatch(response.Header.Get("Location"))
	nodeID := match[1]
	log.Printf("Node added with ID %s", nodeID)
	// Create SNMP Interfaces
	for _, intf := range node.SNMPInterfaces {
		if err = api.SetSnmpInterface(nodeID, &intf); err != nil {
			return err
		}
	}
	// Create IP Interfaces
	for _, intf := range node.IPInterfaces {
		if err = api.SetIPInterface(nodeID, &intf); err != nil {
			return err
		}
	}
	// Create Categories
	for _, cat := range node.Categories {
		if err = api.AddCategory(nodeID, &cat); err != nil {
			return err
		}
	}
	// Create Metadata
	for _, meta := range node.Meta {
		if err = api.SetNodeMetadata(nodeID, meta); err != nil {
			return err
		}
	}
	return nil
}

func (api nodesAPI) DeleteNode(nodeCriteria string) error {
	if err := api.isCriteriaValid(nodeCriteria); err != nil {
		return err
	}
	return api.rest.Delete("/api/v2/nodes/" + nodeCriteria)
}

func (api nodesAPI) GetNodeMetadata(nodeCriteria string) ([]model.MetaData, error) {
	if err := api.isCriteriaValid(nodeCriteria); err != nil {
		return nil, err
	}
	return api.getMetadata("/api/v2/nodes/" + nodeCriteria)
}

func (api nodesAPI) SetNodeMetadata(nodeCriteria string, meta model.MetaData) error {
	if err := api.isCriteriaValid(nodeCriteria); err != nil {
		return err
	}
	return api.setMetadata("/api/v2/nodes/"+nodeCriteria, meta)
}

func (api nodesAPI) DeleteNodeMetadata(nodeCriteria string, context string, key string) error {
	if err := api.isCriteriaValid(nodeCriteria); err != nil {
		return err
	}
	return api.deleteMetadata("/api/v2/nodes/"+nodeCriteria, context, key)
}

func (api nodesAPI) GetIPInterfaces(nodeCriteria string) (*model.OnmsIPInterfaceList, error) {
	if err := api.isCriteriaValid(nodeCriteria); err != nil {
		return nil, err
	}
	url := fmt.Sprintf("/api/v2/nodes/%s/ipinterfaces?limit=%d&offset=%d&orderBy=ipAddress", nodeCriteria, defaultLimit, 0)
	bytes, err := api.rest.Get(url)
	if err != nil {
		return nil, err
	}
	list := &model.OnmsIPInterfaceList{}
	if bytes != nil && len(bytes) > 0 {
		if err = json.Unmarshal(bytes, &list); err != nil {
			return nil, err
		}
	}
	if list.TotalCount > list.Count {
		pages := list.TotalCount / defaultLimit
		if list.TotalCount%defaultLimit > 0 {
			pages++
		}
		m := sync.Mutex{}
		wg := &sync.WaitGroup{}
		for i := 1; i < pages; i++ {
			wg.Add(1)
			go func(page int, wg *sync.WaitGroup) {
				defer wg.Done()
				url := fmt.Sprintf("/api/v2/nodes/%s/ipinterfaces?limit=%d&offset=%d&orderBy=ipAddress", nodeCriteria, defaultLimit, defaultLimit*page)
				if bytes, err = api.rest.Get(url); err != nil {
					return
				}
				temp := &model.OnmsIPInterfaceList{}
				if bytes != nil && len(bytes) > 0 {
					if err = json.Unmarshal(bytes, temp); err != nil {
						return
					}
				}
				m.Lock()
				list.Interfaces = append(list.Interfaces, temp.Interfaces...)
				m.Unlock()
			}(i, wg)
		}
		wg.Wait()
		list.Count = len(list.Interfaces)
	}
	return list, nil
}

func (api nodesAPI) GetIPInterface(nodeCriteria string, ipAddress string) (*model.OnmsIPInterface, error) {
	if err := api.isCriteriaValid(nodeCriteria); err != nil {
		return nil, err
	}
	bytes, err := api.rest.Get("/api/v2/nodes/" + nodeCriteria + "/ipinterfaces/" + ipAddress)
	if err != nil {
		return nil, err
	}
	intf := &model.OnmsIPInterface{}
	if err = json.Unmarshal(bytes, &intf); err != nil {
		return nil, err
	}
	return intf, nil
}

func (api nodesAPI) SetIPInterface(nodeCriteria string, intf *model.OnmsIPInterface) error {
	if err := api.isCriteriaValid(nodeCriteria); err != nil {
		return err
	}
	if err := intf.Validate(); err != nil {
		return err
	}
	ip := intf.ExtractBasic()
	if intf.IfIndex > 0 {
		snmp, err := api.GetSnmpInterface(nodeCriteria, intf.IfIndex)
		if err != nil {
			return fmt.Errorf("Cannot find SNMP Interface with ifIndex %d on Node %s", intf.IfIndex, nodeCriteria)
		}
		log.Printf("Associating SNMP interface with ID %d and ifIndex %d to %s", snmp.ID, snmp.IfIndex, ip.IPAddress)
		ip.SNMPInterface = snmp.ExtractBasic()
	}
	log.Printf("Adding IP Interface %s", intf.IPAddress)
	jsonBytes, err := json.Marshal(ip)
	if err != nil {
		return err
	}
	if err = api.rest.Post("/api/v2/nodes/"+nodeCriteria+"/ipinterfaces", jsonBytes); err != nil {
		return err
	}
	// Add Services
	for _, svc := range intf.Services {
		if err = api.SetMonitoredService(nodeCriteria, intf.IPAddress, &svc); err != nil {
			return err
		}
	}
	// Create Metadata
	for _, meta := range intf.Meta {
		if err = api.SetIPInterfaceMetadata(nodeCriteria, intf.IPAddress, meta); err != nil {
			return err
		}
	}
	return nil
}

func (api nodesAPI) DeleteIPInterface(nodeCriteria string, ipAddress string) error {
	if err := api.isCriteriaValid(nodeCriteria); err != nil {
		return err
	}
	return api.rest.Delete("/api/v2/nodes/" + nodeCriteria + "/ipinterfaces/" + ipAddress)
}

func (api nodesAPI) GetIPInterfaceMetadata(nodeCriteria string, ipAddress string) ([]model.MetaData, error) {
	if err := api.isCriteriaValid(nodeCriteria); err != nil {
		return nil, err
	}
	return api.getMetadata("/api/v2/nodes/" + nodeCriteria + "/ipinterfaces/" + ipAddress)
}

func (api nodesAPI) SetIPInterfaceMetadata(nodeCriteria string, ipAddress string, meta model.MetaData) error {
	if err := api.isCriteriaValid(nodeCriteria); err != nil {
		return err
	}
	return api.setMetadata("/api/v2/nodes/"+nodeCriteria+"/ipinterfaces/"+ipAddress, meta)
}

func (api nodesAPI) DeleteIPInterfaceMetadata(nodeCriteria string, ipAddress string, context string, key string) error {
	if err := api.isCriteriaValid(nodeCriteria); err != nil {
		return err
	}
	return api.deleteMetadata("/api/v2/nodes/"+nodeCriteria+"/ipinterfaces/"+ipAddress, context, key)
}

func (api nodesAPI) GetSnmpInterfaces(nodeCriteria string) (*model.OnmsSnmpInterfaceList, error) {
	if err := api.isCriteriaValid(nodeCriteria); err != nil {
		return nil, err
	}
	url := fmt.Sprintf("/api/v2/nodes/%s/snmpinterfaces?limit=%d&offset=%d&orderBy=ifName", nodeCriteria, defaultLimit, 0)
	bytes, err := api.rest.Get(url)
	if err != nil {
		return nil, err
	}
	list := &model.OnmsSnmpInterfaceList{}
	if bytes != nil && len(bytes) > 0 {
		if err = json.Unmarshal(bytes, &list); err != nil {
			return nil, err
		}
	}
	if list.TotalCount > list.Count {
		pages := list.TotalCount / defaultLimit
		if list.TotalCount%defaultLimit > 0 {
			pages++
		}
		m := sync.Mutex{}
		wg := &sync.WaitGroup{}
		for i := 1; i < pages; i++ {
			wg.Add(1)
			go func(page int, wg *sync.WaitGroup) {
				defer wg.Done()
				url := fmt.Sprintf("/api/v2/nodes/%s/snmpinterfaces?limit=%d&offset=%d&orderBy=ifName", nodeCriteria, defaultLimit, defaultLimit*page)
				if bytes, err = api.rest.Get(url); err != nil {
					return
				}
				temp := &model.OnmsSnmpInterfaceList{}
				if bytes != nil && len(bytes) > 0 {
					if err = json.Unmarshal(bytes, temp); err != nil {
						return
					}
				}
				m.Lock()
				list.Interfaces = append(list.Interfaces, temp.Interfaces...)
				m.Unlock()
			}(i, wg)
		}
		wg.Wait()
		list.Count = len(list.Interfaces)
	}
	return list, nil
}

func (api nodesAPI) GetSnmpInterface(nodeCriteria string, ifIndex int) (*model.OnmsSnmpInterface, error) {
	if err := api.isCriteriaValid(nodeCriteria); err != nil {
		return nil, err
	}
	bytes, err := api.rest.Get("/api/v2/nodes/" + nodeCriteria + "/snmpinterfaces/" + strconv.Itoa(ifIndex))
	if err != nil {
		return nil, err
	}
	intf := &model.OnmsSnmpInterface{}
	if err = json.Unmarshal(bytes, &intf); err != nil {
		return nil, err
	}
	return intf, nil
}

// There are issues with the APIv2 which is why the APIv1 is used
func (api nodesAPI) SetSnmpInterface(nodeCriteria string, intf *model.OnmsSnmpInterface) error {
	if err := api.isCriteriaValid(nodeCriteria); err != nil {
		return err
	}
	if err := intf.Validate(); err != nil {
		return err
	}
	log.Printf("Adding SNMP Interface with index %d", intf.IfIndex)
	dataBytes, err := xml.Marshal(intf.ExtractBasic())
	if err != nil {
		return err
	}
	response, err := api.rest.PostRaw("/rest/nodes/"+nodeCriteria+"/snmpinterfaces", dataBytes, "application/xml")
	if err != nil {
		return err
	}
	return api.rest.IsValid(response)
}

func (api nodesAPI) DeleteSnmpInterface(nodeCriteria string, ifIndex int) error {
	if err := api.isCriteriaValid(nodeCriteria); err != nil {
		return err
	}
	return api.rest.Delete("/api/v2/nodes/" + nodeCriteria + "/snmpinterfaces/" + strconv.Itoa(ifIndex))
}

func (api nodesAPI) GetMonitoredServices(nodeCriteria string, ipAddress string) (*model.OnmsMonitoredServiceList, error) {
	if err := api.isCriteriaValid(nodeCriteria); err != nil {
		return nil, err
	}
	bytes, err := api.rest.Get("/api/v2/nodes/" + nodeCriteria + "/ipinterfaces/" + ipAddress + "/services?limit=0")
	if err != nil {
		return nil, err
	}
	list := &model.OnmsMonitoredServiceList{}
	if bytes != nil && len(bytes) > 0 {
		if err = json.Unmarshal(bytes, &list); err != nil {
			return nil, err
		}
	}
	return list, nil
}

func (api nodesAPI) GetMonitoredService(nodeCriteria string, ipAddress string, service string) (*model.OnmsMonitoredService, error) {
	if err := api.isCriteriaValid(nodeCriteria); err != nil {
		return nil, err
	}
	bytes, err := api.rest.Get("/api/v2/nodes/" + nodeCriteria + "/ipinterfaces/" + ipAddress + "/services/" + service)
	if err != nil {
		return nil, err
	}
	svc := &model.OnmsMonitoredService{}
	if err = json.Unmarshal(bytes, &svc); err != nil {
		return nil, err
	}
	return svc, nil
}

func (api nodesAPI) SetMonitoredService(nodeCriteria string, ipAddress string, svc *model.OnmsMonitoredService) error {
	if err := api.isCriteriaValid(nodeCriteria); err != nil {
		return err
	}
	if err := svc.Validate(); err != nil {
		return err
	}
	jsonBytes, err := json.Marshal(svc)
	if err != nil {
		return err
	}
	return api.rest.Post("/api/v2/nodes/"+nodeCriteria+"/ipinterfaces/"+ipAddress+"/services", jsonBytes)
}

func (api nodesAPI) DeleteMonitoredService(nodeCriteria string, ipAddress string, service string) error {
	if err := api.isCriteriaValid(nodeCriteria); err != nil {
		return err
	}
	return api.rest.Delete("/api/v2/nodes/" + nodeCriteria + "/ipinterfaces/" + ipAddress + "/services/" + service)
}

func (api nodesAPI) GetMonitoredServiceMetadata(nodeCriteria string, ipAddress string, service string) ([]model.MetaData, error) {
	if err := api.isCriteriaValid(nodeCriteria); err != nil {
		return nil, err
	}
	return api.getMetadata("/api/v2/nodes/" + nodeCriteria + "/ipinterfaces/" + ipAddress + "/services/" + service)
}

func (api nodesAPI) SetMonitoredServiceMetadata(nodeCriteria string, ipAddress string, service string, meta model.MetaData) error {
	if err := api.isCriteriaValid(nodeCriteria); err != nil {
		return err
	}
	return api.setMetadata("/api/v2/nodes/"+nodeCriteria+"/ipinterfaces/"+ipAddress+"/services/"+service, meta)
}

func (api nodesAPI) DeleteMonitoredServiceMetadata(nodeCriteria string, ipAddress string, service string, context string, key string) error {
	if err := api.isCriteriaValid(nodeCriteria); err != nil {
		return err
	}
	return api.deleteMetadata("/api/v2/nodes/"+nodeCriteria+"/ipinterfaces/"+ipAddress+"/services/"+service, context, key)
}

func (api nodesAPI) GetCategories(nodeCriteria string) ([]model.OnmsCategory, error) {
	if err := api.isCriteriaValid(nodeCriteria); err != nil {
		return nil, err
	}
	bytes, err := api.rest.Get("/api/v2/nodes/" + nodeCriteria + "/categories?limit=0")
	if err != nil {
		return nil, err
	}
	list := &model.OnmsCategoryList{}
	if bytes != nil && len(bytes) > 0 {
		if err = json.Unmarshal(bytes, &list); err != nil {
			return nil, err
		}
	}
	return list.Categories, nil
}

func (api nodesAPI) AddCategory(nodeCriteria string, category *model.OnmsCategory) error {
	if err := api.isCriteriaValid(nodeCriteria); err != nil {
		return err
	}
	jsonBytes, err := json.Marshal(category)
	if err != nil {
		return err
	}
	return api.rest.Post("/api/v2/nodes/"+nodeCriteria+"/categories", jsonBytes)
}

func (api nodesAPI) DeleteCategory(nodeCriteria string, category string) error {
	if err := api.isCriteriaValid(nodeCriteria); err != nil {
		return err
	}
	return api.rest.Delete("/api/v2/nodes/" + nodeCriteria + "/categories/" + category)
}

func (api nodesAPI) GetAssetRecord(nodeCriteria string) (*model.OnmsAssetRecord, error) {
	if err := api.isCriteriaValid(nodeCriteria); err != nil {
		return nil, err
	}
	bytes, err := api.rest.Get("/api/v2/nodes/" + nodeCriteria)
	if err != nil {
		return nil, err
	}
	if bytes != nil && len(bytes) > 0 {
		n := &model.OnmsNode{}
		if err = json.Unmarshal(bytes, n); err != nil {
			return nil, err
		}
		return n.AssetRecord, nil
	}
	return nil, fmt.Errorf("Not Found")
}

func (api nodesAPI) SetAssetField(nodeCriteria string, field string, value string) error {
	if err := api.isCriteriaValid(nodeCriteria); err != nil {
		return err
	}
	data := url.Values{}
	data.Set(field, value)
	return api.rest.Put("/rest/nodes/"+nodeCriteria+"/assetRecord", []byte(data.Encode()), "application/x-www-form-urlencoded")
}

func (api nodesAPI) getMetadata(baseURL string) ([]model.MetaData, error) {
	bytes, err := api.rest.Get(baseURL + "/metadata")
	if err != nil {
		return nil, err
	}
	list := &model.MetaDataList{}
	json.Unmarshal(bytes, &list)
	return list.Metadata, nil
}

func (api nodesAPI) setMetadata(baseURL string, meta model.MetaData) error {
	if err := meta.Validate(); err != nil {
		return err
	}
	jsonBytes, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	return api.rest.Post(baseURL+"/metadata", jsonBytes)
}

func (api nodesAPI) deleteMetadata(baseURL string, context string, key string) error {
	return api.rest.Delete(baseURL + "/metadata/" + context + "/" + key)
}

func (api nodesAPI) searchNodes(fiqlFilter string) (*model.OnmsNodeList, error) {
	bytes, err := api.rest.Get("/api/v2/nodes?limit=0&_s=" + fiqlFilter)
	if err != nil {
		return nil, err
	}
	list := &model.OnmsNodeList{}
	if bytes != nil && len(bytes) > 0 {
		if err = json.Unmarshal(bytes, list); err != nil {
			return nil, err
		}
	}
	return list, nil
}

func (api nodesAPI) isCriteriaValid(nodeCriteria string) error {
	if _, err := strconv.Atoi(nodeCriteria); err == nil {
		return nil
	} else {
		parts := strings.Split(nodeCriteria, ":")
		if len(parts) == 2 {
			return nil
		}

	}
	return fmt.Errorf("Invalid node criteria. It should be either the node ID or the foreignSource:foreignID combination.")
}
