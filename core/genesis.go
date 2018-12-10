// Copyright 2014 The go-ethereum Authors
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

package core

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/allsportschain/go-allsportschain/common"
	"github.com/allsportschain/go-allsportschain/common/hexutil"
	"github.com/allsportschain/go-allsportschain/common/math"
	"github.com/allsportschain/go-allsportschain/core/rawdb"
	"github.com/allsportschain/go-allsportschain/core/state"
	"github.com/allsportschain/go-allsportschain/core/types"
	"github.com/allsportschain/go-allsportschain/socdb"
	"github.com/allsportschain/go-allsportschain/log"
	"github.com/allsportschain/go-allsportschain/params"
	"github.com/allsportschain/go-allsportschain/rlp"
	"sort"
)

//go:generate gencodec -type Genesis -field-override genesisSpecMarshaling -out gen_genesis.go
//go:generate gencodec -type GenesisAccount -field-override genesisAccountMarshaling -out gen_genesis_account.go

var errGenesisNoConfig = errors.New("genesis has no chain configuration")

// Genesis specifies the header fields, state of a genesis block. It also defines hard
// fork switch-over blocks through the chain configuration.
type Genesis struct {
	Config     *params.ChainConfig `json:"config"`
	Nonce      uint64              `json:"nonce"`
	Timestamp  uint64              `json:"timestamp"`
	ExtraData  []byte              `json:"extraData"`
	GasLimit   uint64              `json:"gasLimit"   gencodec:"required"`
	Difficulty *big.Int            `json:"difficulty" gencodec:"required"`
	Mixhash    common.Hash         `json:"mixHash"`
	Coinbase   common.Address      `json:"coinbase"`
	Alloc      GenesisAlloc        `json:"alloc"      gencodec:"required"`

	// These fields are used for consensus tests. Please don't use them
	// in actual genesis blocks.
	Number     uint64      `json:"number"`
	GasUsed    uint64      `json:"gasUsed"`
	ParentHash common.Hash `json:"parentHash"`
}

// GenesisAlloc specifies the initial state that is part of the genesis block.
type GenesisAlloc map[common.Address]GenesisAccount

func (ga *GenesisAlloc) UnmarshalJSON(data []byte) error {
	m := make(map[common.UnprefixedAddress]GenesisAccount)
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	*ga = make(GenesisAlloc)
	for addr, a := range m {
		(*ga)[common.Address(addr)] = a
	}
	return nil
}

// GenesisAccount is an account in the state of the genesis block.
type GenesisAccount struct {
	Code       []byte                      `json:"code,omitempty"`
	Storage    map[common.Hash]common.Hash `json:"storage,omitempty"`
	Balance    *big.Int                    `json:"balance" gencodec:"required"`
	Nonce      uint64                      `json:"nonce,omitempty"`
	PrivateKey []byte                      `json:"secretKey,omitempty"` // for tests
}

// field type overrides for gencodec
type genesisSpecMarshaling struct {
	Nonce      math.HexOrDecimal64
	Timestamp  math.HexOrDecimal64
	ExtraData  hexutil.Bytes
	GasLimit   math.HexOrDecimal64
	GasUsed    math.HexOrDecimal64
	Number     math.HexOrDecimal64
	Difficulty *math.HexOrDecimal256
	Alloc      map[common.UnprefixedAddress]GenesisAccount
}

type genesisAccountMarshaling struct {
	Code       hexutil.Bytes
	Balance    *math.HexOrDecimal256
	Nonce      math.HexOrDecimal64
	Storage    map[storageJSON]storageJSON
	PrivateKey hexutil.Bytes
}

// storageJSON represents a 256 bit byte array, but allows less than 256 bits when
// unmarshaling from hex.
type storageJSON common.Hash

func (h *storageJSON) UnmarshalText(text []byte) error {
	text = bytes.TrimPrefix(text, []byte("0x"))
	if len(text) > 64 {
		return fmt.Errorf("too many hex characters in storage key/value %q", text)
	}
	offset := len(h) - len(text)/2 // pad on the left
	if _, err := hex.Decode(h[offset:], text); err != nil {
		fmt.Println(err)
		return fmt.Errorf("invalid hex storage key/value %q", text)
	}
	return nil
}

