package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

var SOLANA_DENOM_LABEL = map[string]string{
	"denom": "sol",
}

var (
	Stakedtotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "solana_staked_total",
			Help:        "Total activated staked of the network",
			ConstLabels: SOLANA_DENOM_LABEL,
		},
		nil,
	)
	ValidatorCommission = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "solana_validator_commission_rate",
			Help: "Commission rate of the validator",
		},
		[]string{"validator_address"},
	)
	ValidatorDelegatorCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "solana_validator_delegators_count",
			Help: "Number of delegators per validator",
		},
		[]string{"validator_address"},
	)
	ValidatorEpochCredits = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "solana_validator_epoch_credits",
			Help:        "Credits earned by validator by the end of each epoch",
			ConstLabels: SOLANA_DENOM_LABEL,
		},
		[]string{"validator_address", "epoch"},
	)
	ValidatorStaked = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "solana_validator_staked",
			Help:        "Activated stake per validator",
			ConstLabels: SOLANA_DENOM_LABEL,
		},
		[]string{"validator_address"},
	)
	ValidatorStakedRanking = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "solana_validator_staked_ranking",
			Help: "Activated stake per validator",
		},
		[]string{"validator_address"},
	)
	ValidatorInflationReward = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "solana_validator_inflation_reward",
			Help:        "Reward earned by validator by the end of each epoch",
			ConstLabels: SOLANA_DENOM_LABEL,
		},
		[]string{"validator_address", "epoch"},
	)
	ValidatorInflationRewardCurrentEpoch = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "solana_validator_inflation_reward_current_epoch",
			Help:        "Reward earned by validator by the end of each epoch",
			ConstLabels: SOLANA_DENOM_LABEL,
		},
		[]string{"validator_address"},
	)

	AvailableBalance = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "solana_available_balance",
			Help:        "Solana available balance",
			ConstLabels: SOLANA_DENOM_LABEL,
		},
		[]string{"delegator_address"},
	)
	DelegatorInflationReward = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "solana_delegator_inflation_reward",
			Help:        "Reward earned by delegator by the end of each epoch",
			ConstLabels: SOLANA_DENOM_LABEL,
		},
		[]string{"delegator_address", "epoch"},
	)
	DelegatorInflationRewardCurrentEpoch = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "solana_delegator_inflation_reward_current_epoch",
			Help:        "Reward earned by delegator by the end of each epoch",
			ConstLabels: SOLANA_DENOM_LABEL,
		},
		[]string{"delegator_address"},
	)

	// represents number of errors while collecting chain stats
	// collector label is used to determine which collector to debug
	ErrorGauge = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "solana_exporter_error_count",
			Help: "Total errors while collecting chain stats",
		},
		[]string{"collector"},
	)
)

func init() {
	prometheus.MustRegister(
		Stakedtotal,
		ValidatorCommission,
		ValidatorDelegatorCount,
		ValidatorEpochCredits,
		ValidatorStaked,
		ValidatorStakedRanking,
		ValidatorInflationRewardCurrentEpoch,
		AvailableBalance,
		DelegatorInflationReward,
		DelegatorInflationRewardCurrentEpoch,
		ErrorGauge,
	)
}
