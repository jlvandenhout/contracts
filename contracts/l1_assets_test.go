package contracts

import (
	"jlvandenhout/contracts/utils"
	"math/big"
	"testing"

	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/iotaledger/wasp/packages/vm/core/evm/iscmagic"
	"github.com/stretchr/testify/require"
)

func TestL1Assets(t *testing.T) {
	// Set up chain
	env := solo.New(t)
	chain := env.NewChain()

	// Set up accounts
	owner := utils.NewEVMAccount(chain, 1_000_000)
	sender := utils.NewEVMAccount(chain, 0)
	receiver := utils.NewL1Account(chain, 0, 0)

	// Deploy contracts
	sandbox, err := Sandbox.OnChain(chain)
	require.NoError(t, err)

	contract, _, err := L1Assets.Deploy(chain, owner, big.NewInt(0))
	require.NoError(t, err)

	// Base tokens
	chain.GetL2FundsFromFaucet(sender.AgentID, 10_000)

	// Native tokens
	foundrySerialNumber, nativeTokenID, err := chain.NewFoundryParams(10).CreateFoundry()
	require.NoError(t, err)

	err = chain.MintTokens(foundrySerialNumber, 10, chain.OriginatorPrivateKey)
	require.NoError(t, err)

	err = chain.SendFromL2ToL2AccountNativeTokens(nativeTokenID, sender.AgentID, 10, chain.OriginatorPrivateKey)
	require.NoError(t, err)

	// NFT
	nft, _, err := env.MintNFTL1(chain.OriginatorPrivateKey, chain.OriginatorAddress, []byte("L1 NFT"))
	require.NoError(t, err)

	err = chain.DepositNFT(nft, sender.AgentID, chain.OriginatorPrivateKey)
	require.NoError(t, err)

	// Construct assets to send to the receiver
	// NOTE: Make sure the base tokens cover the required storage deposit
	assets := chain.L2Assets(sender.AgentID)

	assetsToWithdraw := isc.NewEmptyAssets()
	assetsToWithdraw.BaseTokens = 5000
	assetsToWithdraw.NativeTokens = assets.NativeTokens
	assetsToWithdraw.NFTs = assets.NFTs

	// wrappedL2Assets := iscmagic.WrapISCAssets(chain.L2Assets(sender.AgentID))
	wrappedL2Assets := iscmagic.WrapISCAssets(assetsToWithdraw)
	wrappedL1Address := iscmagic.WrapL1Address(receiver.Address)

	// Allow the contract to take the assets
	_, err = sandbox.Call(sender, "allow", big.NewInt(0), contract.Address, wrappedL2Assets)
	require.NoError(t, err)

	// Transfer the assets
	_, err = contract.Call(sender, "withdraw", big.NewInt(0), wrappedL1Address)
	require.NoError(t, err)

	env.AssertL1BaseTokens(receiver.Address, assetsToWithdraw.BaseTokens)
	for _, nativeToken := range assetsToWithdraw.NativeTokens {
		env.AssertL1NativeTokens(receiver.Address, nativeToken.ID, nativeToken.Amount)
	}
	for _, nft := range assetsToWithdraw.NFTs {
		env.HasL1NFT(receiver.Address, &nft)
	}
}
