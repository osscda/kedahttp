package cli

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "capps",
		Short: "Auto-scaling containers",
		Long:  `This project implements a prototype of auto-scaling containers on either Kubernetes or ACI`,
	}
)

// Execute executes the root command.
func Execute() error {
	rootCmd.AddCommand(newUndeployCmd())
	rootCmd.AddCommand(newDeployCmd())
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			er(err)
		}

		// Search for config files named "cappsconfig" first in local dir, then $HOME/.capps then $HOME
		viper.AddConfigPath(".")
		viper.AddConfigPath(home + "/.capps/")
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("cappsconfig")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
