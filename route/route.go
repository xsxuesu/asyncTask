/**
 *	Author Fanxu(746439274@qq.com)
 */
package route

import (
	"net/http"
	"strings"
	"log"
	"asyncTask/control"
	"asyncTask/helpers"
	"asyncTask/model"
	"reflect"
)

func Run(w http.ResponseWriter, r *http.Request){
	helpers.Try(func() {
		uriArr := strings.Split(r.RequestURI,"/")
		log.Println("request:", r.RequestURI)
		if len(uriArr) < 3 {
			panic(model.CustomException{404,"Method Not Found!"})
		}
		controller := uriArr[1]
		params := strings.Split(uriArr[2],"?")
		params1 := strings.Split(params[0],"#")
		method := params1[0]
		auth := helpers.CheckAuth(r)
		if !auth{
			panic(model.CustomException{401,"没有权限操作!"})
		}
		fun := controller + "_" + method
		control.Execute(w,r,helpers.StrFirstToUpper(fun))
	}, func(i interface{}) {
		log.Println(i)
		code := 500
		var msg string
		switch i.(type) {
		case string:
			msg = i.(string)
			break
		case model.CustomException:
			customErr := reflect.ValueOf(i).Interface().(model.CustomException)
			code = customErr.Code
			msg = customErr.Message
			break
		case error:
			err := reflect.ValueOf(i).Interface().(error)
			msg = err.Error()
			break
		default:
			msg = "System Error"
			break
		}
		helpers.ResponseJson(code, msg, helpers.EmptyResult, w)
	})
}