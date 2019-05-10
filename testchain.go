// Copyright © 2019 Weald Technology Trading
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testchain

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

// TestChain is a test blockchain structure
type TestChain struct {
	Chain        *backends.SimulatedBackend
	Accounts     []common.Address
	PrivateKeys  map[common.Address]*ecdsa.PrivateKey
	TransactOpts map[common.Address]*bind.TransactOpts
}

func _bigInt(input string) *big.Int {
	res, _ := new(big.Int).SetString(input, 10)
	return res
}

// CreateSignedTransaction creates a signed transaction
func (tc *TestChain) CreateSignedTransaction(from common.Address, to *common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	nonce, err := tc.Chain.PendingNonceAt(context.Background(), from)
	if err != nil {
		return nil, errors.New("failed to obtain pending nonce")
	}

	msg := ethereum.CallMsg{From: from, To: to, Value: value, Data: data}
	gasLimit, err := tc.Chain.EstimateGas(context.Background(), msg)
	if err != nil {
		return nil, errors.New("failed to obtain gas limit")
	}

	gasPrice := big.NewInt(1)

	// Create the transaction
	var tx *types.Transaction
	if to == nil {
		tx = types.NewContractCreation(nonce, value, gasLimit, gasPrice, data)
	} else {
		tx = types.NewTransaction(nonce, *to, value, gasLimit, gasPrice, data)
	}

	// Sign the transaction and return it
	return types.SignTx(tx, types.MakeSigner(params.AllEthashProtocolChanges, nil), tc.PrivateKeys[from])
}

// NewTestChain creates a new test blockchain
func NewTestChain() *TestChain {
	alloc := core.GenesisAlloc{}
	accounts := make([]common.Address, 0)
	privateKeys := make(map[common.Address]*ecdsa.PrivateKey)
	transactOpts := make(map[common.Address]*bind.TransactOpts)

	// Create some accounts with funds
	for i := 0; i < 128; i++ {
		privateKey, _ := crypto.GenerateKey()
		opts := bind.NewKeyedTransactor(privateKey)
		privateKeys[opts.From] = privateKey
		transactOpts[opts.From] = opts
		accounts = append(accounts, opts.From)
		alloc[opts.From] = core.GenesisAccount{Balance: _bigInt("100000000000000000000000000")}
	}

	// Create the chain
	chain := backends.NewSimulatedBackend(alloc, 8000000)

	return &TestChain{
		Chain:        chain,
		Accounts:     accounts,
		PrivateKeys:  privateKeys,
		TransactOpts: transactOpts,
	}
}

// DeployContract deploys a contract to the blockchain, returning the address of the contract.
func (tc *TestChain) DeployContract(from common.Address, value *big.Int, data []byte) (*common.Address, error) {
	signedTx, err := tc.CreateSignedTransaction(from, nil, value, data)
	if err != nil {
		return nil, err
	}

	err = tc.Chain.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, err
	}
	tc.Chain.Commit()

	receipt, err := tc.Chain.TransactionReceipt(context.Background(), signedTx.Hash())
	if err != nil {
		return nil, err
	}

	if receipt.Status == 0 {
		return nil, errors.New("contract deployment failed")
	}

	return &receipt.ContractAddress, err
}
