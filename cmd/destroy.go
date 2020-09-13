package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wanhuasong/vss/server"
	"log"
)

func initDestroyCmd() {
	rootCmd.AddCommand(destroyCmd)
	destroyCmd.Flags().StringVar(&option.InstanceID, "instance_id", "", "instance id")
}

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "destroy a cloud server.",
	Long:  "destroy a cloud server.",
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
		if err = validateDestroyOptions(); err != nil {
			return
		}
		switch option.Provider {
		case server.ProviderAliyun:
			client := server.NewAliyunClient(config.AliyunAccessKeyID, config.AliyunAccessSecret)
			err = client.DestroyInstance(option.InstanceID)
		case server.ProviderVultr:
			client := server.NewVultrClient(config.VultrAPIKey)
			err = client.DestroyInstance(option.InstanceID)
		default:
			err = fmt.Errorf("provider not support: %s", option.Provider)
		}
		if err == nil {
			log.Printf("instance destroy success: %s", option.InstanceID)
		}
		return
	},
}

func validateDestroyOptions() (err error) {
	if option.Provider == "" {
		err = errors.New("provider cannot be null")
		return
	}
	if option.InstanceID == "" {
		err = errors.New("instance id cannot be null")
		return
	}
	return nil
}
