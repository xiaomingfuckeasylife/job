package main

import (
	"job/conf"
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"encoding/hex"
	"job/db"
	"job/cron"
)

var (
	GET_TX_URL = conf.Config.ChainApi.GetTransactionByHash
	TX_PERIOD  = conf.Config.Job.TxPeriod
	FEE_PERIOD = conf.Config.Job.FeePeriod
)

func main() {

	dia := db.Dialect{}
	dia.Create()
	defer dia.Close()

	cron.AddScheduleBySec(TX_PERIOD,
		func() {
			//TODO get txId from tables
			processTx("42a1ea8e584448673dd23fc91abc9a1adaf4ccbb1fb557bad82dd8bed73c6bc4",&dia)
		})

	cron.AddScheduleByHours(FEE_PERIOD, func() {
		processFee()
	})

}

//TODO
func processFee(){

}

func processTx(txId string , dia *db.Dialect) (bool , error){
	resp , err := http.Get(GET_TX_URL+txId)
	if err != nil {
		log.Fatalf("Error Fetching %s , %v \n", GET_TX_URL+txId, err)
	}
	body , err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var txInfo map[string]interface{}
	json.Unmarshal(body,&txInfo)
	rstMap := txInfo["Result"].(map[string]interface{})
	blockHash := rstMap["blockhash"].(string)
	timestamp := fmt.Sprintf("%.0f",rstMap["time"])
	attrsArr := rstMap["attributes"].([]interface{})
	data ,err := hex.DecodeString(attrsArr[0].(map[string]interface{})["data"].(string))
	if err != nil {
		return false , err
	}
	memo := string(data)
	// TODO blockheight get from getBlockbyhash
	blockHeight := "0"
	sql := "insert into tx_info (txid,height,memo,timestamp,blockhash) values('"+txId + "','"+blockHeight+"','"+memo+"',"+timestamp+",'"+blockHash+"')"
	_ , err = dia.Save(sql)
	if err != nil{
		log.Fatal(err)
	}

	return true, nil
}

