package HttpOperate

import (
	"io/ioutil"
	"encoding/json"
	"net/http"
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"
)
type ResponseInfo struct {
	Result string `json:"result"`
	data []string `json:"data"`
}

func PostOperationLog(logArray []string, zeusIP string, zeusPort string) (string, error) {

	url := "http://" + zeusIP + ":" + zeusPort + "/manage/control/operationLog/save"
	jsonReq, _ := json.Marshal(logArray)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonReq))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求zeus错误！ " + time.Now().Format("2006-01-02 15:04:05 PM"))
		return "", err
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Status:", resp.Status)
	statusCodeStr := strings.Split(resp.Status, " ")[0]
	statusCode,err := strconv.Atoi(statusCodeStr)
	//fmt.Println("StatusCode:", statusCode)
	if statusCode > 300 || statusCode < 200 {
		fmt.Println("请求zeus响应异常！，状态码为 " + resp.Status + "! " + time.Now().Format("2006-01-02 15:04:05 PM"))
		return "", err
	}
	//fmt.Println("response Headers:", resp.Header)
	//fmt.Println("response Body:", string(respBody))
	var response ResponseInfo
	err = json.Unmarshal([]byte(respBody), &response)
	if err != nil {
		fmt.Println("respBody转为ResponseInfo结构体错误！ " + time.Now().Format("2006-01-02 15:04:05 PM"))
		return "", err
	}

	return response.Result, nil
}