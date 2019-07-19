package main

import (
	"fmt"
	"os"

	"github.com/OpenNMS/onmsctl/cli/events"
	"github.com/OpenNMS/onmsctl/cli/info"
	"github.com/OpenNMS/onmsctl/cli/provisioning"
	"github.com/OpenNMS/onmsctl/cli/snmp"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/urfave/cli"
)

func main() {
	var app = cli.NewApp()
	initCliInfo(app)
	initCliFlags(app)
	initCliCommands(app)

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
}

func initCliInfo(app *cli.App) {
	app.Name = "onmsctl"
	app.Usage = "A CLI to manage OpenNMS"
	app.Author = "Alejandro Galue"
	app.Email = "agalue@opennms.org"
	app.Version = "1.0.0"
	app.EnableBashCompletion = true
}

func initCliFlags(app *cli.App) {
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "url",
			Value: rest.Instance.URL,
			Usage: "OpenNMS Base URL",
		},
		cli.StringFlag{
			Name:  "user",
			Value: rest.Instance.Username,
			Usage: "OpenNMS Username (with ROLE_REST or ROLE_ADMIN)",
		},
		cli.StringFlag{
			Name:  "passwd",
			Value: rest.Instance.Password,
			Usage: "OpenNMS User's Password",
		},
		cli.BoolFlag{
			Name:        "insecure-https",
			Destination: &rest.Instance.InsecureSkipVerify,
			Usage:       "to skip HTTPS certificate validation (for self-signed certificates)",
		},
	}
}

func initCliCommands(app *cli.App) {
	app.Commands = []cli.Command{
		provisioning.CliCommand,
		snmp.CliCommand,
		info.CliCommand,
		events.CliCommand,
	}
}
