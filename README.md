# go-testchain

[![Tag](https://img.shields.io/github/tag/wealdtech/go-testchain.svg)](https://github.com/wealdtech/go-testchain/releases/)
[![License](https://img.shields.io/github/license/wealdtech/go-testchain.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/wealdtech/go-testchain?status.svg)](https://godoc.org/github.com/wealdtech/go-testchain)
[![Travis CI](https://img.shields.io/travis/wealdtech/go-testchain.svg)](https://travis-ci.org/wealdtech/go-testchain)
[![codecov.io](https://img.shields.io/codecov/c/github/wealdtech/go-testchain.svg)](https://codecov.io/github/wealdtech/go-testchain)

Go module that creates an in-memory Ethereum blockchain to allow for easy testing of Go Ethereum client code.

## Table of Contents

- [Install](#install)
- [Usage](#usage)
- [Maintainers](#maintainers)
- [Contribute](#contribute)
- [License](#license)

## Install

`go-testchain` is a standard Go module which can be installed with:

```sh
go get github.com/wealdtech/go-testchain
```

## Usage

`go-testchain` provides functions that create and populate an in-memory Ethereum blockchain that can be used to test your code without requiring access to an external blockchain.

The blockchain is created with the `NewTestChain()` function, which returns a structure that contains the following fields:

  - `Chain` the blockchain itself.  This contains a subset of a full chain functionality; details in [the documentation]( https://godoc.org/github.com/ethereum/go-ethereum/accounts/abi/bind/backends).  This implements `bind.ContractBackend` and can be used in most places that an `ethclient.Client` would be, for example when instantiating an `abigen`erated contract
  - `Accounts` an array of 128 addresses.  These are each populated with 100,000,000 Ether
  - `PrivateKeys` a map from account address to private key for the accounts in `Accounts`
  - `TransactOpts` a map from account address to transaction options for the accounts in `Accounts`.  The transaction options are commonly used when sending contract transactions

### Example

```go
package main
  
import (
    "context"
    "fmt"
    "math/big"
    "os"
    "testing"

    ethereum "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/common"
    string2eth "github.com/wealdtech/go-string2eth"
    testchain "github.com/wealdtech/go-testchain"
    erc777 "github.com/wealdtech/go-token/contracts/erc777"
)

func TestMain(m *testing.M) {
    setup()
    os.Exit(m.Run())
}

var tc *testchain.TestChain
var erc777Addr *common.Address

func setup() {
    tc = testchain.NewTestChain()
    // Deploy the ERC-1820 registry
    testchain.DeployERC1820(tc)

    // Deploy an ERC-777 token, created by tc.Accounts[127]
    var err error
    erc777Addr, err = testchain.DeployERC777(tc, big.NewInt(1), tc.Accounts[127])
    if err != nil {
        panic(err)
    }
}

func TestBalance(t *testing.T) {
    // Obtain an account balance
    result, err := tc.Chain.BalanceAt(context.Background(), tc.Accounts[0], nil)
    if err != nil {
        t.Errorf("Error: %v", err)
    }
    fmt.Printf("Balance of %v is %v\n", tc.Accounts[0].Hex(), string2eth.WeiToString(result, true))
}

func TestTotalSupply(t *testing.T) {
    // Call the totalSupply() method of the ERC-777 token by hand
    data := []byte{0x18, 0x16, 0x0d, 0xdd}
    msg := ethereum.CallMsg{
        From: tc.Accounts[0],
        To:   erc777Addr,
        Data: data,
    }
    result, err := tc.Chain.CallContract(context.Background(), msg, nil)
    if err != nil {
        t.Errorf("Error: %v", err)
    }
    fmt.Printf("Total supply is %v\n", new(big.Int).SetBytes(result))
}

func TestTotalSupplyContract(t *testing.T) {
    // NewContract is part of a Go file generated by `abigen` with an ERC-777 ABI
    token, err := erc777.NewContract(*erc777Addr, tc.Chain)
    if err != nil {
        t.Errorf("Error: %v", err)
    }
    totalSupply, err := token.TotalSupply(nil)
    if err != nil {
        t.Errorf("Error: %v", err)
    }
    fmt.Printf("Total supply is %v\n", totalSupply)
}

func TestSend(t *testing.T) {
    // NewContract is part of a Go file generated by `abigen` with an ERC-777 ABI
    token, err := erc777.NewContract(*erc777Addr, tc.Chain)
    if err != nil {
        t.Errorf("Error: %v", err)
    }

    oldSenderBalance, err := token.BalanceOf(nil, tc.Accounts[127])
    if err != nil {
        t.Errorf("Error: %v", err)
    }

    oldRecipientBalance, err := token.BalanceOf(nil, tc.Accounts[0])
    if err != nil {
        t.Errorf("Error: %v", err)
    }

    sendAmount := big.NewInt(1000000000)
    _, err = token.Send(tc.TransactOpts[tc.Accounts[127]], tc.Accounts[0], sendAmount, nil)
    if err != nil {
        t.Errorf("Error: %v", err)
    }
    // Need to tell the test chain to commit any pending transactions
    tc.Chain.Commit()

    newSenderBalance, err := token.BalanceOf(nil, tc.Accounts[127])
    if err != nil {
        t.Errorf("Error: %v", err)
    }
    senderBalanceChange := new(big.Int).Sub(oldSenderBalance, newSenderBalance)
    if senderBalanceChange.Cmp(sendAmount) != 0 {
        t.Errorf("Sender balance change incorrect")
    }

    newRecipientBalance, err := token.BalanceOf(nil, tc.Accounts[0])
    if err != nil {
        t.Errorf("Error: %v", err)
    }
    recipientBalanceChange := new(big.Int).Sub(newRecipientBalance, oldRecipientBalance)
    if recipientBalanceChange.Cmp(sendAmount) != 0 {
        t.Errorf("Recipient balance change incorrect")
    }
}
```

## Maintainers

Jim McDonald: [@mcdee](https://github.com/mcdee).

## Contribute

Contributions welcome. Please check out [the issues](https://github.com/wealdtech/go-testchain/issues).

## License

[Apache-2.0](LICENSE) © 2019 Weald Technology Trading Ltd