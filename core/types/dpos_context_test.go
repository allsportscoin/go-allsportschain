package types

import (
	"testing"

	"github.com/allsportschain/go-allsportschain/common"
	"github.com/allsportschain/go-allsportschain/socdb"
	"github.com/stretchr/testify/assert"
	"github.com/allsportschain/go-allsportschain/params"
	"math/big"
	"github.com/allsportschain/go-allsportschain/trie"
	"errors"
	"fmt"
	"strconv"
)

func TestDposContextSnapshot(t *testing.T) {
	config := &params.ChainConfig{
		MultiVoteBlock:      big.NewInt(0),
	}
	header := &Header{
		Number: big.NewInt(1000),
	}

	db := socdb.NewMemDatabase()
	dposContext, err := NewDposContext(db)
	assert.Nil(t, err)

	snapshot := dposContext.Snapshot()
	assert.Equal(t, dposContext.Root(), snapshot.Root())
	assert.NotEqual(t, dposContext, snapshot)

	// change dposContext
	assert.Nil(t, dposContext.BecomeCandidate(config, header, common.HexToAddress("0x44d1ce0b7cb3588bca96151fe1bc05af38f91b6c")))
	assert.NotEqual(t, dposContext.Root(), snapshot.Root())

	// revert snapshot
	dposContext.RevertToSnapShot(snapshot)
	assert.Equal(t, dposContext.Root(), snapshot.Root())
	assert.NotEqual(t, dposContext, snapshot)
}

func TestMultiVote(t *testing.T) {
	config := &params.ChainConfig{
		MultiVoteBlock:      big.NewInt(0),
	}
	header := &Header{
		Number: big.NewInt(1000),
	}

	delegateAddr := common.HexToAddress("0xb040353ec0f2c113d5639444f7253681aecda1f8")
	candidate1Addr := common.HexToAddress("0x44d1ce0b7cb3588bca96151fe1bc05af38f91b6e")
	candidate2Addr := common.HexToAddress("0xa60a3886b552ff9992cfcd208ec1152079e046c2")
	candidate3Addr := common.HexToAddress("0x4e080e49f62694554871e669aeb4ebe17c4a9670")
	voteMap := map[common.Address]common.Address{
		candidate1Addr: delegateAddr,
		candidate2Addr: delegateAddr,
		candidate3Addr: delegateAddr,
	}
	db := socdb.NewMemDatabase()
	dposContext, err := NewDposContext(db)
	assert.Nil(t, err)
	for candidate, elector := range voteMap {
		assert.Nil(t, dposContext.BecomeCandidate(config,header,candidate))
		assert.Nil(t, dposContext.Delegate(config, header, elector, candidate))
	}

	voteIterator := trie.NewIterator(dposContext.VoteTrie().NodeIterator(nil))
	existVote := voteIterator.Next()
	for existVote {
		assert.Equal(t,append(voteMap[common.BytesToAddress(voteIterator.Value)].Bytes(),voteIterator.Value...),voteIterator.Key)
		existVote = voteIterator.Next()
	}

	delegateIterator := trie.NewIterator(dposContext.DelegateTrie().NodeIterator(nil))
	existDelegate := delegateIterator.Next()
	for existDelegate {
		candidate := delegateIterator.Key[:len(delegateIterator.Key)/2]
		assert.Equal(t,append(candidate,voteMap[common.BytesToAddress(candidate)].Bytes()...),delegateIterator.Key)
		existDelegate = delegateIterator.Next()
	}

	voteIterator = trie.NewIterator(dposContext.VoteTrie().NodeIterator(nil))
	assert.Equal(t, voteIterator.NextPrefixCount(delegateAddr.Bytes()), int64(len(voteMap)))

	assert.Nil(t, dposContext.UnDelegate(config, header, delegateAddr, candidate1Addr))
	voteIterator = trie.NewIterator(dposContext.VoteTrie().NodeIterator(nil))
	assert.Equal(t, voteIterator.NextPrefixCount(delegateAddr.Bytes()), int64(len(voteMap)-1))

	assert.Nil(t, dposContext.UnDelegate(config, header, delegateAddr, candidate2Addr))
	voteIterator = trie.NewIterator(dposContext.VoteTrie().NodeIterator(nil))
	assert.Equal(t, voteIterator.NextPrefixCount(delegateAddr.Bytes()), int64(len(voteMap)-2))

}

