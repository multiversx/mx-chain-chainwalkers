package blockparser

import (
	"github.com/ElrondNetwork/elrond-go/core/indexer"
)

func convertIndexerTxToParserTx(tx indexer.Transaction) Transaction {
	return Transaction{
		Hash:          tx.Hash,
		MBHash:        tx.MBHash,
		BlockHash:     tx.BlockHash,
		Nonce:         tx.Nonce,
		Round:         tx.Round,
		Value:         tx.Value,
		Receiver:      tx.Receiver,
		Sender:        tx.Sender,
		ReceiverShard: tx.ReceiverShard,
		SenderShard:   tx.SenderShard,
		GasPrice:      tx.GasPrice,
		GasLimit:      tx.GasLimit,
		Data:          tx.Data,
		Signature:     tx.Signature,
		Timestamp:     tx.Timestamp,
		Status:        tx.Status,
	}
}

func convertIndexerTxsToParserTxs(txs []indexer.Transaction) []Transaction {
	parserTxs := make([]Transaction, 0)
	for _, tx := range txs {
		parserTxs = append(parserTxs, convertIndexerTxToParserTx(tx))
	}

	return parserTxs
}

func convertIndexerBlockToParserBlock(block indexer.Block, hash string) Block {
	// TODO: Check if we can drop "validators" field - is it required?
	// TODO: Check if we should also add miniBlocksHashes field.
	// TODO: Check why MiniBlocks field is nil.

	return Block{
		Nonce:         block.Nonce,
		Round:         block.Round,
		Epoch:         block.Epoch,
		Hash:          hash,
		Proposer:      block.Proposer,
		Validators:    block.Validators,
		PubKeyBitmap:  block.PubKeyBitmap,
		Size:          block.Size,
		Timestamp:     block.Timestamp,
		StateRootHash: block.StateRootHash,
		PrevHash:      block.PrevHash,
		ShardID:       block.ShardID,
		TxCount:       block.TxCount,
		MiniBlocks:    nil,
	}
}
