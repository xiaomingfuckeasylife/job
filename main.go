package main

import (
	"job/db"
	"fmt"
	"Elastos.ELA/common/log"
)

func main() {
	dia := db.Dialect{}
	dia.Create()

	id , err := dia.Save("insert into test(name,age) values ('adf',111)")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)
	list := dia.Query("select * from test")
	for e :=list.Front(); e != nil ; e = e.Next(){
		fmt.Println(e.Value.(map[string]string)["name"])
	}
}
