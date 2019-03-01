// Copyright (c) 2016, Alan Chen
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
//    this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors
//    may be used to endorse or promote products derived from this software
//    without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package common

import (
	"encoding/json"
	"math/big"
	"sync/atomic"
	"github.com/alanchchen/web3go/rlp"
	"io"
)

const (
	hashLength    = 32
	addressLength = 20
)

// Hash ...
type Hash [hashLength]byte

func NewHash(data []byte) (result Hash) {
	copy(result[:], data)
	return result
}

func (hash *Hash) String() string {
	return BytesToHex(hash[:])
}

// Address ...
type Address [addressLength]byte

func NewAddress(data []byte) (result Address) {
	copy(result[:], data)
	return result
}

func (addr *Address) String() string {
	return BytesToHex(addr[:])
}

// SyncStatus ...
type SyncStatus struct {
	Result        bool
	StartingBlock *big.Int
	CurrentBlock  *big.Int
	HighestBlock  *big.Int
}

// TransactionRequest ...
type TransactionRequest struct {
	From     Address  `json:"from"`
	To       Address  `json:"to"`
	Gas      *big.Int `json:"gas"`
	GasPrice *big.Int `json:"gasprice"`
	Value    *big.Int `json:"value"`
	Data     []byte   `json:"data"`
}

func (tx *TransactionRequest) String() string {
	jsonBytes, _ := json.Marshal(tx)
	return string(jsonBytes)
}

// Transaction ...
type Transaction struct {
	Hash             Hash     `json:"hash"`
	Nonce            Hash     `json:"nonce"`
	BlockHash        Hash     `json:"blockHash"`
	BlockNumber      *big.Int `json:"blockNumber"`
	TransactionIndex uint64   `json:"transactionIndex"`
	From             Address  `json:"from"`
	To               Address  `json:"to"`
	Gas              *big.Int `json:"gas"`
	GasPrice         *big.Int `json:"gasprice"`
	Value            *big.Int `json:"value"`
	Data             []byte   `json:"input"`
}

func (tx *Transaction) String() string {
	jsonBytes, _ := json.Marshal(tx)
	return string(jsonBytes)
}

type Topic struct {
	Data []byte
}

type Topics []Topic

// Log ...
type Log struct {
	LogIndex         uint64   `json:"logIndex"`
	BlockNumber      *big.Int `json:"blockNumber"`
	BlockHash        Hash     `json:"blockHash"`
	TransactionHash  Hash     `json:"transactionHash"`
	TransactionIndex uint64   `json:"transactionIndex"`
	Address          Address  `json:"address"`
	Data             []byte   `json:"data"`
	Topics           Topics   `json:"topics"`
}

// TransactionReceipt ...
type TransactionReceipt struct {
	Hash              Hash     `json:"transactionHash"`
	TransactionIndex  uint64   `json:"transactionIndex"`
	BlockNumber       *big.Int `json:"blockNumber"`
	BlockHash         Hash     `json:"blockHash"`
	CumulativeGasUsed *big.Int `json:"cumulativeGasUsed"`
	GasUsed           *big.Int `json:"gasUsed"`
	ContractAddress   Address  `json:"contractAddress"`
	Logs              []Log    `json:"logs"`
}

func (tx *TransactionReceipt) String() string {
	jsonBytes, _ := json.Marshal(tx)
	return string(jsonBytes)
}

// Block ...
type Block struct {
	Number          *big.Int `json:"number"`
	Hash            Hash     `json:"hash"`
	ParentHash      Hash     `json:"parentHash"`
	Nonce           Hash     `json:"nonce"`
	Sha3Uncles      Hash     `json:"sha3Uncles"`
	Bloom           Hash     `json:"logsBloom"`
	TransactionRoot Hash     `json:"transactionsRoot"`
	StateRoot       Hash     `json:"stateRoot"`
	Miner           Address  `json:"miner"`
	Difficulty      *big.Int `json:"difficulty"`
	TotalDifficulty *big.Int `json:"totalDifficulty"`
	ExtraData       Hash     `json:"extraData"`
	Size            *big.Int `json:"size"`
	GasLimit        *big.Int `json:"gasLimit"`
	GasUsed         *big.Int `json:"gasUsed"`
	Timestamp       *big.Int `json:"timestamp"`
	Transactions    []Hash   `json:"transactions"`
	Uncles          []Hash   `json:"uncles"`
	//MinGasPrice     *big.Int `json:"minGasPrice"`
}

type Transactions struct {
	data txdata
	// caches
	hash atomic.Value
	size atomic.Value
	from atomic.Value
}

type txdata struct {
	AccountNonce uint64          `json:"nonce"    gencodec:"required"`
	Price        *big.Int        `json:"gasPrice" gencodec:"required"`
	GasLimit     uint64          `json:"gas"      gencodec:"required"`
	Recipient    *Address `json:"to"       rlp:"nil"` // nil means contract creation
	Amount       *big.Int        `json:"value"    gencodec:"required"`
	Payload      []byte          `json:"input"    gencodec:"required"`

	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`

	// This is only used when marshaling to JSON.
	Hash *Hash `json:"hash" rlp:"-"`
}

func NewTransactions(nonce uint64, to Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *Transactions {
	return newTransactions(nonce, &to, amount, gasLimit, gasPrice, data)
}

// CopyBytes returns an exact copy of the provided bytes.
func CopyBytes(b []byte) (copiedBytes []byte) {
	if b == nil {
		return nil
	}
	copiedBytes = make([]byte, len(b))
	copy(copiedBytes, b)

	return
}

func newTransactions(nonce uint64, to *Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *Transactions {
	if len(data) > 0 {
		data = CopyBytes(data)
	}
	d := txdata{
		AccountNonce: nonce,
		Recipient:    to,
		Payload:      data,
		Amount:       new(big.Int),
		GasLimit:     gasLimit,
		Price:        new(big.Int),
		V:            new(big.Int),
		R:            new(big.Int),
		S:            new(big.Int),
	}
	if amount != nil {
		d.Amount.Set(amount)
	}
	if gasPrice != nil {
		d.Price.Set(gasPrice)
	}

	return &Transactions{data: d}
}

// EncodeRLP implements rlp.Encoder
func (tx *Transactions) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &tx.data)
}