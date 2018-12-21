package dpos

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
		"sort"

	"github.com/allsportschain/go-allsportschain/common"
	"github.com/allsportschain/go-allsportschain/core/state"
	"github.com/allsportschain/go-allsportschain/core/types"
		"github.com/allsportschain/go-allsportschain/log"
	"github.com/allsportschain/go-allsportschain/trie"
	"github.com/allsportschain/go-allsportschain/params"
)

type EpochContext struct {
	TimeStamp   int64
	DposContext *types.DposContext
	statedb     *state.StateDB
}

// countVotes
func (ec *EpochContext) countVotes() (votes map[common.Address]*big.Int, err error) {
	votes = map[common.Address]*big.Int{}
	delegateTrie := ec.DposContext.DelegateTrie()
	candidateTrie := ec.DposContext.CandidateTrie()
	statedb := ec.statedb

	iterCandidate := trie.NewIterator(candidateTrie.NodeIterator(nil))
	existCandidate := iterCandidate.Next()
	if !existCandidate {
		return votes, errors.New("no candidates")
	}
	for existCandidate {
		candidate := iterCandidate.Value
		candidateAddr := common.BytesToAddress(candidate)
		delegateIterator := trie.NewIterator(delegateTrie.NodeIterator(candidate))

		existDelegator := delegateIterator.NextPrefix(candidate)
		if !existDelegator {
			votes[candidateAddr] = new(big.Int)
			existCandidate = iterCandidate.Next()
			continue
		}
		for existDelegator {
			delegator := delegateIterator.Value
			score, ok := votes[candidateAddr]
			if !ok {
				score = new(big.Int)
			}
			delegatorAddr := common.BytesToAddress(delegator)
			weight := statedb.GetBalance(delegatorAddr)
			score.Add(score, weight)
			votes[candidateAddr] = score
			existDelegator = delegateIterator.NextPrefix(candidate)
		}
		existCandidate = iterCandidate.Next()
	}
	return votes, nil
}

func (ec *EpochContext) kickoutValidator(config *params.ChainConfig, header *types.Header, epoch int64) error {
	validators, err := ec.DposContext.GetValidators()
	if err != nil {
		return fmt.Errorf("failed to get validator: %s", err)
	}
	if len(validators) == 0 {
		return errors.New("no validator could be kickout")
	}

	epochDuration := types.EpochInterval
	// First epoch duration may lt epoch interval,
	// while the first block time wouldn't always align with epoch interval,
	// so caculate the first epoch duartion with first block time instead of epoch interval,
	// prevent the validators were kickout incorrectly.
	if ec.TimeStamp-timeOfFirstBlock < types.EpochInterval {
		epochDuration = ec.TimeStamp - timeOfFirstBlock
	}

	needKickoutValidators := common.SortableAddresses{}
	for _, validator := range validators {
		key := make([]byte, 8)
		binary.BigEndian.PutUint64(key, uint64(epoch))
		key = append(key, validator.Bytes()...)
		cnt := int64(0)
		if cntBytes := ec.DposContext.MintCntTrie().Get(key); cntBytes != nil {
			cnt = int64(binary.BigEndian.Uint64(cntBytes))
		}
		if cnt < epochDuration/types.BlockInterval/ types.MaxValidatorSize /2 {
			// not active validators need kickout
			needKickoutValidators = append(needKickoutValidators, &common.SortableAddress{Address:validator, Weight:big.NewInt(cnt)})
		}
	}
	// no validators need kickout
	needKickoutValidatorCnt := len(needKickoutValidators)
	if needKickoutValidatorCnt <= 0 {
		return nil
	}
	sort.Sort(sort.Reverse(needKickoutValidators))

	candidateCount := 0
	iter := trie.NewIterator(ec.DposContext.CandidateTrie().NodeIterator(nil))
	for iter.Next() {
		candidateCount++
		if candidateCount >= needKickoutValidatorCnt+types.SafeSize {
			break
		}
	}

	for i, validator := range needKickoutValidators {
		// ensure candidate count greater than or equal to safeSize
		if candidateCount <= types.SafeSize {
			log.Info("No more candidate can be kickout", "prevEpochID", epoch, "candidateCount", candidateCount, "needKickoutCount", len(needKickoutValidators)-i)
			return nil
		}

		if err := ec.DposContext.KickoutCandidate(config, header, validator.Address); err != nil {
			return err
		}
		// if kickout success, candidateCount minus 1
		candidateCount--
		log.Info("Kickout candidate", "prevEpochID", epoch, "candidate", validator.Address.String(), "mintCnt", validator.Weight.String())
	}
	return nil
}

