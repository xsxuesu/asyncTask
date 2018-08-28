package helpers

import (
	"io/ioutil"
	"strings"
	"strconv"
	"fmt"
)
type MemStatus struct {
	All  string `json:"all"`
	Used string `json:"used"`
	Free string `json:"free"`
}
type CpuStatus struct {
	Usage string `json:"usage"`
	All string `json:"all"`
	Count int `json:"count"`
}
var CpuUsage CpuStatus
var Mem MemStatus
func GetCPUSample() (idle, total uint64) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range(lines) {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				total += val // tally up all the numbers to get total ticks
				if i == 4 {  // idle is the 5th field in the cpu line
					idle = val
				}
			}
			return
		}
	}
	return
}
func GetCPUInfo() (idle, total uint64) {
	CpuUsage.Count = 0
	all := 0.00
	contents, err := ioutil.ReadFile("/proc/cpuinfo")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range(lines) {
		fields := strings.Fields(line)
		if len(fields) < 1 {
			continue
		}
		if fields[1] == "MHz" {
			CpuUsage.Count += 1
			hz ,_:= strconv.ParseFloat(fields[3], 64)
			all += hz
		}
	}
	CpuUsage.All = fmt.Sprintf("%6.2fMHz",all)
	return
}

func GetSysSample() (idle, total uint64) {
	contents, err := ioutil.ReadFile("/tmp/baasstat")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range(lines) {
		fields := strings.Fields(line)
		if len(fields) < 1{
			continue
		}
		if fields[1] == "Mem"{
			usage ,_:=strconv.ParseFloat(fields[3], 64)
			//内存
			Mem.All = fmt.Sprintf("%6.2fmb",usage / 1024)
			usage ,_=strconv.ParseFloat(fields[7], 64)
			Mem.Used = fmt.Sprintf("%6.2fmb",usage / 1024)
			usage ,_=strconv.ParseFloat(fields[5], 64)
			Mem.Free = fmt.Sprintf("%6.2fmb",usage / 1024)
		}
		if fields[0] == "%Cpu(s):"{
			usage ,_:=strconv.ParseFloat(fields[1], 64)
			//CPU
			CpuUsage.Usage = fmt.Sprintf("%6.2f",usage) + "%"
		}
	}
	return
}
type DockerStatus struct {
	Name string `json:"name"`
	ContainerId string `json:"container_id"`
	CpuUsage string `json:"cpuUsage"`
	MemUsage string `json:"memUsage"`
	Status int `json:"status"`
}
var DockerUsage map[string]*DockerStatus
func GetDockerSample() (idle, total uint64) {
	contents, err := ioutil.ReadFile("/tmp/baasdockerstat")
	if err != nil {
		return
	}
	//销毁之前的内存
	for k, _ := range DockerUsage {
		delete(DockerUsage,k)
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range(lines) {
		fields := strings.Fields(line)
		if len(fields) < 1{
			continue
		}
		if "PIDS" == fields[len(fields)-1]{
			continue
		}
		DockerUsage[fields[0]] = &DockerStatus{
			"",
			"",
			fields[1],
			fields[7],
			1,
		}
	}
	return
}

