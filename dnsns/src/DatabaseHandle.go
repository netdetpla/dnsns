package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	//"fmt"
)

type RightRecord struct {
	domain      string
	rightAs     []string
	rightCNames []string
}

func SplitAs(AsStr string) (As []string) {
	TempAs := strings.Split(AsStr, ",")
	for _, a := range TempAs {
		As = append(As, strings.Split(a, "-")[0])
	}
	return
}

func SplitCNames(CNamesStr string) (CNames []string) {
	CNames = strings.Split(CNamesStr, ",")
	return
}

func getRightValue(domains []string) (rightRecords []*RightRecord, err error) {
	//TODO 备份
	db, err := sql.Open(
		"mysql",
		"root:123456@tcp(10.96.129.133:3306)/cncert_initiative_probe_system?timeout=20s")
		//"zyq:123456@tcp(10.96.129.6:3306)/cncert_initiative_probe_system")
	db.SetMaxOpenConns(50);
	db.SetMaxIdleConns(1);
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	defer db.Close()
	//查询数据
	querySQL := "SELECT dns_domain_name, ns_A FROM domain_library WHERE domain_name=?"
	//fmt.Println(len(domains))
	for _, domain := range domains {
		if len(domain) == 0 {
			continue
		}
		rows, err := db.Query(querySQL, domain)
		if err != nil {
			return nil, err
		}
		if !rows.Next() {
			continue
		}
		var (
			rightCNamesStr string
			rightAsStr     string
		)
		err = rows.Scan(&rightCNamesStr, &rightAsStr)
		if err != nil {
			return nil, err
		}
		rightRecord := &RightRecord{
			domain:      domain,
			rightAs:     SplitAs(rightAsStr),
			rightCNames: SplitCNames(rightCNamesStr),
		}
		rightRecords = append(rightRecords, rightRecord)
		rows.Close()
	}
	return
}
