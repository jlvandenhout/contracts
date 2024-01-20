package test

import (
	"math/big"
	"testing"

	"github.com/iotaledger/wasp/packages/solo"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	env := solo.New(t)
	chain := env.NewChain()
	creator, _ := chain.NewEthereumAccountWithL2Funds()

	// deploy solidity `storage` contract, with 42
	chain.DeployEVMContract(creator, Storage.Abi, Storage.Bytecode, big.NewInt(0), uint32(42))

	// call EVM contract's `retrieve` view, get 42
	res, err := chain.CallView("Storage", "retrieve")
	require.NoError(t, err)
	t.Log(res)
}
