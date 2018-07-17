package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"container/list"
	"log"
	"errors"
	"time"
)


type Dialect struct {
	db *sql.DB
	driver string
	source string
	maxConns int
	maxIdles int
	connMaxLifeTime time.Duration
}

func (dia *Dialect) Create(driver string,source string) error{
	db, err := sql.Open(driver,
		source)
	if err != nil {
		return err
	}
	if dia.maxConns == 0 {
		db.SetMaxOpenConns(5 )
	}else {
		db.SetMaxOpenConns(dia.maxConns)
	}
	if dia.maxIdles == 0 {
		db.SetMaxIdleConns(3 )
	}else {
		db.SetMaxIdleConns(dia.maxIdles)
	}
	if dia.connMaxLifeTime == 0 {
		db.SetConnMaxLifetime( time.Second * 14440 )
	}else {
		db.SetConnMaxLifetime(dia.connMaxLifeTime)
	}
	dia.driver = driver
	dia.source = source
	dia.db = db
	return nil
}

func (dia *Dialect) isConnected() error {
	error := dia.db.Ping()
	if error != nil {
		log.Printf("%T || %s \n" ,error, error.Error())
		return error
	}
	return nil
}

func (dia *Dialect) Begin() (tx *sql.Tx,err error){
	tx , err = dia.db.Begin()
	log.Printf(" begin transaction %v\n",tx)
	return
}

func (dia *Dialect) Commit(tx *sql.Tx) error {
	log.Printf(" commit transaction %v\n",tx)
	return tx.Commit()
}

func (dia *Dialect) Rollback(tx *sql.Tx) error {
	log.Printf(" rollback transaction %v\n",tx)
	return tx.Rollback()
}

func (dia *Dialect) ExecTx(sqlStr string , tx *sql.Tx) (int64, error) {
	log.Printf("sql with tx : %s\n",sqlStr)
	if tx == nil {
		return -1 , errors.New("tx can not be nil ")
	}
	stmt, err := tx.Prepare(sqlStr)
	if err != nil {
		return -1, err
	}

	result, err := stmt.Exec()
	if err != nil {
		return -1, err
	}
	id, err := result.LastInsertId();
	if err != nil {
		return -1, err
	}
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (dia *Dialect) Exec(sqlStr string) (int64, error) {
	log.Printf("sql : %s\n",sqlStr)
	if dia.db == nil || dia.isConnected() != nil{
		if dia.driver != "" && dia.source != "" {
			log.Println("start a new db instance ,close the old one ")
			err := dia.Create(dia.driver,dia.source)
			if err != nil {
				return -1 , err
			}
		}else{
			return -1 , errors.New("db is nil or closed")
		}
	}
	stmt, err := dia.db.Prepare(sqlStr)
	if err != nil {
		return -1, err
	}

	result, err := stmt.Exec()
	if err != nil {
		return -1, err
	}
	id, err := result.LastInsertId();
	if err != nil {
		return -1, err
	}
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (dia *Dialect) Query(sqlStr string) (*list.List, error) {
	log.Printf("sql : %s\n",sqlStr)
	if dia.db == nil || dia.isConnected() != nil{
		if dia.driver != "" && dia.source != "" {
			log.Println("start a new db instance , close the old one ")
			err := dia.Create(dia.driver,dia.source)
			if err != nil {
				return nil , err
			}
		}else{
			return nil , errors.New("db is nil or closed")
		}
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

func (dia *Dialect) Close() error{
	return dia.db.Close()
}

func (dia *Dialect) SetMaxOpenConnections(maxConn int) error {
	if maxConn <= 0 {
		return errors.New("maxConn must bigger than 0")
	}
	dia.maxConns = maxConn
	return nil
}

func (dia *Dialect) SetMaxIdles(maxIdles int)  error{
	if maxIdles <= 0 {
		return errors.New("maxIdles must bigger than 0")
	}
	dia.maxIdles = maxIdles
	return nil
}

func (dia *Dialect) SetConnMaxLifeTime(connMaxLifeTime int)  error{
	if connMaxLifeTime <= 0 {
		return errors.New("connMaxLifeTime must bigger than 0")
	}
	dia.maxIdles = connMaxLifeTime
	return nil
}
