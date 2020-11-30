package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

// ConnectionConfig represents a single proxy connection from the config file
type ConnectionConfig struct {
	SourcePort Port `yaml:"sourcePort"`
	DstPort    Port `yaml:"dstPort"`

	/*blocksAfterNAttempts bool
	nAttemps             int
	useSSL             bool
	SSLCertificatePath string*/
}

func (c *ConnectionConfig) String() string {
	return fmt.Sprintf("%s->%s", c.SourcePort.String(), c.DstPort.String())
}

// Config is the complete config of the proxy
type Config struct {
	Version     string             `yaml:"version"`
	Connections []ConnectionConfig `yaml:"connections"`
}

func (c *Config) String() string {
	sb := strings.Builder{}
	sb.WriteString(c.Version + "\n")
	for _, val := range c.Connections {
		sb.WriteString(val.String() + "\n")
	}
	return sb.String()
}

// LoadConfig creates and returns a Config structure that represents the configuration for the proxy server and its ports.
func LoadConfig(configPath string) (*Config, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	fmt.Println(config)

	return &config, nil
}
