package core

import (
	"bytes"
	"encoding/gob"
	"log"
	"math/big"

	"github.com/dsxg666/snakecoin/common"
)

type State struct {
	Nonce    uint64      // The number of transactions completed by the wallet
	Balance  *big.Int    // Current Account Balance
	Storage  common.Hash // Contract Code Storage Hash
	CodeHash []byte      // Contract Code Hash
}

func NewState() *State {
	return &State{Balance: big.NewInt(100)}
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