func (h storageJSON) MarshalText() ([]byte, error) {
	return hexutil.Bytes(h[:]).MarshalText()
}

// GenesisMismatchError is raised when trying to overwrite an existing
// genesis block with an incompatible one.
type GenesisMismatchError struct {
	Stored, New common.Hash
}

func (e *GenesisMismatchError) Error() string {
	return fmt.Sprintf("database already contains an incompatible genesis block (have %x, new %x)", e.Stored[:8], e.New[:8])
}

// SetupGenesisBlock writes or updates the genesis block in db.
// The block that will be used is:
//
//                          genesis == nil       genesis != nil
//                       +------------------------------------------
//     db has no genesis |  main-net default  |  genesis
//     db has genesis    |  from DB           |  genesis (if compatible)
//
// The stored chain configuration will be updated if it is compatible (i.e. does not
// specify a fork block below the local head block). In case of a conflict, the
// error is a *params.ConfigCompatError and the new, unwritten config is returned.
//
// The returned chain configuration is never nil.
func SetupGenesisBlock(db socdb.Database, genesis *Genesis) (*params.ChainConfig, common.Hash, error) {
	if genesis != nil && genesis.Config == nil {
		return params.AllSochashProtocolChanges, common.Hash{}, errGenesisNoConfig
	}

	// Just commit the new block if there is no stored genesis block.
	stored := rawdb.ReadCanonicalHash(db, 0)
	if (stored == common.Hash{}) {
		if genesis == nil {
			log.Info("Writing default main-net genesis block")
			genesis = DefaultGenesisBlock()
		} else {
			log.Info("Writing custom genesis block")
		}
		block, err := genesis.Commit(db)
		return genesis.Config, block.Hash(), err
	}

	// Check whether the genesis block is already written.
	if genesis != nil {
		hash := genesis.ToBlock(nil).Hash()
		if hash != stored {
			return genesis.Config, hash, &GenesisMismatchError{stored, hash}
		}
	}

	// Get the existing chain configuration.
	log.Info("stored is "+stored.String())
	newcfg := genesis.configOrDefault(stored)
	storedcfg := rawdb.ReadChainConfig(db, stored)
	if storedcfg == nil {
		log.Warn("Found genesis block without chain config")
		rawdb.WriteChainConfig(db, stored, newcfg)
		return newcfg, stored, nil
	}
	// Special case: don't change the existing config of a non-mainnet chain if no new
	// config is supplied. These chains would get AllProtocolChanges (and a compat error)
	// if we just continued here.
	if genesis == nil && stored != params.MainnetGenesisHash {
		return storedcfg, stored, nil
	}

	// Check config compatibility and write the config. Compatibility errors
	// are returned to the caller unless we're already at block zero.
	height := rawdb.ReadHeaderNumber(db, rawdb.ReadHeadHeaderHash(db))
	if height == nil {
		return newcfg, stored, fmt.Errorf("missing block number for head header hash")
	}
	compatErr := storedcfg.CheckCompatible(newcfg, *height)
	if compatErr != nil && *height != 0 && compatErr.RewindTo != 0 {
		return newcfg, stored, compatErr
	}
	rawdb.WriteChainConfig(db, stored, newcfg)
	return newcfg, stored, nil
}

func (g *Genesis) configOrDefault(ghash common.Hash) *params.ChainConfig {
	switch {
	case g != nil:
		log.Info("ChainConfig is genesis")
		return g.Config
	case ghash == params.MainnetGenesisHash:
		log.Info("ChainConfig is MainnetChainConfig")
		return params.MainnetChainConfig
	case ghash == params.TestnetGenesisHash:
		log.Info("ChainConfig is TestnetChainConfig")
		return params.TestnetChainConfig
	default:
		log.Info("ChainConfig is AllSochashProtocolChanges")
		return params.AllSochashProtocolChanges
	}
}

