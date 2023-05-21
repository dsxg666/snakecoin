package server

import (
	"context"
	"encoding/hex"
	"github.com/dsxg666/snakecoin/consensus/pow"
	"log"
	"math/big"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/pebble"
	"github.com/dsxg666/snakecoin/common"
	"github.com/dsxg666/snakecoin/core"
	"github.com/dsxg666/snakecoin/db"
	"github.com/dsxg666/snakecoin/grpc/pb"
	"github.com/dsxg666/snakecoin/mpt"
	"github.com/dsxg666/snakecoin/rlp"
	"github.com/dsxg666/snakecoin/wallet"
)

type Server struct {
	MptDB   *pebble.DB
	ChainDB *pebble.DB
	TxDB    *pebble.DB
	pb.UnimplementedRpcServer
}

func (s *Server) NewAccount(ctx context.Context, req *pb.NewAccountReq) (*pb.NewAccountResp, error) {
	w := wallet.NewWallet()
	path := db.KeystorePath + "/" + w.Address.Hex()
	w.StoreKey(path, []byte(req.Password))

	// Save state to mptrie
	mptBytes := db.Get([]byte("latest"), s.MptDB)
	var e []interface{}
	err := rlp.DecodeBytes(mptBytes, &e)
	if err != nil {
		log.Panic("Failed to DecodeBytes:", err)
	}
	trie := mpt.NewTrieWithDecodeData(e)
	st := core.NewState()
	err = trie.Put(w.Address.Bytes(), st.Serialize())
	if err != nil {
		log.Panic("Failed to Put:", err)
	}
	db.Set([]byte("latest"), mpt.Serialize(trie.Root), s.MptDB)
	return &pb.NewAccountResp{Account: w.Address.Hex()}, nil
}

func (s *Server) GetInfoByBlockNum(ctx context.Context, req *pb.GetInfoByBlockNumReq) (*pb.GetInfoByBlockNumResp, error) {
	lastBlockHashBytes := db.Get([]byte("latest"), s.ChainDB)
	lastBlockBytes := db.Get(lastBlockHashBytes, s.ChainDB)
	lastBlock := core.DeserializeBlock(lastBlockBytes)
	num := lastBlock.Header.Number

	in, _ := new(big.Int).SetString(req.Num, 10)
	if in.Cmp(common.Big0) == -1 || num.Cmp(in) == -1 {
		return nil, nil
	} else {
		var temp *core.Block
		temp = lastBlock
		for i := 0; i < int(new(big.Int).Sub(num, in).Int64()); i++ {
			prevBlockHash := db.Get(temp.Header.PrevBlockHash.Bytes(), s.ChainDB)
			temp = core.DeserializeBlock(prevBlockHash)
		}
		return &pb.GetInfoByBlockNumResp{
			Number:         temp.Header.Number.String(),
			Nonce:          temp.Header.Nonce.String(),
			Time:           common.TimestampToTime(int64(temp.Header.Time)),
			Txs:            strconv.Itoa(len(temp.Body.Txs)),
			Reward:         temp.Header.Reward.String(),
			Difficulty:     temp.Header.Difficulty.String(),
			Coinbase:       temp.Header.Coinbase.Hex(),
			BlockHash:      temp.Header.BlockHash.Hex(),
			PrevBlockHash:  temp.Header.PrevBlockHash.Hex(),
			StateTreeRoot:  temp.Header.PrevBlockHash.Hex(),
			MerkleTreeRoot: temp.Header.MerkleTreeRoot.Hex(),
		}, nil
	}
}

func (s *Server) GetInfoByBlockHash(ctx context.Context, req *pb.GetInfoByBlockHashReq) (*pb.GetInfoByBlockHashResp, error) {
	lastBlockHashBytes := db.Get([]byte("latest"), s.ChainDB)
	lastBlockBytes := db.Get(lastBlockHashBytes, s.ChainDB)
	lastBlock := core.DeserializeBlock(lastBlockBytes)
	var temp *core.Block
	temp = lastBlock

	for i := 0; i < int(lastBlock.Header.Number.Int64())+1; i++ {
		if strings.Compare(req.Hash, hex.EncodeToString(temp.Header.BlockHash.Bytes())) == 0 {
			return &pb.GetInfoByBlockHashResp{
				Number:         temp.Header.Number.String(),
				Nonce:          temp.Header.Nonce.String(),
				Time:           common.TimestampToTime(int64(temp.Header.Time)),
				Txs:            strconv.Itoa(len(temp.Body.Txs)),
				Reward:         temp.Header.Reward.String(),
				Difficulty:     temp.Header.Difficulty.String(),
				Coinbase:       temp.Header.Coinbase.Hex(),
				BlockHash:      temp.Header.BlockHash.Hex(),
				PrevBlockHash:  temp.Header.PrevBlockHash.Hex(),
				StateTreeRoot:  temp.Header.PrevBlockHash.Hex(),
				MerkleTreeRoot: temp.Header.MerkleTreeRoot.Hex(),
			}, nil
		}
		if i != int(lastBlock.Header.Number.Int64()) {
			prevBlockHash := db.Get(temp.Header.PrevBlockHash.Bytes(), s.ChainDB)
			temp = core.DeserializeBlock(prevBlockHash)
		}
	}
	return nil, nil
}

