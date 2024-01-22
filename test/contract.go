package test

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"slices"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/iotaledger/wasp/packages/solo"
)

type Contract struct {
	Abi      string
	Bytecode []byte
}

func NewContractFromArtifact(data []byte) Contract {
	var contract Contract

	if err := json.Unmarshal(data, &contract); err != nil {
		log.Fatal(err)
	}

	return contract
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

func (contract *Contract) Deploy(chain *solo.Chain, creator *ecdsa.PrivateKey, value *big.Int, args ...interface{}) ContractInstance {
	if creator == nil {
		creator, _ = chain.NewEthereumAccountWithL2Funds()
	}

	address, abi := chain.DeployEVMContract(creator, contract.Abi, contract.Bytecode, value, args...)
	return ContractInstance{Abi: abi, Address: address, Chain: chain, Creator: creator}
}

type ContractInstance struct {
	Abi     abi.ABI
	Address common.Address
	Chain   *solo.Chain
	Creator *ecdsa.PrivateKey
}

func (instance *ContractInstance) CallView(caller *ecdsa.PrivateKey, function string, value *big.Int, args ...interface{}) ([]interface{}, error) {
	data, err := instance.Abi.Pack(function, args...)
	if err != nil {
		return nil, err
	}

	callMsg := ethereum.CallMsg{
		To:    &instance.Address,
		Data:  data,
		Value: value,
	}

	if caller != nil {
		callMsg.From = crypto.PubkeyToAddress(caller.PublicKey)
	} else {
		callMsg.From = crypto.PubkeyToAddress(instance.Creator.PublicKey)
	}

	ret, err := instance.Chain.EVM().CallContract(callMsg, nil)
	if err != nil {
		return nil, err
	}

	result, err := instance.Abi.Unpack(function, ret)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (instance *ContractInstance) Call(caller *ecdsa.PrivateKey, function string, value *big.Int, args ...interface{}) (*types.Receipt, error) {
	data, err := instance.Abi.Pack(function, args...)
	if err != nil {
		return nil, err
	}

	callMsg := ethereum.CallMsg{
		To:       &instance.Address,
		Data:     data,
		Value:    value,
		GasPrice: instance.Chain.EVM().GasPrice(),
	}

	var signer *ecdsa.PrivateKey
	if caller != nil {
		signer = caller
	} else {
		signer = instance.Creator
	}

	callMsg.From = crypto.PubkeyToAddress(signer.PublicKey)

	gas, err := instance.Chain.EVM().EstimateGas(callMsg, nil)
	if err != nil {
		return nil, err
	}

	nonce, err := instance.Chain.EVM().TransactionCount(callMsg.From, nil)
	if err != nil {
		return nil, err
	}

	transaction, err := types.SignNewTx(signer, types.NewEIP155Signer(big.NewInt(int64(instance.Chain.EVM().ChainID()))), &types.LegacyTx{
		Nonce:    nonce,
		Gas:      gas,
		GasPrice: callMsg.GasPrice,
		To:       callMsg.To,
		Value:    value,
		Data:     callMsg.Data,
	})
	if err != nil {
		return nil, err
	}

	err = instance.Chain.EVM().SendTransaction(transaction)
	if err != nil {
		return nil, err
	}

	receipt := instance.Chain.EVM().TransactionReceipt(transaction.Hash())

	return receipt, nil
}

func (instance *ContractInstance) EventFromReceipt(event string, receipt *types.Receipt) ([]interface{}, error) {
	topic := instance.Abi.Events[event].ID

	for _, log := range receipt.Logs {
		if slices.Contains(log.Topics, topic) {
			return instance.Abi.Unpack(event, log.Data)
		}
	}

	return nil, fmt.Errorf("event with name %v not found", event)
}
