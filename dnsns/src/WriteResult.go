package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

const ResultPath  = "/tmp/result/"
var resultLine = make(chan string)

func GenerateResultLine(record *Record, taskID string, taskName string) {
	var resultList []string
	detectAsStr := strings.Join(record.detectAs, "+")
	detectCNamesStr := strings.Join(record.detectCNames, "+")
	fmt.Println(record.detectAs)
	fmt.Println(record.detectCNames)
	resultList = append(resultList,
		taskName, taskID, record.domain, record.reServer, detectCNamesStr + "/" + detectAsStr, "\n")
	resultStr := strings.Join(resultList, ";")
	fmt.Println(resultStr)
	resultLine <- resultStr
	return
}

func ControlWriteResultRoutine(tasks *Task) (err error){
	runtime.GOMAXPROCS(runtime.NumCPU())
	for _, record := range tasks.records {
		go GenerateResultLine(record, tasks.taskID, tasks.taskName)
	}
	var resultContent string
	for i := 0; i < len(tasks.records); i++ {
		resultContent += <- resultLine
	}
	close(resultLine)
	err = os.Mkdir(ResultPath, 0777)
	if err != nil && !os.IsExist(err) {
		return
	}
	err = ioutil.WriteFile(ResultPath + tasks.taskID + ".result",
		[]byte(resultContent), 0644)
	return
}
