package conf

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Gateway *Gateway `yaml:"gateway"`
	Client  *Client  `yaml:"client"`
}

type Gateway struct {
	Addr string `yaml:"addr"`
}

type Client struct {
	Auth    Auth    `yaml:"auth"`
	Comment Comment `yaml:"comment"`
}

type Auth struct {
	Addr string `yaml:"addr"`
}

type Comment struct {
	Addr string `yaml:"addr"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
