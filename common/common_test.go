package common

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/urfave/cli"
	"gotest.tools/assert"
)

func TestReadInput(t *testing.T) {
	text, expected := "This is a test", ""

	file, err := ioutil.TempFile("/tmp", "test-")
	assert.NilError(t, err)
	ioutil.WriteFile(file.Name(), []byte(text), 0644)
	inputStream = file // To emulate os.Stdin
	defer os.Remove(file.Name())

	app := cli.NewApp()
	app.Name = "test"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file",
			Usage: "External YAML file (use '-' for STDIN Pipe)",
		},
	}
	app.Action = func(c *cli.Context) error {
		bytes, err := ReadInput(c, 0)
		if err != nil {
			return err
		}
		expected = string(bytes)
		return nil
	}

	err = app.Run([]string{app.Name, "-file", file.Name(), "something"})
	assert.NilError(t, err)
	assert.Equal(t, expected, text)

	err = app.Run([]string{app.Name, "-file", "-", "something"})
	assert.NilError(t, err)
	assert.Equal(t, expected, text)

	err = app.Run([]string{app.Name, text, "something"})
	assert.NilError(t, err)
	assert.Equal(t, expected, text)
}
