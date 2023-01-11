package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/multiversx/mx-chain-chainwalkers/parsing/config"
	"github.com/multiversx/mx-chain-core-go/core"
)

func main() {
	pb, err := initHeightParser()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	height, err := pb.GetHeight()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("{\"height\": %d} \n", height)
}

func initHeightParser() (*ParserHeight, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return &ParserHeight{}, err
	}

	configurationFileName := dir + "/../config.toml"

	cfg, err := loadEconomicsConfig(configurationFileName)
	if err != nil {
		return &ParserHeight{}, err
	}

	pb, err := NewParserHeight(cfg.ElasticSearchConnector)
	if err != nil {
		return &ParserHeight{}, err
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
