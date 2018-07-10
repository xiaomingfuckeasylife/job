package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"container/list"
	"log"
	"errors"
)

type Dialect struct {
	db *sql.DB
}

func (dia *Dialect) Create(driver string,source string) {

	//db, err := sql.Open("mysql",
	//	"root:@tcp(127.0.0.1:3306)/test")
	db, err := sql.Open(driver,
		source)
	if err != nil {
		log.Fatal(err)
	}
	dia.db = db
}

func (dia *Dialect) isConnected() bool {
	error := dia.db.Ping()
	if error != nil {
		return false
	}
	return true
}

func (dia *Dialect) Save(sqlStr string) (int64, error) {
	log.Printf("sql : %s\n",sqlStr)
	if dia.db == nil || !dia.isConnected() {
		return -1 , errors.New("db is nil or closed")
	}
	stmt, err := dia.db.Prepare(sqlStr)
	if err != nil {
		return -1, err
	}

	result, err := stmt.Exec()
	if err != nil {
		log.Fatal(err)
		return -1, err
	}
	id, err := result.LastInsertId();
	if err != nil {
		log.Fatal(err)
		return -1, err
	}
	return id, nil
}

func (dia *Dialect) Query(sqlStr string) (*list.List, error) {
	log.Printf("sql : %s\n",sqlStr)
	if dia.db == nil || !dia.isConnected() {
		return nil , errors.New("db is nil or closed")
	}

	rows, err := dia.db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
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
			return nil, err
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
		return nil, err
	}
	return retList, nil
}

func (dia *Dialect) Close() {
	dia.db.Close()
}
