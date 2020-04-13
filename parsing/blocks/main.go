package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	parsing "github.com/ElrondNetwork/chainwalkers-elrond-go"
	"github.com/ElrondNetwork/chainwalkers-elrond-go/blocks/blockparser"
	"github.com/ElrondNetwork/chainwalkers-elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/core"
)

func main() {
	blocksNonces := parceProgramArguments()

	parseBlocks(blocksNonces)
}

func parceProgramArguments() []uint64 {
	lenArgs := len(os.Args)
	if lenArgs < 2 {
		log.Fatal("invalid arguments")
	}

	blocksNonces := make([]uint64, 0)
	for i := 1; i < lenArgs; i++ {
		nonce, err := strconv.ParseUint(os.Args[i], 10, 64)
		if err != nil {
			log.Fatal("invalid arguments")
		}
		blocksNonces = append(blocksNonces, nonce)
	}

	return blocksNonces
}

func parseBlocks(nonces []uint64) {
	parser, err := initBlocksParser()
	if err != nil {
		log.Fatal("cannot init blocks parser ", err)
	}

	parser.MetaBlocks(nonces)

}

func initBlocksParser() (parsing.ParserBlock, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return nil, err
	}
	configurationFileName := dir + "/../config.toml"
	//configurationFileName := "../config.toml"

	cfg, err := loadEconomicsConfig(configurationFileName)
	if err != nil {
		return nil, err
	}

	pb, err := blockparser.NewParserBlock(cfg.ElasticSearchConnector)
	if err != nil {
		return nil, err
	}

	return pb, nil
}

func loadEconomicsConfig(filepath string) (*config.Config, error) {
	cfg := &config.Config{}
	err := core.LoadTomlFile(cfg, filepath)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
