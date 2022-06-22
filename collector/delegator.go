package collector

import (
	"strconv"
	"sync"

	"github.com/forbole/solana-exporter/types"
	"github.com/prometheus/client_golang/prometheus"
)

// Denom is added to easily integrate with price exporter.
var SOLANA_DENOM_LABEL = map[string]string{
	"denom": "sol",
}

type SolanaDelegatorCollector struct {
	SolanaClient       *types.Client
	DelegatorAddresses []string

	AvailableBalance                     *prometheus.Desc
	DelegatorInflationReward             *prometheus.Desc
	DelegatorInflationRewardCurrentEpoch *prometheus.Desc
}

func NewSolanaDelegatorCollector(solanaclient *types.Client, delegator_addresses []string) *SolanaDelegatorCollector {
	return &SolanaDelegatorCollector{
		SolanaClient:       solanaclient,
		DelegatorAddresses: delegator_addresses,

		AvailableBalance: prometheus.NewDesc(
			"solana_available_balance",
			"Solana available balance",
			[]string{"delegator_address"},
			SOLANA_DENOM_LABEL,
		),
		DelegatorInflationReward: prometheus.NewDesc(
			"solana_delegator_inflation_reward",
			"Reward earned by delegator by the end of each epoch",
			[]string{"delegator_address", "epoch"},
			SOLANA_DENOM_LABEL,
		),
		DelegatorInflationRewardCurrentEpoch: prometheus.NewDesc(
			"solana_delegator_inflation_reward_current_epoch",
			"Reward earned by delegator by the end of each epoch",
			[]string{"delegator_address"},
			SOLANA_DENOM_LABEL,
		),
	}
}

func (collector *SolanaDelegatorCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.AvailableBalance
	ch <- collector.DelegatorInflationReward
	ch <- collector.DelegatorInflationRewardCurrentEpoch
}

func (collector *SolanaDelegatorCollector) Collect(ch chan<- prometheus.Metric) {
	var wg sync.WaitGroup
	for _, address := range collector.DelegatorAddresses {
		wg.Add(1)
		go func(address string) {
			defer wg.Done()
			if balance, err := collector.SolanaClient.GetBalance(address); err != nil {
				ch <- prometheus.NewInvalidMetric(collector.AvailableBalance, err)
			} else {
				ch <- prometheus.MustNewConstMetric(collector.AvailableBalance, prometheus.GaugeValue, types.ConvertLamportToSolana(balance.Value), address)
			}

			if rewards, err := collector.SolanaClient.GetInflationReward(address); err != nil {
				ch <- prometheus.NewInvalidMetric(collector.DelegatorInflationReward, err)
			} else {
				for _, reward := range rewards {
					ch <- prometheus.MustNewConstMetric(collector.DelegatorInflationReward, prometheus.GaugeValue, types.ConvertLamportToSolana(reward.Amount), address, strconv.FormatUint(reward.Epoch, 10))
					ch <- prometheus.MustNewConstMetric(collector.DelegatorInflationRewardCurrentEpoch, prometheus.GaugeValue, types.ConvertLamportToSolana(reward.Amount), address)
				}
			}
		}(address)
	}

	wg.Wait()

}
