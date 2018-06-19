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
	"crypto/ecdsa"

	"github.com/Magicking/faktur-daemon/merkle"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

// Internal receipt format

type Receipt struct {
	Proof      merkle.Branch  `json:"proof"`
	MerkleRoot ethcommon.Hash `json:"merkleRoot"`
	TargetHash ethcommon.Hash `json:"targetHash"`
}

type Faktur struct {
	Receipt
	PrivateKey *ecdsa.PrivateKey
}

// TODO
func NewFaktur(r *Receipt, p *ecdsa.PrivateKey) *Faktur {
	return nil
}

// TODO
func (f *Faktur) Serialize() []byte {
	return nil
}
