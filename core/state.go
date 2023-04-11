package core

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/dsxg666/snakecoin/common"
	"github.com/shopspring/decimal"
)

type State struct {
	Nonce    uint64          // The number of transactions completed by the account
	Balance  decimal.Decimal // Current Account Balance
	Storage  common.Hash     // Contract Code Storage Hash
	CodeHash []byte          // Contract Code Hash
}

func NewState() *State {
	return &State{Balance: decimal.NewFromFloat(10)}
}

func (s *State) Serialize() []byte {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(s)
	if err != nil {
		log.Panic("Failed to Encode:", err)
	}
	return buf.Bytes()
}

func DeserializeState(b []byte) *State {
	var state State
	decoder := gob.NewDecoder(bytes.NewReader(b))
	err := decoder.Decode(&state)
	if err != nil {
		log.Panic("Failed to Decode:", err)
	}
	return &state
}
