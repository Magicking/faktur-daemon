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

/*
Computes a deterministic minimal height merkle tree hash.
If the number of items is not a power of two, some leaves
will be at different levels. Tries to keep both sides of
the tree the same size, but the left may be one greater.

Use this for short deterministic trees, such as the validator list.
For larger datasets, use IAVLTree.

                        *
                       / \
                     /     \
                   /         \
                 /             \
                *               *
               / \             / \
              /   \           /   \
             /     \         /     \
            *       *       *       h6
           / \     / \     / \
          h0  h1  h2  h3  h4  h5

*/

package merkle

import (
	"bytes"

	"golang.org/x/crypto/sha3"
)

type Hashable interface {
	Bytes() []byte
}

type OrderedBytes []Hashable

func (h OrderedBytes) Less(i, j int) bool {
	return bytes.Compare(h[i].Bytes(), h[j].Bytes()) < 0
}

func (h OrderedBytes) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h OrderedBytes) Len() int {
	return len(h)
}

// Merkle Tree Hash Proof version
func SimpleHashFromTwoHashes(left []byte, right []byte) []byte {
	buffer := [][]byte{left, right}
	hashed := sha3.Sum256(bytes.Join(buffer, nil))

	return hashed[:]
}

func SimpleHashFromHashes(hashes [][]byte) []byte {
	// Recursive impl.
	switch len(hashes) {
	case 0:
		return nil
	case 1:
		return hashes[0]
	default:
		left := SimpleHashFromHashes(hashes[:(len(hashes)+1)/2])
		right := SimpleHashFromHashes(hashes[(len(hashes)+1)/2:])
		return SimpleHashFromTwoHashes(left, right)
	}
}

//--------------------------------------------------------------------------------

type Leaf struct {
	Left  []byte `json:"left,omitempty"`  // Hashes from leaf's sibling to a root's child.
	Right []byte `json:"right,omitempty"` // Hashes from leaf's sibling to a root's child.
}

func (l *Leaf) Serialize() string {
	if l.Left != nil && l.Right != nil {
		panic("Both left and right leaf set")
	}
	if l.Left == nil && l.Right == nil {
		panic("Both left and right leaf unset")
	}
	if l.Left != nil {
		return string(l.Left)
	}
	return string(l.Right)
}

type Branch struct {
	AuditPath []Leaf // Hashes from leaf's sibling to a root's child.
}

func (b *Branch) String() (ret string) {
	var buffer bytes.Buffer
	for i, r := range b.AuditPath {
		if r.Left != nil {
			buffer.WriteRune('0')
			buffer.Write(r.Left)
			continue
		}
		if r.Write != nil {
			buffer.WriteRune('1')
			buffer.Write(r.Right)
			continue
		}
		panic("Both left and right leaf unset")
	}
	return ret.String()
}

func (b *Branch) Serialize() (ret string) {
	for _, _ = range b.AuditPath {
		//set bitmap
		return "TODO1" //TODO
	}
	return "TODO2"
}

// proofs[0] is the proof for items[0].
func MerkleTreeHashProofsFromHashables(items []Hashable) (rootHash []byte, proofs []*Branch) {
	trails, rootSPN := trailsFromHashables(items)
	rootHash = rootSPN.Hash
	proofs = make([]*Branch, len(items))
	for i, trail := range trails {
		proofs[i] = &Branch{
			AuditPath: trail.FlattenAunts(),
		}
	}
	return
}

// Verify that leafHash is a leaf hash of the simple-merkle-tree
// which hashes to rootHash.
func Verify(targetHash []byte, proof Branch, rootHash []byte) bool {
	computedHash := computeHashFromAunts(targetHash, proof)
	if computedHash == nil {
		return false
	}
	if !bytes.Equal(computedHash, rootHash) {
		return false
	}
	return true
}

// Use the leafHash and innerHashes to get the root merkle hash.
// If the length of the innerHashes slice isn't exactly correct, the result is nil.
func computeHashFromAunts(targetHash []byte, proofs Branch) []byte {
	result := targetHash
	for _, proof := range proofs.AuditPath {
		if proof.Left != nil {
			result = SimpleHashFromTwoHashes(proof.Left, result)
		} else if proof.Right != nil {
			result = SimpleHashFromTwoHashes(result, proof.Right)
		} else {
			panic("Left or right should be set for leaf proof")
		}
	}
	return result
}

// Helper structure to construct merkle proof.
// The node and the tree is thrown away afterwards.
// Exactly one of node.Left and node.Right is nil, unless node is the root, in which case both are nil.
// node.Parent.Hash = hash(node.Hash, node.Right.Hash) or
// 									  hash(node.Left.Hash, node.Hash), depending on whether node is a left/right child.
type SimpleProofNode struct {
	Hash   []byte
	Parent *SimpleProofNode
	Left   *SimpleProofNode // Left sibling  (only one of Left,Right is set)
	Right  *SimpleProofNode // Right sibling (only one of Left,Right is set)
}

// Starting from a leaf SimpleProofNode, FlattenAunts() will return
// the inner hashes for the item corresponding to the leaf.
func (spn *SimpleProofNode) FlattenAunts() (leaflist []Leaf) {
	// Nonrecursive impl.
	for spn != nil {
		var lv Leaf
		if spn.Left != nil {
			lv = Leaf{Left: spn.Left.Hash}
		} else if spn.Right != nil {
			lv = Leaf{Right: spn.Right.Hash}
		} else {
			break
		}
		spn = spn.Parent
		leaflist = append(leaflist, lv)
	}
	return leaflist
}

// trails[0].Hash is the leaf hash for items[0].
// trails[i].Parent.Parent....Parent == root for all i.
func trailsFromHashables(items []Hashable) (trails []*SimpleProofNode, root *SimpleProofNode) {
	// Recursive impl.
	switch len(items) {
	case 0:
		return nil, nil
	case 1:
		trail := &SimpleProofNode{items[0].Bytes(), nil, nil, nil}
		return []*SimpleProofNode{trail}, trail
	default:
		lefts, leftRoot := trailsFromHashables(items[:(len(items)+1)/2])
		rights, rightRoot := trailsFromHashables(items[(len(items)+1)/2:])
		rootHash := SimpleHashFromTwoHashes(leftRoot.Hash, rightRoot.Hash)
		root := &SimpleProofNode{rootHash, nil, nil, nil}
		leftRoot.Parent = root
		leftRoot.Right = rightRoot
		rightRoot.Parent = root
		rightRoot.Left = leftRoot
		return append(lefts, rights...), root
	}
}
