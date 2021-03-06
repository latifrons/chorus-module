// Copyright 2014 The chorus Authors
// This file is part of the chorus library.
//
// The chorus library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The chorus library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the chorus library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/Baptist-Publication/chorus-module/lib/eth/common"
	"github.com/Baptist-Publication/chorus-module/lib/eth/common/hexutil"
	"github.com/Baptist-Publication/chorus-module/lib/eth/rlp"
)

var errMissingLogFields = errors.New("missing required JSON log fields")

// Log represents a contract log event. These events are generated by the LOG opcode and
// stored/indexed by the node.
type Log struct {
	// Consensus fields.
	Address common.Address // address of the contract that generated the event
	Topics  []common.Hash  // list of topics provided by the contract.
	Data    []byte         // supplied by the contract, usually ABI-encoded

	// Derived fields. These fields are filled in by the node
	// but not secured by consensus.
	BlockNumber uint64      // block in which the transaction was included
	TxHash      common.Hash // hash of the transaction
	TxIndex     uint        // index of the transaction in the block
	BlockHash   common.Hash // hash of the block in which the transaction was included
	Index       uint        // index of the log in the receipt

	// The Removed field is true if this log was reverted due to a chain reorganisation.
	// You must pay attention to this field if you receive logs through a filter query.
	Removed bool
}

type rlpLog struct {
	Address common.Address
	Topics  []common.Hash
	Data    []byte
}

type rlpStorageLog struct {
	Address     common.Address
	Topics      []common.Hash
	Data        []byte
	BlockNumber uint64
	TxHash      common.Hash
	TxIndex     uint
	BlockHash   common.Hash
	Index       uint
}

type jsonLog struct {
	Address     *common.Address `json:"address"`
	Topics      *[]common.Hash  `json:"topics"`
	Data        *hexutil.Bytes  `json:"data"`
	BlockNumber *hexutil.Uint64 `json:"blockNumber"`
	TxIndex     *hexutil.Uint   `json:"transactionIndex"`
	TxHash      *common.Hash    `json:"transactionHash"`
	BlockHash   *common.Hash    `json:"blockHash"`
	Index       *hexutil.Uint   `json:"logIndex"`
	Removed     bool            `json:"removed"`
}

// EncodeRLP implements rlp.Encoder.
func (l *Log) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, rlpLog{Address: l.Address, Topics: l.Topics, Data: l.Data})
}

// DecodeRLP implements rlp.Decoder.
func (l *Log) DecodeRLP(s *rlp.Stream) error {
	var dec rlpLog
	err := s.Decode(&dec)
	if err == nil {
		l.Address, l.Topics, l.Data = dec.Address, dec.Topics, dec.Data
	}
	return err
}

func (l *Log) String() string {
	return fmt.Sprintf(`log: %x %x %x %x %d %x %d`, l.Address, l.Topics, l.Data, l.TxHash, l.TxIndex, l.BlockHash, l.Index)
}

// MarshalJSON implements json.Marshaler.
func (l *Log) MarshalJSON() ([]byte, error) {
	jslog := &jsonLog{
		Address: &l.Address,
		Topics:  &l.Topics,
		Data:    (*hexutil.Bytes)(&l.Data),
		TxIndex: (*hexutil.Uint)(&l.TxIndex),
		TxHash:  &l.TxHash,
		Index:   (*hexutil.Uint)(&l.Index),
		Removed: l.Removed,
	}
	// Set block information for mined logs.
	if (l.BlockHash != common.Hash{}) {
		jslog.BlockHash = &l.BlockHash
		jslog.BlockNumber = (*hexutil.Uint64)(&l.BlockNumber)
	}
	return json.Marshal(jslog)
}

// UnmarshalJSON implements json.Umarshaler.
func (l *Log) UnmarshalJSON(input []byte) error {
	var dec jsonLog
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.Address == nil || dec.Topics == nil || dec.Data == nil ||
		dec.TxIndex == nil || dec.TxHash == nil || dec.Index == nil {
		return errMissingLogFields
	}
	declog := Log{
		Address: *dec.Address,
		Topics:  *dec.Topics,
		Data:    *dec.Data,
		TxHash:  *dec.TxHash,
		TxIndex: uint(*dec.TxIndex),
		Index:   uint(*dec.Index),
		Removed: dec.Removed,
	}
	// Block information may be missing if the log is received through
	// the pending log filter, so it's handled specially here.
	if dec.BlockHash != nil && dec.BlockNumber != nil {
		declog.BlockHash = *dec.BlockHash
		declog.BlockNumber = uint64(*dec.BlockNumber)
	}
	*l = declog
	return nil
}

// LogForStorage is a wrapper around a Log that flattens and parses the entire content of
// a log including non-consensus fields.
type LogForStorage Log

// EncodeRLP implements rlp.Encoder.
func (l *LogForStorage) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, rlpStorageLog{
		Address:     l.Address,
		Topics:      l.Topics,
		Data:        l.Data,
		BlockNumber: l.BlockNumber,
		TxHash:      l.TxHash,
		TxIndex:     l.TxIndex,
		BlockHash:   l.BlockHash,
		Index:       l.Index,
	})
}

// DecodeRLP implements rlp.Decoder.
func (l *LogForStorage) DecodeRLP(s *rlp.Stream) error {
	var dec rlpStorageLog
	err := s.Decode(&dec)
	if err == nil {
		*l = LogForStorage{
			Address:     dec.Address,
			Topics:      dec.Topics,
			Data:        dec.Data,
			BlockNumber: dec.BlockNumber,
			TxHash:      dec.TxHash,
			TxIndex:     dec.TxIndex,
			BlockHash:   dec.BlockHash,
			Index:       dec.Index,
		}
	}
	return err
}
