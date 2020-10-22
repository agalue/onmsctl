package search

import (
	"encoding/json"
	"fmt"

	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

// Entities list of valid searchable entities
var Entities = &model.EnumValue{
	Enum: []string{"nodes", "events", "alarms", "outages"},
}

// CliCommand the CLI command to provide search capabilities information
var CliCommand = cli.Command{
	Name:  "search",
	Usage: "Search OpenNMS database",
	Flags: []cli.Flag{
		cli.GenericFlag{
			Name:     "entity, e",
			Value:    Entities,
			Usage:    "The severity of the event: " + Entities.EnumAsString(),
			Required: true,
		},
		cli.StringFlag{
			Name:  "filter, f",
			Usage: "The filter to apply in FIQL format",
		},
		cli.IntFlag{
			Name:  "limit, l",
			Usage: "The amount of entities per query",
			Value: 10,
		},
		cli.IntFlag{
			Name:  "offset, o",
			Usage: "The starting entity index (for pagination)",
			Value: 0,
		},
	},
	Action: func(c *cli.Context) error {
		entity := c.String("entity")
		if entity == "" {
			return fmt.Errorf("Entity required; options: %s", Entities.EnumAsString())
		}
		url := fmt.Sprintf("/api/v2/%s?limit=%d&offset=%d", entity, c.Int("limit"), c.Int("offset"))
		filter := c.String("filter")
		if filter != "" {
			url += "&_s=" + filter
		}
		jsonBytes, err := rest.Instance.Get(url)
		if err != nil {
			return err
		}
		if len(jsonBytes) == 0 {
			fmt.Printf("There is no data for %s\n", entity)
			return nil
		}
		var data interface{}
		switch entity {
		case "nodes":
			data = &model.OnmsNodeList{}
		case "events":
			data = &model.OnmsEventList{}
		case "alarms":
			data = &model.OnmsAlarmList{}
		case "outages":
			data = &model.OnmsOutageList{}
		}
		err = json.Unmarshal(jsonBytes, data)
		if err != nil {
			return err
		}
		yamlBytes, _ := yaml.Marshal(data)
		fmt.Println(string(yamlBytes))
		return nil
	},
}
