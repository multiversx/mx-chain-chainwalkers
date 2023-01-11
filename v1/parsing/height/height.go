package main

import (
	"github.com/multiversx/mx-chain-chainwalkers/parsing"
	"github.com/multiversx/mx-chain-chainwalkers/parsing/config"
	"github.com/multiversx/mx-chain-chainwalkers/parsing/server"
)

type ParserHeight struct {
	dg parsing.ParserHeight
}

func NewParserHeight(cfg config.ElasticSearchConfig) (*ParserHeight, error) {
	esClient, err := server.NewElasticServer(cfg)
	if err != nil {
		return &ParserHeight{}, err
	}

	return &ParserHeight{
		dg: esClient,
	}, nil
}

func (ph *ParserHeight) GetHeight() (uint64, error) {
	height, err := ph.dg.Height()
	if err != nil {
		return 0, err
	}

	return height, nil
}
