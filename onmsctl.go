package main

import (
	"fmt"
	"os"

	"github.com/OpenNMS/onmsctl/cli/daemon"
	"github.com/OpenNMS/onmsctl/cli/events"
	"github.com/OpenNMS/onmsctl/cli/info"
	"github.com/OpenNMS/onmsctl/cli/nodes"
	"github.com/OpenNMS/onmsctl/cli/profiles"
	"github.com/OpenNMS/onmsctl/cli/provisioning"
	"github.com/OpenNMS/onmsctl/cli/resources"
	"github.com/OpenNMS/onmsctl/cli/search"
	"github.com/OpenNMS/onmsctl/cli/snmp"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/urfave/cli"
)

var (
	version = "v1.0.0-beta3"
)

func main() {
	var app = cli.NewApp()
	initCliInfo(app)
	initCliFlags(app)
	initCliCommands(app)

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}

func initCliInfo(app *cli.App) {
	app.Name = "onmsctl"
	app.Usage = "A CLI to manage OpenNMS"
	app.Author = "Alejandro Galue"
	app.Email = "agalue@opennms.org"
	app.Version = version
	app.EnableBashCompletion = true
}

func initCliFlags(app *cli.App) {
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "url",
			Value:       rest.Instance.URL,
			Destination: &rest.Instance.URL,
			Usage:       "OpenNMS Base URL",
		},
		cli.StringFlag{
			Name:        "user, u",
			Value:       rest.Instance.Username,
			Destination: &rest.Instance.Username,
			Usage:       "OpenNMS Username (with ROLE_REST or ROLE_ADMIN)",
		},
		cli.StringFlag{
			Name:        "passwd, p",
			Value:       rest.Instance.Password,
			Destination: &rest.Instance.Password,
			Usage:       "OpenNMS User's Password",
		},
		cli.IntFlag{
			Name:        "timeout, t",
			Value:       rest.Instance.Timeout,
			Destination: &rest.Instance.Timeout,
			Usage:       "Connection Timeout in Seconds",
		},
		cli.BoolFlag{
			Name:        "insecure, k",
			Destination: &rest.Instance.Insecure,
			Usage:       "Skips HTTPS certificate validation (e.x. self-signed certificates)",
		},
		cli.BoolFlag{
			Name:        "debug, d",
			Destination: &rest.Instance.Debug,
			Usage:       "Enable DEBUG for HTTP requests",
		},
	}
}

func initCliCommands(app *cli.App) {
	app.Commands = []cli.Command{
		info.CliCommand,
		provisioning.CliCommand,
		nodes.CliCommand,
		snmp.CliCommand,
		events.CliCommand,
		daemon.CliCommand,
		resources.CliCommand,
		search.CliCommand,
		profiles.CliCommand,
	}
}
