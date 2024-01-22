package test

import (
	"math/big"
	"testing"

	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntropy(t *testing.T) {
	env := solo.New(t)
	chain := env.NewChain()
	contract := Entropy.Deploy(chain, nil, big.NewInt(0))

	receipt, err := contract.Call(nil, "emitEntropy", nil)
	require.NoError(t, err)

	result, err := contract.EventFromReceipt("EntropyEvent", receipt)
	require.NoError(t, err)
	assert.NotEqualValues(t, result[0], hashing.NilHash)
}
