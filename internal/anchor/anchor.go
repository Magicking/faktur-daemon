// Copyright 2018 SixUnDeuxZero
//
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

package anchor

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	cmn "github.com/Magicking/faktur-daemon/common"
)

type Anchor struct {
	key       *ecdsa.PrivateKey
	from      common.Address
	lastNonce uint64
	chainId   *big.Int
}

func NewAnchor(_key *ecdsa.PrivateKey, _lastNonce uint64, _chainId *big.Int) *Anchor {
	_from := crypto.PubkeyToAddress(_key.PublicKey)
	a := Anchor{
		key:       _key,
		from:      _from,
		lastNonce: _lastNonce,
		chainId:   _chainId,
	}
	return &a
}

func (a *Anchor) SendWithValueMessage(ctx context.Context, to common.Address, value *big.Int, data []byte) (common.Hash, error) {
	nc := cmn.ClientFromContext(ctx)
	auth := bind.NewKeyedTransactor(a.key)

	if value == nil {
		value = new(big.Int)
	}
	gasPrice, err := nc.SuggestGasPrice(ctx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to suggest gas price: %v", err)
	}
	_to := to
	gasLimit, err := nc.EstimateGas(ctx, ethereum.CallMsg{
		From:     a.from,
		To:       &_to,
		Gas:      0,
		GasPrice: gasPrice,
		Value:    value,
		Data:     data,
	})
	if err != nil {
		return common.Hash{}, fmt.Errorf("Could not estimate gas: %v", err)
	}
	rawTx := types.NewTransaction(a.lastNonce, to, value, gasLimit, gasPrice, data)
	signedTx, err := auth.Signer(types.NewEIP155Signer(a.chainId), a.from, rawTx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("SendWithValueMessage: %v", err)
	}
	err = nc.SendTransaction(ctx, signedTx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("SendTransaction: %v", err)
	}
	a.lastNonce++
	return signedTx.Hash(), nil
}

func (a *Anchor) updateNonce(ctx context.Context) error {
	nc := cmn.ClientFromContext(ctx)
	nonce, err := nc.NonceAt(ctx, a.from, nil)
	if err != nil {
		return err
	}
	if a.lastNonce >= nonce {
		return nil
	}
	log.Println("Nonce updated to", nonce)
	a.lastNonce = nonce
	return nil
}

func (a *Anchor) Run(ctx context.Context, contractAddress common.Address, c chan common.Hash) {
	if a.lastNonce == 0 {
		if err := a.updateNonce(ctx); err != nil {
			log.Fatalf("Could not obtain fresh nonce: %v", err)
		}
	}

	go func() {
		ticker := time.NewTicker(time.Duration(10 * time.Second))
		for {
			select {
			case <-ticker.C:
				// Get RETRY
				roots := []common.Hash{}
				for _, root := range roots {
					txHash, err := a.SendWithValueMessage(ctx, contractAddress, new(big.Int), root.Bytes())
					if err != nil {
						// TODO Save merkleroot to database with state RETRY
						log.Printf("Could not sent transaction for hash %v: %v", root.Hex(), err)
						continue
					}
					log.Printf("Transaction sent: %v", txHash.Hex())
					// TODO Save merkleroot to database with state WAITING CONFIRMATION
				}
				// GET WAITING CONFIRMATION
				// Check if timeout too old
				// Set to retry if necessary
				a.updateNonce(ctx)
			}
			//SEND
		}
	}()
	for e := range c {
		// TODO Save merkleroot to database with state NOT SENT
		txHash, err := a.SendWithValueMessage(ctx, contractAddress, new(big.Int), e.Bytes())
		if err != nil {
			// TODO Save merkleroot to database with state RETRY
			log.Printf("Could not sent transaction for hash %v: %v", e.Hex(), err)
			continue
		}
		log.Printf("Transaction sent: %v", txHash.Hex())
		// TODO Save merkleroot to database with state WAITING CONFIRMATION
	}
}
