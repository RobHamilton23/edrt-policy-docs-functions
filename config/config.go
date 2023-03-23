package config

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

//go:embed config.yaml
var config string

func GetConfig() map[string]interface{} {
	var configMap map[string]interface{}
	yaml.Unmarshal([]byte(config), &configMap)
	return configMap
}
