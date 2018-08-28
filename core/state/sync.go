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

package state

import (
	"bytes"

	"github.com/allsportschain/go-allsportschain/common"
	"github.com/allsportschain/go-allsportschain/rlp"
	"github.com/allsportschain/go-allsportschain/trie"
	"github.com/allsportschain/go-allsportschain/log"
	"fmt"
)

// NewStateSync create a new state trie download scheduler.
func NewStateSync(root common.Hash, database trie.DatabaseReader) *trie.Sync {
	var syncer *trie.Sync
	callback := func(leaf []byte, parent common.Hash) error {
		var obj Account
		log.Debug(fmt.Sprintf("ltf_NewStateSync  %v %v\n", leaf, root.String()))
		if err := rlp.Decode(bytes.NewReader(leaf), &obj); err != nil {
			log.Debug(fmt.Sprintf("ltf_NewStateSync err  %v %v\n", err, root.String()))
			return err
		}
		log.Debug(fmt.Sprintf("ltf_NewStateSync done  %v %v\n", obj, root.String()))
		syncer.AddSubTrie(obj.Root, 64, parent, nil)
		syncer.AddRawEntry(common.BytesToHash(obj.CodeHash), 64, parent)
		return nil
	}
	syncer = trie.NewSync(root, database, callback)
	return syncer
}
// NewStateSync create a new state trie download scheduler.
func NewDopsStateSync(root common.Hash, database trie.DatabaseReader) *trie.Sync {
	var syncer *trie.Sync
	syncer = trie.NewSync(root, database, nil)
	return syncer
}
