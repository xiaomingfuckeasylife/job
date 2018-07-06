package conf

import (
	"io/ioutil"
	"os"
	"bytes"
	"encoding/json"
	"log"
)

var (
	DEFAULT_FILE = "./config.json"
	Config       *config
)

type config struct {
	ChainApi      chainApi `json:"ChainApi"`
	Job           job      `json:"Job"`
	Fee           fee      `json:"Fee"`
	InitialHeight int64
}

type fee struct {
	FeeAMT        float64 `json:"FeeAMT"`
	FeeNum        float64 `json:"FeeNum"`
	SenderPubAddr string  `json:"SenderPubAddr"`
	SenderPrivKey string  `json:"SenderPrivKey"`
}

type job struct {
	TxPeriod  int64
	FeePeriod int64
}

type chainApi struct {
	GetBestHeight        string `json:"GetBestHeight"`
	GetBlockByHeight     string `json:"GetBlockByHeight"`
	GetBlockByHash       string `json:"GetBlockByHash"`
	GetTransactionByHash string `json:"GetTransactionByHash"`
	SendTransfer         string `json:"SendTransfer"`
}

func init() {
	buf, err := ioutil.ReadFile(DEFAULT_FILE)
	if err != nil {
		log.Fatalf("Error init config file %v \n", err)
		os.Exit(-1)
	}
	buf = bytes.TrimPrefix(buf, []byte("\xef\xbb\xbf"))
	Config = &config{}
	err = json.Unmarshal(buf, Config)
	if err != nil {
		log.Fatalf("not a valid json config file %v \n", err)
		os.Exit(-1)
	}
}
