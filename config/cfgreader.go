package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

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
		StartYear     int      `yaml:"start_year"`
		EndYear       int      `yaml:"end_year"`
		PocDir        string   `yaml:"poc_dir"`
	}

	POC struct {
		DownloadPOC bool `yaml:"downloadPOC"`
	}

	Exploitdb struct {
		EdbPocDir string `yaml:"poc_dir"`
	}

	PacketStorm struct {
		PacketstormPocDir string `yaml:"poc_dir"`
	}

	Email struct {
		SmtpServer string   `yaml:"smtp_server"`
		SmtpPort   string   `yaml:"smtp_port"`
		Username   string   `yaml:"username"`
		Password   string   `yaml:"password"`
		From       string   `yaml:"from"`
		To         []string `yaml:"to"`
	}
}

func getConfigPath() string {
	dir, _ := filepath.Abs(filepath.Dir(filepath.Join("")))
	configPath := filepath.Join(dir, "config", "config.yaml")
	return configPath
}

func LoadConfig() *Config {
	cfgPath := getConfigPath()
	yamlFile, err := ioutil.ReadFile(cfgPath)
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
