package utils

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type config struct {
	Project struct {
		Name string `yaml:"name"`
		Port string `yaml:"port"`
	} `yaml:"project"`
	Database struct {
		Host string `yaml:"host"`
		User string `yaml:"user"`
		Pswd string `yaml:"pswd"`
		Dbnm string `yaml:"dbnm"`
		Port string `yaml:"port"`
	} `yaml:"database"`
}

type projectConfig struct {
	Name string
	Port string
}

type databaseConfig struct {
	Host string
	User string
	Pswd string
	Dbnm string
	Port string
}

type service interface {
	ProjectConfigs() (*projectConfig, error)
	DatabaseConfigs() (*databaseConfig, error)
}

type configsImpl struct{}

func (c *configsImpl) getConfig() (*config, error) {
	con := &config{}
	path, _ := os.Getwd()
	yamlFile, err := ioutil.ReadFile(fmt.Sprintf("%s/configs/config.yaml", path))
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, con)
	if err != nil {
		return nil, err
	}

	return con, nil
}

func (c *configsImpl) ProjectConfigs() (*projectConfig, error) {
	config, err := c.getConfig()
	if err != nil {
		return nil, err
	}
	return &projectConfig{
		Name: config.Project.Name,
		Port: config.Project.Port,
	}, nil
}

func (c *configsImpl) DatabaseConfigs() (*databaseConfig, error) {
	config, err := c.getConfig()
	if err != nil {
		return nil, err
	}
	return &databaseConfig{
		Host: config.Database.Host,
		User: config.Database.User,
		Pswd: config.Database.Pswd,
		Dbnm: config.Database.Dbnm,
		Port: config.Database.Port,
	}, nil
}

func NewConfigs() service {
	return &configsImpl{}
}
