package cmd

import (
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hsowan-me/vss/server"
	"github.com/spf13/cobra"
	"log"
)

func initListCmd() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVar(&option.RegionID, "region_id", "", "option for aliyun")
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list cloud servers.",
	Long:  "list cloud servers.",
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
		if err = validateListOptions(); err != nil {
			return
		}
		if option.Provider == server.ProviderAliyun {
			var instances []ecs.Instance
			instances, err = server.ListAliyunInstances(config, option)
			if err != nil {
				return
			}
			for i, instance := range instances {
				log.Printf("%d. public ip address: %s, status: %s", i, instance.PublicIpAddress.IpAddress[0], instance.Status)
			}
		}
		return nil
	},
}

func validateListOptions() error {
	if option.Provider == server.ProviderAliyun {
		if option.RegionID == "" {
			return errors.New("region_id cannot be nul")
		}
		regionAvailable, err := isRegionAvailable(config, option.RegionID)
		if err != nil {
			return err
		}
		if !regionAvailable {
			return errors.New("region_id is not available")
		}
		return nil
	} else {
		return fmt.Errorf("provider not support: %s", option.Provider)
	}
}
