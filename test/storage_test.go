package test

import (
	"testing"

	"github.com/iotaledger/wasp/packages/solo"
	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {
	env := solo.New(t)
	chain := env.NewChain()
	creator, _ := chain.NewEthereumAccountWithL2Funds()

	// deploy solidity `storage` contract, with 42
	contract := Storage.Deploy(creator, chain, uint32(42))

	// call EVM contract's `retrieve` view, get 42
	result := contract.CallView(nil, "retrieve", nil)
	assert.Equal(t, uint32(42), result[0])

	// call EVM contract's `store`, set 20
	contract.Call(nil, "store", nil, uint32(20))

	// call EVM contract's `retrieve` view, get 20
	result = contract.CallView(nil, "retrieve", nil)
	assert.Equal(t, uint32(20), result[0])
}
