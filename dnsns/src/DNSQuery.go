package main

import (
	"fmt"
	"github.com/miekg/dns"
	"github.com/weppos/publicsuffix-go/publicsuffix"
	"runtime"
	"strings"
)

var quit = make(chan error)
var ctrl = make(chan int, 50)

func ParseRR(rrs []dns.RR) (as []string, cNames []string) {
	for _, rr := range rrs {
		rrElements := strings.Split(rr.String(), "\t")
		fmt.Println(rrElements)
		if len(rrElements) == 5 {
			if rrElements[3] == "CNAME" {
				cName := string([]rune(rrElements[4])[:len(rrElements[4])-1])
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

func ParseRRNS(rrs []dns.RR) (nss []string) {
	for _, rr := range rrs {
		rrElements := strings.Split(rr.String(), "\t")
		if len(rrElements) == 5 {
			if rrElements[3] == "NS" {
				nss = append(nss, rrElements[4][:len(rrElements[4])-1])
			}
		}
	}
	return
}

func ParseRRA(rrs []dns.RR) (as []string) {
	for _, rr := range rrs {
		rrElements := strings.Split(rr.String(), "\t")
		if len(rrElements) == 5 {
			if rrElements[3] == "A" {
				as = append(as, rrElements[4])
			}
		}
	}
	return
}

func SendDNSQuery(record *Record) {
	fmt.Println("t start")
	//探测NS
	m := new(dns.Msg)
	domain, err := publicsuffix.Domain(record.rightRecord.domain)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeNS)
	errCount := 3
NSStart:
	in, err := dns.Exchange(m, record.reServer+":53")
	if err != nil {
		if strings.Index(err.Error(), "timeout") >= 0 {
			if errCount == 0 {
				record.timeoutFlag = true
			} else {
				errCount--
				goto NSStart
			}
		}
	} else {
		record.timeoutFlag = false
		record.detectCNames = ParseRRNS(in.Answer)
	}
	fmt.Println("DNS NS")
	fmt.Println(record.detectCNames)
	//探测NS A
	for _, ns := range record.detectCNames {
		m2 := new(dns.Msg)
		m2.SetQuestion(dns.Fqdn(ns), dns.TypeA)
		errCount := 3
	AStart:
		in2, err := dns.Exchange(m2, record.reServer+":53")
		if err != nil {
			if errCount == 0 {
				continue
			} else {
				errCount--
				goto AStart
			}
		} else {
			fmt.Println(in2.Answer)
			record.detectAs = append(record.detectAs, ParseRRA(in2.Answer)...)
			fmt.Println("DNS A")
			fmt.Println(record.detectAs)
		}
	}
	<- ctrl
	fmt.Println("t quit")
	quit <- nil
	return
}

func ControlDNSQueryRoutine(tasks *Task) (err error) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println(len(tasks.records))
	for _, record := range tasks.records {
		fmt.Println("ctrl start")
		ctrl <- 1
		go SendDNSQuery(record)
	}
	for i := 0; i < len(tasks.records); i++ {
		err = <-quit
		fmt.Println("ctrl quit")
	}
	close(ctrl)
	close(quit)
	return err
}
