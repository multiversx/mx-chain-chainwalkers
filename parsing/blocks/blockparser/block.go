package blockparser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	parsing "github.com/ElrondNetwork/chaimwalkers-elrong-go"
	"github.com/ElrondNetwork/chaimwalkers-elrong-go/config"
	"github.com/ElrondNetwork/chaimwalkers-elrong-go/server"
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

func (pb *ParserBlock) MetaBlock(nonce uint64) {
	mb, err := pb.prepareBlock(nonce)
	if err != nil {
		log.Fatalf("Fatal error: %s", err)
	}

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(mb); err != nil {
		log.Fatalf("Error encoding blocks: %s", err)
	}

	fmt.Println(string(buf.Bytes()))
}

func (pb *ParserBlock) prepareBlock(nonce uint64) (MetaBlock, error) {
	metaBlock, err := pb.dg.GetMetaBlock(nonce)
	if err != nil {
		return MetaBlock{}, err
	}

	baseBlock := convertIndexerBlockToParserBlock(metaBlock)
	metachainBlock := MetaBlock{}
	metachainBlock.Block = baseBlock

	metachainBlock.MiniBlocks = pb.getMiniBlocksAndTxs(metaBlock.MiniBlocksHashes)

	shardBlocks := make([]ShardBlock, 0)
	for _, shardBlockHash := range metaBlock.NotarizedBlocksHashes {
		shardBlock, err := pb.dg.GetShardBlockByHash(shardBlockHash)
		if err != nil {
			continue
		}

		mbs := pb.getMiniBlocksAndTxs(shardBlock.MiniBlocksHashes)

		sb := convertIndexerBlockToParserBlock(shardBlock)
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
		if err != nil || len(txs) == 0 {
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
