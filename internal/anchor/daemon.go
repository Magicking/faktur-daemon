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
	"sort"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Magicking/faktur-daemon/internal/db"
	"github.com/Magicking/faktur-daemon/merkle"
	"github.com/ethereum/go-ethereum/common"
)

func AnchorDaemon(ctx context.Context, hashC chan []merkle.Hashable, merkleRootC chan common.Hash, timeout time.Duration) {

	ticker := time.NewTicker(timeout)
	var hashs []merkle.Hashable
	for {
		select {
		case hash := <-hashC:
			// TODO(6120) emit preReceipt
			hashs = append(hashs, hash...)
		case <-ticker.C:
			if len(hashs) == 0 {
				continue
			}
			if len(hashs) == 1 {
				hashs = append(hashs, hashs[0])
			}
			log.Printf("Hashs length: %v", len(hashs))
			// For Merkle tree order stability
			sort.Sort(merkle.OrderedBytes(hashs))
			merkleRoot, receipts := merkle.MerkleTreeHashProofsFromHashables(hashs)
			root := common.BytesToHash(merkleRoot)
			// save to database with merkleroot as key
			for i, e := range receipts {
				// Save Receipt
				err := db.SaveReceipt(ctx, e, hashs[i], root)
				if err != nil {
					log.WithFields(log.Fields{
						"hash": common.BytesToHash(hashs[i].Bytes()).Hex(),
						"root": root.Hex(),
					}).Warn(err)
				}
			}
			// Send merkleRoot to blockchain
			merkleRootC <- root
			log.Printf("Hash Root: %v\nReceipts len: %v", root.Hex(), len(receipts))
			hashs = nil
		}
	}

}
