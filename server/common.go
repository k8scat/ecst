package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
)

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

func ValidateConfig(config *Config, provider string) (err error) {
	switch provider {
	case ProviderAliyun:
		if config.AliyunAccessKeyID == "" {
			err = errors.New("aliyun_access_key_id cannot be null")
		}
		if config.AliyunAccessSecret == "" {
			err = errors.New("aliyun_access_secret cannot be null")
		}
	case ProviderVultr:
		if config.VultrAPIKey == "" {
			err = errors.New("vultr_api_key cannot be null")
		}
	default:
		err = fmt.Errorf("provider not support: %s", provider)
	}
	return
}
