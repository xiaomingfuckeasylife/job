package main

import (
	"job/conf"
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"job/db"
	"job/cron"
	"strings"
	"bytom/src/github.com/bytom/errors"
	"strconv"
	"encoding/hex"
	"sync"
	"math"
)

var (
	GET_BEST_HEIGHT_URL             = conf.Config.ChainApi.GetBestHeight
	GET_BLOCK_BY_HEIGHT_URL         = conf.Config.ChainApi.GetBlockByHeight
	GET_BLOCK_BY_HASH               = conf.Config.ChainApi.GetBlockByHash
	GET_TX_URL                      = conf.Config.ChainApi.GetTransactionByHash
	SEND_TRANSFER_URL               = conf.Config.ChainApi.SendTransfer
	TX_PERIOD                       = conf.Config.Job.TxPeriod
	FEE_PERIOD                      = conf.Config.Job.FeePeriod
	FEE_NUM                         = conf.Config.Fee.FeeNum
	FEE_AMT                         = conf.Config.Fee.FeeAMT
	SENDER_PUBADDR                  = conf.Config.Fee.SenderPubAddr
	SENDER_PRIVKEY                  = conf.Config.Fee.SenderPrivKey
	SELA                    float64 = 100000000
	InitialHeight                   = conf.Config.InitialHeight
)

var syncTxlock sync.RWMutex
var feeLock sync.RWMutex
var syncHeight int

func main() {

	dia := db.Dialect{}
	dia.Create(conf.Config.DriverName,conf.Config.DataSourceName)
	defer dia.Close()

	cron.AddScheduleBySec(TX_PERIOD, func() {
		syncTxlock.Lock()
		log.Println("sync block start")
		defer syncTxlock.Unlock()
		_, err := processTx(&dia)
		if err != nil {
			log.Fatal(err)
		}
	})

	cron.AddScheduleByHours(FEE_PERIOD, func() {
		// TODO get addresses from tables
		feeLock.Lock()
		log.Println("send Fee start")
		defer feeLock.Unlock()
		ri := make([]map[string]string, 1)
		m := make(map[string]string)
		m["address"] = "EKDb9T8hDgT5CwrvxRuoCeKN3WcAKCShB2"
		m["amount"] = fmt.Sprintf("%.8f", math.Round(FEE_AMT*FEE_NUM*SELA)/SELA)
		ri[0] = m

		b, err := json.Marshal(ri)
		if err != nil {
			log.Fatal("json marshal error ", err)
		}
		fmt.Println(string(b))
		processFee(string(b))
	})

}

func processFee(receivInfo string) (bool, error) {
	body := `{
			"Action":"transfer",
			"Version":"1.0.0",
			"Data":
				{"senderAddr":"` + SENDER_PUBADDR + `",
   				 "senderPrivateKey":"` + SENDER_PRIVKEY + `",
				 "memo":"chinajoy fee transfer",
				 "receiver":` + receivInfo + `
				}
			}`
	r := strings.NewReader(body)
	rsp, err := http.Post(SEND_TRANSFER_URL, "application/json", r)
	if err != nil {
		return false, err
	}
	bytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return false, err
	}
	log.Printf("ret Msg : %s \n", string(bytes))
	var ret map[string]interface{}
	err = json.Unmarshal(bytes, &ret)
	if err != nil {
		return false, err
	}
	if ret["error"].(float64) != 0 {
		return false, errors.New(" transfer api error ")
	}
	// TODO update table status
	return true, nil
}

func processTx(dia *db.Dialect) (bool, error) {
	var start int
	var end int

	resp, err := http.Get(GET_BEST_HEIGHT_URL)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	result := make(map[string]interface{})
	json.Unmarshal(body, &result)
	end = int(result["Result"].(float64))
	if syncHeight == 0 {
		list, err := dia.Query(" select height from tx_info order by height desc limit 1")
		if err != nil {
			return false, err
		}
		if list.Len() == 0 {
			start = int(InitialHeight)
		} else {
			m := list.Front().Value.(map[string]string)
			startStr := m["height"]
			start, err = strconv.Atoi(startStr)
			if err != nil {
				return false, err
			}
			start = start + 1
		}
	} else {
		start = syncHeight + 1
	}
	if start >= end+1 {
		log.Println("no block need to sync")
		return true, nil
	}
	for height := start; height < end+1; height++ {
		log.Printf("sync height : %d \n", height)
		resp, err := http.Get(GET_BLOCK_BY_HEIGHT_URL + strconv.Itoa(height))
		if err != nil {
			return false, err
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, err
		}
		var blockInfo map[string]interface{}
		json.Unmarshal(body, &blockInfo)
		rstMap := blockInfo["Result"].(map[string]interface{})
		blockHash := rstMap["hash"].(string)
		timestamp := fmt.Sprintf("%.0f", rstMap["time"])
		txArr := rstMap["tx"].([]interface{})
		blockHeight := int(rstMap["height"].(float64))
		for i := 0; i < len(txArr); i++ {
			txInfo := txArr[i].(map[string]interface{})
			txId := txInfo["txid"].(string)
			if int(txInfo["type"].(float64)) != 2 {
				continue;
			}
			attributes := txInfo["attributes"].([]interface{})
			memoByte, err := hex.DecodeString(attributes[0].(map[string]interface{})["data"].(string))
			if err != nil {
				return false, err
			}
			memo := string(memoByte)
			log.Printf("memo info %s \n", memo)
			if !isValid(memo) {
				continue
			}
			//if !strings.HasPrefix(memo,"chinajoy") {
			//	continue
			//}
			sql := "insert into tx_info (txid,height,memo,timestamp,blockhash) values('" + txId + "'," + strconv.Itoa(blockHeight) + ",'" + memo + "'," + timestamp + ",'" + blockHash + "')"
			_, err = dia.Save(sql)
			if err != nil {
				return false, err
			}
			// TODO update table status
		}
		if height == end && height > syncHeight {
			syncHeight = height
		}
	}

	return true, nil
}

func isValid(memo string) bool {
	if (len(memo) > 0) {
		for i := 0; i < len(memo); i++ {
			if (memo[i] > 102 || (memo[i] < 97 && memo[i] > 57) || memo[i] < 48) {
				return false
			}
		}
	}
	return true
}
