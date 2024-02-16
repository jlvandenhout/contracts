package utils

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/solo"
)

type L1Account struct {
	AgentID *isc.AddressAgentID
	Address iotago.Address
	Keys    *cryptolib.KeyPair
}

type EVMAccount struct {
	AgentID *isc.EthereumAddressAgentID
	Address common.Address
	Keys    *ecdsa.PrivateKey
}

func NewL1Account(chain *solo.Chain, l1BaseTokens uint64, l2BaseTokens uint64) L1Account {
	keys, address := chain.Env.NewKeyPair()
	agentID := isc.NewAddressAgentID(address)

	if l1BaseTokens > 0 {
		chain.Env.GetFundsFromFaucet(address, l1BaseTokens)
	}

	if l2BaseTokens > 0 {
		chain.GetL2FundsFromFaucet(agentID, l2BaseTokens)
	}

	account := L1Account{
		AgentID: agentID,
		Address: address,
		Keys:    keys,
	}

	return account
}

func NewEVMAccount(chain *solo.Chain, l2BaseTokens uint64) EVMAccount {
	keys, address := chain.NewEthereumAccountWithL2Funds(l2BaseTokens)
	agentID := isc.NewEthereumAddressAgentID(chain.ChainID, address)

	account := EVMAccount{
		AgentID: agentID,
		Address: address,
		Keys:    keys,
	}

	return account
}
