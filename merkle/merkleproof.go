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

package merkle

import (
	"encoding/hex"

	"golang.org/x/crypto/sha3"
)

type Hashitem []byte

func (hi Hashitem) Hash() []byte {
	ret := sha3.Sum256(hi)
	return ret[:]
}

func (hi Hashitem) Bytes() []byte {
	return []byte(hi)
}

type AnchorPoint struct {
	SourceID string `json:"sourceId"`
	Type     string `json:"type"`
}

type Chainpoint struct {
	Context    string        `json:"@context"`
	Anchors    []AnchorPoint `json:"anchors"`
	MerkleRoot string        `json:"merkleRoot"`
	Proof      Branch        `json:"proof"`
	TargetHash string        `json:"targetHash"`
	Type       string        `json:"type"`
}

func NewChainpoints(items []Hashable) ([]Chainpoint, []byte) {
	rootH, proofs := MerkleTreeHashProofsFromHashables(items)
	if len(proofs) != len(items) {
		panic("Not all items were entered into merkle tree")
	}
	unanchoredReceipts := make([]Chainpoint, len(items))
	for i, v := range items {
		unanchoredReceipts[i].Type = "ChainpointSHA3-256v2"
		unanchoredReceipts[i].Context = "https://w3id.org/chainpoint/v2"
		unanchoredReceipts[i].Proof = *proofs[i]
		unanchoredReceipts[i].MerkleRoot = hex.EncodeToString(rootH)
		unanchoredReceipts[i].TargetHash = hex.EncodeToString(v.Bytes())
	}
	return unanchoredReceipts, rootH
}

func (cp *Chainpoint) Verify() bool {
	targetHash, err := hex.DecodeString(cp.TargetHash)
	if err != nil {
		return false
	}
	root, err := hex.DecodeString(cp.MerkleRoot)
	if err != nil {
		return false
	}
	return Verify(targetHash, cp.Proof, root)
}
