package contracts

import (
	"jlvandenhout/contracts/utils"
	"math/big"
	"testing"

	"github.com/iotaledger/wasp/packages/solo"
	"github.com/iotaledger/wasp/packages/vm/core/evm/iscmagic"
	"github.com/stretchr/testify/require"
)

func TestL2Balance(t *testing.T) {
	// Set up chain
	env := solo.New(t)
	chain := env.NewChain()

	// Set up accounts
	owner := utils.NewEVMAccount(chain, 1_000_000)
	user := utils.NewL1Account(chain, 10_000_000, 0)

	// Deploy contracts
	sandbox, err := Accounts.OnChain(chain)
	require.NoError(t, err)

	contract, _, err := L2Balance.Deploy(chain, owner, big.NewInt(0))
	require.NoError(t, err)

	// Native tokens
	foundrySerialNumber, nativeTokenID, err := chain.NewFoundryParams(10).CreateFoundry()
	require.NoError(t, err)

	err = chain.MintTokens(foundrySerialNumber, 10, chain.OriginatorPrivateKey)
	require.NoError(t, err)

	err = chain.SendFromL2ToL2AccountNativeTokens(nativeTokenID, user.AgentID, 10, chain.OriginatorPrivateKey)
	require.NoError(t, err)

	// Wrap
	wrappedNativeTokenID := iscmagic.WrapNativeTokenID(nativeTokenID)
	wrappedAgentID := iscmagic.WrapISCAgentID(user.AgentID)

	// Test
	a, err := sandbox.CallView(owner, "getL2BalanceNativeTokens", big.NewInt(0), wrappedNativeTokenID, wrappedAgentID)
	require.NoError(t, err)

	b, err := contract.CallView(owner, "getBalance", big.NewInt(0), wrappedNativeTokenID, wrappedAgentID)
	require.NoError(t, err)

	balance := chain.L2NativeTokens(user.AgentID, nativeTokenID)

	require.Equal(t, new(big.Int).SetBytes(a), balance)
	require.Equal(t, new(big.Int).SetBytes(b), balance)
}
