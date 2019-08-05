package provisioning

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/OpenNMS/onmsctl/common"
	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/urfave/cli"

	"gopkg.in/yaml.v2"
)

// PoliciesCliCommand the CLI command configuration for managing foreign source detectors
var PoliciesCliCommand = cli.Command{
	Name:     "policy",
	Usage:    "Manage foreign source policies",
	Category: "Foreign Source Definitions",
	Subcommands: []cli.Command{
		{
			Name:      "list",
			Usage:     "List all the policy from a given foreign source definition",
			ArgsUsage: "<foreignSource>",
			Action:    listPolicies,
		},
		{
			Name:      "enumerate",
			ShortName: "enum",
			Usage:     "Enumerate the list of available policy classes",
			Action:    enumeratePolicyClasses,
		},
		{
			Name:      "describe",
			ShortName: "desc",
			Usage:     "Describe a given policy class",
			ArgsUsage: "<policyName|ClassName>",
			Action:    describePolicyClass,
		},
		{
			Name:      "get",
			Usage:     "Gets a policy from a given foreign source definition",
			ArgsUsage: "<foreignSource> <policyName|className>",
			Action:    getPolicy,
		},
		{
			Name:      "set",
			Usage:     "Adds or update a policy for a given foreign source definition",
			ArgsUsage: "<foreignSource> <policyName> <className>",
			Action:    setPolicy,
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name:  "parameter, p",
					Usage: "A policy parameter (e.x. -p 'matchBehavior=ALL_PARAMETERS')",
				},
			},
		},
		{
			Name:   "apply",
			Usage:  "Creates or updates a policy from a external YAML file",
			Action: applyPolicy,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file, f",
					Usage: "External YAML file (use '-' for STDIN Pipe)",
				},
			},
			ArgsUsage: "<foreignSource> <yaml>",
		},
		{
			Name:      "delete",
			ShortName: "del",
			Usage:     "Deletes an existing policy from a given foreign source definition",
			ArgsUsage: "<foreignSource> <policyName>",
			Action:    deletePolicy,
		},
	},
}

func listPolicies(c *cli.Context) error {
	fsDef, err := GetForeignSourceDef(c)
	if err != nil {
		return err
	}
	writer := common.NewTableWriter()
	fmt.Fprintln(writer, "Policy Name\tPolicy Class")
	for _, policy := range fsDef.Policies {
		fmt.Fprintf(writer, "%s\t%s\n", policy.Name, policy.Class)
	}
	writer.Flush()
	return nil
}

func enumeratePolicyClasses(c *cli.Context) error {
	policies, err := getPolicies()
	if err != nil {
		return err
	}
	writer := common.NewTableWriter()
	fmt.Fprintln(writer, "Policy Name\tPolicy Class")
	for _, plugin := range policies.Plugins {
		fmt.Fprintf(writer, "%s\t%s\n", plugin.Name, plugin.Class)
	}
	writer.Flush()
	return nil
}

func describePolicyClass(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Policy name or class required")
	}
	src := c.Args().Get(0)
	policies, err := getPolicies()
	if err != nil {
		return err
	}
	for _, plugin := range policies.Plugins {
		if plugin.Class == src || plugin.Name == src {
			data, _ := yaml.Marshal(&plugin)
			fmt.Println(string(data))
			return nil
		}
	}
	return fmt.Errorf("Cannot find policy for %s", src)
}

func getPolicy(c *cli.Context) error {
	fsDef, err := GetForeignSourceDef(c)
	if err != nil {
		return err
	}
	src := c.Args().Get(1)
	if src == "" {
		return fmt.Errorf("Policy name or class required")
	}
	for _, policy := range fsDef.Policies {
		if policy.Class == src || policy.Name == src {
			data, _ := yaml.Marshal(&policy)
			fmt.Println(string(data))
			return nil
		}
	}
	return fmt.Errorf("Cannot find policy for %s", src)
}

func setPolicy(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Foreign source name, policy name and class required")
	}
	foreignSource := c.Args().Get(0)
	if !RequisitionExists(foreignSource) {
		return fmt.Errorf("Foreign source %s doesn't exist", foreignSource)
	}
	policyName := c.Args().Get(1)
	if policyName == "" {
		return fmt.Errorf("Policy name required")
	}
	policyClass := c.Args().Get(2)
	if policyClass == "" {
		return fmt.Errorf("Policy class required")
	}
	policy := model.Policy{Name: policyName, Class: policyClass}
	params := c.StringSlice("parameter")
	for _, p := range params {
		data := strings.Split(p, "=")
		param := model.Parameter{Key: data[0], Value: data[1]}
		policy.Parameters = append(policy.Parameters, param)
	}
	policies, err := getPolicies()
	if err != nil {
		return err
	}
	err = isPolicyValid(policy, policies)
	if err != nil {
		return err
	}
	jsonBytes, _ := json.Marshal(policy)
	return rest.Instance.Post("/rest/foreignSources/"+foreignSource+"/policies", jsonBytes)
}

func applyPolicy(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Foreign source name required")
	}
	foreignSource := c.Args().Get(0)
	if !RequisitionExists(foreignSource) {
		return fmt.Errorf("Foreign source %s doesn't exist", foreignSource)
	}
	data, err := common.ReadInput(c, 1)
	if err != nil {
		return err
	}
	policy := &model.Policy{}
	yaml.Unmarshal(data, policy)
	policies, err := getPolicies()
	if err != nil {
		return err
	}
	err = isPolicyValid(*policy, policies)
	if err != nil {
		return err
	}
	fmt.Printf("Updating policy %s...\n", policy.Name)
	jsonBytes, _ := json.Marshal(policy)
	return rest.Instance.Post("/rest/foreignSources/"+foreignSource+"/policies", jsonBytes)
}

func deletePolicy(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Foreign source name and policy name required")
	}
	foreignSource := c.Args().Get(0)
	if !RequisitionExists(foreignSource) {
		return fmt.Errorf("Foreign source %s doesn't exist", foreignSource)
	}
	detector := c.Args().Get(1)
	if detector == "" {
		return fmt.Errorf("Policy name required")
	}
	return rest.Instance.Delete("/rest/foreignSources/" + foreignSource + "/policies/" + detector)
}

func getPolicies() (model.PluginList, error) {
	detectors := model.PluginList{}
	jsonData, err := rest.Instance.Get("/rest/foreignSourcesConfig/policies")
	if err != nil {
		return detectors, fmt.Errorf("Cannot retrieve policy list")
	}
	json.Unmarshal(jsonData, &detectors)
	return detectors, nil
}

func isPolicyValid(policy model.Policy, config model.PluginList) error {
	if err := policy.IsValid(); err != nil {
		return err
	}
	plugin := config.FindPlugin(policy.Class)
	if plugin == nil {
		return fmt.Errorf("Cannot find policy with class %s", policy.Class)
	}
	if err := plugin.VerifyParameters(policy.Parameters); err != nil {
		return err
	}
	return nil
}
