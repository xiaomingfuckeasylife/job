package db

import (
	"testing"
	"log"
	"fmt"
	"github.com/xiaomingfuckeasylife/job/conf"
)

func TestDb(t *testing.T){

	dia := Dialect{}
	dia.Create(conf.Config.DriverName,conf.Config.DataSourceName)
	id , err := dia.Exec("update test set age = 1111 where id = 1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)
	list , err := dia.Query("select * from test")
	if err != nil {
		log.Fatal(err)
	}
	for e :=list.Front(); e != nil ; e = e.Next(){
		fmt.Println(e.Value.(map[string]string)["name"])
	}
	fmt.Println("Connected ? %t",dia.isConnected())
	dia.Close()
	fmt.Println("Connected ? %t",dia.isConnected())
	//dia.Save("insert into test(name,age) values ('test',111)")
	//dia.Query("select * from test")
	fmt.Println("Connected ? %t",dia.isConnected())
	dia.Close()
	fmt.Println("Connected ? %t",dia.isConnected())

}

func Test_dbTx(t *testing.T){
	dia := Dialect{}
	dia.Create("mysql","root:@tcp(127.0.0.1:3306)/test")
	defer dia.Close()
	tx , err := dia.Begin()
	if err != nil {
		fmt.Printf("%v",err)
		return
	}
	_ , err = dia.ExecTx("update test set age = 1001 ",tx)
	if err != nil {
		fmt.Printf("%v",err)
		return
	}
	_ , err = dia.ExecTx("update test set age = 1001 ",tx)
	if err != nil {
		fmt.Printf("%v",err)
		return
	}
	_ , err = dia.ExecTx("update test set age = 1002 ",tx)
	if err != nil {
		fmt.Printf("%v",err)
		return
	}
	_ , err = dia.ExecTx("update test set age = 1003 ",tx)
	if err != nil {
		fmt.Printf("%v",err)
		return
	}
	_ , err = dia.ExecTx("update test set age = 1004 ",tx)
	if err != nil {
		fmt.Printf("%v",err)
		return
	}
	_ , err = dia.ExecTx("update test set age = 1005 ",tx)
	if err != nil {
		fmt.Printf("%v",err)
		return
	}
	_ , err = dia.ExecTx("update test set age = 1006 ",tx)
	if err != nil {
		fmt.Printf("%v",err)
		return
	}
	//dia.Rollback(tx)
	//if err != nil {
	//	fmt.Printf("%v",err)
	//	return
	//}
	err = dia.Commit(tx)
	if err != nil {
		fmt.Printf("%v",err)
		return
	}

	_, err = dia.Exec("update test set age = 1007")
	if err != nil {
		fmt.Printf("%v",err)
		return
	}
}
