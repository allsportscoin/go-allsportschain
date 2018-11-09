// Copyright 2016 The go-ethereum Authors
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

// Package les implements the Light Ethereum Subprotocol.
package les

import (
	"fmt"
	"sync"
	"time"

	"github.com/allsportschain/go-allsportschain/accounts"
	"github.com/allsportschain/go-allsportschain/common"
	"github.com/allsportschain/go-allsportschain/common/hexutil"
	"github.com/allsportschain/go-allsportschain/consensus"
	"github.com/allsportschain/go-allsportschain/consensus/dpos"
	"github.com/allsportschain/go-allsportschain/core"
	"github.com/allsportschain/go-allsportschain/core/bloombits"
	"github.com/allsportschain/go-allsportschain/core/rawdb"
	"github.com/allsportschain/go-allsportschain/core/types"
	"github.com/allsportschain/go-allsportschain/soc"
	"github.com/allsportschain/go-allsportschain/soc/downloader"
	"github.com/allsportschain/go-allsportschain/soc/filters"
	"github.com/allsportschain/go-allsportschain/soc/gasprice"
	"github.com/allsportschain/go-allsportschain/socdb"
	"github.com/allsportschain/go-allsportschain/event"
	"github.com/allsportschain/go-allsportschain/internal/socapi"
	"github.com/allsportschain/go-allsportschain/light"
	"github.com/allsportschain/go-allsportschain/log"
	"github.com/allsportschain/go-allsportschain/node"
	"github.com/allsportschain/go-allsportschain/p2p"
	"github.com/allsportschain/go-allsportschain/p2p/discv5"
	"github.com/allsportschain/go-allsportschain/params"
	rpc "github.com/allsportschain/go-allsportschain/rpc"
)

type LightAllsportschain struct {
	config *soc.Config

	odr         *LesOdr
	relay       *LesTxRelay
	chainConfig *params.ChainConfig
	// Channel for shutting down the service
	shutdownChan chan bool
	// Handlers
	peers           *peerSet
	txPool          *light.TxPool
	blockchain      *light.LightChain
	protocolManager *ProtocolManager
	serverPool      *serverPool
	reqDist         *requestDistributor
	retriever       *retrieveManager
	// DB interfaces
	chainDb socdb.Database // Block chain database

	bloomRequests                              chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer, chtIndexer, bloomTrieIndexer *core.ChainIndexer

	ApiBackend *LesApiBackend

	eventMux       *event.TypeMux
	engine         consensus.Engine
	accountManager *accounts.Manager

	networkId     uint64
	netRPCService *socapi.PublicNetAPI

	wg sync.WaitGroup
}

func New(ctx *node.ServiceContext, config *soc.Config) (*LightAllsportschain, error) {
	chainDb, err := soc.CreateDB(ctx, config, "lightchaindata")
	if err != nil {
		return nil, err
	}
	chainConfig, genesisHash, genesisErr := core.SetupGenesisBlock(chainDb, config.Genesis)
	if _, isCompat := genesisErr.(*params.ConfigCompatError); genesisErr != nil && !isCompat {
		return nil, genesisErr
	}
	log.Info("Initialised chain configuration", "chainConfig", chainConfig, "config", config)

	peers := newPeerSet()
	quitSync := make(chan struct{})

	lsoc := &LightAllsportschain{
		config:           config,
		chainConfig:      chainConfig,
		chainDb:          chainDb,
		eventMux:         ctx.EventMux,
		peers:            peers,
		reqDist:          newRequestDistributor(peers, quitSync),
		accountManager:   ctx.AccountManager,
		engine:           dpos.New(chainConfig.Dpos, chainDb),
		shutdownChan:     make(chan bool),
		networkId:        config.NetworkId,
		bloomRequests:    make(chan chan *bloombits.Retrieval),
		bloomIndexer:     soc.NewBloomIndexer(chainDb, light.BloomTrieFrequency),
		chtIndexer:       light.NewChtIndexer(chainDb, true),
		bloomTrieIndexer: light.NewBloomTrieIndexer(chainDb, true),
	}

	lsoc.relay = NewLesTxRelay(peers, lsoc.reqDist)
	lsoc.serverPool = newServerPool(chainDb, quitSync, &lsoc.wg)
	lsoc.retriever = newRetrieveManager(peers, lsoc.reqDist, lsoc.serverPool)
	lsoc.odr = NewLesOdr(chainDb, lsoc.chtIndexer, lsoc.bloomTrieIndexer, lsoc.bloomIndexer, lsoc.retriever)
	if lsoc.blockchain, err = light.NewLightChain(lsoc.odr, lsoc.chainConfig, lsoc.engine); err != nil {
		return nil, err
	}
	lsoc.bloomIndexer.Start(lsoc.blockchain)
	// Rewind the chain in case of an incompatible config upgrade.
	if compat, ok := genesisErr.(*params.ConfigCompatError); ok {
		log.Warn("Rewinding chain to upgrade configuration", "err", compat)
		lsoc.blockchain.SetHead(compat.RewindTo)
		rawdb.WriteChainConfig(chainDb, genesisHash, chainConfig)
	}

	lsoc.txPool = light.NewTxPool(lsoc.chainConfig, lsoc.blockchain, lsoc.relay)
	if lsoc.protocolManager, err = NewProtocolManager(lsoc.chainConfig, true, ClientProtocolVersions, config.NetworkId, lsoc.eventMux, lsoc.engine, lsoc.peers, lsoc.blockchain, nil, chainDb, lsoc.odr, lsoc.relay, lsoc.serverPool, quitSync, &lsoc.wg); err != nil {
		return nil, err
	}
	lsoc.ApiBackend = &LesApiBackend{lsoc, nil}
	gpoParams := config.GPO
	if gpoParams.Default == nil {
		gpoParams.Default = config.GasPrice
	}
	lsoc.ApiBackend.gpo = gasprice.NewOracle(lsoc.ApiBackend, gpoParams)
	return lsoc, nil
}

