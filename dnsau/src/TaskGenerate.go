package main

import (
	"encoding/base64"
	"io/ioutil"
	"strings"
	"fmt"
)

const ConfPath = "/tmp/conf/busi.conf"

type Record struct {
	rightRecord *RightRecord
	detectAs []string
	detectCNames []string
	authServerNames string
	timeoutFlag bool
	result string
	compareType string
}

type Task struct {
	taskID string
	taskName string
	uuid string
	records []*Record
}

func GetTaskConfig() (task *Task, err error) {
	task = new(Task)
	taskConfigBase64, err := ioutil.ReadFile(ConfPath)
	if err != nil {
        return nil, err
    }
	taskConfigB, err := base64.StdEncoding.DecodeString(string(taskConfigBase64))
	if err != nil {
        return nil, err
    }
	taskConfig := strings.Split(string(taskConfigB), ";")
	fmt.Println(taskConfig)

	task.taskID = taskConfig[1]

	//组合域名、正确值
	domains := taskConfig[4:len(taskConfig)-1]
	rightRecords, err := getRightValue(domains)
	if err != nil {
        return nil, err
    }
	for _, rightRecord := range rightRecords {
		record := new(Record)
		record.rightRecord = rightRecord
		task.records = append(task.records, record)
	}
	task.taskName = taskConfig[2]
	task.uuid = taskConfig[len(taskConfig)-1]
	return
}
