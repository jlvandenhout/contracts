package utils

import (
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
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/solo"
)

type Contract struct {
	ABI      abi.ABI
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
		ABI      abi.ABI
		Bytecode string
	}

	err := json.Unmarshal(data, &fields)
	if err != nil {
		return err
	}

	contract.ABI = fields.ABI
	contract.Bytecode = common.FromHex(strings.TrimSpace(fields.Bytecode))
	return nil
}

func (contract *Contract) Deploy(chain *solo.Chain, owner EVMAccount, value *big.Int, args ...interface{}) (*ContractInstance, *types.Receipt, error) {
	nonce := chain.Nonce(owner.AgentID)

	bytecode := contract.Bytecode
	arguments, err := contract.ABI.Pack("", args...)
	if err != nil {
		return nil, nil, err
	}

	data := make([]byte, 0, len(bytecode)+len(arguments))
	data = append(data, bytecode...)
	data = append(data, arguments...)

	callMessage := ethereum.CallMsg{
		From:  owner.Address,
		To:    nil, // contract creation
		Value: value,
		Data:  data,
	}

	gas, err := chain.EVM().EstimateGas(callMessage, nil)
	if err != nil {
		return nil, nil, err
	}

	signer, err := chain.EVM().Signer()
	if err != nil {
		return nil, nil, err
	}

	transaction, err := types.SignTx(
		types.NewContractCreation(
			nonce,
			callMessage.Value,
			gas,
			chain.EVM().GasPrice(),
			data,
		),
		signer,
		owner.Keys,
	)
	if err != nil {
		return nil, nil, err
	}

	err = chain.EVM().SendTransaction(transaction)
	if err != nil {
		return nil, nil, err
	}

	receipt := chain.EVM().TransactionReceipt(transaction.Hash())

	address := crypto.CreateAddress(owner.Address, nonce)

	agentID := isc.NewEthereumAddressAgentID(chain.ID(), address)
	return &ContractInstance{ABI: contract.ABI, Address: address, AgentID: agentID, Chain: chain}, receipt, nil
}

type CoreContract struct {
	ABI     abi.ABI
	Address common.Address
}

func NewCoreContractFromABIAndAddress(marshalledABI string, address common.Address) CoreContract {
	unmarshalledABI, err := abi.JSON(strings.NewReader(marshalledABI))
	if err != nil {
		log.Fatal(err)
	}

	return CoreContract{ABI: unmarshalledABI, Address: address}
}

func (coreContract *CoreContract) OnChain(chain *solo.Chain) (*ContractInstance, error) {
	agentID := isc.NewEthereumAddressAgentID(chain.ChainID, coreContract.Address)

	return &ContractInstance{ABI: coreContract.ABI, Address: coreContract.Address, AgentID: agentID, Chain: chain}, nil
}

type ContractInstance struct {
	ABI     abi.ABI
	Address common.Address
	AgentID isc.AgentID
	Chain   *solo.Chain
}

func (instance *ContractInstance) CallView(caller EVMAccount, function string, value *big.Int, args ...interface{}) ([]byte, error) {
	data, err := instance.ABI.Pack(function, args...)
	if err != nil {
		return nil, err
	}

	callMsg := ethereum.CallMsg{
		From:  caller.Address,
		To:    &instance.Address,
		Data:  data,
		Value: value,
	}

	ret, err := instance.Chain.EVM().CallContract(callMsg, nil)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (instance *ContractInstance) EstimateGas(caller EVMAccount, function string, value *big.Int, args ...interface{}) (uint64, error) {
	data, err := instance.ABI.Pack(function, args...)
	if err != nil {
		return 0, err
	}

	callMsg := ethereum.CallMsg{
		From:     caller.Address,
		To:       &instance.Address,
		Data:     data,
		Value:    value,
		GasPrice: instance.Chain.EVM().GasPrice(),
	}

	gas, err := instance.Chain.EVM().EstimateGas(callMsg, nil)
	if err != nil {
		return 0, err
	}

	return gas, nil
}

func (instance *ContractInstance) Call(caller EVMAccount, function string, value *big.Int, args ...interface{}) (*types.Receipt, error) {
	data, err := instance.ABI.Pack(function, args...)
	if err != nil {
		return nil, err
	}

	callMsg := ethereum.CallMsg{
		From:     caller.Address,
		To:       &instance.Address,
		Data:     data,
		Value:    value,
		GasPrice: instance.Chain.EVM().GasPrice(),
	}

	gas, err := instance.Chain.EVM().EstimateGas(callMsg, nil)
	if err != nil {
		return nil, err
	}

	nonce, err := instance.Chain.EVM().TransactionCount(callMsg.From, nil)
	if err != nil {
		return nil, err
	}

	signer, err := instance.Chain.EVM().Signer()
	if err != nil {
		return nil, err
	}

	transaction, err := types.SignNewTx(caller.Keys, signer, &types.LegacyTx{
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

func (instance *ContractInstance) EventFromReceipt(event string, receipt *types.Receipt, v interface{}) error {
	topic := instance.ABI.Events[event].ID

	for _, log := range receipt.Logs {
		if slices.Contains(log.Topics, topic) {
			return instance.ABI.UnpackIntoInterface(v, event, log.Data)
		}
	}

	return fmt.Errorf("event with name %v not found", event)
}
