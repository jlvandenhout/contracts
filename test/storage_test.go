package test

import (
	"math/big"
	"testing"

	"github.com/iotaledger/wasp/packages/solo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	env := solo.New(t)
	chain := env.NewChain()
	contract := Storage.Deploy(chain, nil, big.NewInt(0), uint32(42))

	result, err := contract.CallView(nil, "retrieve", big.NewInt(0))
	require.NoError(t, err)
	assert.Equal(t, result[0], uint32(42))

	receipt, err := contract.Call(nil, "store", big.NewInt(0), uint32(20))
	require.NoError(t, err)

	result, err = contract.EventFromReceipt("Stored", receipt)
	require.NoError(t, err)
	assert.Equal(t, result[0], uint32(20))

	result, err = contract.CallView(nil, "retrieve", big.NewInt(0))
	require.NoError(t, err)
	assert.Equal(t, result[0], uint32(20))
}
