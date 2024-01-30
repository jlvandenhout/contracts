package test

import (
	"math/big"
	"testing"

	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestL1Funds(t *testing.T) {
	env := solo.New(t)
	chain := env.NewChain()
	contract := L1Funds.Deploy(chain, nil, big.NewInt(0))
	keys, address := env.NewKeyPairWithFunds()

	metadata := []byte("metadata")
	nft, _ := env.MintNFTL1(keys, address, metadata)
	assets := isc.NewAssets(1000, nil, nft)
	chain.SendFromL1ToL2Account(1000, http.ResponseWriter, r *http.Request
}
