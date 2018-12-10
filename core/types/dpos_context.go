package types

import (
	"bytes"
	"errors"
	"github.com/allsportschain/go-allsportschain/common"
	"github.com/allsportschain/go-allsportschain/crypto/sha3"
	"github.com/allsportschain/go-allsportschain/rlp"
	"github.com/allsportschain/go-allsportschain/trie"
	"github.com/allsportschain/go-allsportschain/socdb"
	"fmt"
	"encoding/binary"
	"github.com/allsportschain/go-allsportschain/params"
)

type DposContext struct {
	epochTrie     *trie.Trie
	delegateTrie  *trie.Trie
	voteTrie      *trie.Trie
	candidateTrie *trie.Trie
	mintCntTrie   *trie.Trie

	db socdb.Database
}

var (
	epochPrefix     = []byte("epoch-")
	delegatePrefix  = []byte("delegate-")
	votePrefix      = []byte("vote-")
	candidatePrefix = []byte("candidate-")
	mintCntPrefix   = []byte("mintCnt-")
)
const (
	 MaxVoteCandidateNum = 30
)

func NewEpochTrie(root common.Hash, db socdb.Database) (*trie.Trie, error) {
	//return trie.NewTrieWithPrefix(root, epochPrefix, db)
	triedb := trie.NewDatabase(db)
	return trie.New(root, triedb)
}

func NewDelegateTrie(root common.Hash, db socdb.Database) (*trie.Trie, error) {
	//return trie.NewTrieWithPrefix(root, delegatePrefix, db)
	triedb := trie.NewDatabase(db)
	return trie.New(root, triedb)
}

func NewVoteTrie(root common.Hash, db socdb.Database) (*trie.Trie, error) {
	//return trie.NewTrieWithPrefix(root, votePrefix, db)
	triedb := trie.NewDatabase(db)
	return trie.New(root, triedb)
}

func NewCandidateTrie(root common.Hash, db socdb.Database) (*trie.Trie, error) {
	//return trie.NewTrieWithPrefix(root, candidatePrefix, db)
	triedb := trie.NewDatabase(db)
	return trie.New(root, triedb)
}

func NewMintCntTrie(root common.Hash, db socdb.Database) (*trie.Trie, error) {
	//return trie.NewTrieWithPrefix(root, mintCntPrefix, db)
	triedb := trie.NewDatabase(db)
	return trie.New(root, triedb)
}

func NewDposContext(db socdb.Database) (*DposContext, error) {
	epochTrie, err := NewEpochTrie(common.Hash{}, db)
	if err != nil {
		return nil, err
	}
	delegateTrie, err := NewDelegateTrie(common.Hash{}, db)
	if err != nil {
		return nil, err
	}
	voteTrie, err := NewVoteTrie(common.Hash{}, db)
	if err != nil {
		return nil, err
	}
	candidateTrie, err := NewCandidateTrie(common.Hash{}, db)
	if err != nil {
		return nil, err
	}
	mintCntTrie, err := NewMintCntTrie(common.Hash{}, db)
	if err != nil {
		return nil, err
	}
	return &DposContext{
		epochTrie:     epochTrie,
		delegateTrie:  delegateTrie,
		voteTrie:      voteTrie,
		candidateTrie: candidateTrie,
		mintCntTrie:   mintCntTrie,
		db:            db,
	}, nil
}

func NewDposContextFromProto(db socdb.Database, ctxProto *DposContextProto) (*DposContext, error) {
	epochTrie, err := NewEpochTrie(ctxProto.EpochHash, db)
	if err != nil {
		return nil, err
	}
	delegateTrie, err := NewDelegateTrie(ctxProto.DelegateHash, db)
	if err != nil {
		return nil, err
	}
	voteTrie, err := NewVoteTrie(ctxProto.VoteHash, db)
	if err != nil {
		return nil, err
	}
	candidateTrie, err := NewCandidateTrie(ctxProto.CandidateHash, db)
	if err != nil {
		return nil, err
	}
	mintCntTrie, err := NewMintCntTrie(ctxProto.MintCntHash, db)
	if err != nil {
		return nil, err
	}
	return &DposContext{
		epochTrie:     epochTrie,
		delegateTrie:  delegateTrie,
		voteTrie:      voteTrie,
		candidateTrie: candidateTrie,
		mintCntTrie:   mintCntTrie,
		db:            db,
	}, nil
}

func (d *DposContext) Copy() *DposContext {
	epochTrie := *d.epochTrie
	delegateTrie := *d.delegateTrie
	voteTrie := *d.voteTrie
	candidateTrie := *d.candidateTrie
	mintCntTrie := *d.mintCntTrie
	return &DposContext{
		epochTrie:     &epochTrie,
		delegateTrie:  &delegateTrie,
		voteTrie:      &voteTrie,
		candidateTrie: &candidateTrie,
		mintCntTrie:   &mintCntTrie,
	}
}

