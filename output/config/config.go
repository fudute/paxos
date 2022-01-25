package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Cluster Cluster `yaml:"cluster"`
}
type Cluster struct {
	Acceptors []*Node `yaml:"acceptors"`
	Proposers []*Node `yaml:"proposers"`
	Learners  []*Node `yaml:"learners"`
}

type Node struct {
	Name string `yaml:"name"`
	Addr string `yaml:"addr"`
}

func LoadConfig() (*Config, error) {
	conf := Config{}

	pwd, _ := os.Getwd()
	fmt.Printf("pwd: %v\n", pwd)
	file := os.Getenv("CLUSTER_CONFIG")
	if file == "" {
		file = "config/config.yaml"
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(data, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