func (s *Server) NewTx(ctx context.Context, req *pb.NewTxReq) (*pb.NewTxResp, error) {
	password := req.Password
	from := req.From
	to := req.To
	amount := req.Amount
	// isExist?
	if !AccountIsExist(from) {
		return &pb.NewTxResp{State: "0"}, nil
	}
	if !AccountIsExist(to) {
		return &pb.NewTxResp{State: "1"}, nil
	}
	// getWallet
	path := db.KeystorePath + "/" + from
	w := wallet.LoadWallet(path, []byte(password), from)
	if w == nil {
		return &pb.NewTxResp{State: "2"}, nil
	}
	// isFull?
	ok, loc := core.TxIsFull(s.TxDB)
	if ok {
		return &pb.NewTxResp{State: "3"}, nil
	}
	// getBalance
	mptBytes := db.Get([]byte("latest"), s.MptDB)
	var e []interface{}
	err := rlp.DecodeBytes(mptBytes, &e)
	if err != nil {
		log.Panic("Failed to DecodeBytes:", err)
	}
	trie := mpt.NewTrieWithDecodeData(e)

	fromByte := common.Hex2Bytes(from[2:])
	fromStateByte, _ := trie.Get(fromByte)
	fromState := core.DeserializeState(fromStateByte)
	fromBalance := fromState.Balance

	val, _ := new(big.Int).SetString(amount, 10)
	if fromBalance.Cmp(val) == -1 {
		return &pb.NewTxResp{State: "4"}, nil
	}

	fromState.Balance.Sub(fromBalance, val)
	fromState.Nonce += 1
	trie.Update(fromByte, fromState.Serialize())

	toByte := common.Hex2Bytes(to[2:])
	toStateByte, _ := trie.Get(toByte)
	toState := core.DeserializeState(toStateByte)
	toState.Balance.Add(toState.Balance, val)
	trie.Update(toByte, toState.Serialize())
	db.Set([]byte("latest"), mpt.Serialize(trie.Root), s.MptDB)

	tx := core.NewTransaction(val, uint64(time.Now().Unix()),
		common.BytesToAddress(fromByte), common.BytesToAddress(toByte), wallet.EncodePubKey(w.PubKey))
	txHash := tx.Hash()
	tx.TxHash.SetBytes(txHash)
	sign := w.Sign(txHash)
	tx.Signature = sign
	core.PushTxToPool(loc, tx, s.TxDB)
	return &pb.NewTxResp{State: "5"}, nil
}

func AccountIsExist(address string) bool {
	files, _ := filepath.Glob(db.KeystorePath + "/*")
	for i := 0; i < len(files); i++ {
		if strings.Compare(files[i][14:], address) == 0 {
			return true
		}
	}
	return false
}

func (s *Server) GetTxPool(ctx context.Context, req *pb.GetTxPoolReq) (*pb.GetTxPoolResp, error) {
	_, loc := core.TxIsFull(s.TxDB)
	if loc[0] == 0 {
		return nil, nil
	}
	var txs []*pb.GetTxPoolResp_Tx
	for i := 0; i < int(loc[0]); i++ {
		txBytes := db.Get([]byte{byte(i)}, s.TxDB)
		tx := core.DeserializeTx(txBytes)
		txs = append(txs, &pb.GetTxPoolResp_Tx{
			TxHash: tx.TxHash.Hex(),
			From:   tx.From.Hex(),
			To:     tx.To.Hex(),
			Amount: tx.Value.String(),
			Time:   common.TimestampToTime(int64(tx.Time)),
		})
	}
	return &pb.GetTxPoolResp{Txs: txs}, nil
}

