package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wanhuasong/vss/server"
	"github.com/wanhuasong/vss/utils"
	"io/ioutil"
	"log"
	"os"
	"path"
)

var optionFile string

func initCreateCmd() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVar(&optionFile, "option", "", "option file used to create a server")
	createCmd.Flags().StringVar(&option.ScriptFile, "script_file", "", "run script file after instance start")
	// aliyun flags
	createCmd.Flags().StringVar(&option.RegionID, "region_id", "", "option for aliyun")
	createCmd.Flags().StringVar(&option.ImageID, "image_id", "", "option for aliyun")
	createCmd.Flags().StringVar(&option.InstanceType, "instance_type", "", "option for aliyun")
	createCmd.Flags().StringVar(&option.SecurityGroupID, "security_group_id", "", "option for aliyun")
	createCmd.Flags().StringVar(&option.VSwitchID, "v_switch_id", "", "option for aliyun")
	// vultr flags
	createCmd.Flags().StringVar(&option.DCID, "dcid", "", "option for vultr")
	createCmd.Flags().StringVar(&option.VPSPlanID, "vps_plan_id", "", "option for vultr")
	createCmd.Flags().StringVar(&option.OSID, "osid", "", "option for vultr")
	createCmd.Flags().StringVar(&option.FirewallGroupID, "firewall_group_id", "", "option for vultr")
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
		if err = server.ValidateConfig(config, option.Provider); err != nil {
			return
		}
		var instance *server.Instance
		switch option.Provider {
		case server.ProviderAliyun:
			client := server.NewAliyunClient(config.AliyunAccessKeyID, config.AliyunAccessSecret)
			instance, err = client.CreateInstance(option.RegionID, option.ImageID, option.InstanceType, option.SecurityGroupID, option.VSwitchID)
		case server.ProviderVultr:
			client := server.NewVultrClient(config.VultrAPIKey)
			instance, err = client.CreateInstance(option.DCID, option.VPSPlanID, option.OSID, option.FirewallGroupID)
		default:
			err = fmt.Errorf("provider not support: %s", option.Provider)
		}
		if err != nil {
			return
		}
		log.Printf("instance id: %s", instance.ID)
		log.Printf("instance ip: %s", instance.PublicIP)
		log.Printf("instance password: %s", instance.Password)
		log.Printf("instance status: %s", instance.Status)
		if option.ScriptFile != "" {
			var f os.FileInfo
			f, err = os.Stat(option.ScriptFile)
			if err != nil {
				return
			}
			if f.IsDir() {
				err = fmt.Errorf("script file is not a regular file: %s", option.ScriptFile)
				return
			}
			log.Printf("start run script file remote: %s", option.ScriptFile)
			filename := path.Base(option.ScriptFile)
			remotePath := path.Join(utils.DefaultDir, filename)
			err = utils.RunScriptRemote(utils.DefaultSSHUser, instance.Password, instance.PublicIP, utils.DefaultSSHPort, option.ScriptFile, remotePath, utils.DefaultPermissions)
		}
		return
	},
}

func validateCreateOptions() (err error) {
	if option.Provider == "" {
		err = errors.New("provider cannot be null")
		return
	}
	if optionFile != "" {
		err = initOption()
		if err != nil {
			return
		}
	}
	switch option.Provider {
	case server.ProviderAliyun:
		if option.RegionID == "" {
			return errors.New("region_id cannot be null")
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
	case server.ProviderVultr:
		if option.DCID == "" {
			return errors.New("dcid cannot be null")
		}
		if option.VPSPlanID == "" {
			return errors.New("vps_plan_id cannot be null")
		}
		if option.OSID == "" {
			return errors.New("osid cannot be null")
		}
	default:
		return fmt.Errorf("provider not support: %s", option.Provider)
	}
	return
}

func initOption() error {
	f, err := os.Stat(optionFile)
	if err != nil {
		return err
	}
	if f.IsDir() {
		return fmt.Errorf("%s is not a regular file", optionFile)
	}
	b, err := ioutil.ReadFile(optionFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, option)
}
