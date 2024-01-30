package test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/stretchr/testify/require"
)

func TestL1Funds(t *testing.T) {
	env := solo.New(t)
	chain := env.NewChain()
	contract := L1Funds.Deploy(chain, nil, big.NewInt(0))
	L1Keys, L1Address := env.NewKeyPairWithFunds()
	L2Keys, L2Address := solo.NewEthereumAccount()

	metadata := []byte("metadata")
	nft, _, err := env.MintNFTL1(L1Keys, L1Address, metadata)
	require.NoError(t, err)

	agentId := isc.NewEthereumAddressAgentID(chain.ChainID, L2Address)

	assets := isc.NewAssets(1000, nil, nft.ID)
	chain.SendFromL1ToL2Account(2600, assets, agentId, L1Keys)

	receipt, err := contract.Call(L2Keys, "deposit", big.NewInt(0), assets)
	require.NoError(t, err)

	result, err := contract.EventFromReceipt("BaseTokenEvent", receipt)
	require.NoError(t, err)

	fmt.Println(result)

	result, err = contract.EventFromReceipt("NFTEvent", receipt)
	require.NoError(t, err)

	fmt.Println(result)
}
