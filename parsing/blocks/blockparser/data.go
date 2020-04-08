package blockparser

import "time"

type Block struct {
	Nonce         uint64        `json:"nonce"`
	Round         uint64        `json:"round"`
	Epoch         uint32        `json:"epoch"`
	Hash          string        `json:"hash"`
	Proposer      uint64        `json:"proposer"`
	Validators    []uint64      `json:"validators"`
	PubKeyBitmap  string        `json:"pubKeyBitmap"`
	Size          int64         `json:"size"`
	Timestamp     time.Duration `json:"timestamp"`
	StateRootHash string        `json:"stateRootHash"`
	PrevHash      string        `json:"prevHash"`
	ShardID       uint32        `json:"shardId"`
	TxCount       uint32        `json:"txCount"`
	MiniBlocks    []MiniBlock   `json:"miniBlocks"`
}

type MetaBlock struct {
	Block
	ShardBlocks []ShardBlock `json:"shardBlocks"`
}

type ShardBlock struct {
	Block
}

type MiniBlock struct {
	Hash            string        `json:"hash"`
	SenderShardId   uint32        `json:"senderShardID"`
	ReceiverShardId uint32        `json:"receiverShardID"`
	Transactions    []Transaction `json:"transactions"`
}

type Transaction struct {
	Hash          string        `json:"hash"`
	MBHash        string        `json:"miniBlockHash"`
	BlockHash     string        `json:"blockHash"`
	Nonce         uint64        `json:"nonce"`
	Round         uint64        `json:"round"`
	Value         string        `json:"value"`
	Receiver      string        `json:"receiver"`
	Sender        string        `json:"sender"`
	ReceiverShard uint32        `json:"receiverShard"`
	SenderShard   uint32        `json:"senderShard"`
	GasPrice      uint64        `json:"gasPrice"`
	GasLimit      uint64        `json:"gasLimit"`
	Data          string        `json:"data"`
	Signature     string        `json:"signature"`
	Timestamp     time.Duration `json:"timestamp"`
	Status        string        `json:"status"`
}
