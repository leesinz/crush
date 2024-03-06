package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

var configPath = filepath.Join(getConfigDir(), "config", "config.yaml")

type Config struct {
	Database struct {
		DBPort     int    `yaml:"db_port"`
		DBUsername string `yaml:"db_username"`
		DBPassword string `yaml:"db_password"`
		Name       string `yaml:"name"`
	}
	Github struct {
		GithubToken   string   `yaml:"github_token"`
		BlacklistUser []string `yaml:"blacklist"`
	}
	MSF struct {
		MsfDir string `yaml:"msf_dir"`
	}

	Vulhub struct {
		VulhubDir string `yaml:"vulhub_dir"`
	}
	Email struct {
		SMTP_SERVER string   `yaml:"smtp_server"`
		SMTP_PORT   string   `yaml:"smtp_port"`
		Username    string   `yaml:"username"`
		Password    string   `yaml:"password"`
		From        string   `yaml:"from"`
		To          []string `yaml:"to"`
	}
}

func getConfigDir() string {
	dir, _ := filepath.Abs(filepath.Dir(filepath.Join("")))
	return dir
}

func LoadConfig() *Config {
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("Error reading config file: %v", err))
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(fmt.Sprintf("Error parsing config file: %v", err))
	}

	return &config
}