func TestMultiVoteHasVoted(t *testing.T) {
	config := &params.ChainConfig{
		MultiVoteBlock: big.NewInt(0),
	}
	header := &Header{
		Number: big.NewInt(1000),
	}


	delegateAddr := common.HexToAddress("0xb040353ec0f2c113d5639444f7253681aecda1f8")
	candidateAddr := common.HexToAddress("0x4e080e49f62694554871e669aeb4ebe17c4a9670")
	db := socdb.NewMemDatabase()
	dposContext, err := NewDposContext(db)
	assert.Nil(t, err)
	assert.Nil(t, dposContext.BecomeCandidate(config, header, candidateAddr))
	assert.Nil(t, dposContext.Delegate(config, header, delegateAddr, candidateAddr))
	err = dposContext.Delegate(config, header, delegateAddr, candidateAddr)
	assert.EqualError(t,errors.New(candidateAddr.String() + " Has already been voted"),err.Error())
}

func TestMultiVoteExceedMaxVoteCandidateNum(t *testing.T) {
	config := &params.ChainConfig{
		MultiVoteBlock: big.NewInt(0),
	}
	header := &Header{
		Number: big.NewInt(1000),
	}


	delegateAddr := common.HexToAddress("0xb040353ec0f2c113d5639444f7253681aecda1f8")
	candidateAddr := common.HexToAddress("0x4e080e49f62694554871e669aeb4ebe17c4a9670")
	db := socdb.NewMemDatabase()
	dposContext, err := NewDposContext(db)
	assert.Nil(t, err)

	for i := 0; i < MaxVoteCandidateNum; i++ {
		tmpCandidateAddr := common.BytesToAddress([]byte("addr" + strconv.Itoa(i)))
		assert.Nil(t, dposContext.BecomeCandidate(config, header, tmpCandidateAddr))
		assert.Nil(t, dposContext.Delegate(config, header, delegateAddr, tmpCandidateAddr))
	}
	voteIterator := trie.NewIterator(dposContext.VoteTrie().NodeIterator(delegateAddr.Bytes()))
	existVoteCount := voteIterator.NextPrefixCount(delegateAddr.Bytes())
	assert.Nil(t, dposContext.BecomeCandidate(config, header, candidateAddr))
	err = dposContext.Delegate(config, header, delegateAddr, candidateAddr)
	expectedErr := errors.New(fmt.Sprintf("%v has already voted %v votes, Can't exceed %v votes.", delegateAddr.String(), existVoteCount, MaxVoteCandidateNum))
	assert.EqualError(t,expectedErr,err.Error())
}

func TestMultiVoteInvalidCandidate(t *testing.T) {
	config := &params.ChainConfig{
		MultiVoteBlock: big.NewInt(0),
	}
	header := &Header{
		Number: big.NewInt(1000),
	}

	delegateAddr := common.HexToAddress("0xb040353ec0f2c113d5639444f7253681aecda1f8")
	candidateAddr := common.HexToAddress("0x4e080e49f62694554871e669aeb4ebe17c4a9670")
	db := socdb.NewMemDatabase()
	dposContext, err := NewDposContext(db)
	assert.Nil(t, err)
	err = dposContext.Delegate(config, header, delegateAddr, candidateAddr)
	expectedErr := errors.New(candidateAddr.String() + " is invalid candidate")
	assert.EqualError(t,expectedErr,err.Error())
}