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
	"time"

	"github.com/Magicking/faktur-daemon/backends"
	"github.com/Magicking/faktur-daemon/common"
	"github.com/Magicking/faktur-daemon/internal/anchor"
	"github.com/Magicking/faktur-daemon/internal/db"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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
		go func() {
			time.Sleep(2 * time.Second)
			log.Fatal("Forced interrupt")
		}()
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
	var nonce uint64
	if opts.CacheNonce != "" {
		_nonce, err := ioutil.ReadFile(opts.CacheNonce)
		if err != nil {
			log.Println("ReadFile", err)
		} else {
			nonce = big.NewInt(0).SetBytes(_nonce).Uint64()
		}
	}
	chainId := new(big.Int)
	chainId.SetString(opts.ChainId, 10)
	contractAddress := ethcommon.HexToAddress(opts.ContractAddress)
	log.Printf("Anchor version 0.1")
	log.Printf("RPC URL %v", opts.RpcURL)
	log.Printf("Chain ID %v", chainId.String())
	log.Printf("Nonce %v", nonce)
	log.Printf("Contract is at: %v", contractAddress.Hex())
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
	db.MigrateDatabase(ctx)
	common.NewGethClientToContext(ctx, opts.RpcURL)

	// Setup ethereum transaction sender/signer
	key, err := crypto.HexToECDSA(opts.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	anc := anchor.NewAnchor(key, nonce, chainId)

	// TODO Get timeout from Smart Contract
	timeout := time.Duration(1 * time.Second)
	go anchor.AnchorDaemon(ctx, hashC, merkleRootC, timeout)
	// Handle gathering hash
	// On timeout (period/10) produce merkleRoot and send that to merkleRoot chan
	// Send signed preReceipt
	// go start receiver
	http, err := backends.NewHTTP(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if err := http.Init(ctx, hashC); err != nil {
		log.Fatal(err)
	}
	go http.Run(ctx)
	// go start preReceipt submitter
	// go start merkleRoot submitter
	go anc.Run(ctx, contractAddress, merkleRootC)
	<-done
}
