package contracts

import (
	"jlvandenhout/contracts/artifacts"
	"jlvandenhout/contracts/utils"

	"github.com/iotaledger/wasp/packages/vm/core/evm/iscmagic"
)

var (
	L1Assets = utils.NewContractFromArtifact(artifacts.L1Assets)
)

var (
	Sandbox = utils.NewCoreContractFromABIAndAddress(iscmagic.SandboxABI, iscmagic.Address)
)