func lesTopic(genesisHash common.Hash, protocolVersion uint) discv5.Topic {
	var name string
	switch protocolVersion {
	case lpv1:
		name = "LES"
	case lpv2:
		name = "LES2"
	default:
		panic(nil)
	}
	return discv5.Topic(name + "@" + common.Bytes2Hex(genesisHash.Bytes()[0:8]))
}

type LightDummyAPI struct{}

// Etherbase is the address that mining rewards will be send to
func (s *LightDummyAPI) Socerbase() (common.Address, error) {
	return common.Address{}, fmt.Errorf("not supported")
}

// Coinbase is the address that mining rewards will be send to (alias for Etherbase)
func (s *LightDummyAPI) Coinbase() (common.Address, error) {
	return common.Address{}, fmt.Errorf("not supported")
}

// Hashrate returns the POW hashrate
func (s *LightDummyAPI) Hashrate() hexutil.Uint {
	return 0
}

// Mining returns an indication if this node is currently mining.
func (s *LightDummyAPI) Mining() bool {
	return false
}

// APIs returns the collection of RPC services the ethereum package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
func (s *LightAllsportschain) APIs() []rpc.API {
	return append(socapi.GetAPIs(s.ApiBackend), []rpc.API{
		{
			Namespace: "soc",
			Version:   "1.0",
			Service:   &LightDummyAPI{},
			Public:    true,
		}, {
			Namespace: "soc",
			Version:   "1.0",
			Service:   downloader.NewPublicDownloaderAPI(s.protocolManager.downloader, s.eventMux),
			Public:    true,
		}, {
			Namespace: "soc",
			Version:   "1.0",
			Service:   filters.NewPublicFilterAPI(s.ApiBackend, true),
			Public:    true,
		}, {
			Namespace: "net",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		},
	}...)
}

func (s *LightAllsportschain) ResetWithGenesisBlock(gb *types.Block) {
	s.blockchain.ResetWithGenesisBlock(gb)
}

func (s *LightAllsportschain) BlockChain() *light.LightChain      { return s.blockchain }
func (s *LightAllsportschain) TxPool() *light.TxPool              { return s.txPool }
func (s *LightAllsportschain) Engine() consensus.Engine           { return s.engine }
func (s *LightAllsportschain) LesVersion() int                    { return int(s.protocolManager.SubProtocols[0].Version) }
func (s *LightAllsportschain) Downloader() *downloader.Downloader { return s.protocolManager.downloader }
func (s *LightAllsportschain) EventMux() *event.TypeMux           { return s.eventMux }

// Protocols implements node.Service, returning all the currently configured
// network protocols to start.
func (s *LightAllsportschain) Protocols() []p2p.Protocol {
	return s.protocolManager.SubProtocols
}

// Start implements node.Service, starting all internal goroutines needed by the
// Ethereum protocol implementation.
func (s *LightAllsportschain) Start(srvr *p2p.Server) error {
	s.startBloomHandlers()
	log.Warn("Light client mode is an experimental feature")
	s.netRPCService = socapi.NewPublicNetAPI(srvr, s.networkId)
	// clients are searching for the first advertised protocol in the list
	protocolVersion := AdvertiseProtocolVersions[0]
	s.serverPool.start(srvr, lesTopic(s.blockchain.Genesis().Hash(), protocolVersion))
	s.protocolManager.Start(s.config.LightPeers)
	return nil
}

// Stop implements node.Service, terminating all internal goroutines used by the
// Ethereum protocol.
func (s *LightAllsportschain) Stop() error {
	s.odr.Stop()
	if s.bloomIndexer != nil {
		s.bloomIndexer.Close()
	}
	if s.chtIndexer != nil {
		s.chtIndexer.Close()
	}
	if s.bloomTrieIndexer != nil {
		s.bloomTrieIndexer.Close()
	}
	s.blockchain.Stop()
	s.protocolManager.Stop()
	s.txPool.Stop()

	s.eventMux.Stop()

	time.Sleep(time.Millisecond * 200)
	s.chainDb.Close()
	close(s.shutdownChan)

	return nil
}
