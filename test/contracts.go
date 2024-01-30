package test

import "jlvandenhout/contracts/artifacts"

var (
	Storage = NewContractFromArtifact(artifacts.Storage)
	Entropy = NewContractFromArtifact(artifacts.Entropy)
	L1Funds = NewContractFromArtifact(artifacts.L1Funds)
)
