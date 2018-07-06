package main

import (
	"testing"
	"job/db"
	"log"
	"fmt"
)

func Test_processTx(t *testing.T){
	b , err := processTx(&db.Dialect{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(b)
}


func Test_other(t *testing.T){
	fmt.Println('a','f','0','9')
}
