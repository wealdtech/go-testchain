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
	"math/big"
	"testing"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestERC777(t *testing.T) {
	tc := NewTestChain()
	err := DeployERC1820(tc)
	require.Nil(t, err, "failed to deploy ERC-1820 contract")
	contractAddress, err := DeployERC777(tc, big.NewInt(1000000000), tc.Accounts[127])
	require.Nil(t, err, "failed to deploy ERC-777 contract")

	// Fetch the total supply to ensure this really has been deployed
	data := []byte{0x18, 0x16, 0x0d, 0xdd}
	msg := ethereum.CallMsg{
		From: tc.Accounts[0],
		To:   contractAddress,
		Data: data,
	}
	result, err := tc.Chain.CallContract(context.Background(), msg, nil)
	require.Nil(t, err, "failed to obtain total supply")
	expected, _ := new(big.Int).SetString("100000000000000000000000000", 10)
	totalSupply := new(big.Int).SetBytes(result)
	assert.Equal(t, expected, totalSupply)

	// Fetch the granularity to ensure the input value of 1 has been set
	data = []byte{0x55, 0x6f, 0x0d, 0xc7}
	msg = ethereum.CallMsg{
		From: tc.Accounts[0],
		To:   contractAddress,
		Data: data,
	}
	result, err = tc.Chain.CallContract(context.Background(), msg, nil)
	expected = _bigInt("1000000000")
	granularity := new(big.Int).SetBytes(result)
	require.Nil(t, err, "failed to obtain granularity")
	assert.Equal(t, expected, granularity)
}
