package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wanhuasong/vss/server"
	"log"
)

const (
	listTypeInstance  = "instance"
	listTypeRegion    = "region"
	listTypeDCID      = "dcid"
	listTypeOSID      = "osid"
	listTypeVPSPlanID = "vps_plan_id"
)

func initListCmd() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVar(&option.RegionID, "region_id", "", "option for aliyun")
	listCmd.Flags().StringVar(&option.ListType, "list_type", "instance",
		"list type, support server|region|dcid|osid|vpsplanid, default instance")
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
		switch option.ListType {
		case listTypeInstance:
			var instances []*server.Instance
			switch option.Provider {
			case server.ProviderAliyun:
				client := server.NewAliyunClient(config.AliyunAccessKeyID, config.AliyunAccessSecret)
				instances, err = client.ListInstances(option.RegionID)
			case server.ProviderVultr:
				client := server.NewVultrClient(config.VultrAPIKey)
				instances, err = client.ListInstances()
			}
			if err != nil {
				return
			}
			for _, instance := range instances {
				log.Printf("instance id: %s, public ip address: %s, status: %s", instance.ID, instance.PublicIP, instance.Status)
			}
		}

		return
	},
}

func validateListOptions() (err error) {
	if option.ListType == listTypeInstance && option.Provider == "" {
		err = errors.New("provider cannot be null while list type is instance")
		return
	}
	if option.ListType == "" {
		err = errors.New("list type cannot be null")
		return
	}
	switch option.Provider {
	case server.ProviderAliyun:
		if option.RegionID == "" {
			err = errors.New("region_id cannot be nul")
			return
		}
	case server.ProviderVultr:
		return
	default:
		err = fmt.Errorf("provider not support: %s", option.Provider)
	}
	return
}