// ToBlock creates the genesis block and writes state of a genesis specification
// to the given database (or discards it if nil).
func (g *Genesis) ToBlock(db socdb.Database) *types.Block {
	if db == nil {
		db = socdb.NewMemDatabase()
	}
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(db))

	size := len(g.Alloc)
	txs := make(types.TxByNonce, 0, size)
	signer := types.MakeSigner(g.Config, big.NewInt(0))
	for addr, account := range g.Alloc {
		statedb.AddBalance(addr, account.Balance)
		statedb.SetCode(addr, account.Code)
		statedb.SetNonce(addr, uint64(0))
		for key, value := range account.Storage {
			statedb.SetState(addr, key, value)
		}
		tx := types.NewTransaction(account.Nonce, addr, account.Balance,types.Normal,1, big.NewInt(1),nil)
		tx, _ = tx.WithSignature(signer, common.Hex2Bytes("9bea4c4daac7c7c52e093e6a4c35dbbcf8856f1af7b059ba20253e70848d094f8a8fae537ce25ed8cb5af9adac3f141af69bd515bd2ba031522df09b97dd72b100"))
		txs = append(txs, tx)
	}

	sort.Sort(txs)
	root := statedb.IntermediateRoot(false)

	head := &types.Header{
		Number:     new(big.Int).SetUint64(g.Number),
		Nonce:      types.EncodeNonce(g.Nonce),
		Time:       new(big.Int).SetUint64(g.Timestamp),
		ParentHash: g.ParentHash,
		Extra:      g.ExtraData,
		GasLimit:   g.GasLimit,
		GasUsed:    g.GasUsed,
		Difficulty: g.Difficulty,
		MixDigest:  g.Mixhash,
		Coinbase:   g.Coinbase,
		Root:       root,
	}
	// add dposcontext
	dposContext := initGenesisDposContext(g,head, db)
	head.DposContext = dposContext.ToProto()

	if g.GasLimit == 0 {
		head.GasLimit = params.GenesisGasLimit
	}
	if g.Difficulty == nil {
		head.Difficulty = params.GenesisDifficulty
	}
	statedb.Commit(false)
	statedb.Database().TrieDB().Commit(root, true)

	block := types.NewBlock(head, txs, nil, nil)
	block.DposContext = dposContext

	return block
}

// Commit writes the block and state of a genesis specification to the database.
// The block is committed as the canonical head block.
func (g *Genesis) Commit(db socdb.Database) (*types.Block, error) {
	block := g.ToBlock(db)
	// add dposcontext
	if _, err := block.DposContext.CommitTo(); err != nil {
		return nil, err
	}

	if block.Number().Sign() != 0 {
		return nil, fmt.Errorf("can't commit genesis block with number > 0")
	}

	rawdb.WriteTd(db, block.Hash(), block.NumberU64(), g.Difficulty)
	rawdb.WriteBlock(db, block)
	rawdb.WriteReceipts(db, block.Hash(), block.NumberU64(), nil)
	rawdb.WriteCanonicalHash(db, block.Hash(), block.NumberU64())
	rawdb.WriteHeadBlockHash(db, block.Hash())
	rawdb.WriteHeadHeaderHash(db, block.Hash())
	rawdb.WriteTxLookupEntries(db, block)

	config := g.Config
	if config == nil {
		config = params.AllSochashProtocolChanges
	}
	rawdb.WriteChainConfig(db, block.Hash(), config)
	return block, nil
}

// MustCommit writes the genesis block and state to db, panicking on error.
// The block is committed as the canonical head block.
func (g *Genesis) MustCommit(db socdb.Database) *types.Block {
	block, err := g.Commit(db)
	if err != nil {
		panic(err)
	}
	return block
}

// GenesisBlockForTesting creates and writes a block in which addr has the given wei balance.
func GenesisBlockForTesting(db socdb.Database, addr common.Address, balance *big.Int) *types.Block {
	g := Genesis{Alloc: GenesisAlloc{addr: {Balance: balance}}}
	return g.MustCommit(db)
}

// DefaultGenesisBlock returns the Ethereum main net genesis block.
func DefaultGenesisBlock() *Genesis {
	return &Genesis{
		Config:     params.MainnetChainConfig,
		Nonce:      66,
		ExtraData:  hexutil.MustDecode("0x11bbe8db4e347b4e8c937c1c8370e4b5ed33adb3db69cbdb7a38e1e50b1b82fa"),
		GasLimit:   1342177280,
		Difficulty: big.NewInt(131072),
		//Alloc:      decodePrealloc(mainnetAllocData),
		Alloc:		defaultMainNetGennesisAlloc(),
	}
}

