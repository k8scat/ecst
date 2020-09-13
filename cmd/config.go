package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

func initConfigCmd() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().StringVar(&config.AliyunAccessKeyID, "aliyun_access_key_id", config.AliyunAccessKeyID, "aliyun accessKeyID")
	configCmd.Flags().StringVar(&config.AliyunAccessSecret, "aliyun_access_secret", config.AliyunAccessSecret, "aliyun accessSecret")
	configCmd.Flags().StringVar(&config.VultrAPIKey, "vultr_api_key", config.VultrAPIKey, "vultr apiKey")
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "cloud server provider config.",
	Long:  "cloud server provider config.",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		defer func() {
			if p := recover(); p != nil {
				switch p := p.(type) {
				case error:
					err = p
				default:
					err = fmt.Errorf("%s", p)
				}
			}
		}()
		if err = config.WriteFile(cfgFile); err != nil {
			return
		}
		log.Println("config saved.")
		return
	},
}