func (d *DposContext) Root() (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, d.epochTrie.Hash())
	rlp.Encode(hw, d.delegateTrie.Hash())
	rlp.Encode(hw, d.candidateTrie.Hash())
	rlp.Encode(hw, d.voteTrie.Hash())
	rlp.Encode(hw, d.mintCntTrie.Hash())
	hw.Sum(h[:0])
	return h
}

func (d *DposContext) Snapshot() *DposContext {
	return d.Copy()
}

func (d *DposContext) RevertToSnapShot(snapshot *DposContext) {
	d.epochTrie = snapshot.epochTrie
	d.delegateTrie = snapshot.delegateTrie
	d.candidateTrie = snapshot.candidateTrie
	d.voteTrie = snapshot.voteTrie
	d.mintCntTrie = snapshot.mintCntTrie
}

func (d *DposContext) FromProto(dcp *DposContextProto) error {
	var err error
	d.epochTrie, err = NewEpochTrie(dcp.EpochHash, d.db)
	if err != nil {
		return err
	}
	d.delegateTrie, err = NewDelegateTrie(dcp.DelegateHash, d.db)
	if err != nil {
		return err
	}
	d.candidateTrie, err = NewCandidateTrie(dcp.CandidateHash, d.db)
	if err != nil {
		return err
	}
	d.voteTrie, err = NewVoteTrie(dcp.VoteHash, d.db)
	if err != nil {
		return err
	}
	d.mintCntTrie, err = NewMintCntTrie(dcp.MintCntHash, d.db)
	return err
}

type DposContextProto struct {
	EpochHash     common.Hash `json:"epochRoot"        gencodec:"required"`
	DelegateHash  common.Hash `json:"delegateRoot"     gencodec:"required"`
	CandidateHash common.Hash `json:"candidateRoot"    gencodec:"required"`
	VoteHash      common.Hash `json:"voteRoot"         gencodec:"required"`
	MintCntHash   common.Hash `json:"mintCntRoot"      gencodec:"required"`
}

func (d *DposContext) ToProto() *DposContextProto {
	return &DposContextProto{
		EpochHash:     d.epochTrie.Hash(),
		DelegateHash:  d.delegateTrie.Hash(),
		CandidateHash: d.candidateTrie.Hash(),
		VoteHash:      d.voteTrie.Hash(),
		MintCntHash:   d.mintCntTrie.Hash(),
	}
}

func (p *DposContextProto) Root() (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, p.EpochHash)
	rlp.Encode(hw, p.DelegateHash)
	rlp.Encode(hw, p.CandidateHash)
	rlp.Encode(hw, p.VoteHash)
	rlp.Encode(hw, p.MintCntHash)
	hw.Sum(h[:0])
	return h
}

func (d *DposContext) KickoutCandidate(config *params.ChainConfig, header *Header, candidateAddr common.Address) error {
	candidate := candidateAddr.Bytes()
	err := d.candidateTrie.TryDelete(candidate)
	if err != nil {
		if _, ok := err.(*trie.MissingNodeError); !ok {
			return err
		}
	}
	iter := trie.NewIterator(d.delegateTrie.NodeIterator(candidate))
	for iter.NextPrefix(candidate) {
		delegator := iter.Value
		key := append(candidate, delegator...)
		err = d.delegateTrie.TryDelete(key)
		if err != nil {
			if _, ok := err.(*trie.MissingNodeError); !ok {
				return err
			}
		}
		voteKey := []byte{}
		if config.IsMultiVote(header.Number) {
			voteKey = append(delegator, candidate...)
		}else{
			voteKey = delegator
		}

		v, err := d.voteTrie.TryGet(voteKey)
		if err != nil {
			if _, ok := err.(*trie.MissingNodeError); !ok {
				return err
			}
		}
		if err == nil && bytes.Equal(v, candidate) {
			err = d.voteTrie.TryDelete(voteKey)
			if err != nil {
				if _, ok := err.(*trie.MissingNodeError); !ok {
					return err
				}
			}
		}
	}
	return nil
}

func (d *DposContext) BecomeCandidate(config *params.ChainConfig, header *Header, candidateAddr common.Address) error {
	candidate := candidateAddr.Bytes()
	return d.candidateTrie.TryUpdate(candidate, candidate)
}