// DefaultTestnetGenesisBlock returns the Ropsten network genesis block.
func DefaultTestnetGenesisBlock() *Genesis {
	return &Genesis{
		Config:     params.TestnetChainConfig,
		Nonce:      66,
		ExtraData:  hexutil.MustDecode("0x11bbe8db4e347b4e8c937c1c8370e4b5ed33adb3db69cbdb7a38e1e50b1b82fa"),
		GasLimit:   1342177280,
		Difficulty: big.NewInt(131072),
		Alloc:      defaultTestNetGennesisAlloc(),
	}
}

func defaultTestNetGennesisAlloc() map[common.Address]GenesisAccount {
	alloc := map[common.Address]GenesisAccount{
		common.HexToAddress("0x90ae4a42d524506f99249e5fc10d948c4e07f441"): {Balance: big.NewInt(2).Mul(big.NewInt(1e+8),big.NewInt(1e+18)), Nonce:0},
		common.HexToAddress("0x97cf512dc01011c3e4926c80b12d55609729bc4a"): {Balance: big.NewInt(2).Mul(big.NewInt(1e+8),big.NewInt(1e+18)), Nonce:1},
		common.HexToAddress("0xaaf44b8cdb34c41b17bcdb6dedd34bd5c775f9d7"): {Balance: big.NewInt(2).Mul(big.NewInt(1e+8),big.NewInt(1e+18)), Nonce:2},
		common.HexToAddress("0x7e3a758190beba57902b5b08b59f15a102e53e67"): {Balance: big.NewInt(2).Mul(big.NewInt(1e+8),big.NewInt(1e+18)), Nonce:3},
		common.HexToAddress("0xe72239a57f06079b1c849d90a4c606e0ff1e3cad"): {Balance: big.NewInt(2).Mul(big.NewInt(1e+8),big.NewInt(1e+18)), Nonce:4},
		common.HexToAddress("0x6034094ff39f12786f8d5f45ae1ece5ec6b83064"): {Balance: big.NewInt(2).Mul(big.NewInt(1e+8),big.NewInt(1e+18)), Nonce:5},
		common.HexToAddress("0x6c18f4f165572afa4068dfa3ce537c4e22575144"): {Balance: big.NewInt(2).Mul(big.NewInt(1e+8),big.NewInt(1e+18)), Nonce:6},
		common.HexToAddress("0xf97e86587b04c6f7a033fb365a8413e2e1af1f3e"): {Balance: big.NewInt(1).Mul(big.NewInt(1e+8),big.NewInt(1e+18)), Nonce:7},

	}
	return alloc
}

func defaultMainNetGennesisAlloc() map[common.Address]GenesisAccount {
	alloc := map[common.Address]GenesisAccount{
		common.HexToAddress("0x90ae4a42d524506f99249e5fc10d948c4e07f441"): {Balance: big.NewInt(2).Mul(big.NewInt(2e+8),big.NewInt(1e+18)), Nonce:0},
		common.HexToAddress("0x97cf512dc01011c3e4926c80b12d55609729bc4a"): {Balance: big.NewInt(2).Mul(big.NewInt(2e+8),big.NewInt(1e+18)), Nonce:1},
		common.HexToAddress("0xaaf44b8cdb34c41b17bcdb6dedd34bd5c775f9d7"): {Balance: big.NewInt(2).Mul(big.NewInt(2e+8),big.NewInt(1e+18)), Nonce:2},
		common.HexToAddress("0x7e3a758190beba57902b5b08b59f15a102e53e67"): {Balance: big.NewInt(2).Mul(big.NewInt(2e+8),big.NewInt(1e+18)), Nonce:3},
		common.HexToAddress("0xe72239a57f06079b1c849d90a4c606e0ff1e3cad"): {Balance: big.NewInt(2).Mul(big.NewInt(2e+8),big.NewInt(1e+18)), Nonce:4},
		common.HexToAddress("0x6034094ff39f12786f8d5f45ae1ece5ec6b83064"): {Balance: big.NewInt(2).Mul(big.NewInt(2e+8),big.NewInt(1e+18)), Nonce:5},
		common.HexToAddress("0x6c18f4f165572afa4068dfa3ce537c4e22575144"): {Balance: big.NewInt(2).Mul(big.NewInt(2e+8),big.NewInt(1e+18)), Nonce:6},
		common.HexToAddress("0xf97e86587b04c6f7a033fb365a8413e2e1af1f3e"): {Balance: big.NewInt(1).Mul(big.NewInt(2e+8),big.NewInt(1e+18)), Nonce:7},

	}
	return alloc
}

