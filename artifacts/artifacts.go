package artifacts

import (
	_ "embed"
)

var (
	//go:embed contracts/Storage.sol/Storage.json
	Storage []byte
	//go:embed contracts/Entropy.sol/Entropy.json
	Entropy []byte
)
