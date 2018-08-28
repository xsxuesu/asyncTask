/**
 *	Author Fanxu(746439274@qq.com)
 */
package control

import (
	"net/http"
	"asyncTask/helpers"
	"strings"
	"asyncTask/model"
)
type SystemDockerPerformanceResult struct {
	Name string `json:"name"`
	ContainerId string `json:"container_id"`
}
func (this * BaseControl) SystemInfo(w http.ResponseWriter, r *http.Request){
	cmd := "cat /proc/version"
	result,err := helpers.ExecShell(cmd)
	if err != nil {
		panic(model.CustomException{402, err.Error()})
	}
	var sysVer int
	if strings.Contains(result,"ubuntu") {
		sysVer = 0
	} else if strings.Contains(result,"centos"){
		sysVer = 1
	} else{
		sysVer = 2
	}
	helpers.ResponseJson(200,"ok",map[string] interface{}{ "type":sysVer , "version":result },w)
	return
}
type SystemPerformanceParam struct {
	Dockers []string `json:"dockers"`
}
func (this * BaseControl) SystemPerformance(w http.ResponseWriter, r *http.Request){
	//body, err := ioutil.ReadAll(r.Body)
	//if err != nil{
	//	panic(model.CustomException{402, err.Error()})
	//}
	//var p SystemPerformanceParam
	//fmt.Println(string(body))
	//err = json.Unmarshal(body,&p)
	//if err != nil{
	//	panic(model.CustomException{402, "dockers params : "+err.Error()})
	//}
	//contents , err := helpers.ExecShell("docker ps -a")
	//if err != nil{
	//	panic(model.CustomException{400, err.Error()})
	//}
	//lines := strings.Split(string(contents), "\n")
	//i := 0
	//dockers := make([]*SystemDockerPerformanceResult,len(p.Dockers))
	//for _, line := range(lines) {
	//	fields := strings.Fields(line)
	//	if len(fields) < 1{
	//		continue
	//	}
	//	if "NAMES" == fields[len(fields)-1]{
	//		continue
	//	}
	//	fmt.Println(p.Dockers)
	//	if !helpers.HasElement(p.Dockers,fields[len(fields)-1]){
	//		continue
	//	}
	//	dockers[i] = &SystemDockerPerformanceResult{
	//		fields[len(fields)-1],
	//		fields[0],
	//	}
	//	i++
	//}
	helpers.ResponseJson(200,"ok",map[string] interface{}{
		"cpu":helpers.CpuUsage,"mem":map[string]helpers.MemStatus{
			"usage":helpers.Mem,
		} },w)
	return
}