// DefaultRinkebyGenesisBlock returns the Rinkeby network genesis block.
func DefaultRinkebyGenesisBlock() *Genesis {
	return &Genesis{
		Config:     params.RinkebyChainConfig,
		Timestamp:  1492009146,
		ExtraData:  hexutil.MustDecode("0x52657370656374206d7920617574686f7269746168207e452e436172746d616e42eb768f2244c8811c63729a21a3569731535f067ffc57839b00206d1ad20c69a1981b489f772031b279182d99e65703f0076e4812653aab85fca0f00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
		GasLimit:   4700000,
		Difficulty: big.NewInt(1),
		Alloc:      decodePrealloc(rinkebyAllocData),
	}
}

// DeveloperGenesisBlock returns the 'geth --dev' genesis block. Note, this must
// be seeded with the
func DeveloperGenesisBlock(period uint64, faucet common.Address) *Genesis {
	// Override the default period to the user requested one
	config := *params.AllCliqueProtocolChanges
	config.Clique.Period = period

	// Assemble and return the genesis with the precompiles and faucet pre-funded
	return &Genesis{
		Config:     &config,
		ExtraData:  append(append(make([]byte, 32), faucet[:]...), make([]byte, 65)...),
		GasLimit:   6283185,
		Difficulty: big.NewInt(1),
		Alloc: map[common.Address]GenesisAccount{
			common.BytesToAddress([]byte{1}): {Balance: big.NewInt(1)}, // ECRecover
			common.BytesToAddress([]byte{2}): {Balance: big.NewInt(1)}, // SHA256
			common.BytesToAddress([]byte{3}): {Balance: big.NewInt(1)}, // RIPEMD
			common.BytesToAddress([]byte{4}): {Balance: big.NewInt(1)}, // Identity
			common.BytesToAddress([]byte{5}): {Balance: big.NewInt(1)}, // ModExp
			common.BytesToAddress([]byte{6}): {Balance: big.NewInt(1)}, // ECAdd
			common.BytesToAddress([]byte{7}): {Balance: big.NewInt(1)}, // ECScalarMul
			common.BytesToAddress([]byte{8}): {Balance: big.NewInt(1)}, // ECPairing
			faucet: {Balance: new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(9))},
		},
	}
}

func decodePrealloc(data string) GenesisAlloc {
	var p []struct{ Addr, Balance *big.Int }
	if err := rlp.NewStream(strings.NewReader(data), 0).Decode(&p); err != nil {
		panic(err)
	}
	ga := make(GenesisAlloc, len(p))
	for _, account := range p {
		ga[common.BigToAddress(account.Addr)] = GenesisAccount{Balance: account.Balance}
	}
	return ga
}

func initGenesisDposContext(g *Genesis,header *types.Header, db socdb.Database) *types.DposContext {
	dc, err := types.NewDposContextFromProto(db, &types.DposContextProto{})
	if err != nil {
		return nil
	}
	if g.Config != nil && g.Config.Dpos != nil && g.Config.Dpos.Validators != nil {
		dc.SetValidators(g.Config.Dpos.Validators)
		for _, validator := range g.Config.Dpos.Validators {
			dc.BecomeCandidate(g.Config,header,validator)
			dc.Delegate(g.Config,header,validator, validator)
		}
		log.Info("Will change multi-vote on block number : "+ g.Config.MultiVoteBlock.String())
	}

	return dc
}