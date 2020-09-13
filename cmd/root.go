package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/wanhuasong/vss/server"
	"io/ioutil"
	"os"
	"runtime/debug"
)

var (
	cfgFile string
	config  *server.Config
	option  *server.Option

	rootCmd = &cobra.Command{
		Use:   "vss COMMAND [args...]",
		Short: "cloud server tool.",
		Long:  "cloud server tool which supports multiple platforms.",
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		debug.PrintStack()
		panic(err)
	}
}

func init() {
	config = new(server.Config)
	option = new(server.Option)
	initConfig()
	// global flags
	rootCmd.PersistentFlags().StringVar(&option.Provider, "provider", "", "server provider")
	// init commands
	initConfigCmd()
	initCreateCmd()
	initDestroyCmd()
	initListCmd()
}

func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	cfgFile = fmt.Sprintf("%s/.vss/config.json", home)
	fi, err := os.Stat(cfgFile)
	if err == nil && !fi.IsDir() {
		b, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			panic(err)
		}
		if err = json.Unmarshal(b, config); err != nil {
			panic(err)
		}
	}
}
