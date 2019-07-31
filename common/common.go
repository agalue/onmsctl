package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/tabwriter"

	"github.com/OpenNMS/onmsctl/rest"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

// TableWriterOutput the default output for table writers
var TableWriterOutput = os.Stdout

// Reads YAML configuration from file and place it on a target object
func init() {
	configFile := getConfigFile()
	if fileExists(configFile) {
		data, _ := ioutil.ReadFile(configFile)
		err := yaml.Unmarshal(data, &rest.Instance)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: cannot read configuration file %s; %s\n", configFile, err)
			os.Exit(1)
		}
	}
}

// NewTableWriter creates a new table writer
func NewTableWriter() *tabwriter.Writer {
	return tabwriter.NewWriter(TableWriterOutput, 0, 8, 1, '\t', tabwriter.AlignRight)
}

// ReadInput reads data from a file specified on the CLI context
func ReadInput(c *cli.Context, dataIndex int) ([]byte, error) {
	var data []byte
	ymlFile := c.String("file")
	if ymlFile == "" {
		arg := c.Args().Get(dataIndex)
		if arg == "" {
			return nil, fmt.Errorf("YAML content cannot be empty")
		}
		data = []byte(arg)
	} else if ymlFile == "-" { // TODO Does this work on Windows ?
		fi, err := os.Stdin.Stat()
		if err != nil {
			return nil, err
		}
		if fi.Mode()&os.ModeNamedPipe == 0 {
			return nil, fmt.Errorf("There is no YAML content on STDIN pipe")
		}
		data, _ = ioutil.ReadAll(os.Stdin)
	} else {
		if fileExists(ymlFile) {
			data, _ = ioutil.ReadFile(ymlFile)
		} else {
			return nil, fmt.Errorf("YAML file %s doesn't exist", ymlFile)
		}
	}
	return data, nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getConfigFile() string {
	homeDir, _ := os.UserHomeDir()
	configFile := homeDir + string(os.PathSeparator) + ".onms" + string(os.PathSeparator) + "config.yaml"
	return getEnv("ONMSCONFIG", configFile)
}
