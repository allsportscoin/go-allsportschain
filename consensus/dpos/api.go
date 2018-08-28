// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package dpos

import (
	"github.com/allsportschain/go-allsportschain/common"
	"github.com/allsportschain/go-allsportschain/consensus"
	"github.com/allsportschain/go-allsportschain/core/types"
	"github.com/allsportschain/go-allsportschain/rpc"
	"math/big"
	"github.com/allsportschain/go-allsportschain/trie"
	"fmt"
	"bytes"
	"github.com/allsportschain/go-allsportschain/rlp"
	"github.com/allsportschain/go-allsportschain/log"
)

// API is a user facing RPC API to allow controlling the delegate and voting
// mechanisms of the delegated-proof-of-stake
type API struct {
	chain consensus.ChainReader
	dpos  *Dpos
}

// GetValidators retrieves the list of the validators at specified block
func (api *API) GetValidators(number *rpc.BlockNumber) ([]common.Address, error) {
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	if header == nil {
		return nil, errUnknownBlock
	}

	epochTrie, err := types.NewEpochTrie(header.DposContext.EpochHash, api.dpos.db)
	if err != nil {
		return nil, err
	}

	validators := []common.Address{}
	key := []byte("validator")
	validatorsRLP := epochTrie.Get(key)
	if err := rlp.DecodeBytes(validatorsRLP, &validators); err != nil {
		return nil, fmt.Errorf("failed to decode validators: %s", err)
	}
	return validators, nil
}
// SetValidators retrieves the list of the validators at specified block
func (api *API) SetValidators(number *rpc.BlockNumber, validators []common.Address) error {
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	if header == nil {
		return errUnknownBlock
	}

	epochTrie, err := types.NewEpochTrie(header.DposContext.EpochHash, api.dpos.db)
	if err != nil {
		return err
	}
	dposContext := types.DposContext{}
	dposContext.SetEpoch(epochTrie)
	dposContext.SetValidators(validators)
	return nil
}

// GetConfirmedBlockNumber retrieves the latest irreversible block
func (api *API) GetConfirmedBlockNumber() (*big.Int, error) {
	var err error
	header := api.dpos.confirmedBlockHeader
	if header == nil {
		header, err = api.dpos.loadConfirmedBlockHeader(api.chain)
		if err != nil {
			return nil, err
		}
	}
	return header.Number, nil
}

//api for get all candidates form candidate trie
func (api * API) GetCandidates(number *rpc.BlockNumber) ([]common.Address, error) {
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		log.Error("latestedbalocknumber")
		header = api.chain.CurrentHeader()
	} else {
		log.Error("has number")
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	if header == nil {
		return nil, errUnknownBlock
	}

	candidateTrie, err := types.NewCandidateTrie(header.DposContext.CandidateHash, api.dpos.db)
	if err != nil {
		return nil, err
	}

	candidates := []common.Address{}
	iter := trie.NewIterator(candidateTrie.NodeIterator(nil))
	for iter.Next() {
		candidates = append(candidates, common.BytesToAddress(iter.Value))
	}

	return candidates, nil
}

//api for get candidate vote addr form vote trie
func (api * API) GetAddrVote(candidate common.Address) ([]common.Address, error) {
	header := api.chain.CurrentHeader()
	voteTrie, err := types.NewVoteTrie(header.DposContext.VoteHash, api.dpos.db)

	if err != nil {
		return []common.Address{}, err
	}

	delegateList := make([]common.Address, 0)
	voteIter := trie.NewIterator(voteTrie.NodeIterator(nil))
	for voteIter.Next() {
		dele := voteIter.Key
		cand := voteIter.Value
		if bytes.Equal(cand,candidate.Bytes()) {
			delegateList = append(delegateList,common.BytesToAddress(dele))
		}
	}
	return delegateList, nil
}
