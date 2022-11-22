package collector

import (
	"log"
	"sort"
	"strconv"
	"sync"

	"github.com/forbole/solana-exporter/types"
)

func (collector *SolanaCollector) CollectValidatorStats() {
	voteAccounts, err := collector.SolanaClient.GetVoteAccounts()

	if err != nil {
		ErrorGauge.WithLabelValues("get_vote_accounts").Inc()
		log.Print(err)
	}

	// Sort to get validator ranking.
	sort.Slice(voteAccounts.Current, func(i, j int) bool {
		return voteAccounts.Current[i].ActivatedStake > voteAccounts.Current[j].ActivatedStake
	})

	var stakedTotal uint64
	for index, account := range voteAccounts.Current {
		stakedTotal += account.ActivatedStake

		if _, ok := collector.ValidatorAddresses[account.VotePubkey]; ok {
			ValidatorCommission.WithLabelValues(account.VotePubkey).Set(float64(account.Commission))
			ValidatorStaked.WithLabelValues(account.VotePubkey).Set(types.ConvertLamportToSolana(account.ActivatedStake))
			ValidatorStakedRanking.WithLabelValues(account.VotePubkey).Set(float64(index + 1))

			// Reward handler
			for _, credits := range account.EpochCredits {
				// epochCredits: <array> [epoch, credits, previousCredits]
				epoch := credits[0]
				earnedCedit := credits[1] - credits[2]
				ValidatorEpochCredits.WithLabelValues(account.VotePubkey, strconv.Itoa(epoch)).Set(float64(earnedCedit))
			}
		}
	}

	Stakedtotal.WithLabelValues().Set(types.ConvertLamportToSolana(stakedTotal))

	var wg sync.WaitGroup
	for address := range collector.ValidatorAddresses {
		wg.Add(1)
		go func(address string) {
			defer wg.Done()
			if rewards, err := collector.SolanaClient.GetInflationReward(address); err != nil {
				ErrorGauge.WithLabelValues("get_validator_inflation_reward").Inc()
				log.Print(err)
			} else {
				for _, reward := range rewards {
					ValidatorInflationRewardCurrentEpoch.WithLabelValues(address).Set(types.ConvertLamportToSolana(reward.Amount))
					ValidatorInflationReward.WithLabelValues(address, strconv.FormatUint(reward.Epoch, 10)).Set(types.ConvertLamportToSolana(reward.Amount))
				}
			}
			if delegators, err := collector.SolanaClient.GetDelegatorsCount(address); err != nil {
				ErrorGauge.WithLabelValues("get_delegators_count").Inc()
				log.Print(err)
			} else {
				ValidatorDelegatorCount.WithLabelValues(address).Set(float64(len(delegators)))
			}
		}(address)
	}

	wg.Wait()
}
