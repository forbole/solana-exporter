package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/forbole/solana-exporter/types"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	homeDir string
	config  types.Config
)

var rootCmd = &cobra.Command{
	Use:   "solana_exporter",
	Short: "A solana exporter to get validator and delegator balances",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	handleInitError(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&homeDir, "home", "", "Directory for config and data (default is $HOME/.solana_exporter)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if homeDir != "" {
		cfgFile := path.Join(homeDir, "config.yaml")
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		handleInitError(err)
		viper.AddConfigPath(path.Join(home, ".solana_exporter"))
		viper.SetConfigName("config")
	}
}

func handleInitError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
