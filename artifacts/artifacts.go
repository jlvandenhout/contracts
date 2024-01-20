package artifacts

import (
	_ "embed"
)

var (
	//go:embed contracts/Storage.sol/Storage.json
	Storage []byte
)
