package test

import (
	"github.com/urfave/cli"
)

// CreateCli Creates a CLI Application object
func CreateCli(cmd cli.Command) *cli.App {
	var app = cli.NewApp()
	app.Name = "onmsctl"
	app.Commands = []cli.Command{cmd}
	return app
}
