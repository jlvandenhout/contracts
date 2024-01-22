package test

import (
	"testing"

	"github.com/iotaledger/wasp/packages/solo"
)

func TestEntropy(t *testing.T) {
	env := solo.New(t)
	chain := env.NewChain()
	creator, _ := chain.NewEthereumAccountWithL2Funds()

	// deploy solidity `entropy` contract
	contract := Entropy.Deploy(creator, chain)

	// call EVM contract's `emitEntropy`
	contract.Call(nil, "emitEntropy", nil)
}
