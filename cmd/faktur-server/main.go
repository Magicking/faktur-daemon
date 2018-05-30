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

package main

import (
	"context"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"os/signal"
	"syscall"

	"github.com/Magicking/faktur-daemon/backends"
	"github.com/Magicking/faktur-daemon/common"
	"github.com/Magicking/faktur-daemon/internal/anchor"

	ethcommon "github.com/ethereum/go-ethereum/common"
	flags "github.com/jessevdk/go-flags"
)

var opts struct {
	RpcURL          string `long:"rpc-url" env:"RPC_URL" description:"RPC URL for the node"`
	PrivateKey      string `long:"key" required:"true" env:"PRIVATE_KEY" description:"Private key used to sign transaction"`
	ChainId         string `long:"chain-id" required:"true" env:"CHAIN_ID" description:"Ethereum chain id"`
	ContractAddress string `long:"contract" required:"true" env:"CONTRACT" description:"faktur contract adresse"`
	CacheNonce      string `long:"cache-nonce" env:"CACHE_NONCE" description:"Path to cached nonce"`
	DbDSN           string `long:"db-dsn" env:"DB_DSN" required:"true" description:"Database DSN (e.g: /tmp/test.sqlite)"`
}

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// TODO reload config ?
	go func() {
		<-sigs
		done <- true
	}()

	parser := flags.NewParser(nil, flags.Default)
	parser.ShortDescription = "faktur"
	parser.LongDescription = "faktur daemon\n"

	_, err := parser.AddGroup("Services", "Services", &opts)
	if err != nil {
		log.Fatalf("Could not add group services: %v", err)
	}
	_, err = parser.AddGroup("HTTP Backend", "HTTP Backend", &backends.HTTPopts)
	if err != nil {
		log.Fatalf("Could not add group HTTP Backend: %v", err)
	}
	if _, err := parser.Parse(); err != nil {
		log.Fatalf("Could not parse arguments: %v", err)
	}
	var nonce *big.Int
	if opts.CacheNonce != "" {
		_nonce, err := ioutil.ReadFile(opts.CacheNonce)
		if err != nil {
			log.Println("ReadFile", err)
		} else {
			nonce = big.NewInt(0).SetBytes(_nonce)
		}
	}
	log.Printf("Anchor version 0.1")
	log.Printf("RPC URL %v", opts.RpcURL)
	log.Printf("Chain ID %v", opts.ChainId)
	log.Printf("Nonce %v", nonce)
	log.Printf("Contract is at: %v", opts.ContractAddress)
	/*
		var buf []byte
		if err == nil {
			buf = big.NewInt(int64(tx.Nonce()) + 1).Bytes()
		}
		if opts.CacheNonce != "" {
			if err := ioutil.WriteFile(opts.CacheNonce, buf, 0666); err != nil {
				log.Println(err)
			}
		}*/

	// create channel hash receiver
	hashC := make(chan ethcommon.Hash, 1)
	// TODO create channel preReceipt receiver
	// create channel merkleRoot receiver
	merkleRootC := make(chan ethcommon.Hash, 1)
	ctx := common.InitContext(context.Background())
	common.NewDBToContext(ctx, opts.DbDSN)
	common.NewGethClienToContext(ctx, opts.RpcURL)
	// TODO Get timeout from Smart Contract
	timeout := big.NewInt(10)
	go anchor.AnchorDaemon(hashC, merkleRootC, timeout)
	// Handle gathering hash
	// On timeout (period/10) produce merkleRoot and send that to merkleRoot chan
	// Send signed preReceipt
	// go start receiver
	// go start preReceipt submitter
	// go start merkleRoot submitter
	<-done
}
