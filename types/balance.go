package types

type BalanceResponse struct {
	Context Context `json:"context"`
	Value   uint64
}
