package config

import (
    "os"
    "encoding/json"
)

type ServerConfig struct {
	Host string `json:"host"`
	Port string   `json:"port"`
}

func Load(configPath string) (*ServerConfig,error) {
    var serverConfig ServerConfig
	data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, err
    }

    err = json.Unmarshal(data, &serverConfig)
    if err != nil {
        return nil, err
    }
    return &serverConfig, nil
}
