package main

import (
	"fmt"
	"strconv"
	"time"
	"os"
)

func main() {
	randString := ""
	_ = SendUDP("", randString, "start")
	err := os.Mkdir(AppstatusPath, 0777)
	if err != nil && !os.IsExist(err) {
		fmt.Println(err.Error())
		os.Exit(10)
	}
	err = os.Mkdir(LogPath, 0777)
	if err != nil && !os.IsExist(err) {
		WriteError2Appstatus(err.Error(), 9)
	}
	//网络检查
	netCheckFlag, err := NetCheck()
	if err != nil || !netCheckFlag {
		ConnectFail()
		WriteError2Appstatus("Can not connect to the Internet.", 2)
	}
	//任务开始
	TaskStart()
	//读取配置
	GetConf()
	tasks, err := GetTaskConfig()
	if err != nil {
		GetConfFail()
		WriteError2Appstatus(err.Error(), 3)
	}
	GetConfSuccess()
	startTime := time.Now().Unix()
	_ = SendUDP(tasks.taskID, tasks.subID, "run")
	//任务执行
	_ = SendUDP(tasks.taskID, randString, "run")
	TaskRun()
	err = ControlDNSQueryRoutine(tasks)
	if err != nil {
		TaskRunFail()
		WriteError2Appstatus(err.Error(), 1)
	}
	ControlCompareRoutine(tasks)
	TaskRunSuccess()
	//进度
	process := len(tasks.records) / 30
	final_process := len(tasks.records) % 30
	for i := 0; i < process; i++ {
		err = SendProcess(tasks.taskID, tasks.uuid, "DomainInfo", 30, false)
		if err != nil {
			WriteResultFail()
			WriteError2Appstatus(err.Error(), 1)
		}
		time.Sleep(time.Duration(1 * time.Second))
	}
	err = SendProcess(tasks.taskID, tasks.uuid, "DomainInfo", final_process, true)
		if err != nil {
			WriteResultFail()
			WriteError2Appstatus(err.Error(), 1)
		}
	//写结果
	WriteResult()
	err = ControlWriteResultRoutine(tasks)
	if err != nil {
		WriteResultFail()
		WriteError2Appstatus(err.Error(), 1)
	}
	WriteResultSuccess()
	endTime := time.Now().Unix()
	duration := endTime - startTime
	_ = SendUDP(tasks.taskID, tasks.subID, "len: "+strconv.Itoa(len(tasks.records))+"; duration: "+strconv.FormatInt(duration, 10))
	//写状态文件
	WriteSuccess2Appstatus()
}
