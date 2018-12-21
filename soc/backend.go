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

// Package eth implements the Ethereum protocol.
package soc

import (
	"errors"
	"fmt"
	"math/big"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/allsportschain/go-allsportschain/accounts"
	"github.com/allsportschain/go-allsportschain/common"
	"github.com/allsportschain/go-allsportschain/common/hexutil"
	"github.com/allsportschain/go-allsportschain/consensus"
	"github.com/allsportschain/go-allsportschain/consensus/dpos"
	"github.com/allsportschain/go-allsportschain/consensus/sochash"
	"github.com/allsportschain/go-allsportschain/core"
	"github.com/allsportschain/go-allsportschain/core/bloombits"
	"github.com/allsportschain/go-allsportschain/core/rawdb"
	"github.com/allsportschain/go-allsportschain/core/types"
	"github.com/allsportschain/go-allsportschain/core/vm"
	"github.com/allsportschain/go-allsportschain/soc/downloader"
	"github.com/allsportschain/go-allsportschain/soc/filters"
	"github.com/allsportschain/go-allsportschain/soc/gasprice"
	"github.com/allsportschain/go-allsportschain/socdb"
	"github.com/allsportschain/go-allsportschain/event"
	"github.com/allsportschain/go-allsportschain/internal/socapi"
	"github.com/allsportschain/go-allsportschain/log"
	"github.com/allsportschain/go-allsportschain/miner"
	"github.com/allsportschain/go-allsportschain/node"
	"github.com/allsportschain/go-allsportschain/p2p"
	"github.com/allsportschain/go-allsportschain/params"
	"github.com/allsportschain/go-allsportschain/rlp"
	"github.com/allsportschain/go-allsportschain/rpc"
)

type LesServer interface {
	Start(srvr *p2p.Server)
	Stop()
	Protocols() []p2p.Protocol
	SetBloomBitsIndexer(bbIndexer *core.ChainIndexer)
}

// Allsportschain implements the Allsportschain full node service.
type Allsportschain struct {
	config      *Config
	chainConfig *params.ChainConfig

	// Channel for shutting down the service
	shutdownChan chan bool // Channel for shutting down the Allsportschain

	// Handlers
	txPool          *core.TxPool
	blockchain      *core.BlockChain
	protocolManager *ProtocolManager
	lesServer       LesServer

	// DB interfaces
	chainDb socdb.Database // Block chain database

	eventMux       *event.TypeMux
	engine         consensus.Engine
	accountManager *accounts.Manager

	bloomRequests chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer  *core.ChainIndexer             // Bloom indexer operating during block imports

	APIBackend *SocAPIBackend

	miner     *miner.Miner
	gasPrice  *big.Int
	validator common.Address
	socerbase common.Address

	networkID     uint64
	netRPCService *socapi.PublicNetAPI

	lock sync.RWMutex // Protects the variadic fields (e.g. gas price and socerbase)
}

func (s *Allsportschain) AddLesServer(ls LesServer) {
	s.lesServer = ls
	ls.SetBloomBitsIndexer(s.bloomIndexer)
}

