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
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestERC1820(t *testing.T) {
	tc := NewTestChain()
	err := DeployERC1820(tc)
	require.Nil(t, err, "failed to deploy ERC-1820 contract")

	// Ensure the contract is present at the expected address
	erc1820 := common.HexToAddress("1820a4B7618BdE71Dce8cdc73aAB6C95905faD24")
	code, err := tc.Chain.CodeAt(context.Background(), erc1820, nil)
	require.Nil(t, err, "failed to obtain code")
	assert.True(t, len(code) > 0)
}
