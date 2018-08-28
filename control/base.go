/**
 *	Author Fanxu(746439274@qq.com)
 */

package control

import (
	"net/http"
	"reflect"
	"fmt"
	"asyncTask/model"
)
//定义控制器函数Map类型，便于后续快捷使用
type ControllerMapsType map[string]reflect.Value
var ControllerMaps = make(ControllerMapsType,0)
type BaseControl struct {
	NeedLogin bool
	Login func() bool
	User map[string] interface{}
}
type RpcServer struct {
}

func Init(){
	var routers BaseControl
	//不传入地址就只能反射Routers静态定义的方法
	vf := reflect.ValueOf(&routers)
	vft := vf.Type()
	//读取方法数量
	mNum := vf.NumMethod()
	//遍历路由器的方法，并将其存入控制器映射变量中
	for i := 0; i < mNum; i++ {
		mName := vft.Method(i).Name
		fmt.Println("index:", i, " MethodName:", mName)
		ControllerMaps[mName] = vf.Method(i) //<<<
	}
}

func Execute(w http.ResponseWriter, r *http.Request,fun string){
	if _, ok := ControllerMaps[fun]; !ok {
		//不存在
		panic(model.CustomException{404,"Method Not Found!"})
	}
	//演示
	params := make([]reflect.Value,2)
	params[0] = reflect.ValueOf(w)
	params[1] = reflect.ValueOf(r)
	//创建带调用方法时需要传入的参数列表
	//使用方法名字符串调用指定方法
	ControllerMaps[fun].Call(params)
}