package parsing

import (
	"github.com/multiversx/mx-chain-es-indexer-go/data"
)

type DataGetter interface {
	GetTransactionsByMbHash(hash string) ([]data.Transaction, error)
	GetMetaBlock(nonce uint64) (data.Block, string, error)
	GetShardBlockByHash(hash string) (data.Block, string, error)
}

type ParserBlock interface {
	MetaBlocks(nonces []uint64)
}

type ParserHeight interface {
	Height() (uint64, error)
}
