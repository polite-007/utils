package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
	"time"
)

type fofa struct{}

type fofaResult struct {
	FofaQuery        string
	State            bool
	ResultNum        string
	NewFofaQuery     string
	IsHistoryAppName bool
	RawBody          string
}

var Fofa = &fofa{}

func (f *fofaResult) checkResultStatus() {
	if strings.Contains(f.RawBody, "查询语法错误") && strings.Contains(f.RawBody, "820000") {
		f.State = false
	} else if strings.Contains(f.RawBody, "规则不存在") && strings.Contains(f.RawBody, "811001") {
		f.State = false
	} else if strings.Contains(f.RawBody, "\"error\":false,") {
		f.State = true
	} else {
		f.State = false
	}
}

func (f *fofaResult) updateQuery() error {
	re := regexp.MustCompile(`"query":(.*?),`)
	if !re.MatchString(f.RawBody) {
		return errors.New("无法匹配到查询语法错误")
	}
	resQuery := strings.ReplaceAll(strings.Trim(re.FindStringSubmatch(f.RawBody)[1], "\""), "\\", "") + "\""
	if f.FofaQuery == resQuery {
		f.IsHistoryAppName = false
	} else {
		f.IsHistoryAppName = true
	}
	f.NewFofaQuery = strings.ReplaceAll(strings.Trim(re.FindStringSubmatch(f.RawBody)[1], "\""), "\\", "") + "\""
	return nil
}

func (f *fofaResult) extractResultSize() error {
	re := regexp.MustCompile(`"size":(.*?),`)
	if !re.MatchString(f.RawBody) {
		return errors.New("无法匹配到查询语法错误")
	}
	f.ResultNum = re.FindStringSubmatch(f.RawBody)[1]
	return nil
}

func (f *fofa) FofaSearchAll(fofaQuery string, proxy string) (*fofaResult, error) {
	client := &HttpClient{
		Timeout: 25 * time.Second,
		Proxy:   proxy,
		Url:     "https://fofa.info/api/v1/search/all?&key=12ca5f1dfc5490ff9ac4a018b25aeec1&qbase64=" + base64.StdEncoding.EncodeToString([]byte(fofaQuery)),
	}
	res, err := client.Get()
	if err != nil {
		log.Fatalf("Error on request for %s: %v\n", fofaQuery, err)
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error on request for %s: %v\n", fofaQuery, err)
		return nil, err
	}
	defer res.Body.Close()

	fofa_result := &fofaResult{
		FofaQuery: fofaQuery,
		RawBody:   string(resBody),
	}
	fofa_result.checkResultStatus()

	if fofa_result.State {
		fofa_result.extractResultSize()
		fofa_result.updateQuery()
	}
	return fofa_result, nil
}

func (f *fofa) Fofasearch(fofaquerys []string, proxy string) []string {
	resultPrints := []string{}
	for i, fofaquery := range fofaquerys {
		fofaquery = fofaquery
		resultPrint, err := f.FofaSearchAll(fofaquery, proxy)
		if err != nil {
			resultPrints = append(resultPrints, err.Error())
			fmt.Printf("%d:%s\n", i, err.Error())
			continue
		}
		if !resultPrint.State {
			resultPrints = append(resultPrints, "false")
			fmt.Printf("%d:%s\n", i, "false")
			continue
		}
		resultPrints = append(resultPrints, resultPrint.ResultNum)
		fmt.Printf("%d:%s, num:%s\n", i, "true", resultPrint.ResultNum)
	}
	return resultPrints
}
