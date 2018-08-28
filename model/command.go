package model

type Command struct {

	Action int `json:"action"` //操作

	Time string `json:"time"` //操作时间

	Data map[string] interface{} //传来的 各种数据

}

type ReconnectData struct{
	Index int64
}