// New creates a new Allsportschain object (including the
// initialisation of the common Allsportschain object)
func New(ctx *node.ServiceContext, config *Config) (*Allsportschain, error) {
	if config.SyncMode == downloader.LightSync {
		return nil, errors.New("can't run soc.Allsportschain in light sync mode, use les.LightAllsportschain")
	}
	if !config.SyncMode.IsValid() {
		return nil, fmt.Errorf("invalid sync mode %d", config.SyncMode)
	}
	chainDb, err := CreateDB(ctx, config, "chaindata")
	if err != nil {
		return nil, err
	}
	chainConfig, genesisHash, genesisErr := core.SetupGenesisBlock(chainDb, config.Genesis)
	if _, ok := genesisErr.(*params.ConfigCompatError); genesisErr != nil && !ok {
		return nil, genesisErr
	}
	log.Info("Initialised chain configuration", "config", chainConfig)

	soc := &Allsportschain{
		config:         config,
		chainDb:        chainDb,
		chainConfig:    chainConfig,
		eventMux:       ctx.EventMux,
		accountManager: ctx.AccountManager,
		engine:         CreateConsensusEngine(ctx, &config.Sochash, chainConfig, chainDb),
		shutdownChan:   make(chan bool),
		networkID:      config.NetworkId,
		gasPrice:       config.GasPrice,
		socerbase:      config.Socerbase,
		bloomRequests:  make(chan chan *bloombits.Retrieval),
		bloomIndexer:   NewBloomIndexer(chainDb, params.BloomBitsBlocks),
	}

	soc.setAuthorize()
	log.Info("Initialising Allsportschain protocol", "versions", ProtocolVersions, "network", config.NetworkId)

	if !config.SkipBcVersionCheck {
		bcVersion := rawdb.ReadDatabaseVersion(chainDb)
		if bcVersion != core.BlockChainVersion && bcVersion != 0 {
			return nil, fmt.Errorf("Blockchain DB version mismatch (%d / %d). Run gsoc upgradedb.\n", bcVersion, core.BlockChainVersion)
		}
		rawdb.WriteDatabaseVersion(chainDb, core.BlockChainVersion)
	}
	var (
		vmConfig    = vm.Config{EnablePreimageRecording: config.EnablePreimageRecording}
		cacheConfig = &core.CacheConfig{Disabled: config.NoPruning, TrieNodeLimit: config.TrieCache, TrieTimeLimit: config.TrieTimeout}
	)
	soc.blockchain, err = core.NewBlockChain(chainDb, cacheConfig, soc.chainConfig, soc.engine, vmConfig)
	if err != nil {
		return nil, err
	}
	// Rewind the chain in case of an incompatible config upgrade.
	if compat, ok := genesisErr.(*params.ConfigCompatError); ok {
		log.Warn("Rewinding chain to upgrade configuration", "err", compat)
		soc.blockchain.SetHead(compat.RewindTo)
		rawdb.WriteChainConfig(chainDb, genesisHash, chainConfig)
	}
	soc.bloomIndexer.Start(soc.blockchain)

	if config.TxPool.Journal != "" {
		config.TxPool.Journal = ctx.ResolvePath(config.TxPool.Journal)
	}
	soc.txPool = core.NewTxPool(config.TxPool, soc.chainConfig, soc.blockchain)

	if soc.protocolManager, err = NewProtocolManager(soc.chainConfig, config.SyncMode, config.NetworkId, soc.eventMux, soc.txPool, soc.engine, soc.blockchain, chainDb); err != nil {
		return nil, err
	}
	soc.chainConfig.IpcPath = fmt.Sprintf("%s/%s",config.DataDir, config.IPCPath)
	soc.chainConfig.NetworkId = config.NetworkId
	soc.miner = miner.New(soc, soc.chainConfig, soc.EventMux(), soc.engine)
	soc.miner.SetExtra(makeExtraData(config.ExtraData))

	soc.APIBackend = &SocAPIBackend{soc, nil}
	gpoParams := config.GPO
	if gpoParams.Default == nil {
		gpoParams.Default = config.GasPrice
	}
	soc.APIBackend.gpo = gasprice.NewOracle(soc.APIBackend, gpoParams)

	return soc, nil
}

func makeExtraData(extra []byte) []byte {
	if len(extra) == 0 {
		// create default extradata
		extra, _ = rlp.EncodeToBytes([]interface{}{
			uint(params.VersionMajor<<16 | params.VersionMinor<<8 | params.VersionPatch),
			"gsoc",
			runtime.Version(),
			runtime.GOOS,
		})
	}
	if uint64(len(extra)) > params.MaximumExtraDataSize {
		log.Warn("Miner extra data exceed limit", "extra", hexutil.Bytes(extra), "limit", params.MaximumExtraDataSize)
		extra = nil
	}
	return extra
}

// CreateDB creates the chain database.
func CreateDB(ctx *node.ServiceContext, config *Config, name string) (socdb.Database, error) {
	db, err := ctx.OpenDatabase(name, config.DatabaseCache, config.DatabaseHandles)
	if err != nil {
		return nil, err
	}
	if db, ok := db.(*socdb.LDBDatabase); ok {
		db.Meter("soc/db/chaindata/")
	}
	return db, nil
}

