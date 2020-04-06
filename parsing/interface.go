package parsing

import "github.com/ElrondNetwork/elrond-go/core/indexer"

type DataGetter interface {
	GetTransactionsByMbHash(hash string) ([]indexer.Transaction, error)
	GetMetaBlock(nonce uint64) (indexer.Block, string, error)
	GetShardBlockByHash(hash string) (indexer.Block, string, error)
}

type ParserBlock interface {
	MetaBlocks(nonces []uint64)
}

type ParserHeight interface {
	Height() (uint64, error)
}