func (ec *EpochContext) lookupValidator(now int64) (validator common.Address, err error) {
	offset := now % types.EpochInterval
	if offset%types.BlockInterval != 0 {
		return common.Address{}, ErrInvalidMintBlockTime
	}
	offset /= types.BlockInterval

	validators, err := ec.DposContext.GetValidators()
	if err != nil {
		return common.Address{}, err
	}
	validatorSize := len(validators)
	if validatorSize == 0 {
		return common.Address{}, errors.New("failed to lookup validator")
	}
	offset %= int64(validatorSize)
	return validators[offset], nil
}


func (ec *EpochContext) tryElect(config *params.ChainConfig, header,genesis, parent *types.Header) error {

	if config.MultiVoteBlock != nil && config.MultiVoteBlock.Cmp(genesis.Number) !=0 && config.MultiVoteBlock.Cmp(header.Number) == 0{
		voteIterator := trie.NewIterator(ec.DposContext.VoteTrie().NodeIterator(nil))
		existVote := voteIterator.Next()
		for existVote {
			ec.DposContext.VoteTrie().Delete(voteIterator.Key)
			ec.DposContext.VoteTrie().TryUpdate(append(voteIterator.Key, voteIterator.Value...),voteIterator.Value)
			existVote = voteIterator.Next()
		}
		log.Info("change multi-vote on block number : "+ config.MultiVoteBlock.String())
	}

	genesisEpoch := genesis.Time.Int64() / types.EpochInterval
	prevEpoch := parent.Time.Int64() / types.EpochInterval
	currentEpoch := ec.TimeStamp / types.EpochInterval

	prevEpochIsGenesis := prevEpoch == genesisEpoch
	if prevEpochIsGenesis && prevEpoch < currentEpoch {
		prevEpoch = currentEpoch - 1
	}
	prevEpochBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(prevEpochBytes, uint64(prevEpoch))
	//iter := trie.NewIterator(ec.DposContext.MintCntTrie().PrefixIterator(prevEpochBytes))
	iter := trie.NewIterator(ec.DposContext.MintCntTrie().NodeIterator(prevEpochBytes))
	for i := prevEpoch; i < currentEpoch; i++ {
		iBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(iBytes, uint64(i))
		// if prevEpoch is not genesis, kickout not active candidate
		if !prevEpochIsGenesis && iter.NextPrefix(iBytes) {
			if err := ec.kickoutValidator(config, header, i); err != nil {
				return err
			}
		}
		votes, err := ec.countVotes()
		if err != nil {
			return err
		}
		candidates := common.SortableAddresses{}
		for candidate, cnt := range votes {
			candidates = append(candidates, &common.SortableAddress{Address:candidate, Weight:cnt})
		}
		if len(candidates) < types.SafeSize {
			return errors.New("too few candidates")
		}
		sort.Sort(candidates)
		if len(candidates) > types.MaxValidatorSize {
			candidates = candidates[:types.MaxValidatorSize]
		}

		sortedValidators := make([]common.Address, 0)
		for _, candidate := range candidates {
			sortedValidators = append(sortedValidators, candidate.Address)
		}

		ec.DposContext.SetValidators(sortedValidators)
		log.Info("Come to new epoch", "prevEpoch", i, "nextEpoch", i+1)
	}
	return nil
}