func (d *DposContext) Delegate(config *params.ChainConfig, header *Header, delegatorAddr, candidateAddr common.Address) error {
	delegator, candidate := delegatorAddr.Bytes(), candidateAddr.Bytes()

	// the candidate must be candidate
	candidateInTrie, err := d.candidateTrie.TryGet(candidate)
	if err != nil {
		return err
	}
	if candidateInTrie == nil {
		return errors.New(candidateAddr.String() + " is invalid candidate")
	}

	if config.IsMultiVote(header.Number) {
		// judge candidate exists
		oldCandidate, err := d.voteTrie.TryGet(append(delegator, candidate...))
		if err != nil {
			if _, ok := err.(*trie.MissingNodeError); !ok {
				return err
			}
		}
		if oldCandidate != nil {
			return errors.New(candidateAddr.String() + " Has already been voted")
		}

		// judge count
		voteIterator := trie.NewIterator(d.voteTrie.NodeIterator(delegator))
		existVoteCount := voteIterator.NextPrefixCount(delegator)
		if existVoteCount >= MaxVoteCandidateNum {
			return errors.New(fmt.Sprintf("%v has already voted %v votes, Can't exceed %v votes.", delegatorAddr.String(), existVoteCount, MaxVoteCandidateNum))
		}

		//vote
		if err = d.delegateTrie.TryUpdate(append(candidate, delegator...), delegator); err != nil {
			return err
		}
		if err := d.voteTrie.TryUpdate(append(delegator, candidate...), candidate); err != nil {
			return err
		}
	}else{
		// delete old candidate if exists
		oldCandidate, err := d.voteTrie.TryGet(delegator)
		if err != nil {
			if _, ok := err.(*trie.MissingNodeError); !ok {
				return err
			}
		}
		if oldCandidate != nil {
			d.delegateTrie.Delete(append(oldCandidate, delegator...))
		}
		if err = d.delegateTrie.TryUpdate(append(candidate, delegator...), delegator); err != nil {
			return err
		}
		if err := d.voteTrie.TryUpdate(delegator, candidate); err != nil {
			return err
		}
	}
	return nil
}

func (d *DposContext) UnDelegate(config *params.ChainConfig, header *Header, delegatorAddr, candidateAddr common.Address) error {
	delegator, candidate := delegatorAddr.Bytes(), candidateAddr.Bytes()

	// the candidate must be candidate
	candidateInTrie, err := d.candidateTrie.TryGet(candidate)
	if err != nil {
		return err
	}
	if candidateInTrie == nil {
		return errors.New("invalid candidate to undelegate")
	}

	voteKey := []byte{}
	if config.IsMultiVote(header.Number) {
		voteKey = append(delegator, candidate...)
	}else{
		voteKey = delegator
	}

	oldCandidate, err := d.voteTrie.TryGet(voteKey)
	if err != nil {
		return err
	}
	if !bytes.Equal(candidate, oldCandidate) {
		return errors.New("mismatch candidate to undelegate")
	}

	if err = d.delegateTrie.TryDelete(append(candidate, delegator...)); err != nil {
		return err
	}
	if err = d.voteTrie.TryDelete(voteKey); err != nil {
		return err
	}

	return nil
}

//TODO js api
func (d *DposContext) Prods(config *params.ChainConfig, header *Header, delegatorAddr common.Address, candidateAddrList []common.Address) error {

	if !config.IsMultiVote(header.Number) {
		return errors.New("Don't support multi-vote before block number:"+config.MultiVoteBlock.String())
	}
	// judge count
	if len(candidateAddrList) > MaxVoteCandidateNum {
		return errors.New(fmt.Sprintf("Can't exceed %v candidates", MaxVoteCandidateNum))
	}

	// the candidate must be candidate
	for _, candidateAddr := range candidateAddrList {
		candidate := candidateAddr.Bytes()
		candidateInTrie, err := d.candidateTrie.TryGet(candidate)
		if err != nil {
			return err
		}
		if candidateInTrie == nil {
			return errors.New(candidateAddr.String() +" is not a candidate")
		}
	}

	delegator := delegatorAddr.Bytes()

	// delete old candidate
	voteIterator := trie.NewIterator(d.voteTrie.NodeIterator(delegator))
	existVote := voteIterator.NextPrefix(delegator)
	for existVote {
		d.voteTrie.Delete(voteIterator.Key)
		d.delegateTrie.Delete(append(voteIterator.Value, delegator...))
		existVote = voteIterator.NextPrefix(delegator)
	}

	// vote
	for _, candidateAddr := range candidateAddrList {
		candidate := candidateAddr.Bytes()
		if err := d.delegateTrie.TryUpdate(append(candidate, delegator...), delegator); err != nil {
			return err
		}
		if err := d.voteTrie.TryUpdate(append(delegator, candidate...), candidate); err != nil {
			return err
		}
	}
	return nil
}

