package main

import (
	"errors"
	"regexp"
)

type quakeResult struct {
	QuakeQuery string
	State      string
	ResultNum  string
	RawBody    string
}

func (q *quakeResult) checkResultStatus(resBody string) (string, error) {
	re := regexp.MustCompile(`"message":"(.*?)",`)
	if !re.MatchString(resBody) {
		return "", errors.New("解析错误")
	}
	return re.FindStringSubmatch(resBody)[1], nil
}

func (q *quakeResult) extractResultSize(resBody string) (string, error) {
	re := regexp.MustCompile(`"total":(.*?)}}}`)
	if !re.MatchString(resBody) {
		return "", errors.New("解析错误")
	}
	return re.FindStringSubmatch(resBody)[1], nil
}

//func handleBody(resBody string) (*quakeResult, error) {
//	queryResult := &quakeResult{}
//	state, err := checkResultStatus(resBody)
//	if err != nil {
//		return nil, err
//	}
//	if state != "Successful." {
//		queryResult.State = state
//		return queryResult, nil
//	}
//	num, _ := extractResultSize(resBody)
//	queryResult.State = state
//	queryResult.Num = num
//	return queryResult, nil
//}
//
//func QuakeSearchAll(quakeQuery string) (*quakeResult, error) {
//	bodyNew, _ := json.Marshal(map[string]string{
//		"query":        quakeQuery,
//		"start":        strconv.Itoa(0),
//		"size":         strconv.Itoa(10),
//		"ignore_cache": "False",
//		"latest":       "True",
//	})
//	client := HttpClient{
//		Timeout: 25 * time.Second,
//		Proxy:   "http://127.0.0.1:8080",
//		Url:     "https://quake.360.net/api/v3/search/quake_service",
//		Header: map[string]string{
//			"X-QuakeToken": "4903fae5-455b-482c-a520-89160bceacf0",
//			"Content-Type": "application/json",
//		},
//		Body: string(bodyNew),
//	}
//	res, err := client.Post()
//	resbdoy, err := io.ReadAll(res.Body)
//	if err != nil {
//		return nil, err
//	}
//	defer res.Body.Close()
//	query_Result := &quakeResult{}
//	query_Result.checkResultStatus(string(resbdoy))
//
//	return handleBody(string(resbdoy))
//}
