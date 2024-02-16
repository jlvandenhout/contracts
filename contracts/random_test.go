package contracts

import (
	"jlvandenhout/contracts/utils"
	"math/big"
	"testing"

	"github.com/iotaledger/wasp/packages/solo"
	"github.com/stretchr/testify/require"
)

func Setup(t *testing.T) (*utils.ContractInstance, utils.EVMAccount) {
	// Set up chain
	env := solo.New(t)
	chain := env.NewChain()

	// Set up accounts
	user := utils.NewEVMAccount(chain, 1_000_000)

	// Deploy contracts
	contract, _, err := Random.Deploy(chain, user, big.NewInt(0))
	require.NoError(t, err)

	return contract, user
}

func TestRandomByte(t *testing.T) {
	contract, user := Setup(t)

	receipt, err := contract.Call(user, "getByte", big.NewInt(0))
	require.NoError(t, err)

	var value byte
	err = contract.EventFromReceipt("Byte", receipt, &value)
	require.NoError(t, err)

	require.NotZero(t, value)
}

func TestRandomUInt64(t *testing.T) {
	contract, user := Setup(t)

	receipt, err := contract.Call(user, "getUInt64", big.NewInt(0))
	require.NoError(t, err)

	var value uint64
	err = contract.EventFromReceipt("UInt64", receipt, &value)
	require.NoError(t, err)

	require.NotZero(t, value)
}

func TestRandomBytesNotEqual(t *testing.T) {
	contract, user := Setup(t)

	receipt, err := contract.Call(user, "getBytes", big.NewInt(0), big.NewInt(100))
	require.NoError(t, err)

	var a []byte
	err = contract.EventFromReceipt("Bytes", receipt, &a)
	require.NoError(t, err)

	receipt, err = contract.Call(user, "getBytes", big.NewInt(0), big.NewInt(100))
	require.NoError(t, err)

	var b []byte
	err = contract.EventFromReceipt("Bytes", receipt, &b)
	require.NoError(t, err)

	require.NotEqual(t, a, b)
}

func TestRandomBytesLength(t *testing.T) {
	contract, user := Setup(t)

	lengths := []int64{0, 1, 32, 64, 100, 500}

	for _, length := range lengths {
		receipt, err := contract.Call(user, "getBytes", big.NewInt(0), big.NewInt(length))
		require.NoError(t, err)

		var value []byte
		err = contract.EventFromReceipt("Bytes", receipt, &value)
		require.NoError(t, err)

		require.Len(t, value, int(length))
	}
}
