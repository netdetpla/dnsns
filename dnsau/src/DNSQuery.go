package main

import (
	"github.com/miekg/dns"
	"runtime"
	"strings"
	"fmt"
)

var quit = make(chan error)
var ctrl = make(chan int, 100)

func ParseRR(rrs []dns.RR) (as []string, cNames []string) {
	for _, rr := range rrs {
		rrElements := strings.Split(rr.String(), "\t")
		fmt.Println(rrElements)
		if len(rrElements) == 5 {
			if rrElements[3] == "CNAME" {
				cName := string([]rune(rrElements[4])[:len(rrElements[4]) - 1])
				cNames = append(cNames, cName)
			} else if rrElements[3] == "A" {
				as = append(as, rrElements[4])
			}
		}
	}
	InfoLog("ParseRR")
	fmt.Println(as)
	fmt.Println(cNames)
	return
}

func SendDNSQuery(record *Record) {
    i:=0
    Loop:
    var dig Dig
    rsp,err := dig.Trace(record.rightRecord.domain)
    if err != nil{
		if i < 3 {
            i++
			goto Loop
        }
        record.timeoutFlag = true
    }
	if len(rsp)!= 0 {
		rsps := rsp[len(rsp)-1]
		record.timeoutFlag = false
		record.detectAs, record.detectCNames = ParseRR(rsps.Msg.Answer)
		record.authServerNames = rsps.Server[:len(rsps.Server)-1] + ":" + rsps.ServerIP[1:len(rsps.ServerIP)-4]
	}else{
		record.timeoutFlag = true
	}
	
	<- ctrl
	quit <- nil
    return
}

func ControlDNSQueryRoutine(tasks *Task) (err error){
	runtime.GOMAXPROCS(runtime.NumCPU())
	for _, record := range tasks.records {
		ctrl <- 1
		go SendDNSQuery(record)
	}
	for i := 0; i < len(tasks.records); i++ {
		err = <- quit
	}
	close(ctrl)
	close(quit)
	return err
}