func (d *DposContext) CommitTo() (*DposContextProto, error) {
//func (d *DposContext) CommitTo(dbw trie.DatabaseWriter) (*DposContextProto, error) {
	//epochRoot, err := d.epochTrie.CommitTo(dbw)
	epochRoot, err := d.epochTrie.Commit(nil)
	if err != nil {
		return nil, err
	}
	d.epochTrie.DbCommit(epochRoot, true)
	//delegateRoot, err := d.delegateTrie.CommitTo(dbw)
	delegateRoot, err := d.delegateTrie.Commit(nil)
	if err != nil {
		return nil, err
	}
	d.delegateTrie.DbCommit(delegateRoot, true)
	//voteRoot, err := d.voteTrie.CommitTo(dbw)
	voteRoot, err := d.voteTrie.Commit(nil)
	if err != nil {
		return nil, err
	}
	d.voteTrie.DbCommit(voteRoot, true)
	//candidateRoot, err := d.candidateTrie.CommitTo(dbw)
	candidateRoot, err := d.candidateTrie.Commit(nil)
	if err != nil {
		return nil, err
	}
	d.candidateTrie.DbCommit(candidateRoot, true)
	//mintCntRoot, err := d.mintCntTrie.CommitTo(dbw)
	mintCntRoot, err := d.mintCntTrie.Commit(nil)
	if err != nil {
		return nil, err
	}
	d.mintCntTrie.DbCommit(mintCntRoot, true)
	return &DposContextProto{
		EpochHash:     epochRoot,
		DelegateHash:  delegateRoot,
		VoteHash:      voteRoot,
		CandidateHash: candidateRoot,
		MintCntHash:   mintCntRoot,
	}, nil
}

func (d *DposContext) CandidateTrie() *trie.Trie          { return d.candidateTrie }
func (d *DposContext) DelegateTrie() *trie.Trie           { return d.delegateTrie }
func (d *DposContext) VoteTrie() *trie.Trie               { return d.voteTrie }
func (d *DposContext) EpochTrie() *trie.Trie              { return d.epochTrie }
func (d *DposContext) MintCntTrie() *trie.Trie            { return d.mintCntTrie }
func (d *DposContext) DB() socdb.Database                 { return d.db }
func (dc *DposContext) SetEpoch(epoch *trie.Trie)         { dc.epochTrie = epoch }
func (dc *DposContext) SetDelegate(delegate *trie.Trie)   { dc.delegateTrie = delegate }
func (dc *DposContext) SetVote(vote *trie.Trie)           { dc.voteTrie = vote }
func (dc *DposContext) SetCandidate(candidate *trie.Trie) { dc.candidateTrie = candidate }
func (dc *DposContext) SetMintCnt(mintCnt *trie.Trie)     { dc.mintCntTrie = mintCnt }

func (dc *DposContext) GetValidators() ([]common.Address, error) {
	var validators []common.Address
	key := []byte("validator")
	validatorsRLP := dc.epochTrie.Get(key)
	if err := rlp.DecodeBytes(validatorsRLP, &validators); err != nil {
		return nil, fmt.Errorf("failed to decode validators: %s", err)
	}
	return validators, nil
}

func (dc *DposContext) SetValidators(validators []common.Address) error {
	key := []byte("validator")
	validatorsRLP, err := rlp.EncodeToBytes(validators)
	if err != nil {
		return fmt.Errorf("failed to encode validators to rlp bytes: %s", err)
	}
	dc.epochTrie.Update(key, validatorsRLP)
	return nil
}


// update counts in MintCntTrie for the miner of newBlock
func  (dc *DposContext) UpdateMintCnt(parentBlockTime, currentBlockTime int64, validator common.Address, epochInterval int64 ) {
	currentMintCntTrie := dc.MintCntTrie()
	currentEpoch := parentBlockTime / epochInterval
	currentEpochBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(currentEpochBytes, uint64(currentEpoch))

	cnt := int64(1)
	newEpoch := currentBlockTime / epochInterval
	// still during the currentEpochID
	if currentEpoch == newEpoch {

		// when current is not genesis, read last count from the MintCntTrie
		cntBytes := currentMintCntTrie.Get(append(currentEpochBytes, validator.Bytes()...))

		// not the first time to mint
		if cntBytes != nil {
			cnt = int64(binary.BigEndian.Uint64(cntBytes)) + 1
		}
	}

	newCntBytes := make([]byte, 8)
	newEpochBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(newEpochBytes, uint64(newEpoch))
	binary.BigEndian.PutUint64(newCntBytes, uint64(cnt))
	dc.MintCntTrie().TryUpdate(append(newEpochBytes, validator.Bytes()...), newCntBytes)
}

//get all candidates form candidate trie
func (dc * DposContext) GetCandidates() ([]common.Address, error) {
	var candidates []common.Address
	iter := trie.NewIterator(dc.candidateTrie.NodeIterator(nil))
	for iter.Next() {
		candidates = append(candidates, common.BytesToAddress(iter.Value))
	}
	return candidates, nil
}


