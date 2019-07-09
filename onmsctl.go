package main

import (
	"fmt"
	"os"

	"github.com/OpenNMS/onmsctl/cli/provisioning"
	"github.com/OpenNMS/onmsctl/common"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/urfave/cli"
)

var app = cli.NewApp()

func main() {
	common.ReadConfig(&rest.Instance)
	initCliInfo()
	initCliFlags()
	initCliCommands()

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
}

func initCliInfo() {
	app.Name = "onmsctl"
	app.Usage = "A CLI to manage OpenNMS"
	app.Author = "Alejandro Galue"
	app.Email = "agalue@opennms.org"
	app.Version = "1.0.0"
}

func initCliFlags() {
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
	}
}

func initCliCommands() {
	app.Commands = []cli.Command{
		provisioning.CliCommand,
	}
}
