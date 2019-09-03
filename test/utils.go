package test

import (
	"github.com/urfave/cli"
)

// PoliciesJSON A JSON representation of a policy list (for testing purposes)
var PoliciesJSON = `
{
	"plugins": [
		{
			"name": "Set Node Category",
			"class": "org.opennms.netmgt.provision.persist.policies.NodeCategorySettingPolicy",
			"parameters": [
				{
					"key": "category",
					"required": true,
					"options": []
				},
				{
					"key": "matchBehavior",
					"required": true,
					"options": [
						"ALL_PARAMETERS",
						"ANY_PARAMETER",
						"NO_PARAMETERS"
					]
				},
				{
					"key": "foreignId",
					"required": false,
					"options": []
				},
				{
					"key": "foreignSource",
					"required": false,
					"options": []
				},
				{
					"key": "label",
					"required": false,
					"options": []
				},
				{
					"key": "labelSource",
					"required": false,
					"options": []
				},
				{
					"key": "netBiosDomain",
					"required": false,
					"options": []
				},
				{
					"key": "netBiosName",
					"required": false,
					"options": []
				},
				{
					"key": "operatingSystem",
					"required": false,
					"options": []
				},
				{
					"key": "sysContact",
					"required": false,
					"options": []
				},
				{
					"key": "sysDescription",
					"required": false,
					"options": []
				},
				{
					"key": "sysLocation",
					"required": false,
					"options": []
				},
				{
					"key": "sysName",
					"required": false,
					"options": []
				},
				{
					"key": "sysObjectId",
					"required": false,
					"options": []
				},
				{
					"key": "type",
					"required": false,
					"options": []
				}
			]
		}
	],
	"count": 1,
	"totalCount": 1,
	"offset": 0
}
`

// DetectorsJSON A JSON representation of a detector list (for testing purposes)
var DetectorsJSON = `
{
  "plugins": [
    {
      "name": "ICMP",
      "class": "org.opennms.netmgt.provision.detector.icmp.IcmpDetector",
      "parameters": [
        {
          "key": "allowFragmentation",
          "required": false,
          "options": []
        },
        {
          "key": "dscp",
          "required": false,
          "options": []
        },
        {
          "key": "ipMatch",
          "required": false,
          "options": []
        },
        {
          "key": "port",
          "required": false,
          "options": []
        },
        {
          "key": "retries",
          "required": false,
          "options": []
        },
        {
          "key": "serviceName",
          "required": false,
          "options": []
        },
        {
          "key": "timeout",
          "required": false,
          "options": []
        }
      ]
    },
    {
      "name": "SNMP",
      "class": "org.opennms.netmgt.provision.detector.snmp.SnmpDetector",
      "parameters": [
        {
          "key": "forceVersion",
          "required": false,
          "options": []
        },
        {
          "key": "hex",
          "required": false,
          "options": []
        },
        {
          "key": "ipMatch",
          "required": false,
          "options": []
        },
        {
          "key": "isTable",
          "required": false,
          "options": []
        },
        {
          "key": "matchType",
          "required": false,
          "options": []
        },
        {
          "key": "oid",
          "required": false,
          "options": []
        },
        {
          "key": "port",
          "required": false,
          "options": []
        },
        {
          "key": "retries",
          "required": false,
          "options": []
        },
        {
          "key": "serviceName",
          "required": false,
          "options": []
        },
        {
          "key": "timeout",
          "required": false,
          "options": []
        },
        {
          "key": "vbvalue",
          "required": false,
          "options": []
        }
      ]
    }
	],
	"count": 2,
	"totalCount": 2,
	"offset": 0
}
`

// CreateCli Creates a CLI Application object
func CreateCli(cmd cli.Command) *cli.App {
	var app = cli.NewApp()
	app.Name = "onmsctl"
	app.Commands = []cli.Command{cmd}
	return app
}