func (s *Server) GetInfoByTxHash(ctx context.Context, req *pb.GetInfoByTxHashReq) (*pb.GetInfoByTxHashResp, error) {
	lastBlockHashBytes := db.Get([]byte("latest"), s.ChainDB)
	lastBlockBytes := db.Get(lastBlockHashBytes, s.ChainDB)
	lastBlock := core.DeserializeBlock(lastBlockBytes)
	var temp *core.Block
	temp = lastBlock

	for i := 0; i < int(lastBlock.Header.Number.Int64())+1; i++ {
		for i := 0; i < len(temp.Body.Txs); i++ {
			if strings.Compare(req.GetHash(), hex.EncodeToString(temp.Body.Txs[i].TxHash.Bytes())) == 0 {
				return &pb.GetInfoByTxHashResp{
					TxHash: temp.Body.Txs[i].TxHash.Hex(),
					From:   temp.Body.Txs[i].From.Hex(),
					To:     temp.Body.Txs[i].To.Hex(),
					Amount: temp.Body.Txs[i].Value.String(),
					Time:   common.TimestampToTime(int64(temp.Body.Txs[i].Time)),
					Block:  temp.Header.Number.String(),
				}, nil
			}
		}
		if i != int(lastBlock.Header.Number.Int64()) {
			prevBlockHash := db.Get(temp.Header.PrevBlockHash.Bytes(), s.ChainDB)
			temp = core.DeserializeBlock(prevBlockHash)
		}
	}
	return nil, nil
}

func (s *Server) Mine(ctx context.Context, req *pb.MineReq) (*pb.MineResp, error) {
	chainDB := s.ChainDB
	lastBlockHash := db.Get([]byte("latest"), chainDB)
	lastBlockBytes := db.Get(lastBlockHash, chainDB)
	lastBlock := core.DeserializeBlock(lastBlockBytes)
	nonce, diff := pow.Pow(lastBlock.Header.Difficulty, pow.CombinedData(
		lastBlock.Header.Number.Bytes(),
		pow.ToBytes(lastBlock.Header.Time),
		lastBlock.Header.Coinbase.Bytes(),
		lastBlock.Header.PrevBlockHash.Bytes(),
		lastBlock.Header.MerkleTreeRoot.Bytes(),
		lastBlock.Header.StateTreeRoot.Bytes(),
	))
	db.Set([]byte("difficulty"), []byte(diff.String()), chainDB)
	return &pb.MineResp{Nonce: nonce.String()}, nil
}

func (s *Server) NewBlock(ctx context.Context, req *pb.NewBlockReq) (*pb.NewBlockResp, error) {
	nonceIn := req.GetNonce()
	miner := req.GetMiner()
	mptDB := s.MptDB
	txDB := s.TxDB
	chainDB := s.ChainDB
	diffByte := db.Get([]byte("difficulty"), chainDB)
	diff, _ := new(big.Int).SetString(string(diffByte), 10)
	_, loc := core.TxIsFull(txDB)
	if loc[0] == 0 {
		return &pb.NewBlockResp{State: "0"}, nil
	}
	if !AccountIsExist(miner) {
		return &pb.NewBlockResp{State: "1"}, nil
	}

	lastBlockHash := db.Get([]byte("latest"), chainDB)
	lastBlockBytes := db.Get(lastBlockHash, chainDB)
	lastBlock := core.DeserializeBlock(lastBlockBytes)
	nonce, _ := new(big.Int).SetString(nonceIn, 10)

	if !pow.Mine(lastBlock.Header.Difficulty, nonce, pow.CombinedData(
		lastBlock.Header.Number.Bytes(),
		pow.ToBytes(lastBlock.Header.Time),
		lastBlock.Header.Coinbase.Bytes(),
		lastBlock.Header.PrevBlockHash.Bytes(),
		lastBlock.Header.MerkleTreeRoot.Bytes(),
		lastBlock.Header.StateTreeRoot.Bytes(),
	)) {
		return &pb.NewBlockResp{State: "2"}, nil
	}

	// Verity Transaction
	var txs []*core.Transaction
	for i := 0; i < int(loc[0]); i++ {
		txBytes := db.Get([]byte{byte(i)}, txDB)
		tx := core.DeserializeTx(txBytes)
		if !wallet.Verity(tx.Hash(), tx.Signature, wallet.DecodePubKey(tx.PubKey)) {
			return &pb.NewBlockResp{State: "3"}, nil
		}
		tx.State = 1
		txs = append(txs, tx)
		db.Set([]byte{byte(i)}, tx.Serialize(), txDB)
	}

	accBytes := common.Hex2Bytes(miner[2:])
	// Update state tree
	mptBytes := db.Get([]byte("latest"), mptDB)
	var e []interface{}
	err := rlp.DecodeBytes(mptBytes, &e)
	if err != nil {
		log.Panic("Failed to DecodeBytes: ", err)
	}
	trie := mpt.NewTrieWithDecodeData(e)
	stateBytes, _ := trie.Get(accBytes)
	state := core.DeserializeState(stateBytes)
	i, _ := new(big.Int).SetString("20000000000000000000", 10)
	state.Balance = state.Balance.Add(state.Balance, i)
	trie.Update(accBytes, state.Serialize())
	db.Set([]byte("latest"), mpt.Serialize(trie.Root), mptDB)
	// NewBlock
	core.NewBlock2(i, nonce, diff, lastBlock.Header.Number, lastBlock.Header.BlockHash.Bytes(), common.BytesToAddress(accBytes), chainDB, mptDB, txs)

	return &pb.NewBlockResp{State: "4"}, nil
}

