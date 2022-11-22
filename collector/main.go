package collector

import "github.com/forbole/solana-exporter/types"

type SolanaCollector struct {
	SolanaClient       *types.Client
	DelegatorAddresses []string
	ValidatorAddresses map[string]struct{}
}

func NewSolanaCollector(solanaClient *types.Client, delegatorAddresses []string, validatorAddresses map[string]struct{}) SolanaCollector {
	return SolanaCollector{
		SolanaClient:       solanaClient,
		DelegatorAddresses: delegatorAddresses,
		ValidatorAddresses: validatorAddresses,
	}
}

func (c *SolanaCollector) Collect() {
	c.CollectValidatorStats()
	c.CollectDelegatorAddressStats()
}
