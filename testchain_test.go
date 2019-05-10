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
	"bytes"
	"context"
	"math/big"
	"testing"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTestChain(t *testing.T) {
	tc := NewTestChain()
	balance, err := tc.Chain.BalanceAt(context.Background(), tc.Accounts[0], nil)
	require.Nil(t, err, "failed to obtain balance")
	expected, _ := new(big.Int).SetString("100000000000000000000000000", 10)
	assert.Equal(t, expected, balance)
}

func TestDeployContracT(t *testing.T) {
	tc := NewTestChain()

	data := common.FromHex("608060405234801561001057600080fd5b50610108806100206000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c806319ff1d2114602d575b600080fd5b603360a5565b6040805160208082528351818301528351919283929083019185019080838360005b83811015606b5781810151838201526020016055565b50505050905090810190601f16801560975780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b60408051808201909152600d81527f48656c6c6f2c20776f726c64210000000000000000000000000000000000000060208201529056fea165627a7a72305820025f93d88ad4ed442aaaed9b4041425029b6a951425bb127494ab09293568c5c0029")
	contractAddress, err := tc.DeployContract(tc.Accounts[0], big.NewInt(0), data)
	require.Nil(t, err, "failed to deploy contract")

	data = []byte{0x19, 0xff, 0x1d, 0x21}
	msg := ethereum.CallMsg{
		From: tc.Accounts[0],
		To:   contractAddress,
		Data: data,
	}
	result, err := tc.Chain.CallContract(context.Background(), msg, nil)
	require.Nil(t, err, "failed to obtain decimals")
	expected := common.FromHex("0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000d48656c6c6f2c20776f726c642100000000000000000000000000000000000000")
	assert.True(t, bytes.Compare(result, expected) == 0)
}
