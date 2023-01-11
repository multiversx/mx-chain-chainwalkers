package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/multiversx/mx-chain-chainwalkers/parsing/config"
	"github.com/multiversx/mx-chain-core-go/core"
	indexerData "github.com/multiversx/mx-chain-es-indexer-go/data"
)

type object = map[string]interface{}

type elasticServer struct {
	client *elasticsearch.Client
}

func NewElasticServer(cfg config.ElasticSearchConfig) (*elasticServer, error) {
	elasticCfg := elasticsearch.Config{
		Addresses: []string{
			cfg.URL,
		},
		Username: cfg.Username,
		Password: cfg.Password,
	}

	// Instantiate a new Elasticsearch client object instance
	client, err := elasticsearch.NewClient(elasticCfg)
	if err != nil {
		return &elasticServer{}, err
	}

	return &elasticServer{
		client: client,
	}, nil
}

func createQueryBlock(nonce uint64, shardId uint32) map[string]interface{} {
	query := object{
		"query": object{
			"bool": object{
				"must": []interface{}{
					object{
						"match": object{
							"nonce": fmt.Sprintf("%d", nonce),
						},
					},
					object{
						"match": object{
							"shardId": fmt.Sprintf("%d", shardId),
						},
					},
				},
			},
		},
	}

	return query
}

func (es *elasticServer) getBlock(query object) (indexerData.Block, string, error) {
	buff, err := encodeQuery(query)
	if err != nil {
		return indexerData.Block{}, "", err
	}

	res, err := es.client.Search(
		es.client.Search.WithContext(context.Background()),
		es.client.Search.WithIndex("blocks"),
		es.client.Search.WithBody(&buff),
		es.client.Search.WithTrackTotalHits(true),
		es.client.Search.WithPretty(),
	)
	if err != nil {
		return indexerData.Block{}, "", err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if res.IsError() {
		return indexerData.Block{}, "", fmt.Errorf("error response: %s", res)
	}

	// TODO: check why the field "miniBlocksHashes" is not populated.

	var r object
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return indexerData.Block{}, "", fmt.Errorf("cannot decode response %s", err.Error())
	}

	h1 := r["hits"].(object)["hits"].([]interface{})
	if len(h1) == 0 {
		return indexerData.Block{}, "", fmt.Errorf("cannot find block in database")
	}

	h2 := h1[0].(object)["_source"]
	bbb, _ := json.Marshal(h2)
	var block indexerData.Block
	err = json.Unmarshal(bbb, &block)
	if err != nil {
		return indexerData.Block{}, "", fmt.Errorf("cannot unmarshal blocks")
	}

	h3 := h1[0].(object)["_id"]
	blockHash := fmt.Sprint(h3)

	return block, blockHash, nil
}

func (es *elasticServer) getTxByMbHash(hash string) ([]indexerData.Transaction, error) {
	query := object{
		"query": object{
			"match": object{
				"miniBlockHash": hash,
			},
		},
	}

	buff, err := encodeQuery(query)
	if err != nil {
		return nil, err
	}

	res, err := es.client.Search(
		es.client.Search.WithSize(1000),
		es.client.Search.WithScroll(time.Minute),
		es.client.Search.WithContext(context.Background()),
		es.client.Search.WithIndex("transactions"),
		es.client.Search.WithBody(&buff),
		es.client.Search.WithTrackTotalHits(true),
		es.client.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("error response: %s", res)
	}

	defer func() {
		_ = res.Body.Close()
	}()

	var r object
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	txs := make([]indexerData.Transaction, 0)
	hits := r["hits"].(object)
	txs = append(txs, formatTxs(hits)...)

	// use scroll because cannot get more than 10k of transactions
	scrollID := r["_scroll_id"].(string)
	for {
		rScroll, err := es.getScrollResponse(scrollID)
		if err != nil {
			return nil, err
		}

		hits := rScroll["hits"].(object)
		if len(hits["hits"].([]interface{})) < 1 {
			break
		}

		txs = append(txs, formatTxs(hits)...)
	}

	return txs, nil
}

func (es *elasticServer) getScrollResponse(scrollID string) (object, error) {
	resScroll, err := es.client.Scroll(
		es.client.Scroll.WithScrollID(scrollID),
		es.client.Scroll.WithScroll(time.Minute),
	)
	if err != nil {
		return nil, err
	}
	if resScroll.IsError() {
		return nil, fmt.Errorf("error response: %s", resScroll)
	}

	defer func() {
		_ = resScroll.Body.Close()
	}()

	var rScroll object
	if err := json.NewDecoder(resScroll.Body).Decode(&rScroll); err != nil {
		return nil, err
	}

	return rScroll, nil
}

func formatTxs(data object) []indexerData.Transaction {
	var err error

	txs := make([]indexerData.Transaction, 0)
	for _, h1 := range data["hits"].([]interface{}) {
		h2 := h1.(object)["_source"]
		h3 := h1.(object)["_id"]

		var tx indexerData.Transaction
		bbb, _ := json.Marshal(h2)
		err = json.Unmarshal(bbb, &tx)
		if err != nil {
			continue
		}

		tx.Hash = fmt.Sprintf("%s", h3)
		txs = append(txs, tx)
	}
	return txs
}

func encodeQuery(query object) (bytes.Buffer, error) {
	var buff bytes.Buffer
	if err := json.NewEncoder(&buff).Encode(query); err != nil {
		return bytes.Buffer{}, fmt.Errorf("error encoding query: %s", err.Error())
	}

	return buff, nil
}

func (es *elasticServer) GetTransactionsByMbHash(hash string) ([]indexerData.Transaction, error) {
	return es.getTxByMbHash(hash)
}

func (es *elasticServer) GetMetaBlock(nonce uint64) (indexerData.Block, string, error) {
	query := createQueryBlock(nonce, core.MetachainShardId)
	return es.getBlock(query)
}

func (es *elasticServer) GetShardBlockByHash(hash string) (indexerData.Block, string, error) {
	query := object{
		"query": object{
			"match": object{
				"_id": hash,
			},
		},
	}

	return es.getBlock(query)
}

func (es *elasticServer) Height() (uint64, error) {
	query := object{
		"size": 1,
		"query": object{
			"match": object{
				"shardId": fmt.Sprintf("%d", core.MetachainShardId),
			},
		},
		"sort": object{
			"nonce": object{
				"order": "desc",
			},
		},
	}

	buff, err := encodeQuery(query)
	if err != nil {
		return 0, err
	}

	res, err := es.client.Search(
		es.client.Search.WithContext(context.Background()),
		es.client.Search.WithIndex("blocks"),
		es.client.Search.WithBody(&buff),
		es.client.Search.WithTrackTotalHits(true),
		es.client.Search.WithPretty(),
	)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if res.IsError() {
		return 0, fmt.Errorf("error response: %s", res)
	}

	var r object
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return 0, fmt.Errorf("cannot decode response %s", err.Error())
	}

	h1 := r["hits"].(object)["hits"].([]interface{})
	if len(h1) == 0 {
		return 0, fmt.Errorf("cannot find blocks in database")
	}

	h2 := h1[0].(object)["_source"]

	bbb, _ := json.Marshal(h2)
	var block indexerData.Block
	err = json.Unmarshal(bbb, &block)
	if err != nil {
		return 0, fmt.Errorf("cannot unmarshal blocks")
	}

	return block.Nonce, nil
}
