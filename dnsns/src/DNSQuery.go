package main

import (
	"fmt"
	"github.com/miekg/dns"
	"runtime"
	"strings"
)

var quit = make(chan error)

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
				nss = append(nss, rrElements[4])
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
	//探测NS
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(record.rightRecord.domain), dns.TypeNS)
	in, err := dns.Exchange(m, record.reServer+":53")
	if err != nil {
		if strings.Index(err.Error(), "timeout") >= 0 {
			record.timeoutFlag = true
		} else {
			quit <- err
			return
		}
	} else {
		record.timeoutFlag = false
		record.detectCNames = ParseRRNS(in.Answer)
	}
	//探测NS A
	for _, ns := range record.detectCNames {
		m2 := new(dns.Msg)
		m2.SetQuestion(dns.Fqdn(ns), dns.TypeA)
		in, err := dns.Exchange(m, record.reServer+":53")
		if err != nil {
			continue
		} else {
			record.detectAs = append(record.detectAs, ParseRRA(in.Answer)...)
		}
	}
	quit <- nil
	return
}

func ControlDNSQueryRoutine(tasks *Task) (err error) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	for _, record := range tasks.records {
		go SendDNSQuery(record)
	}
	for i := 0; i < len(tasks.records); i++ {
		err = <-quit
	}
	close(quit)
	return err
}
