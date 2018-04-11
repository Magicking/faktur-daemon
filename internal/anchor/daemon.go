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
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func AnchorDaemon(hashC chan common.Hash, merkleRootC chan common.Hash, timeout *big.Int) {

	/*	// For Merkle tree order stability
		sort.Strings(voteHashs)
		hashs := make([]merkle.Hashable, len(voteHashs))
		for i, e := range voteHashs {
			h := sha3.New256()
			voteHash := common.StringToHash(e)
			//sum256 the file
			if _, err := io.Copy(h, bytes.NewBuffer(voteHash.Bytes())); err != nil {
				err_str := fmt.Sprintf("ListVotes(%s) error: %v", params.BallotID, err)
				return operations.NewScrutinCloseDefault(500).WithPayload(&models.Error{
					Message: &err_str})
			}
			// hash to insert
			hash := h.Sum(nil)
			hashs[i] = merkle.Hashitem(hash[:])
		}
		if len(voteHashs) == 1 {
			hashs = append(hashs, hashs[0])
		}
		receipts, merkleRoot := merkle.NewChainpoints(hashs)
		// Send transaction
		txhash, err := SendWithValueMessage(ctx, _from, nil, merkleRoot)
		if err != nil {
			err_str := fmt.Sprintf("Sendata(%s) error: %v", merkleRoot, err)
			return operations.NewScrutinCloseDefault(500).WithPayload(&models.Error{
				Message: &err_str})
		}
		now := time.Now()
		for i, v := range receipts {
			//Discard second entry in tree if there is only one vote
			if len(voteHashs) == 1 && i > 0 {
				break
			}
			//fill receipt
			v.Anchors = []merkle.AnchorPoint{merkle.AnchorPoint{SourceID: txhash.String(), Type: "ETHData"}}
		}
	*/
}
