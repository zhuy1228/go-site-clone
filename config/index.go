package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Port            int         `yaml:"port"`
	AppName         string      `yaml:"appName"`
	SiteFileDir     string      `yaml:"siteFileDir"`
	PackSiteFileDir string      `yaml:"packSiteFileDir"`
	Nginx           NginxConfig `yaml:"nginx"`
}

type NginxConfig struct {
	NginxPath     string `yaml:"nginxPath"`
	NginxConfPath string `yaml:"nginxConfPath"`
}

func LoadConfig() (*AppConfig, error) {
	// 读取文件内容
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}

	// 解析 YAML
	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// 更新配置（写回文件）
func SaveConfig(cfg *AppConfig) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	// 写回文件，0644 表示文件权限
	return os.WriteFile("config.yaml", data, 0644)
}
