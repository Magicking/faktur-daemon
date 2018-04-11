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

package common

import (
	"context"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jinzhu/gorm"
)

type key int

var ctxValKey key = 0
var dbKey key = 1
var ethClientKey key = 2

type ctxValues struct {
	m map[key]interface{}
}

func (v ctxValues) Set(k key, val interface{}) {
	v.m[k] = val
}

func (v ctxValues) Get(k key) interface{} {
	val, ok := v.m[k]
	if !ok {
		log.Fatalf("Could not find key: %v", k)
	}
	return val
}

func NewCtxValues() *ctxValues {
	mm := make(map[key]interface{})
	cv := &ctxValues{
		m: mm,
	}
	return cv
}

func getContextValue(ctx context.Context, k key) interface{} {
	v, ok := ctx.Value(ctxValKey).(*ctxValues)
	if !ok {
		log.Fatalf("Could not obtain map context values")
	}
	return v.Get(k)
}

func setContextValue(ctx context.Context, k key, val interface{}) {
	v, ok := ctx.Value(ctxValKey).(*ctxValues)
	if !ok {
		log.Fatalf("Could not obtain map context values")
	}
	v.Set(k, val)
}

func InitContext(ctx context.Context) context.Context {
	values := NewCtxValues()
	return context.WithValue(ctx, ctxValKey, values)
}

func NewDBToContext(ctx context.Context, dbDsn string) {
	db, err := InitDatabase(dbDsn)
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}
	setContextValue(ctx, dbKey, db)
}

func DBFromContext(ctx context.Context) *gorm.DB {
	key := dbKey
	ret, ok := getContextValue(ctx, key).(*gorm.DB)
	if !ok {
		log.Fatalf("Could not cast context with key %d", key)
	}
	return ret
}

func NewGethClienToContext(ctx context.Context, rpc_url string) {
	var client *ethclient.Client
	var err error
	for i := 1; i < 10; i++ {
		client, err = ethclient.Dial(rpc_url)
		if err == nil || i == 10 {
			break
		}
		sleep := (2 << uint(i)) * time.Second
		log.Printf("Could not connect to Ethereum: %v", err)
		log.Printf("Waiting %v before retry", sleep)
		time.Sleep(sleep)
	}
	if err != nil {
		log.Fatalf("Could not initialize client context: %v", err)
	}
	// TODO Check needed protocol dependancies
	setContextValue(ctx, ethClientKey, client)
}

func ClientFromContext(ctx context.Context) *ethclient.Client {
	key := ethClientKey
	ret, ok := getContextValue(ctx, key).(*ethclient.Client)
	if !ok {
		log.Fatalf("Could not cast context with key %d", key)
	}
	return ret
}

/*
func NewSchedulerToContext(ctx context.Context, tick time.Duration) {
	c := NewScheduler(ctx, tick)
	setContextValue(ctx, schedulerKey, c)
}

func SchedulerChanFromContext(ctx context.Context) chan callback {
	key := schedulerKey
	ret, ok := getContextValue(ctx, key).(chan callback)
	if !ok {
		log.Fatalf("Could not cast context with key %d", key)
	}
	return ret
}*/
