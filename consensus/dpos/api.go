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
	"github.com/allsportschain/go-allsportschain/common/hexutil"
	"encoding/binary"
		"sort"
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
	log.Debug("number:","number",*number)
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
func (api * API) GetDelegatesByCandidate(candidate common.Address,number *rpc.BlockNumber) ([]common.Address, error) {
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	if header == nil {
		return nil, errUnknownBlock
	}
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
func (api * API) GetCandidatesByDelegate(delegate common.Address, number *rpc.BlockNumber) ([]common.Address, error) {
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	if header == nil {
		return nil, errUnknownBlock
	}
	voteTrie, err := types.NewVoteTrie(header.DposContext.VoteHash, api.dpos.db)
	if err != nil {
		return []common.Address{}, err
	}

	candidateList := make([]common.Address, 0)
	candidateIterator := trie.NewIterator(voteTrie.NodeIterator(delegate.Bytes()))
	for candidateIterator.NextPrefix(delegate.Bytes()) {
		candidate := candidateIterator.Value
		candidateList = append(candidateList,common.BytesToAddress(candidate))
	}
	return candidateList, nil
}

func (api * API) GetAddrIsCandidate(addr common.Address,number *rpc.BlockNumber) (bool, error) {
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	if header == nil {
		return false, errUnknownBlock
	}
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

func (api * API) GetTotalWei(number *rpc.BlockNumber) (*big.Int, error) {
	blockNumber := int64(0)
	if number == nil || *number == rpc.LatestBlockNumber {
		blockNumber = api.chain.CurrentHeader().Number.Int64()
	} else {
		blockNumber = number.Int64()
	}
	increaseWei := big.NewInt(0).Mul(defaultBlockReward,big.NewInt(blockNumber))
	totalWei := big.NewInt(0).Add(socInitCount,increaseWei)
	return totalWei, nil
}

func (api * API) GetEpochInterval(number *rpc.BlockNumber) ( *big.Int, error) {
	return big.NewInt(epochInterval), nil
}

func (api * API) GetValidatorsMintCount(number *rpc.BlockNumber) (map[string]interface{}, error) {
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	if header == nil {
		return nil, errUnknownBlock
	}
	timestamp := header.Time.Uint64()
	epochIntervalBig,_ := api.GetEpochInterval(number)
	epochIntervalUint64 := epochIntervalBig.Uint64()


	mintCountList := []int64{}
	currentEpoch := timestamp / epochIntervalUint64
	epochProgress := timestamp % epochIntervalUint64
	fields := map[string]interface{}{
		"validatorsList": []common.Address{},
		"mintCountList": mintCountList,
		"epochInterval":epochIntervalUint64,
		"currentEpoch":currentEpoch,
		"epochProgress":epochProgress,
	}

	validators,err := api.GetValidators(number)
	if err != nil {
		return fields, err
	}
	fields["validatorsList"] = validators


	MintCntTrie, err := types.NewMintCntTrie(header.DposContext.MintCntHash, api.dpos.db)
	if err != nil {
		return fields, err
	}

	for _, validator := range validators {
		key := make([]byte, 8)
		binary.BigEndian.PutUint64(key, uint64(currentEpoch))
		key = append(key, validator.Bytes()...)
		cnt := int64(0)
		if cntBytes := MintCntTrie.Get(key); cntBytes != nil {
			cnt = int64(binary.BigEndian.Uint64(cntBytes))
		}
		mintCountList = append(mintCountList,cnt)
	}
	fields["mintCountList"] = mintCountList
	return fields,nil
}

func (api * API) GetCandidatesAndVoteCountTopN(topN hexutil.Uint64, number *rpc.BlockNumber) ( map[string]interface{}, error) {
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	if header == nil {
		return nil, errUnknownBlock
	}

	candidatesList := []common.Address{}
	voteCountList := []*big.Int{}
	totalVoteCount := big.NewInt(0)
	fields := map[string]interface{}{
		"candidatesList": candidatesList,
		"voteCountList": voteCountList,
		"totalVoteCount":totalVoteCount,
	}
	retCandidates := common.SortableAddresses{}

	statedb, err := api.chain.StateAt(header.Root)
	if statedb == nil || err != nil {
		return fields, err
	}

	dposContext, err := types.NewDposContextFromProto(api.dpos.db, &types.DposContextProto{
		EpochHash:     header.DposContext.EpochHash,
		DelegateHash:  header.DposContext.DelegateHash,
		CandidateHash: header.DposContext.CandidateHash,
		VoteHash:      header.DposContext.VoteHash,
		MintCntHash:  header.DposContext.MintCntHash,
	})
	if err != nil {
		return fields, err
	}

	epochContext := &EpochContext{
		DposContext: dposContext,
		statedb:     statedb,
	}
	votes, err := epochContext.countVotes()
	if err != nil {
		return fields, err
	}

	for candidate, cnt := range votes {
		retCandidates = append(retCandidates, &common.SortableAddress{ Address:candidate, Weight:cnt})
	}

	sort.Sort(retCandidates)
	if topN > 0 && hexutil.Uint64(len(retCandidates)) > topN {
		retCandidates = retCandidates[:topN]
	}
	for _, sortableAddress := range retCandidates{
		candidatesList = append(candidatesList,sortableAddress.Address)
		voteCountList = append(voteCountList,sortableAddress.Weight)
	}

	delegateMap := map[common.Address]bool{}
	delegateIterator := trie.NewIterator(dposContext.DelegateTrie().NodeIterator(nil))
	existDelegator := delegateIterator.Next()

	for existDelegator {
		delegator := delegateIterator.Value
		delegatorAddr := common.BytesToAddress(delegator)
		if delegateMap[delegatorAddr] == false {
			count := statedb.GetBalance(delegatorAddr)
			totalVoteCount.Add(totalVoteCount, count)
			delegateMap[delegatorAddr] = true
		}
		existDelegator = delegateIterator.Next()
	}

	fields["candidatesList"] = candidatesList
	fields["voteCountList"] = voteCountList
	fields["totalVoteCount"] = totalVoteCount
	return fields,nil
}

func (api *API) GetAddrVoteCount(address common.Address, number *rpc.BlockNumber) ( *big.Int, error) {
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	if header == nil {
		return nil, errUnknownBlock
	}
	totalVoteCount := big.NewInt(0)

	delegateTrie, err := types.NewDelegateTrie(header.DposContext.DelegateHash, api.dpos.db)
	if err != nil {
		return totalVoteCount, err
	}
	statedb, err := api.chain.StateAt(header.Root)
	if statedb == nil || err != nil {
		return totalVoteCount, err
	}

	delegateIterator := trie.NewIterator(delegateTrie.NodeIterator(address.Bytes()))
	existDelegator := delegateIterator.NextPrefix(address.Bytes())

	for existDelegator {
		delegator := delegateIterator.Value
		delegatorAddr := common.BytesToAddress(delegator)
		count := statedb.GetBalance(delegatorAddr)
		totalVoteCount.Add(totalVoteCount, count)
		existDelegator = delegateIterator.NextPrefix(address.Bytes())
	}

	return totalVoteCount, nil
}