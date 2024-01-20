package test

import (
	_ "embed"
	"encoding/json"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type Contract struct {
	Abi      string
	Bytecode []byte
}

func (contract *Contract) UnmarshalJSON(data []byte) error {
	var fields struct {
		Abi      interface{}
		Bytecode string
	}

	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	if abi, err := json.Marshal(fields.Abi); err != nil {
		return err
	} else {
		contract.Abi = string(abi)
	}

	contract.Bytecode = common.FromHex(strings.TrimSpace(fields.Bytecode))

	return nil
}

func NewContractFromArtifact(data []byte) Contract {
	var contract Contract

	if err := json.Unmarshal(data, &contract); err != nil {
		log.Fatal(err)
	}

	return contract
}
