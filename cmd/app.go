package cmd

import (
	"fmt"
	"log"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Use for the config file
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "pman",
		Short: "Go Pass-Man: a simple password manager",
		Long:  `Go Pass-Man: a simple CLI password manager coded in GO.`,
	}
)

// Execute executes the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {

	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pman.yaml)")

}

func initConfig() {
	if cfgFile != "" {
		// use configFile from flag
		viper.SetConfigFile(cfgFile)
	} else {
		// find home directory
		home, err := homedir.Dir()
		if err != nil {
			er(err)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".pman")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func er(msg interface{}) {

	log.Fatal(msg)

}
