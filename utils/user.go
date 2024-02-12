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

type EVM struct {
	AgentID *isc.EthereumAddressAgentID
	Address common.Address
	Keys    *ecdsa.PrivateKey
}

type User struct {
	L1
	EVM
}

type BaseTokens struct {
	L1  uint64
	EVM uint64
}

func NewUser(chain *solo.Chain, baseTokens BaseTokens) User {
	l1Keys, l1Address := chain.Env.NewKeyPair()
	l1AgentID := isc.NewAddressAgentID(l1Address)

	if baseTokens.L1 > 0 {
		chain.Env.GetFundsFromFaucet(l1Address, baseTokens.L1)
	}

	evmKeys, evmAddress := chain.NewEthereumAccountWithL2Funds(baseTokens.EVM)
	evmAgentID := isc.NewEthereumAddressAgentID(chain.ChainID, evmAddress)

	user := User{
		L1{
			AgentID: l1AgentID,
			Address: l1Address,
			Keys:    l1Keys,
		},
		EVM{
			AgentID: evmAgentID,
			Address: evmAddress,
			Keys:    evmKeys,
		},
	}

	return user
}
