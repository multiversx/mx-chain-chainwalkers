package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/ElrondNetwork/chainwalkers-elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/indexer"
	"github.com/elastic/go-elasticsearch/v7"
)

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
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{
						"match": map[string]interface{}{
							"nonce": fmt.Sprintf("%d", nonce),
						},
					},
					map[string]interface{}{
						"match": map[string]interface{}{
							"shardId": fmt.Sprintf("%d", shardId),
						},
					},
				},
			},
		},
	}

	return query
}

func (es *elasticServer) getBlock(query map[string]interface{}) (indexer.Block, string, error) {
	buff, err := encodeQuery(query)
	if err != nil {
		return indexer.Block{}, "", err
	}

	res, err := es.client.Search(
		es.client.Search.WithContext(context.Background()),
		es.client.Search.WithIndex("blocks"),
		es.client.Search.WithBody(&buff),
		es.client.Search.WithTrackTotalHits(true),
		es.client.Search.WithPretty(),
	)
	if err != nil {
		return indexer.Block{}, "", err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return indexer.Block{}, "", fmt.Errorf("error parsing the response body: %s", err.Error())
		} else {
			return indexer.Block{}, "",
				fmt.Errorf("[%s] %s: %s",
					res.Status(),
					e["error"].(map[string]interface{})["type"],
					e["error"].(map[string]interface{})["reason"],
				)
		}
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return indexer.Block{}, "", fmt.Errorf("cannot decode response %s", err.Error())
	}

	h1 := r["hits"].(map[string]interface{})["hits"].([]interface{})
	if len(h1) == 0 {
		return indexer.Block{}, "", fmt.Errorf("cannot find blocks in database")
	}

	h2 := h1[0].(map[string]interface{})["_source"]
	bbb, _ := json.Marshal(h2)
	var block indexer.Block
	err = json.Unmarshal(bbb, &block)
	if err != nil {
		return indexer.Block{}, "", fmt.Errorf("cannot unmarshal blocks")
	}

	h3 := h1[0].(map[string]interface{})["_id"]
	blockHash := fmt.Sprint(h3)

	return block, blockHash, nil
}

func (es *elasticServer) getTxByMbHash(hash string) ([]indexer.Transaction, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"miniBlockHash": hash,
			},
		},
	}

	buff, err := encodeQuery(query)
	if err != nil {
		return nil, err
	}

	res, err := es.client.Search(
		es.client.Search.WithContext(context.Background()),
		es.client.Search.WithIndex("transactions"),
		es.client.Search.WithBody(&buff),
		es.client.Search.WithTrackTotalHits(true),
		es.client.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	txs := make([]indexer.Transaction, 0)
	for _, h1 := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		h2 := h1.(map[string]interface{})["_source"]

		var tx indexer.Transaction
		bbb, _ := json.Marshal(h2)
		err = json.Unmarshal(bbb, &tx)
		if err != nil {
			continue
		}

		txs = append(txs, tx)
	}

	return txs, nil
}

func encodeQuery(query map[string]interface{}) (bytes.Buffer, error) {
	var buff bytes.Buffer
	if err := json.NewEncoder(&buff).Encode(query); err != nil {
		return bytes.Buffer{}, fmt.Errorf("error encoding query: %s", err.Error())
	}

	return buff, nil
}

func (es *elasticServer) GetTransactionsByMbHash(hash string) ([]indexer.Transaction, error) {
	return es.getTxByMbHash(hash)
}

func (es *elasticServer) GetMetaBlock(nonce uint64) (indexer.Block, string, error) {
	query := createQueryBlock(nonce, core.MetachainShardId)
	return es.getBlock(query)
}

func (es *elasticServer) GetShardBlockByHash(hash string) (indexer.Block, string, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"_id": hash,
			},
		},
	}

	return es.getBlock(query)
}

func (es *elasticServer) Height() (uint64, error) {
	query := map[string]interface{}{
		"size": 1,
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"shardId": fmt.Sprintf("%d", core.MetachainShardId),
			},
		},
		"sort": map[string]interface{}{
			"nonce": map[string]interface{}{
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
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return 0, fmt.Errorf("error parsing the response body: %s", err.Error())
		} else {
			return 0,
				fmt.Errorf("[%s] %s: %s",
					res.Status(),
					e["error"].(map[string]interface{})["type"],
					e["error"].(map[string]interface{})["reason"],
				)
		}
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return 0, fmt.Errorf("cannot decode response %s", err.Error())
	}

	h1 := r["hits"].(map[string]interface{})["hits"].([]interface{})
	if len(h1) == 0 {
		return 0, fmt.Errorf("cannot find blocks in database")
	}

	h2 := h1[0].(map[string]interface{})["_source"]

	bbb, _ := json.Marshal(h2)
	var block indexer.Block
	err = json.Unmarshal(bbb, &block)
	if err != nil {
		return 0, fmt.Errorf("cannot unmarshal blocks")
	}

	return block.Nonce, nil
}
