package provisioning

import (
	"fmt"
	"strings"

	"github.com/OpenNMS/onmsctl/common"
	"github.com/OpenNMS/onmsctl/model"
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
	fsDef, err := getFsAPI().GetForeignSourceDef(c.Args().Get(0))
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
	policies, err := getUtilsAPI().GetAvailablePolicies()
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
	plugin, err := getFsAPI().GetPolicyConfig(c.Args().Get(0))
	if err != nil {
		return err
	}
	data, _ := yaml.Marshal(plugin)
	fmt.Println(string(data))
	return nil
}

func getPolicy(c *cli.Context) error {
	detector, err := getFsAPI().GetPolicy(c.Args().Get(0), c.Args().Get(1))
	if err != nil {
		return err
	}
	data, _ := yaml.Marshal(detector)
	fmt.Println(string(data))
	return nil
}

func setPolicy(c *cli.Context) error {
	policy := model.Policy{Name: c.Args().Get(1), Class: c.Args().Get(2)}
	params := c.StringSlice("parameter")
	for _, p := range params {
		data := strings.Split(p, "=")
		param := model.Parameter{Key: data[0], Value: data[1]}
		policy.Parameters = append(policy.Parameters, param)
	}
	return getFsAPI().SetPolicy(c.Args().Get(0), policy)
}

func applyPolicy(c *cli.Context) error {
	data, err := common.ReadInput(c, 1)
	if err != nil {
		return err
	}
	policy := model.Policy{}
	yaml.Unmarshal(data, &policy)
	return getFsAPI().SetPolicy(c.Args().Get(0), policy)
}

func deletePolicy(c *cli.Context) error {
	return getFsAPI().DeletePolicy(c.Args().Get(0), c.Args().Get(1))
}
