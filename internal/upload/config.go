package upload

import (
    "encoding/json"
    "os"
)

type Config struct {
    Region          string `json:"region"`
    Host            string `json:"host"`
    SecretAccessKey string `json:"secretAccessKey"`
    AccessKeyId     string `json:"accessKeyId"`
    Bucket          string `json:"bucket"`
}

func ParseConfig(configDir string) (*Config, error) {
    file, err := os.Open(configDir)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    decoder := json.NewDecoder(file)
    var cfg Config
    err = decoder.Decode(&cfg)
    if err != nil {
        return nil, err
    }

    return &cfg, nil
}
