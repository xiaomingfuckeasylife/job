package db

import (
	"testing"
	"log"
	"fmt"
)

func TestDb(t *testing.T){

	dia := Dialect{}
	dia.Create()
	id , err := dia.Save("insert into test(name,age) values ('test',111)")
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
