/**
 *	Author Fanxu(746439274@qq.com)
 */

package queue

import (
	"time"
	"asyncTask/helpers"
)

//func ServerMonitor(){
//	for true {
//		idle0,total0 := helpers.GetCPUSample()
//		time.Sleep(3 * time.Second)
//		idle1, total1 := helpers.GetCPUSample()
//		idleTicks := float64(idle1 - idle0)
//		totalTicks := float64(total1 - total0)
//		//CpuUsage = 100 * (totalTicks - idleTicks) / totalTicks
//	}
//}

func SystemMonitor(){
	helpers.GetCPUInfo()
	for true {
		helpers.ExecShell("top -b -n 1 -d 1 > /tmp/baasstat")
		//helpers.ExecShell("docker stats --no-stream > /tmp/baasdockerstat")
		time.Sleep(3 * time.Second)
		helpers.GetSysSample()
		//helpers.GetDockerSample()
	}
}