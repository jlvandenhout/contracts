package utils

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/solo"
)

type User struct {
	Address common.Address
	AgentID *isc.EthereumAddressAgentID
	Keys    *ecdsa.PrivateKey
}

func NewUser(chain *solo.Chain, baseTokens uint64) User {
	keys, address := chain.NewEthereumAccountWithL2Funds(baseTokens)
	agentID := isc.NewEthereumAddressAgentID(chain.ChainID, address)

	user := User{
		Keys:    keys,
		Address: address,
		AgentID: agentID,
	}

	return user
}
