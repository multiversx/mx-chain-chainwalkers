package main

import (
	parsing "github.com/ElrondNetwork/chainwalkers-elrond-go"
	"github.com/ElrondNetwork/chainwalkers-elrond-go/config"
	"github.com/ElrondNetwork/chainwalkers-elrond-go/server"
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
