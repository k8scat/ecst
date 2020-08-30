package cmd

import (
	"errors"
	"fmt"
	"github.com/hsowan-me/vss/server"
	"github.com/spf13/cobra"
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
		if option.Provider == server.ProviderAliyun {
			if err = server.DestroyAliyunInstance(config, option); err != nil {
				return
			}
			log.Println("instance destroyed")
			return
		}
		return
	},
}

func validateDestroyOptions() error {
	if option.Provider == "" {
		return errors.New("provider cannot be null")
	}
	if option.Provider == server.ProviderAliyun {
		if option.InstanceID == "" {
			return errors.New("instance id cannot be null")
		}
		return nil
	} else {
		return fmt.Errorf("provider not support: %s", option.Provider)
	}
}
