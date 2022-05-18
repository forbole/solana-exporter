package types

type InflationReward struct {
	Amount        uint64 `json:"amount"`
	Commission    uint8  `json:"commission"`
	EffectiveSlot uint64 `json:"effectiveSlot"`
	Epoch         uint64 `json:"epoch"`
	PostBalance   uint64 `json:"postBalance"`
}

type GetInflationRewardConfig struct {
	Commitment string `json:"commitment,omitempty"`
	Epoch      uint64 `json:"epoch,omitempty"`
}
