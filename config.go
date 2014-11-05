package husky

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

type Config struct {
	config map[string]interface{}
}

func NewConfig() *Config {
	cwd, _ := os.Getwd()
	configFile := path.Join(cwd, "config.yml")

	if _, err := os.Stat(configFile); err != nil {
		panic(fmt.Sprintf("%v", err))
	}

	yamlConfig, err := ioutil.ReadFile(configFile)
	if err != nil {
		// default config
		return &Config{
			config: map[string]interface{}{
				"host": "127.0.0.1",
				"port": 8888,
				"path": "/",
			},
		}
	}

	var config map[string]interface{}
	if err := yaml.Unmarshal(yamlConfig, &config); err != nil {
		panic(fmt.Sprintf("%v\n", err))
	}

	return &Config{
		config: config,
	}
}

func (c *Config) Get(key string) (interface{}, bool) {
	if value, ok := c.config[key]; ok {
		return value, true
	}

	return nil, false
}

func (c *Config) Set(key string, value interface{}) {
	c.config[key] = value
}
