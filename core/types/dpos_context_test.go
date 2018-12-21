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
	"encoding/binary"
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

func TestDposContextDelegateAndUnDelegate(t *testing.T) {
	//not MultiVote
	config := &params.ChainConfig{
		MultiVoteBlock: big.NewInt(1000),
	}
	header := &Header{
		Number: big.NewInt(99),
	}

	delegateAddr := common.HexToAddress("0xb040353ec0f2c113d5639444f7253681aecda1f8")
	candidate1Addr := common.HexToAddress("0x44d1ce0b7cb3588bca96151fe1bc05af38f91b6e")
	voteMap := map[common.Address]common.Address{
		candidate1Addr: delegateAddr,
	}
	db := socdb.NewMemDatabase()
	dposContext, err := NewDposContext(db)
	assert.Nil(t, err)

	// BecomeCandidate and delegate
	for candidate, elector := range voteMap {
		assert.Nil(t, dposContext.BecomeCandidate(config,header,candidate))
		assert.Nil(t, dposContext.Delegate(config, header, elector, candidate))
	}


	voteIterator := trie.NewIterator(dposContext.VoteTrie().NodeIterator(nil))
	existVote := voteIterator.Next()
	for existVote {
		assert.Equal(t,voteMap[common.BytesToAddress(voteIterator.Value)].Bytes(),voteIterator.Key)
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

	//UnDelegate candidate1Addr
	assert.Nil(t, dposContext.UnDelegate(config, header, delegateAddr, candidate1Addr))
	// UnDelegate,then vote Info is none
	voteIterator = trie.NewIterator(dposContext.VoteTrie().NodeIterator(nil))
	assert.Equal(t, voteIterator.NextPrefixCount(delegateAddr.Bytes()), int64(0))
	delegateIter := trie.NewIterator(dposContext.DelegateTrie().NodeIterator(candidate1Addr.Bytes()))
	assert.False(t, delegateIter.NextPrefix(candidate1Addr.Bytes()))

	//delegator delegate to not exist candidate
	assert.NotNil(t, dposContext.Delegate(config, header, delegateAddr, common.HexToAddress("0xab")))
	//undelegator delegate to not exist candidate
	assert.NotNil(t, dposContext.UnDelegate(config, header, delegateAddr, common.HexToAddress("0xab")))
	//undelegator ，Cancel a candidate who didn't vote.
	candidate4Addr := common.HexToAddress("0x4e080e49f62694554871e669aeb4ebe17c4a9671")
	assert.Nil(t, dposContext.BecomeCandidate(config,header,candidate4Addr))
	assert.NotNil(t, dposContext.UnDelegate(config, header, delegateAddr, candidate1Addr))


}

func TestDposContextDelegateAndUnDelegateMultiVote(t *testing.T) {
	//MultiVote
	config := &params.ChainConfig{
		MultiVoteBlock: big.NewInt(0),
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

	//delegate
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

	//UnDelegate candidate1Addr, 1 reduction in the number of votes cast by delegate,
	assert.Nil(t, dposContext.UnDelegate(config, header, delegateAddr, candidate1Addr))
	voteIterator = trie.NewIterator(dposContext.VoteTrie().NodeIterator(nil))
	assert.Equal(t, voteIterator.NextPrefixCount(delegateAddr.Bytes()), int64(len(voteMap)-1))
	delegateIter := trie.NewIterator(dposContext.DelegateTrie().NodeIterator(candidate1Addr.Bytes()))
	assert.False(t, delegateIter.NextPrefix(candidate1Addr.Bytes()))

	//UnDelegate candidate2Addr, 1 reduction in the number of votes cast by delegate,
	assert.Nil(t, dposContext.UnDelegate(config, header, delegateAddr, candidate2Addr))
	voteIterator = trie.NewIterator(dposContext.VoteTrie().NodeIterator(nil))
	assert.Equal(t, voteIterator.NextPrefixCount(delegateAddr.Bytes()), int64(len(voteMap)-2))
	delegateIter = trie.NewIterator(dposContext.DelegateTrie().NodeIterator(candidate2Addr.Bytes()))
	assert.False(t, delegateIter.NextPrefix(candidate2Addr.Bytes()))

	//UnDelegate candidate3Addr, 1 reduction in the number of votes cast by delegate,
	assert.Nil(t, dposContext.UnDelegate(config, header, delegateAddr, candidate3Addr))
	voteIterator = trie.NewIterator(dposContext.VoteTrie().NodeIterator(nil))
	assert.Equal(t, voteIterator.NextPrefixCount(delegateAddr.Bytes()), int64(len(voteMap)-3))
	delegateIter = trie.NewIterator(dposContext.DelegateTrie().NodeIterator(candidate3Addr.Bytes()))
	assert.False(t, delegateIter.NextPrefix(candidate3Addr.Bytes()))


	//delegator delegate to not exist candidate
	assert.NotNil(t, dposContext.Delegate(config, header, delegateAddr, common.HexToAddress("0xab")))
	//undelegator delegate to not exist candidate
	assert.NotNil(t, dposContext.UnDelegate(config, header, delegateAddr, common.HexToAddress("0xab")))
	//undelegator ， a candidate who didn't vote.
	candidate4Addr := common.HexToAddress("0x4e080e49f62694554871e669aeb4ebe17c4a9671")
	assert.Nil(t, dposContext.BecomeCandidate(config,header,candidate4Addr))
	assert.NotNil(t, dposContext.UnDelegate(config, header, delegateAddr, candidate1Addr))
}

func TestDposContextMultiVoteHasVoted(t *testing.T) {
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

func TestDposContextMultiVoteExceedMaxVoteCandidateNum(t *testing.T) {
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

func TestDposContextMultiVoteInvalidCandidate(t *testing.T) {
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
	err = dposContext.UnDelegate(config, header, delegateAddr, candidateAddr)
	expectedErr = errors.New(candidateAddr.String() + " is invalid candidate")
	assert.EqualError(t,expectedErr,err.Error())
}
func TestDposContextBecomeCandidateAndKickoutCandidate(t *testing.T) {
	//MultiVote
	config := &params.ChainConfig{
		MultiVoteBlock: big.NewInt(0),
	}
	header := &Header{
		Number: big.NewInt(1000),
	}
	checkDposContextBecomeCandidateAndKickoutCandidate(t, config, header)

	// no multiVote
	config = &params.ChainConfig{
		MultiVoteBlock: nil,
	}
	header = &Header{
		Number: big.NewInt(990),
	}
	checkDposContextBecomeCandidateAndKickoutCandidate(t, config, header)
}

func checkDposContextBecomeCandidateAndKickoutCandidate(t *testing.T, config *params.ChainConfig, header *Header) {
	candidates := []common.Address{
		common.HexToAddress("0x44d1ce0b7cb3588bca96151fe1bc05af38f91b6e"),
		common.HexToAddress("0xa60a3886b552ff9992cfcd208ec1152079e046c2"),
		common.HexToAddress("0x4e080e49f62694554871e669aeb4ebe17c4a9670"),
	}
	delegateAddress1 :=common.HexToAddress("0x4e080e49f62694554871e669aeb4ebe17c4a9671")
	delegateAddress2 :=common.HexToAddress("0x4e080e49f62694554871e669aeb4ebe17c4a9672")
	db := socdb.NewMemDatabase()
	dposContext, err := NewDposContext(db)
	assert.Nil(t, err)

	//BecomeCandidate
	for _, candidate := range candidates {
		assert.Nil(t, dposContext.BecomeCandidate(config, header, candidate))
		if config.IsMultiVote(header.Number) {
			assert.Nil(t, dposContext.Delegate(config, header, delegateAddress1, candidate))
			assert.Nil(t, dposContext.Delegate(config, header, delegateAddress2, candidate))
		}
		assert.Nil(t, dposContext.Delegate(config, header, candidate, candidate))
	}

	candidateMap := map[common.Address]bool{}
	candidateIter := trie.NewIterator(dposContext.candidateTrie.NodeIterator(nil))
	for candidateIter.Next() {
		candidateMap[common.BytesToAddress(candidateIter.Value)] = true
	}
	voteIter := trie.NewIterator(dposContext.voteTrie.NodeIterator(nil))
	voteMap := map[common.Address]bool{}
	for voteIter.Next() {
		voteMap[common.BytesToAddress(voteIter.Value)] = true
	}
	for _, candidate := range candidates {
		delegateIter := trie.NewIterator(dposContext.delegateTrie.NodeIterator(candidate.Bytes()))
		assert.True(t, delegateIter.Next())
		assert.True(t, candidateMap[candidate])
		assert.True(t, voteMap[candidate])
	}

	// Repeated BecomeCandidate have no effect
	for _, candidate := range candidates {
		assert.Nil(t, dposContext.BecomeCandidate(config, header, candidate))
	}
	candidateMap = map[common.Address]bool{}
	candidateIter = trie.NewIterator(dposContext.candidateTrie.NodeIterator(nil))
	for candidateIter.Next() {
		candidateMap[common.BytesToAddress(candidateIter.Value)] = true
	}
	voteIter = trie.NewIterator(dposContext.voteTrie.NodeIterator(nil))
	voteMap = map[common.Address]bool{}
	for voteIter.Next() {
		voteMap[common.BytesToAddress(voteIter.Value)] = true
	}

	for _, candidate := range candidates {
		delegateIter := trie.NewIterator(dposContext.delegateTrie.NodeIterator(candidate.Bytes()))
		assert.True(t, delegateIter.Next())
		assert.True(t, candidateMap[candidate])
		assert.True(t, voteMap[candidate])
	}


	//KickoutCandidate
	oldValidators,_ := dposContext.GetValidators()
	kickIdx := 1
	assert.Nil(t, dposContext.KickoutCandidate(config, header, candidates[kickIdx]))
	candidateMap = map[common.Address]bool{}
	candidateIter = trie.NewIterator(dposContext.candidateTrie.NodeIterator(nil))
	for candidateIter.Next() {
		candidateMap[common.BytesToAddress(candidateIter.Value)] = true
	}
	voteIter = trie.NewIterator(dposContext.voteTrie.NodeIterator(nil))
	voteMap = map[common.Address]bool{}
	for voteIter.Next() {
		voteMap[common.BytesToAddress(voteIter.Value)] = true
	}
	for i, candidate := range candidates {
		delegateIter := trie.NewIterator(dposContext.delegateTrie.NodeIterator(candidate.Bytes()))
		if i == kickIdx {
			assert.False(t, delegateIter.Next())
			assert.False(t, candidateMap[candidate])
			assert.False(t, voteMap[candidate])
			continue
		}
		assert.True(t, delegateIter.Next())
		assert.True(t, candidateMap[candidate])
		assert.True(t, voteMap[candidate])
	}

	// KickoutCandidate  does not affect validators
	newValidators,_ := dposContext.GetValidators()
	assert.Equal(t,oldValidators,newValidators)

	// Repeated kicks have no effect
	oldValidators,_ = dposContext.GetValidators()
	assert.Nil(t, dposContext.KickoutCandidate(config, header, candidates[kickIdx]))
	candidateMap = map[common.Address]bool{}
	candidateIter = trie.NewIterator(dposContext.candidateTrie.NodeIterator(nil))
	for candidateIter.Next() {
		candidateMap[common.BytesToAddress(candidateIter.Value)] = true
	}
	voteIter = trie.NewIterator(dposContext.voteTrie.NodeIterator(nil))
	voteMap = map[common.Address]bool{}
	for voteIter.Next() {
		voteMap[common.BytesToAddress(voteIter.Value)] = true
	}
	for i, candidate := range candidates {
		delegateIter := trie.NewIterator(dposContext.delegateTrie.NodeIterator(candidate.Bytes()))
		if i == kickIdx {
			assert.False(t, delegateIter.Next())
			assert.False(t, candidateMap[candidate])
			assert.False(t, voteMap[candidate])
			continue
		}
		assert.True(t, delegateIter.Next())
		assert.True(t, candidateMap[candidate])
		assert.True(t, voteMap[candidate])
	}

	// KickoutCandidate  does not affect validators
	newValidators,_ = dposContext.GetValidators()
	assert.Equal(t,oldValidators,newValidators)
}


var (
	MockEpoch = []string{
		"0x44d1ce0b7cb3588bca96151fe1bc05af38f91b6e",
		"0xa60a3886b552ff9992cfcd208ec1152079e046c2",
		"0x4e080e49f62694554871e669aeb4ebe17c4a9670",
		"0xb040353ec0f2c113d5639444f7253681aecda1f8",
		"0x14432e15f21237013017fa6ee90fc99433dec82c",
		"0x9f30d0e5c9c88cade54cd1adecf6bc2c7e0e5af6",
		"0xd83b44a3719720ec54cdb9f54c0202de68f1ebcb",
		"0x56cc452e450551b7b9cffe25084a069e8c1e9441",
		"0xbcfcb3fa8250be4f2bf2b1e70e1da500c668377b",
		"0x9d9667c71bb09d6ca7c3ed12bfe5e7be24e2ffe1",
		"0xabde197e97398864ba74511f02832726edad5967",
		"0x6f99d97a394fa7a623fdf84fdc7446b99c3cb335",
		"0xf78b011e639ce6d8b76f97712118f3fe4a12dd95",
		"0x8db3b6c801dddd624d6ddc2088aa64b5a2493661",
		"0x751b484bd5296f8d267a8537d33f25a848f7f7af",
		"0x646ba1fa42eb940aac67103a71e9a908ef484ec3",
		"0x34d4a8d9f6b53a8f5e674516cb8ad66c843b2801",
		"0x5b76fff970bf8a351c1c9ebfb5e5a9493e956ddd",
		"0x8da3c5aedaf106c61cfee6d8483e1f255fdd60c0",
		"0x2cdbe87a1bd7ee60dd6fe97f7b2d1efbacd5d95d",
		"0x743415d0e979dc6e426bc8189e40beb65bf5ac1d",
	}
)

func mockNewDposContext(db socdb.Database) *DposContext {

	config := &params.ChainConfig{
		MultiVoteBlock:      big.NewInt(0),
	}
	header := &Header{
		Number: big.NewInt(1000),
	}

	dposContext, err := NewDposContextFromProto(db, &DposContextProto{})
	if err != nil {
		return nil
	}
	delegator := []byte{}
	candidate := []byte{}
	addresses := []common.Address{}
	for i := 0; i < MaxValidatorSize; i++ {
		addresses = append(addresses, common.HexToAddress(MockEpoch[i]))
	}
	dposContext.SetValidators(addresses)
	for j := 0; j < len(MockEpoch); j++ {
		delegator = common.HexToAddress(MockEpoch[j]).Bytes()
		candidate = common.HexToAddress(MockEpoch[j]).Bytes()
		dposContext.BecomeCandidate(config,header,common.BytesToAddress(candidate))
		dposContext.Delegate(config,header,common.BytesToAddress(delegator),common.BytesToAddress(candidate))
	}
	return dposContext
}

func setMintCntTrie(epochID int64, candidate common.Address, mintCntTrie *trie.Trie, count int64) {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(epochID))
	cntBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(cntBytes, uint64(count))
	mintCntTrie.TryUpdate(append(key, candidate.Bytes()...), cntBytes)
}

func getMintCnt(epochID int64, candidate common.Address, mintCntTrie *trie.Trie) int64 {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(epochID))
	cntBytes := mintCntTrie.Get(append(key, candidate.Bytes()...))
	if cntBytes == nil {
		return 0
	} else {
		return int64(binary.BigEndian.Uint64(cntBytes))
	}
}

func TestDposContextUpdateMintCnt(t *testing.T) {
	db := socdb.NewMemDatabase()
	dposContext := mockNewDposContext(db)

	// new block still in the same epoch with current block, but newMiner is the first time to mint in the epoch
	lastTime := int64(EpochInterval)

	miner := common.HexToAddress("0xa60a3886b552ff9992cfcd208ec1152079e046c2")
	blockTime := int64(EpochInterval + BlockInterval)

	beforeUpdateCnt := getMintCnt(blockTime/EpochInterval, miner, dposContext.MintCntTrie())
	//updateMintCnt(lastTime, blockTime, miner, dposContext)
	dposContext.UpdateMintCnt(lastTime, blockTime, miner, EpochInterval)
	afterUpdateCnt := getMintCnt(blockTime/EpochInterval, miner, dposContext.MintCntTrie())
	assert.Equal(t, int64(0), beforeUpdateCnt)
	assert.Equal(t, int64(1), afterUpdateCnt)

	// new block still in the same epoch with current block, and newMiner has mint block before in the epoch
	setMintCntTrie(blockTime/EpochInterval, miner, dposContext.MintCntTrie(), int64(1))

	blockTime = EpochInterval + BlockInterval*4

	// currentBlock has recorded the count for the newMiner before UpdateMintCnt
	beforeUpdateCnt = getMintCnt(blockTime/EpochInterval, miner, dposContext.MintCntTrie())
	//updateMintCnt(lastTime, blockTime, miner, dposContext)
	dposContext.UpdateMintCnt(lastTime, blockTime, miner, EpochInterval)
	afterUpdateCnt = getMintCnt(blockTime/EpochInterval, miner, dposContext.MintCntTrie())
	assert.Equal(t, int64(1), beforeUpdateCnt)
	assert.Equal(t, int64(2), afterUpdateCnt)

	// new block come to a new epoch
	blockTime = EpochInterval * 2

	beforeUpdateCnt = getMintCnt(blockTime/EpochInterval, miner, dposContext.MintCntTrie())
	//updateMintCnt(lastTime, blockTime, miner, dposContext)
	dposContext.UpdateMintCnt(lastTime, blockTime, miner, EpochInterval)
	afterUpdateCnt = getMintCnt(blockTime/EpochInterval, miner, dposContext.MintCntTrie())
	assert.Equal(t, int64(0), beforeUpdateCnt)
	assert.Equal(t, int64(1), afterUpdateCnt)
}

func TestDposContextSuffleValidators(t *testing.T) {
	MockAddress := []string{
		"0x44d1ce0b7cb3588bca96151fe1bc05af38f91b6e",
		"0xa60a3886b552ff9992cfcd208ec1152079e046c2",
		"0x4e080e49f62694554871e669aeb4ebe17c4a9670",
		"0xb040353ec0f2c113d5639444f7253681aecda1f8",
		"0x14432e15f21237013017fa6ee90fc99433dec82c",
		"0x9f30d0e5c9c88cade54cd1adecf6bc2c7e0e5af6",
		"0xd83b44a3719720ec54cdb9f54c0202de68f1ebcb",
		"0x56cc452e450551b7b9cffe25084a069e8c1e9441",
		"0xbcfcb3fa8250be4f2bf2b1e70e1da500c668377b",
		"0x9d9667c71bb09d6ca7c3ed12bfe5e7be24e2ffe1",
		"0xabde197e97398864ba74511f02832726edad5967",
		"0x6f99d97a394fa7a623fdf84fdc7446b99c3cb335",
		"0xf78b011e639ce6d8b76f97712118f3fe4a12dd95",
		"0x8db3b6c801dddd624d6ddc2088aa64b5a2493661",
		"0x751b484bd5296f8d267a8537d33f25a848f7f7af",
		"0x646ba1fa42eb940aac67103a71e9a908ef484ec3",
		"0x34d4a8d9f6b53a8f5e674516cb8ad66c843b2801",
		"0x5b76fff970bf8a351c1c9ebfb5e5a9493e956ddd",
		"0x8da3c5aedaf106c61cfee6d8483e1f255fdd60c0",
		"0x2cdbe87a1bd7ee60dd6fe97f7b2d1efbacd5d95d",
		"0x743415d0e979dc6e426bc8189e40beb65bf5ac1d",
	}

	db := socdb.NewMemDatabase()
	dposContext, err := NewDposContext(db)
	assert.Nil(t, err)
	validators := []common.Address{}
	for i :=0; i < MaxValidatorSize && i < len(MockAddress); i++ {
		validators = append(validators,common.HexToAddress(MockAddress[i]))
	}
	assert.Nil(t, dposContext.SetValidators(validators))

	// in same LoopInterval
	header := &Header{
		Time: big.NewInt( LoopInterval - BlockInterval),
	}
	parent := &Header{
		Time: big.NewInt(LoopInterval - BlockInterval*2 ),
	}
	assert.Nil(t, dposContext.SuffleValidators(header,parent))
	suffledValidators , err := dposContext.GetValidators()
	assert.Equal(t, validators,suffledValidators)

	// in different  LoopInterval
	header = &Header{
		Time: big.NewInt( LoopInterval - BlockInterval),
	}
	parent = &Header{
		Time: big.NewInt(LoopInterval + BlockInterval ),
	}
	assert.Nil(t, dposContext.SuffleValidators(header,parent))
	suffledValidators , err = dposContext.GetValidators()
	assert.NotEqual(t, validators,suffledValidators)
}
