package profiles

import (
	"fmt"

	"github.com/OpenNMS/onmsctl/api"
	"github.com/OpenNMS/onmsctl/common"
	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/OpenNMS/onmsctl/services"
	"github.com/urfave/cli"
)

// CliCommand the CLI command to manage server profiles
var CliCommand = cli.Command{
	Name:  "config",
	Usage: "Manage OpenNMS servers configuration profiles",
	Subcommands: []cli.Command{
		{
			Name:   "list",
			Usage:  "Lists all the configuration profiles",
			Action: listConfigProfiles,
		},
		{
			Name:   "set",
			Usage:  "Adds or updates a configuration profile",
			Action: setConfigProfile,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "name",
					Usage:    "Profile name",
					Required: true,
				},
				cli.StringFlag{
					Name:     "url",
					Value:    "http://localhost:8980/opennms",
					Usage:    "OpenNMS Base URL",
					Required: true,
				},
				cli.StringFlag{
					Name:     "user",
					Value:    "admin",
					Usage:    "OpenNMS Username (with ROLE_REST or ROLE_ADMIN)",
					Required: true,
				},
				cli.StringFlag{
					Name:     "passwd",
					Value:    "admin",
					Usage:    "OpenNMS User's Password",
					Required: true,
				},
				cli.IntFlag{
					Name:  "timeout",
					Value: 30,
					Usage: "Connection Timeout in Seconds",
				},
				cli.BoolFlag{
					Name:  "insecure",
					Usage: "To skip TLS Certificate validation",
				},
			},
		},
		{
			Name:      "default",
			Usage:     "Mark an existing configuration profile as default",
			Action:    makeDefaultConfigProfile,
			ArgsUsage: "<name>",
		},
		{
			Name:      "delete",
			Usage:     "Deletes an existing configuration profile",
			Action:    deleteConfigProfile,
			ArgsUsage: "<name>",
		},
	},
}

func listConfigProfiles(c *cli.Context) error {
	cfg, err := getAPI().GetProfilesConfig()
	if err != nil {
		return err
	}
	if cfg == nil || cfg.IsEmpty() {
		fmt.Println("There are no profiles configured")
		return nil
	}
	writer := common.NewTableWriter()
	fmt.Fprintln(writer, "Default\tName\tUser\tURL")
	for _, p := range cfg.Profiles {
		def := ""
		if cfg.Default == p.Name {
			def = "*"
		}
		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", def, p.Name, p.Username, p.URL)
	}
	writer.Flush()
	return nil
}

func setConfigProfile(c *cli.Context) error {
	profile := model.Profile{
		Name:     c.String("name"),
		URL:      c.String("url"),
		Username: c.String("user"),
		Password: c.String("passwd"),
		Timeout:  c.Int("passwd"),
		Insecure: c.Bool("insecure"),
	}
	return getAPI().SetProfile(profile)
}

func makeDefaultConfigProfile(c *cli.Context) error {
	if name := c.Args().First(); name == "" {
		return fmt.Errorf("Profile name required")
	} else {
		return getAPI().SetDefault(name)
	}
}

func deleteConfigProfile(c *cli.Context) error {
	if name := c.Args().First(); name == "" {
		return fmt.Errorf("Profile name required")
	} else {
		return getAPI().DeleteProfile(name)
	}
}

func getAPI() api.ProfilesAPI {
	return services.GetProfilesAPI(rest.Instance)
}
