package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/go/src/pkg/log"
	"container/list"
)

type Dialect struct {
	Db *sql.DB
}

func (dia *Dialect) Create(){

	db, err := sql.Open("mysql",
		"root:@tcp(127.0.0.1:3306)/test")
	if err != nil {
		log.Fatal(err)
	}
	dia.Db = db
}

func (dia *Dialect) isConnected() bool{
	error :=dia.Db.Ping()
	if error != nil {
		return false
	}
	return true
}

func (dia *Dialect) Save(sqlStr string) (int64,error) {
	if dia.Db == nil || !dia.isConnected() {
		dia.Create()
	}
	stmt , err:=dia.Db.Prepare(sqlStr)
	if err != nil {
		return -1 , err
	}

	result , err := stmt.Exec()
	if err != nil {
		log.Fatal(err)
		return -1 , err
	}
	id , err := result.LastInsertId();
	if err != nil {
		log.Fatal(err)
		return -1 , err
	}
	return id , nil
}

func (dia *Dialect) Query(sqlStr string) *list.List{
	if dia.Db == nil || !dia.isConnected() {
		dia.Create()
	}

	rows , err :=dia.Db.Query(sqlStr)
	if err != nil {
		log.Fatal(err.Error())
	}
	columns , err := rows.Columns()
	if err != nil {
		log.Fatal(err.Error())
	}
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	retList := list.New()

	// Fetch rows
	for rows.Next() {
		retMap := make(map[string]string)
		retList.PushBack(retMap)
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			log.Fatal(err.Error()) // proper error handling instead of panic in your app
		}

		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			retMap[columns[i]] = value
		}
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err.Error()) // proper error handling instead of panic in your app
	}
	return retList
}

func (dia *Dialect) Close(){
	dia.Db.Close()
}
