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

// Internal receipt format

type Leaf struct {
	Left  string `json:"left,omitempty"`  // Hashes from leaf's sibling to a root's child.
	Right string `json:"right,omitempty"` // Hashes from leaf's sibling to a root's child.
}

type Receipt struct {
	Proof      []Leaf `json:"proof"`
	MerkleRoot string `json:"merkleRoot"`
	TargetHash string `json:"targetHash"`
}