// CreateConsensusEngine creates the required type of consensus engine instance for an Allsportschain service
func CreateConsensusEngine(ctx *node.ServiceContext, config *sochash.Config, chainConfig *params.ChainConfig, db socdb.Database) consensus.Engine {
	// If proof-of-authority is requested, set it up
	if chainConfig.Dpos != nil {
		return dpos.New(chainConfig.Dpos, db)
	}
	// Otherwise assume proof-of-work
	switch config.PowMode {
	case sochash.ModeFake:
		log.Warn("Sochash used in fake mode")
		return sochash.NewFaker()
	case sochash.ModeTest:
		log.Warn("Sochash used in test mode")
		return sochash.NewTester()
	case sochash.ModeShared:
		log.Warn("Sochash used in shared mode")
		return sochash.NewShared()
	default:
		engine := sochash.New(sochash.Config{
			CacheDir:       ctx.ResolvePath(config.CacheDir),
			CachesInMem:    config.CachesInMem,
			CachesOnDisk:   config.CachesOnDisk,
			DatasetDir:     config.DatasetDir,
			DatasetsInMem:  config.DatasetsInMem,
			DatasetsOnDisk: config.DatasetsOnDisk,
		})
		engine.SetThreads(-1) // Disable CPU mining
		return engine
	}
}

// APIs return the collection of RPC services the socereum package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
func (s *Allsportschain) APIs() []rpc.API {
	apis := socapi.GetAPIs(s.APIBackend)

	// Append any APIs exposed explicitly by the consensus engine
	apis = append(apis, s.engine.APIs(s.BlockChain())...)

	// Append all the local APIs and return
	return append(apis, []rpc.API{
		{
			Namespace: "soc",
			Version:   "1.0",
			Service:   NewPublicAllsportschainAPI(s),
			Public:    true,
		}, {
			Namespace: "soc",
			Version:   "1.0",
			Service:   NewPublicMinerAPI(s),
			Public:    true,
		}, {
			Namespace: "soc",
			Version:   "1.0",
			Service:   downloader.NewPublicDownloaderAPI(s.protocolManager.downloader, s.eventMux),
			Public:    true,
		}, {
			Namespace: "miner",
			Version:   "1.0",
			Service:   NewPrivateMinerAPI(s),
			Public:    false,
		}, {
			Namespace: "soc",
			Version:   "1.0",
			Service:   filters.NewPublicFilterAPI(s.APIBackend, false),
			Public:    true,
		}, {
			Namespace: "admin",
			Version:   "1.0",
			Service:   NewPrivateAdminAPI(s),
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPublicDebugAPI(s),
			Public:    true,
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPrivateDebugAPI(s.chainConfig, s),
		}, {
			Namespace: "net",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		},
	}...)
}

func (s *Allsportschain) ResetWithGenesisBlock(gb *types.Block) {
	s.blockchain.ResetWithGenesisBlock(gb)
}

func (s *Allsportschain) Socerbase() (eb common.Address, err error) {
	s.lock.RLock()
	socerbase := s.socerbase
	s.lock.RUnlock()

	if socerbase != (common.Address{}) {
		return socerbase, nil
	}
	if wallets := s.AccountManager().Wallets(); len(wallets) > 0 {
		if accounts := wallets[0].Accounts(); len(accounts) > 0 {
			socerbase := accounts[0].Address

			s.lock.Lock()
			s.socerbase = socerbase
			s.lock.Unlock()

			log.Info("Socerbase automatically configured", "address", socerbase)
			return socerbase, nil
		}
	}
	return common.Address{}, fmt.Errorf("socerbase must be explicitly specified")
}

// SetSocerbase sets the mining reward address.
func (s *Allsportschain) SetSocerbase(socerbase common.Address) {
	s.lock.Lock()
	s.socerbase = socerbase
	s.lock.Unlock()

	s.miner.SetSocerbase(socerbase)
}

func (s *Allsportschain) Validator() (validator common.Address, err error) {
	s.lock.RLock()
	validator = s.validator
	s.lock.RUnlock()

	if validator != (common.Address{}) {
		return validator, nil
	}
	if wallets := s.AccountManager().Wallets(); len(wallets) > 0 {
		if accounts := wallets[0].Accounts(); len(accounts) > 0 {
			return accounts[0].Address, nil
		}
	}
	return common.Address{}, fmt.Errorf("validator address must be explicitly specified")
}

// set in js console via admin interface or wrapper from cli flags
func (s *Allsportschain) SetValidator(validator common.Address) {
	s.lock.Lock()
	s.validator = validator
	s.lock.Unlock()
}

