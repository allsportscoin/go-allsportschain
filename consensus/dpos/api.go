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
		header = api.chain.CurrentHeader()
	} else {
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

//api for get Delegate addr form delegateTrie
func (api * API) GetDelegatesByCandidate(candidate common.Address) ([]common.Address, error) {
	header := api.chain.CurrentHeader()
	delegateTrie, err := types.NewDelegateTrie(header.DposContext.DelegateHash, api.dpos.db)

	if err != nil {
		return []common.Address{}, err
	}

	delegateList := make([]common.Address, 0)
	delegateIterator := trie.NewIterator(delegateTrie.NodeIterator(candidate.Bytes()))
	for delegateIterator.NextPrefix(candidate.Bytes()) {
		delegate := delegateIterator.Value
		delegateList = append(delegateList,common.BytesToAddress(delegate))
	}
	return delegateList, nil
}

//api for get candidate addr form vote trie
func (api * API) GetCandidatesByDelegate(delegate common.Address) ([]common.Address, error) {
	header := api.chain.CurrentHeader()
	voteTrie, err := types.NewVoteTrie(header.DposContext.VoteHash, api.dpos.db)
	if err != nil {
		return []common.Address{}, err
	}

	candidate, err := voteTrie.TryGet(delegate.Bytes())
	if err != nil {
		return []common.Address{}, err
	}
	candidateList := make([]common.Address, 0)
	if len(candidate) != 0 {
		candidateList = append(candidateList, common.BytesToAddress(candidate))
	}
	return candidateList, nil
}

func (api * API) GetAddrIsCandidate(addr common.Address) (bool, error) {
	header := api.chain.CurrentHeader()
	candidateTrie, err := types.NewCandidateTrie(header.DposContext.CandidateHash, api.dpos.db)
	if err != nil {
		return false, err
	}

	candidate, err := candidateTrie.TryGet(addr.Bytes())
	if err != nil {
		return false, err
	}

	if bytes.Equal(candidate,addr.Bytes()) {
		return true, nil
	}else{
		return false, nil
	}
}
