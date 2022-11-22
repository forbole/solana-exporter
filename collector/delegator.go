package collector

import (
	"log"
	"strconv"
	"sync"

	"github.com/forbole/solana-exporter/types"
)

func (collector *SolanaCollector) CollectDelegatorAddressStats() {
	var wg sync.WaitGroup
	for _, address := range collector.DelegatorAddresses {
		wg.Add(1)
		go func(address string) {
			defer wg.Done()
			if balance, err := collector.SolanaClient.GetBalance(address); err != nil {
				ErrorGauge.WithLabelValues("get_balance").Inc()
				log.Print(err)
			} else {
				AvailableBalance.WithLabelValues(address).Set(types.ConvertLamportToSolana(balance.Value))
			}

			if rewards, err := collector.SolanaClient.GetInflationReward(address); err != nil {
				ErrorGauge.WithLabelValues("get_delegator_inflation_reward").Inc()
				log.Print(err)
			} else {
				for _, reward := range rewards {
					DelegatorInflationReward.WithLabelValues(address, strconv.FormatUint(reward.Epoch, 10)).Set(types.ConvertLamportToSolana(reward.Amount))
					DelegatorInflationRewardCurrentEpoch.WithLabelValues(address).Set(types.ConvertLamportToSolana(reward.Amount))
				}
			}
		}(address)
	}

	wg.Wait()

}
