package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/tabwriter"

	"github.com/OpenNMS/onmsctl/rest"
	"github.com/OpenNMS/onmsctl/services"
	"github.com/urfave/cli"
)

// TableWriterOutput the default output for table writers
var TableWriterOutput = os.Stdout

// Reads YAML configuration from file and place it on a target object
func init() {
	services.GetProfilesAPI(rest.Instance).GetProfilesConfig()
}

// NewTableWriter creates a new table writer
func NewTableWriter() *tabwriter.Writer {
	return tabwriter.NewWriter(TableWriterOutput, 0, 8, 1, '\t', tabwriter.AlignRight)
}

// ReadInput reads data from a file specified on the CLI context
func ReadInput(c *cli.Context, dataIndex int) ([]byte, error) {
	ymlFile := c.String("file")
	if ymlFile == "" {
		arg := c.Args().Get(dataIndex)
		if arg == "" {
			return nil, fmt.Errorf("Content cannot be empty")
		}
		return []byte(arg), nil
	} else if ymlFile == "-" { // TODO Does this work on Windows ?
		fi, err := os.Stdin.Stat()
		if err != nil {
			return nil, err
		}
		if fi.Mode()&os.ModeNamedPipe == 0 {
			return nil, fmt.Errorf("There is no YAML content on STDIN pipe")
		}
		return ioutil.ReadAll(os.Stdin)
	}
	return ioutil.ReadFile(ymlFile)
}
