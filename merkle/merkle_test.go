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
	"bytes"
	"fmt"
	"testing"

	"encoding/hex"
)

func TestMerkle(t *testing.T) {
	H0, _ := hex.DecodeString("bdf8c9bdf076d6aff0292a1c9448691d2ae283f2ce41b045355e2c8cb8e85ef2")
	H1, _ := hex.DecodeString("cb0dbbedb5ec5363e39be9fc43f56f321e1572cfcf304d26fc67cb6ea2e49faf")
	H2, _ := hex.DecodeString("da0ed1fecac504ea4f76d241a45032fa97b9eb692614419a04c9a9c32e39df2d")

	cpProofs, rootP := NewChainpoints([]Hashable{Hashitem(H0), Hashitem(H1), Hashitem(H2)})
	rootH, proofs := ChainpointProofsFromHashables([]Hashable{Hashitem(H0), Hashitem(H1), Hashitem(H2)})
	hash := SimpleHashFromTwoHashes(H0, nil)

	// single item
	ok := Verify(hash, *proofs[0], rootH)
	fmt.Println(hex.EncodeToString(proofs[0].Aunts[0].Left))
	fmt.Println(hex.EncodeToString(proofs[0].Aunts[0].Right))
	fmt.Println(hex.EncodeToString(rootH))
	if ok {
		t.Log("Okey !")
	} else {
		t.Error("Ko !")
	}
	if bytes.Equal(rootP, rootH) {
		t.Log("Okey !")
	} else {
		t.Error("Ko !")
	}
	for i := range cpProofs {
		ok = cpProofs[i].MerkleVerify()
		if ok {
			t.Log("Okey !")
		} else {
			t.Error("Ko !")
		}
	}
}
