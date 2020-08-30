package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
)

const (
	ProviderAliyun              string = "aliyun"
	defaultRegionID             string = "cn-hangzhou"
	aliyunInstanceStatusRunning string = "Running"
	aliyunInstanceStatusStopped string = "Stopped"

	ProviderVultr string = "vultr"
)

type Config struct {
	AliyunAccessKeyID  string `json:"aliyun_access_key_id"`
	AliyunAccessSecret string `json:"aliyun_access_secret"`
	VultrAPIKey        string `json:"vultr_api_key"`
}

type Option struct {
	Provider string
	// Aliyun options
	RegionID        string
	InstanceType    string
	ImageID         string
	InstanceID      string
	SecurityGroupID string
	VSwitchID       string
}

type AliyunInstance struct {
	InstanceID string
	Password   string
	Ips        []string
	Status     string
}

func (c *Config) WriteFile(cfgFile string) error {
	cfgDir := path.Dir(cfgFile)
	fi, err := os.Stat(cfgDir)
	if err != nil || !fi.IsDir() {
		if err = os.MkdirAll(cfgDir, os.ModePerm); err != nil {
			return err
		}
	}
	f, err := os.Create(cfgFile)
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	return encoder.Encode(c)
}

func ValidateConfig(config *Config, provider string) error {
	if provider == ProviderAliyun {
		if config.AliyunAccessKeyID == "" {
			return errors.New("aliyun_access_key_id cannot be null")
		}
		if config.AliyunAccessSecret == "" {
			return errors.New("aliyun_access_secret cannot be null")
		}
		return nil
	} else {
		return fmt.Errorf("provider not support: %s", provider)
	}
}
