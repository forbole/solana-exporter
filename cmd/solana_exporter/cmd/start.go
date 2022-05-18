package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/forbole/solana-exporter/collector"
	"github.com/forbole/solana-exporter/types"
	"github.com/prometheus/client_golang/prometheus"
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

		registry := prometheus.NewPedanticRegistry()
		registry.MustRegister(
			collector.NewSolanaValidatorCollector(solanaClient, types.NewValidatorAddressesMap(config.ValidatorAddresses)),
			collector.NewSolanaDelegatorCollector(solanaClient, config.DelegatorAddresses),
		)

		handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{
			ErrorLog:      log.New(os.Stderr, log.Prefix(), log.Flags()),
			ErrorHandling: promhttp.ContinueOnError,
		})

		http.Handle("/metrics", handler)
		log.Fatal(http.ListenAndServe(config.Port, nil))
		fmt.Printf("Start listening on port %s", config.Port)
		return nil
	},
}
