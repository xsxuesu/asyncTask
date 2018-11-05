package helpers

import (
	"asyncTask/config"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
)

type Db struct {
	Db          *sql.DB
	ParseResult func(query *sql.Rows) ([]byte, error)
	Connect     func() error
	Query       func(sql string) ([]byte, error)
	Execute     func(sql string) error
}

var DbPool struct {
	Dbs       []Db      //数据库连接池
	CurrentDb func() Db //当前可用的数据库连接
}

//初始化mysql 连接池
func InitDbPool() {
	DbPool.Dbs = make([]Db, config.All().MysqlPoolCount)
	for i := 0; i < int(config.All().MysqlPoolCount); i++ {
		err := Reconnect(i)
		if err != nil {
			log.Fatalf(err.Error())
			return
		}
		log.Println("init database connect:", i)
	}
	DbPool.CurrentDb = func() Db {
		index := time.Now().Unix() % int64(len(DbPool.Dbs))
		err := DbPool.Dbs[index].Execute("SET NAMES utf8")
		if err != nil {
			Reconnect(int(index))
		}
		log.Println(index)
		return DbPool.Dbs[index]
	}

}

func Reconnect(i int) error {
	var db Db
	db.Connect = func() error {
		var err error
		db.Db, err = sql.Open("mysql", config.All().Mysql.Conn)
		if err != nil {
			log.Printf(err.Error())
			return err
		}
		return nil
	}
	db.Query = func(sql string) ([]byte, error) {
		if db.Db == nil {
			log.Printf("数据库未连接")
			return nil, errors.New("数据库未连接")
		}
		query, err := db.Db.Query(sql)
		if err != nil {
			log.Printf(err.Error())
			return nil, err
		}
		return db.ParseResult(query)
	}
	db.Execute = func(sql string) error {
		if db.Db == nil {
			log.Printf("数据库未连接")
			return errors.New("数据库未连接")
		}
		_, err := db.Db.Exec(sql)
		if err != nil {
			log.Printf(err.Error())
			return err
		}
		return nil
	}
	db.ParseResult = func(query *sql.Rows) ([]byte, error) {
		column, _ := query.Columns()              //读出查询出的列字段名
		values := make([][]byte, len(column))     //values是每个列的值，这里获取到byte里
		scans := make([]interface{}, len(column)) //因为每次查询出来的列是不定长的，用len(column)定住当次查询的长度
		for j := range values {                   //让每一行数据都填充到[][]byte里面
			scans[j] = &values[j]
		}
		jj := 0
		var rows []interface{}
		for query.Next() { //循环，让游标往下移动
			if err := query.Scan(scans...); err != nil { //query.Scan查询出来的不定长值放到scans[jj] = &values[jj],也就是每行都放在values里
				fmt.Println(err)
				return nil, err
			}
			row := make(map[string]string) //每行数据
			for k, v := range values {     //每行数据是放在values里面，现在把它挪到row里
				key := column[k]
				row[key] = string(v)
			}
			rows = append(rows, row)
			jj++
		}
		return json.Marshal(rows)
	}
	DbPool.Dbs[i] = db
	err := DbPool.Dbs[i].Connect()
	if err != nil {
		return err
	}
	return nil
}
