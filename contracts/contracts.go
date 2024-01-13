package contracts

import (
	"jlvandenhout/contracts/utils"

	"github.com/iotaledger/wasp/packages/vm/core/evm/iscmagic"
)

var ()

var (
	Sandbox = utils.CoreContract{Abi: iscmagic.SandboxABI, Address: iscmagic.Address}
)
