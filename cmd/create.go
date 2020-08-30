package cmd

import (
	"errors"
	"fmt"
	"github.com/hsowan-me/vss/server"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

func initCreateCmd() {
	rootCmd.AddCommand(createCmd)
	// aliyun flags
	createCmd.Flags().StringVar(&option.RegionID, "region_id", "", "option for aliyun")
	createCmd.Flags().StringVar(&option.ImageID, "image_id", "", "option for aliyun")
	createCmd.Flags().StringVar(&option.InstanceType, "instance_type", "", "option for aliyun")
	createCmd.Flags().StringVar(&option.SecurityGroupID, "security_group_id", "", "option for aliyun")
	createCmd.Flags().StringVar(&option.VSwitchID, "v_switch_id", "", "option for aliyun")
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a cloud server.",
	Long:  "create a cloud server.",
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
		if err = validateCreateOptions(); err != nil {
			return
		}
		if option.Provider == server.ProviderAliyun {
			// create instance
			var instance *server.AliyunInstance
			instance, err = server.RunAliyunInstance(config, option)
			if err != nil {
				return
			}
			log.Printf("instance id: %s", instance.InstanceID)
			log.Printf("instance ip: %s", strings.Join(instance.Ips, ","))
			log.Printf("instance password: %s", instance.Password)
			log.Printf("instance status: %s", instance.Status)
		}
		return
	},
}

func validateCreateOptions() error {
	if option.Provider == "" {
		return errors.New("provider cannot be null")
	}
	if option.Provider == server.ProviderAliyun {
		if option.RegionID == "" {
			return errors.New("region_id cannot be null")
		}
		regionAvailable, err := isRegionAvailable(config, option.RegionID)
		if err != nil {
			return err
		}
		if !regionAvailable {
			return errors.New("region_id is not available")
		}
		if option.ImageID == "" {
			return errors.New("image_id cannot be null")
		}
		if option.InstanceType == "" {
			return errors.New("instance_type cannot be null")
		}
		if option.SecurityGroupID == "" {
			return errors.New("security_group_id cannot be null")
		}
		if option.VSwitchID == "" {
			return errors.New("v_switch_id cannot be null")
		}
		return nil
	} else {
		return fmt.Errorf("provider not support: %s", option.Provider)
	}
}
