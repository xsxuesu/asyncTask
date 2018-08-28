/**
 *	Author Fanxu(746439274@qq.com)
 */
package control

import (
	"net/http"
	"encoding/json"
	"fmt"
	"asyncTask/helpers"
	"asyncTask/queue"
	"asyncTask/model"
	"asyncTask/config"
	"reflect"
	"log"
)

func (this * BaseControl) UserIndex(w http.ResponseWriter, r *http.Request){
	//异步任务处理
	var task model.Command
	task.Action = config.ACTION_PRINT
	err := queue.Queue.Worker(task)
	if err != nil{
		//同步返回
		helpers.ResponseJson(404,err.Error(),helpers.EmptyResult,w)
		return
	}
	//同步返回
	helpers.ResponseJson(200,"Hello World",helpers.EmptyResult,w)
}

func (this * BaseControl) UserDb(w http.ResponseWriter, r *http.Request){
	var rows model.TableTestList
	var result map[string] interface{}
	data,err := helpers.DbPool.CurrentDb().Query("SELECT * FROM test")
	if err != nil {
		helpers.ResponseJson(400,"error",result,w)
		return
	}
	json.Unmarshal(data,&rows)
	var re = reflect.ValueOf(rows).Interface()
	log.Println(re)
	d,_ := json.Marshal(re)
	json.Unmarshal(d,&result)
	helpers.ResponseJson(200,"ok",map[string] interface{}{"list":rows},w)
}

func (this RpcServer) UserIndex(params map[string] interface{},result *helpers.JsonResult) error{
	//异步任务添加
	var task model.Command
	task.Action = config.ACTION_PRINT
	err1 := queue.Queue.Worker(task)
	if err1 != nil{
		//同步返回结果
		result.Code = 404
		result.Msg = err1.Error()
		return nil
	}
	//同步返回结果
	result.Code = 200
	result.Msg = "ok"
	b := []byte(`{"list":[]}`)
	var d map[string] interface{}
	err := json.Unmarshal(b, &d)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	result.Data = d
	return nil
}

func (this * BaseControl) UserJson(w http.ResponseWriter, r *http.Request){
	//this.NeedLogin = true
	//ok := helpers.LoginRequired(w,r)
	//if ok != "ok"{
	//	helpers.ResponseError(w)
	//	return
	//}
	b := []byte(`{"name":"wednesday", "age":6, "parents": [ "gomez", "moticia" ]}`)
	var d map[string] interface{}
	err := json.Unmarshal(b, &d)
	if err != nil {
		fmt.Println(err)
		helpers.ResponseError(w)
		return
	}
	helpers.ResponseJson(200,"ok",d,w)
}
