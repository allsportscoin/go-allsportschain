// Copyright 2015 The go-ethereum Authors
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

package soc

import (
	"context"
	"math/big"

	"github.com/allsportschain/go-allsportschain/accounts"
	"github.com/allsportschain/go-allsportschain/common"
	"github.com/allsportschain/go-allsportschain/common/math"
	"github.com/allsportschain/go-allsportschain/core"
	"github.com/allsportschain/go-allsportschain/core/bloombits"
	"github.com/allsportschain/go-allsportschain/core/rawdb"
	"github.com/allsportschain/go-allsportschain/core/state"
	"github.com/allsportschain/go-allsportschain/core/types"
	"github.com/allsportschain/go-allsportschain/core/vm"
	"github.com/allsportschain/go-allsportschain/soc/downloader"
	"github.com/allsportschain/go-allsportschain/soc/gasprice"
	"github.com/allsportschain/go-allsportschain/socdb"
	"github.com/allsportschain/go-allsportschain/event"
	"github.com/allsportschain/go-allsportschain/params"
	"github.com/allsportschain/go-allsportschain/rpc"
)

// SocAPIBackend implements socapi.Backend for full nodes
type SocAPIBackend struct {
	soc *Allsportschain
	gpo *gasprice.Oracle
}

// ChainConfig returns the active chain configuration.
func (b *SocAPIBackend) ChainConfig() *params.ChainConfig {
	return b.soc.chainConfig
}

func (b *SocAPIBackend) CurrentBlock() *types.Block {
	return b.soc.blockchain.CurrentBlock()
}

func (b *SocAPIBackend) SetHead(number uint64) {
	b.soc.protocolManager.downloader.Cancel()
	b.soc.blockchain.SetHead(number)
}

func (b *SocAPIBackend) HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header, error) {
	// Pending block is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block := b.soc.miner.PendingBlock()
		return block.Header(), nil
	}
	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.soc.blockchain.CurrentBlock().Header(), nil
	}
	return b.soc.blockchain.GetHeaderByNumber(uint64(blockNr)), nil
}

func (b *SocAPIBackend) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
    return b.soc.blockchain.GetHeaderByHash(hash), nil
}

func (b *SocAPIBackend) BlockByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Block, error) {
	// Pending block is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block := b.soc.miner.PendingBlock()
		return block, nil
	}
	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.soc.blockchain.CurrentBlock(), nil
	}
	return b.soc.blockchain.GetBlockByNumber(uint64(blockNr)), nil
}

func (b *SocAPIBackend) StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*state.StateDB, *types.Header, error) {
	// Pending state is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block, state := b.soc.miner.Pending()
		return state, block.Header(), nil
	}
	// Otherwise resolve the block number and return its state
	header, err := b.HeaderByNumber(ctx, blockNr)
	if header == nil || err != nil {
		return nil, nil, err
	}
	stateDb, err := b.soc.BlockChain().StateAt(header.Root)
	return stateDb, header, err
}

func (b *SocAPIBackend) GetBlock(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return b.soc.blockchain.GetBlockByHash(hash), nil
}

func (b *SocAPIBackend) GetReceipts(ctx context.Context, hash common.Hash) (types.Receipts, error) {
	if number := rawdb.ReadHeaderNumber(b.soc.chainDb, hash); number != nil {
		return rawdb.ReadReceipts(b.soc.chainDb, hash, *number), nil
	}
	return nil, nil
}

func (b *SocAPIBackend) GetLogs(ctx context.Context, hash common.Hash) ([][]*types.Log, error) {
	number := rawdb.ReadHeaderNumber(b.soc.chainDb, hash)
	if number == nil {
		return nil, nil
	}
	receipts := rawdb.ReadReceipts(b.soc.chainDb, hash, *number)
	if receipts == nil {
		return nil, nil
	}
	logs := make([][]*types.Log, len(receipts))
	for i, receipt := range receipts {
		logs[i] = receipt.Logs
	}
	return logs, nil
}

func (b *SocAPIBackend) GetTd(blockHash common.Hash) *big.Int {
	return b.soc.blockchain.GetTdByHash(blockHash)
}

func (b *SocAPIBackend) GetEVM(ctx context.Context, msg core.Message, state *state.StateDB, header *types.Header, vmCfg vm.Config) (*vm.EVM, func() error, error) {
	state.SetBalance(msg.From(), math.MaxBig256)
	vmError := func() error { return nil }

	context := core.NewEVMContext(msg, header, b.soc.BlockChain(), nil)
	return vm.NewEVM(context, state, b.soc.chainConfig, vmCfg), vmError, nil
}

func (b *SocAPIBackend) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription {
	return b.soc.BlockChain().SubscribeRemovedLogsEvent(ch)
}

func (b *SocAPIBackend) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	return b.soc.BlockChain().SubscribeChainEvent(ch)
}

func (b *SocAPIBackend) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return b.soc.BlockChain().SubscribeChainHeadEvent(ch)
}

func (b *SocAPIBackend) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription {
	return b.soc.BlockChain().SubscribeChainSideEvent(ch)
}

func (b *SocAPIBackend) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return b.soc.BlockChain().SubscribeLogsEvent(ch)
}

func (b *SocAPIBackend) SendTx(ctx context.Context, signedTx *types.Transaction) error {
	return b.soc.txPool.AddLocal(signedTx)
}

func (b *SocAPIBackend) GetPoolTransactions() (types.Transactions, error) {
	pending, err := b.soc.txPool.Pending()
	if err != nil {
		return nil, err
	}
	var txs types.Transactions
	for _, batch := range pending {
		txs = append(txs, batch...)
	}
	return txs, nil
}

func (b *SocAPIBackend) GetPoolTransaction(hash common.Hash) *types.Transaction {
	return b.soc.txPool.Get(hash)
}

func (b *SocAPIBackend) GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error) {
	return b.soc.txPool.State().GetNonce(addr), nil
}

func (b *SocAPIBackend) Stats() (pending int, queued int) {
	return b.soc.txPool.Stats()
}

func (b *SocAPIBackend) TxPoolContent() (map[common.Address]types.Transactions, map[common.Address]types.Transactions) {
	return b.soc.TxPool().Content()
}

func (b *SocAPIBackend) SubscribeNewTxsEvent(ch chan<- core.NewTxsEvent) event.Subscription {
	return b.soc.TxPool().SubscribeNewTxsEvent(ch)
}

func (b *SocAPIBackend) Downloader() *downloader.Downloader {
	return b.soc.Downloader()
}

func (b *SocAPIBackend) ProtocolVersion() int {
	return b.soc.SocVersion()
}

func (b *SocAPIBackend) SuggestPrice(ctx context.Context) (*big.Int, error) {
	return b.gpo.SuggestPrice(ctx)
}

func (b *SocAPIBackend) ChainDb() socdb.Database {
	return b.soc.ChainDb()
}

func (b *SocAPIBackend) EventMux() *event.TypeMux {
	return b.soc.EventMux()
}

func (b *SocAPIBackend) AccountManager() *accounts.Manager {
	return b.soc.AccountManager()
}

func (b *SocAPIBackend) BloomStatus() (uint64, uint64) {
	sections, _, _ := b.soc.bloomIndexer.Sections()
	return params.BloomBitsBlocks, sections
}

func (b *SocAPIBackend) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	for i := 0; i < bloomFilterThreads; i++ {
		go session.Multiplex(bloomRetrievalBatch, bloomRetrievalWait, b.soc.bloomRequests)
	}
}
