package blockparser

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/multiversx/mx-chain-chainwalkers/parsing"
	"github.com/multiversx/mx-chain-chainwalkers/parsing/config"
	"github.com/multiversx/mx-chain-chainwalkers/parsing/server"
)

type ParserBlock struct {
	dg parsing.DataGetter
}

func NewParserBlock(cfg config.ElasticSearchConfig) (*ParserBlock, error) {
	esClient, err := server.NewElasticServer(cfg)
	if err != nil {
		return &ParserBlock{}, err
	}

	return &ParserBlock{
		dg: esClient,
	}, nil
}

func (pb *ParserBlock) MetaBlocks(nonces []uint64) {
	metablocks := make([]MetaBlock, len(nonces))
	for i := 0; i < len(nonces); i++ {
		metaBlock, err := pb.prepareBlock(nonces[i])
		if err != nil {
			log.Fatalf("Fatal error: %s", err)
		}

		metablocks[i] = metaBlock
	}

	buff, err := json.Marshal(&metablocks)
	if err != nil {
		log.Fatalf("Error encoding blocks: %s", err)
	}

	fmt.Println(string(buff))
}

func (pb *ParserBlock) prepareBlock(nonce uint64) (MetaBlock, error) {
	metaBlock, hash, err := pb.dg.GetMetaBlock(nonce)
	if err != nil {
		return MetaBlock{}, err
	}

	baseBlock := convertIndexerBlockToParserBlock(metaBlock, hash)
	metachainBlock := MetaBlock{}
	metachainBlock.Block = baseBlock

	metachainBlock.MiniBlocks = pb.getMiniBlocksAndTxs(metaBlock.MiniBlocksHashes)

	shardBlocks := make([]ShardBlock, 0)
	for _, shardBlockHash := range metaBlock.NotarizedBlocksHashes {
		shardBlock, hash, err := pb.dg.GetShardBlockByHash(shardBlockHash)
		if err != nil {
			continue
		}

		mbs := pb.getMiniBlocksAndTxs(shardBlock.MiniBlocksHashes)

		sb := convertIndexerBlockToParserBlock(shardBlock, hash)
		sb.MiniBlocks = mbs

		s := ShardBlock{}
		s.Block = sb

		shardBlocks = append(shardBlocks, s)
	}

	metachainBlock.ShardBlocks = shardBlocks

	return metachainBlock, nil
}

func (pb *ParserBlock) getMiniBlocksAndTxs(mbHahes []string) []MiniBlock {
	miniblocks := make([]MiniBlock, 0)
	for _, mbHash := range mbHahes {
		txs, err := pb.dg.GetTransactionsByMbHash(mbHash)
		if err != nil {
			log.Fatal("cannot get transactions", err.Error())
		}
		if len(txs) == 0 {
			continue
		}

		miniblocks = append(miniblocks, MiniBlock{
			Hash:            mbHash,
			SenderShardId:   txs[0].SenderShard,
			ReceiverShardId: txs[0].ReceiverShard,
			Transactions:    convertIndexerTxsToParserTxs(txs),
		})
	}

	return miniblocks
}
