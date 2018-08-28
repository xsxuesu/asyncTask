package model

//单条记录
type TableTestRecord struct {
	Id string `json:"id"`
	Name string `json:"name"`
}
//多条列表记录
type TableTestList []TableTestRecord
