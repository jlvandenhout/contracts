package artifacts

import (
	_ "embed"
)

var (
	//go:embed contracts/L1Assets.sol/L1Assets.json
	L1Assets []byte
)
