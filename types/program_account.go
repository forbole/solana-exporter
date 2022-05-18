package types

type ProgramAccountsResponse struct {
	Pubkey  string         `json:"pubkey"`
	Account ProgramAccount `json:"account"`
}

type GetProgramAccountsConfigEncoding string

const (
	// GetProgramAccountsConfigEncodingBase58 limited to Account data of less than 128 bytes
	GetProgramAccountsConfigEncodingBase58     GetProgramAccountsConfigEncoding = "base58"
	GetProgramAccountsConfigEncodingJsonParsed GetProgramAccountsConfigEncoding = "jsonParsed"
	GetProgramAccountsConfigEncodingBase64     GetProgramAccountsConfigEncoding = "base64"
	GetProgramAccountsConfigEncodingBase64Zstd GetProgramAccountsConfigEncoding = "base64+zstd"
)

type GetProgramAccountsConfigFilter struct {
	Memcmp   *Memcmp `json:"memcmp,omitempty"`
	DataSize uint64  `json:"dataSize,omitempty"`
}

type GetProgramAccountsConfig struct {
	Encoding GetProgramAccountsConfigEncoding `json:"encoding,omitempty"`
	Filters  []GetProgramAccountsConfigFilter `json:"filters,omitempty"`
}

type Memcmp struct {
	Offset uint64 `json:"offset"`
	Bytes  string `json:"bytes"`
}

type ProgramAccount struct {
	Lamports   uint64      `json:"lamports"`
	Owner      string      `json:"owner"`
	RentEpoch  uint64      `json:"rentEpoch"`
	Data       interface{} `json:"data"`
	Executable bool        `json:"executable"`
}
