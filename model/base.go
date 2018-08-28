package model

import (
	"reflect"
)

type BaseModel struct {
	Save func(m *BaseModel) bool //保存记录
	One func(m BaseModel) interface{} //单条记录
	All func(m BaseModel) [] interface{} //多条记录
	TableName string
}
var BaseActiveModel BaseModel
func InitBaseModel()  {
	BaseActiveModel.Save = func(m *BaseModel) bool {
		immutable := reflect.TypeOf(m)
		fields := immutable.NumField()
		//sql := fmt.Sprintf("INSERT INTO %s(",m.TableName)
		var fieldSql string
		for i :=0;i<fields ;i++  {
			fieldSql += immutable.Field(i).Name + ","
		}
		return true
	}
}
