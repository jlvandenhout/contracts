package test

import "jlvandenhout/contracts/artifacts"

var (
	Storage = NewContractFromArtifact(artifacts.Storage)
	Entropy = NewContractFromArtifact(artifacts.Entropy)
)