func (s *Allsportschain) setAuthorize() error{
	validator, err := s.Validator()
	if err != nil {
		log.Info("Cannot start mining without validator", "err", err)
		return fmt.Errorf("validator missing: %v", err)
	}

	if dpos, ok := s.engine.(*dpos.Dpos); ok {
		wallet, err := s.accountManager.Find(accounts.Account{Address: validator})
		if wallet == nil || err != nil {
			log.Error("validator account unavailable locally", "err", err)
			return fmt.Errorf("signer missing: %v", err)
		}
		dpos.Authorize(validator, wallet.SignHash)
	}
	return nil
}

func (s *Allsportschain) StartMining(local bool) error {
	sb, err := s.Socerbase()
	if err != nil {
		log.Error("Cannot start mining without socerbase", "err", err)
		return fmt.Errorf("socerbase missing: %v", err)
	}

	err = s.setAuthorize()
	if err != nil {
		log.Error("Cannot start mining without set Authorize", "err", err)
		return err
	}

	if local {
		// If local (CPU) mining is started, we can disable the transaction rejection
		// mechanism introduced to speed sync times. CPU mining on mainnet is ludicrous
		// so none will ever hit this path, whereas marking sync done on CPU mining
		// will ensure that private networks work in single miner mode too.
		atomic.StoreUint32(&s.protocolManager.acceptTxs, 1)
	}
	go s.miner.Start(sb)
	return nil
}

func (s *Allsportschain) StopMining()         { s.miner.Stop() }
func (s *Allsportschain) IsMining() bool      { return s.miner.Mining() }
func (s *Allsportschain) Miner() *miner.Miner { return s.miner }

func (s *Allsportschain) AccountManager() *accounts.Manager  { return s.accountManager }
func (s *Allsportschain) BlockChain() *core.BlockChain       { return s.blockchain }
func (s *Allsportschain) TxPool() *core.TxPool               { return s.txPool }
func (s *Allsportschain) EventMux() *event.TypeMux           { return s.eventMux }
func (s *Allsportschain) Engine() consensus.Engine           { return s.engine }
func (s *Allsportschain) ChainDb() socdb.Database            { return s.chainDb }
func (s *Allsportschain) IsListening() bool                  { return true } // Always listening
func (s *Allsportschain) SocVersion() int                    { return int(s.protocolManager.SubProtocols[0].Version) }
func (s *Allsportschain) NetVersion() uint64                 { return s.networkID }
func (s *Allsportschain) Downloader() *downloader.Downloader { return s.protocolManager.downloader }

// Protocols implements node.Service, returning all the currently configured
// network protocols to start.
func (s *Allsportschain) Protocols() []p2p.Protocol {
	if s.lesServer == nil {
		return s.protocolManager.SubProtocols
	}
	return append(s.protocolManager.SubProtocols, s.lesServer.Protocols()...)
}

// Start implements node.Service, starting all internal goroutines needed by the
// Allsportschain protocol implementation.
func (s *Allsportschain) Start(srvr *p2p.Server) error {
	// Start the bloom bits servicing goroutines
	s.startBloomHandlers()

	// Start the RPC service
	s.netRPCService = socapi.NewPublicNetAPI(srvr, s.NetVersion())

	// Figure out a max peers count based on the server limits
	maxPeers := srvr.MaxPeers
	if s.config.LightServ > 0 {
		if s.config.LightPeers >= srvr.MaxPeers {
			return fmt.Errorf("invalid peer config: light peer count (%d) >= total peer count (%d)", s.config.LightPeers, srvr.MaxPeers)
		}
		maxPeers -= s.config.LightPeers
	}
	// Start the networking layer and the light server if requested
	s.protocolManager.Start(maxPeers)
	if s.lesServer != nil {
		s.lesServer.Start(srvr)
	}
	return nil
}

// Stop implements node.Service, terminating all internal goroutines used by the
// Allsportschain protocol.
func (s *Allsportschain) Stop() error {
	s.bloomIndexer.Close()
	s.blockchain.Stop()
	s.protocolManager.Stop()
	if s.lesServer != nil {
		s.lesServer.Stop()
	}
	s.txPool.Stop()
	s.miner.Stop()
	s.eventMux.Stop()

	s.chainDb.Close()
	close(s.shutdownChan)

	return nil
}
