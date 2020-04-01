package parsing

import "github.com/ElrondNetwork/elrond-go/core/indexer"

type DataGetter interface {
	GetTransactionsByMbHash(hash string) ([]indexer.Transaction, error)
	GetMetaBlock(nonce uint64) (indexer.Block, error)
	GetShardBlockByHash(hash string) (indexer.Block, error)
}

type ParserBlock interface {
	MetaBlock(nonce uint64)
}

type ParserHeight interface {
	Height() (uint64, error)
}
