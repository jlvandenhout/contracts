package utils

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/solo"
)

type L1 struct {
	AgentID *isc.AddressAgentID
	Address iotago.Address
	Keys    *cryptolib.KeyPair
}

type L2 struct {
	AgentID *isc.EthereumAddressAgentID
	Address common.Address
	Keys    *ecdsa.PrivateKey
}

type User struct {
	L1
	L2
}

type BaseTokens struct {
	L1 uint64
	L2 uint64
}

func NewUser(chain *solo.Chain, baseTokens BaseTokens) User {
	l1Keys, l1Address := chain.Env.NewKeyPair()
	l1AgentID := isc.NewAddressAgentID(l1Address)

	if baseTokens.L1 > 0 {
		chain.Env.GetFundsFromFaucet(l1Address, baseTokens.L1)
	}

	l2Keys, l2Address := chain.NewEthereumAccountWithL2Funds(baseTokens.L2)
	l2AgentID := isc.NewEthereumAddressAgentID(chain.ChainID, l2Address)

	user := User{
		L1: L1{
			AgentID: l1AgentID,
			Address: l1Address,
			Keys:    l1Keys,
		},
		L2: L2{
			AgentID: l2AgentID,
			Address: l2Address,
			Keys:    l2Keys,
		},
	}

	return user
}
