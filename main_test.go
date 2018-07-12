package main

import (
	"testing"
	"log"
	"github.com/xiaomingfuckeasylife/job/db"
	"math/rand"
)

func Test_processTx(t *testing.T){
	b , err := processTx(&db.Dialect{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(b)
}


func Test_other(t *testing.T){
	//print(time.Now().Unix())
	rand.Seed(1531308367)
	print(rand.Int31n(10))
}
