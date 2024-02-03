package test

import (
	"math/big"
	"testing"

	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/iotaledger/wasp/packages/vm/core/evm/iscmagic"
	"github.com/stretchr/testify/require"
)

func TestL1Funds(t *testing.T) {
	// Set up
	env := solo.New(t, &solo.InitOptions{AutoAdjustStorageDeposit: true})
	chain := env.NewChain()
	contract := L1Funds.Deploy(chain, nil, big.NewInt(0))

	// Sender account
	senderKeys, senderAddress := chain.NewEthereumAccountWithL2Funds()
	senderAgentID := isc.NewEthereumAddressAgentID(chain.ChainID, senderAddress)

	// Receiver account
	_, receiverAddress := env.NewKeyPairWithFunds()

	// Native tokens
	maxSupply := big.NewInt(10)
	foundrySN, nativeTokenID, err := chain.NewFoundryParams(maxSupply).CreateFoundry()
	require.NoError(t, err)

	err = chain.MintTokens(foundrySN, maxSupply, chain.OriginatorPrivateKey)
	require.NoError(t, err)

	err = chain.SendFromL2ToL2AccountNativeTokens(nativeTokenID, senderAgentID, maxSupply, chain.OriginatorPrivateKey)
	require.NoError(t, err)

	// NFT
	nft, _, err := env.MintNFTL1(chain.OriginatorPrivateKey, chain.OriginatorAddress, []byte("L1 NFT"))
	require.NoError(t, err)

	err = chain.DepositNFT(nft, senderAgentID, chain.OriginatorPrivateKey)
	require.NoError(t, err)

	// Send assets to receiver
	senderAssets := chain.L2Assets(senderAgentID)

	wrappedReceiverAddress := iscmagic.WrapL1Address(receiverAddress)
	wrappedSenderAssets := iscmagic.WrapISCAssets(senderAssets)

	receipt, err := contract.Call(senderKeys, "send", big.NewInt(0), wrappedReceiverAddress, wrappedSenderAssets)
	require.NoError(t, err)

	result, err := contract.EventFromReceipt("Send", receipt)
	require.NoError(t, err)

	t.Log(result)
	t.Log(env.L1BaseTokens(receiverAddress))
	t.Log(env.L1NativeTokens(receiverAddress, nativeTokenID))
	t.Log(env.L1NFTs(receiverAddress))
}
