package cmd

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/forbole/solana-exporter/collector"
	"github.com/forbole/solana-exporter/types"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start exporting solana metrics",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
		err := viper.Unmarshal(&config)

		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		solanaClient := types.NewSolanaClient(config.Node.Address)
		solanaCollector := collector.NewSolanaCollector(solanaClient, config.DelegatorAddresses, types.NewValidatorAddressesMap(config.ValidatorAddresses))

		go func() {
			for {
				solanaCollector.Collect()
				time.Sleep(10 * time.Minute)
			}
		}()

		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(config.Port, nil))
		fmt.Printf("Start listening on port %s", config.Port)
		return nil
	},
}
