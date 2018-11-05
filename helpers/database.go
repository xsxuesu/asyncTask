package helpers

import (
	"asyncTask/config"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
)

var Database struct {
	Db          *sql.DB
	Result      []byte
	ParseResult func(query *sql.Rows)
	Connect     func() error
	Query       func(sql string)
	Execute     func(sql string)
}

//单链接 池
func InitDataBase() error {
	Database.Connect = func() error {
		var err error
		log.Printf(config.All().Mysql.Conn)
		Database.Db, err = sql.Open("mysql", config.All().Mysql.Conn)
		if err != nil {
			log.Printf(err.Error())
			return err
		}
		return nil
	}
	Database.Query = func(sql string) {
		query, err := Database.Db.Query(sql)
		if err != nil {
			log.Printf(err.Error())
			return
		}
		Database.ParseResult(query)
	}
	Database.Execute = func(sql string) {
		_, err := Database.Db.Query(sql)
		if err != nil {
			log.Printf(err.Error())
			return
		}
	}
	Database.ParseResult = func(query *sql.Rows) {
		column, _ := query.Columns()              //读出查询出的列字段名
		values := make([][]byte, len(column))     //values是每个列的值，这里获取到byte里
		scans := make([]interface{}, len(column)) //因为每次查询出来的列是不定长的，用len(column)定住当次查询的长度
		for i := range values {                   //让每一行数据都填充到[][]byte里面
			scans[i] = &values[i]
		}
		i := 0
		results := make(map[int]interface{})
		for query.Next() { //循环，让游标往下移动
			if err := query.Scan(scans...); err != nil { //query.Scan查询出来的不定长值放到scans[i] = &values[i],也就是每行都放在values里
				fmt.Println(err)
				return
			}
			row := make(map[string]string) //每行数据
			for k, v := range values {     //每行数据是放在values里面，现在把它挪到row里
				key := column[k]
				row[key] = string(v)
			}
			results[i] = row //装入结果集中
			i++
		}
		var err error
		Database.Result, err = json.Marshal(results)
		if err != nil { //query.Scan查询出来的不定长值放到scans[i] = &values[i],也就是每行都放在values里
			fmt.Println(err)
			return
		}
		for k, v := range results { //查询出来的数组
			fmt.Println(k, v)
		}
	}
	return nil
}
