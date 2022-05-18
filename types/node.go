package types

type Node struct {
	Address    string `mapstructure:"address"`
	ClientName string `mapstructure:"client_name"`
}

func NewNode(address string, clientName string) *Node {
	return &Node{
		Address:    address,
		ClientName: clientName,
	}
}
