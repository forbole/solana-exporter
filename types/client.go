package types

import (
	"math"

	soljunoClient "github.com/forbole/soljuno/solana/client"
	jsonrpc "github.com/ybbus/jsonrpc/v2"
)

const STAKE_PROGRAM_ID = "Stake11111111111111111111111111111111111111"
const SOLANA_EXPONENT int = 9

type Client struct {
	soljunoClient.Client
}

func NewSolanaClient(endpoint string) *Client {
	rpcClient := jsonrpc.NewClient(endpoint)
	client := soljunoClient.Client{
		RpcClient: rpcClient,
	}
	return &Client{
		Client: client,
	}
}

func (c *Client) GetBalance(address string) (BalanceResponse, error) {
	var balance BalanceResponse
	err := c.Client.RpcClient.CallFor(&balance, "getBalance", address)
	return balance, err
}

func (c *Client) GetInflationReward(address string) ([]InflationReward, error) {
	var reward []InflationReward
	err := c.Client.RpcClient.CallFor(&reward, "getInflationReward", []string{address}, nil)
	return reward, err
}

func (c *Client) GetProgramAccount(programId string, config GetProgramAccountsConfig) ([]ProgramAccountsResponse, error) {
	var programAccount []ProgramAccountsResponse
	err := c.Client.RpcClient.CallFor(&programAccount, "getProgramAccounts", programId, config)
	return programAccount, err
}

// https://stackoverflow.com/questions/70163352/am-i-able-to-get-a-list-of-delegators-by-validator-solana-using-the-json-rpc-a
func (c *Client) GetDelegatorsCount(validatorAddress string) ([]ProgramAccountsResponse, error) {
	validatorFilter := []GetProgramAccountsConfigFilter{
		{
			Memcmp: &Memcmp{
				Offset: 124,
				Bytes:  validatorAddress,
			},
		},
	}
	encoding := GetProgramAccountsConfigEncodingBase64

	getProgramAccountsConfig := GetProgramAccountsConfig{
		Encoding: encoding,
		Filters:  validatorFilter,
	}

	stakeProgram, err := c.GetProgramAccount(STAKE_PROGRAM_ID, getProgramAccountsConfig)

	return stakeProgram, err
}

func ConvertLamportToSolana(n uint64) float64 {
	return float64(n) / math.Pow10(SOLANA_EXPONENT)
}