func (s *Server) GetAllBlock(ctx context.Context, req *pb.GetAllBlockReq) (*pb.GetAllBlockResp, error) {
	chainDB := s.ChainDB
	lastBlockHashBytes := db.Get([]byte("latest"), chainDB)
	lastBlockBytes := db.Get(lastBlockHashBytes, chainDB)
	lastBlock := core.DeserializeBlock(lastBlockBytes)
	var temp *core.Block
	temp = lastBlock
	var bs []*pb.GetAllBlockResp_Block
	for i := 0; i < int(lastBlock.Header.Number.Int64())+1; i++ {
		bs = append(bs, &pb.GetAllBlockResp_Block{
			Number: temp.Header.Number.String(),
			Time:   common.TimestampToTime(int64(temp.Header.Time)),
			Txs:    strconv.Itoa(len(temp.Body.Txs)),
			Reward: temp.Header.Reward.String(),
			Miner:  temp.Header.Coinbase.Hex(),
		})
		if i == int(lastBlock.Header.Number.Int64()) {
			break
		}
		prevBlockHash := db.Get(temp.Header.PrevBlockHash.Bytes(), chainDB)
		temp = core.DeserializeBlock(prevBlockHash)
	}
	return &pb.GetAllBlockResp{Bs: bs}, nil
}

func (s *Server) GetAllTx(ctx context.Context, req *pb.GetAllTxReq) (*pb.GetAllTxResp, error) {
	chainDB := s.ChainDB
	lastBlockHashBytes := db.Get([]byte("latest"), chainDB)
	lastBlockBytes := db.Get(lastBlockHashBytes, chainDB)
	lastBlock := core.DeserializeBlock(lastBlockBytes)
	var temp *core.Block
	temp = lastBlock
	var txs []*pb.GetAllTxResp_Tx
	for i := 0; i < int(lastBlock.Header.Number.Int64())+1; i++ {
		for j := 0; j < len(temp.Body.Txs); j++ {
			txs = append(txs, &pb.GetAllTxResp_Tx{
				TxHash: temp.Body.Txs[j].TxHash.Hex()[2:],
				From:   temp.Body.Txs[j].From.Hex(),
				To:     temp.Body.Txs[j].To.Hex(),
				Amount: temp.Body.Txs[j].Value.String(),
				Time:   common.TimestampToTime(int64(temp.Body.Txs[j].Time)),
				Block:  temp.Header.Number.String(),
			})
		}
		if i == int(lastBlock.Header.Number.Int64()) {
			break
		}
		prevBlockHash := db.Get(temp.Header.PrevBlockHash.Bytes(), chainDB)
		temp = core.DeserializeBlock(prevBlockHash)
	}
	return &pb.GetAllTxResp{Txs: txs}, nil
}

func (s *Server) GetBalance(ctx context.Context, req *pb.GetBalanceReq) (*pb.GetBalanceResp, error) {
	if !AccountIsExist(req.GetAddr()) {
		return nil, nil
	}
	mptBytes := db.Get([]byte("latest"), s.MptDB)
	var e []interface{}
	err := rlp.DecodeBytes(mptBytes, &e)
	if err != nil {
		log.Panic("Failed to DecodeBytes:", err)
	}
	trie := mpt.NewTrieWithDecodeData(e)
	stateB, _ := trie.Get(common.Hex2Bytes(req.GetAddr()[2:]))
	state := core.DeserializeState(stateB)
	return &pb.GetBalanceResp{Balance: state.Balance.String()}, nil
}
