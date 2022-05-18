package types

// Config defines all necessary parameters
type Config struct {
	DelegatorAddresses []string `mapstructure:"delegator_addresses"`
	ValidatorAddresses []string `mapstructure:"validator_addresses"`
	Port               string   `mapstructure:"port"`
	Node               Node     `mapstructure:"node"`
}

// NewConfig builds a new Config instance
func NewConfig(
	delegatorAddresses []string, validatorAddresses []string, port string, nodeCfg Node,
) Config {
	return Config{
		DelegatorAddresses: delegatorAddresses,
		ValidatorAddresses: validatorAddresses,
		Port:               port,
		Node:               nodeCfg,
	}
}

func NewValidatorAddressesMap(validatorAddresses []string) map[string]struct{} {
	addressMap := make(map[string]struct{})
	var s struct{}
	for _, address := range validatorAddresses {
		addressMap[address] = s
	}
	return addressMap
}
