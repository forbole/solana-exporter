package collector

import (
	"sort"
	"strconv"

	"github.com/forbole/solana-exporter/types"
	"github.com/prometheus/client_golang/prometheus"
)

type SolanaValidatorCollector struct {
	SolanaClient       *types.Client
	ValidatorAddresses map[string]struct{}

	Stakedtotal                          *prometheus.Desc
	ValidatorCommission                  *prometheus.Desc
	ValidatorDelegatorCount              *prometheus.Desc
	ValidatorEpochCredits                *prometheus.Desc
	ValidatorInflationReward             *prometheus.Desc
	ValidatorInflationRewardCurrentEpoch *prometheus.Desc
	ValidatorStaked                      *prometheus.Desc
	ValidatorStakedRanking               *prometheus.Desc
}

func NewSolanaValidatorCollector(solanaClient *types.Client, validator_addresses map[string]struct{}) *SolanaValidatorCollector {
	return &SolanaValidatorCollector{
		SolanaClient:       solanaClient,
		ValidatorAddresses: validator_addresses,

		Stakedtotal: prometheus.NewDesc(
			"solana_staked_total",
			"Total activated staked of the network",
			nil,
			SOLANA_DENOM_LABEL,
		),
		ValidatorCommission: prometheus.NewDesc(
			"solana_validator_commission_rate",
			"Commission rate of the validator",
			[]string{"validator_address"},
			nil,
		),
		ValidatorDelegatorCount: prometheus.NewDesc(
			"solana_validator_delegators_count",
			"Number of delegators per validator",
			[]string{"validator_address"},
			nil,
		),
		ValidatorEpochCredits: prometheus.NewDesc(
			"solana_validator_epoch_credits",
			"Credits earned by validator by the end of each epoch",
			[]string{"validator_address", "epoch"},
			SOLANA_DENOM_LABEL,
		),
		ValidatorStaked: prometheus.NewDesc(
			"solana_validator_staked",
			"Activated stake per validator",
			[]string{"validator_address"},
			SOLANA_DENOM_LABEL,
		),
		ValidatorStakedRanking: prometheus.NewDesc(
			"solana_validator_staked_ranking",
			"Activated stake per validator",
			[]string{"validator_address"},
			nil,
		),
		ValidatorInflationReward: prometheus.NewDesc(
			"solana_validator_inflation_reward",
			"Reward earned by validator by the end of each epoch",
			[]string{"validator_address", "epoch"},
			SOLANA_DENOM_LABEL,
		),
		ValidatorInflationRewardCurrentEpoch: prometheus.NewDesc(
			"solana_validator_inflation_reward_current_epoch",
			"Reward earned by validator by the end of each epoch",
			[]string{"validator_address"},
			SOLANA_DENOM_LABEL,
		),
	}
}

func (collector *SolanaValidatorCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.Stakedtotal
	ch <- collector.ValidatorCommission
	ch <- collector.ValidatorDelegatorCount
	ch <- collector.ValidatorEpochCredits
	ch <- collector.ValidatorStaked
	ch <- collector.ValidatorStakedRanking
	ch <- collector.ValidatorInflationReward
	ch <- collector.ValidatorInflationRewardCurrentEpoch
}

func (collector *SolanaValidatorCollector) Collect(ch chan<- prometheus.Metric) {
	voteAccounts, err := collector.SolanaClient.GetVoteAccounts()

	if err != nil {
		ch <- prometheus.NewInvalidMetric(collector.Stakedtotal, err)
		ch <- prometheus.NewInvalidMetric(collector.ValidatorCommission, err)
		ch <- prometheus.NewInvalidMetric(collector.ValidatorStaked, err)
		ch <- prometheus.NewInvalidMetric(collector.ValidatorStakedRanking, err)
	}

	// Sort to get validator ranking.
	sort.Slice(voteAccounts.Current, func(i, j int) bool {
		return voteAccounts.Current[i].ActivatedStake > voteAccounts.Current[j].ActivatedStake
	})

	var stakedTotal uint64
	for index, account := range voteAccounts.Current {
		stakedTotal += account.ActivatedStake

		if _, ok := collector.ValidatorAddresses[account.VotePubkey]; ok {
			ch <- prometheus.MustNewConstMetric(collector.ValidatorCommission, prometheus.GaugeValue, float64(account.Commission), account.VotePubkey)
			ch <- prometheus.MustNewConstMetric(collector.ValidatorStaked, prometheus.GaugeValue, types.ConvertLamportToSolana(account.ActivatedStake), account.VotePubkey)
			ch <- prometheus.MustNewConstMetric(collector.ValidatorStakedRanking, prometheus.GaugeValue, float64(index+1), account.VotePubkey)

			// Reward handler
			for _, credits := range account.EpochCredits {
				// epochCredits: <array> [epoch, credits, previousCredits]
				epoch := credits[0]
				earnedCedit := credits[1] - credits[2]
				ch <- prometheus.MustNewConstMetric(collector.ValidatorEpochCredits, prometheus.GaugeValue, float64(earnedCedit), account.VotePubkey, strconv.Itoa(epoch))
			}
		}
	}
	ch <- prometheus.MustNewConstMetric(collector.Stakedtotal, prometheus.GaugeValue, types.ConvertLamportToSolana(stakedTotal))

	for address := range collector.ValidatorAddresses {
		if delegators, err := collector.SolanaClient.GetDelegatorsCount(address); err != nil {
			ch <- prometheus.NewInvalidMetric(collector.ValidatorDelegatorCount, err)
		} else {
			ch <- prometheus.MustNewConstMetric(collector.ValidatorDelegatorCount, prometheus.GaugeValue, float64(len(delegators)), address)
		}

		if rewards, err := collector.SolanaClient.GetInflationReward(address); err != nil {
			ch <- prometheus.NewInvalidMetric(collector.ValidatorInflationReward, err)
		} else {
			for _, reward := range rewards {
				ch <- prometheus.MustNewConstMetric(collector.ValidatorInflationRewardCurrentEpoch, prometheus.GaugeValue, types.ConvertLamportToSolana(reward.Amount), address)
				ch <- prometheus.MustNewConstMetric(collector.ValidatorInflationReward, prometheus.GaugeValue, types.ConvertLamportToSolana(reward.Amount), address, strconv.FormatUint(reward.Epoch, 10))
			}
		}
	}

}
