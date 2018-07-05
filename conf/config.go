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
	ChainApi chainApi `json:"ChainApi"`
	Job job `json:"Job"`
}

type job struct {
	TxPeriod  int64
	FeePeriod int64
}

type chainApi struct {
	GetBlockByHash  string `json:"GetBlockByHash"`
	GetTransactionByHash   string `json:"GetTransactionByHash"`
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
