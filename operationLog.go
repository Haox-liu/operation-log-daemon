package main

import (
	"os"
	"./SqlOperate"
	"time"
	"./HttpOperate"
	"strconv"
	"strings"
	"fmt"
)

var zeusIP = "10.160.0.85"
var zeusPort = "8085"
//间隔时间
var intervalTime = 10

func main() {
	//getEnv()
	//启用断续器，每隔10秒执行
	//ticker := time.NewTicker(time.Duration(intervalTime) * time.Second)
	//for _=range ticker.C {
		logArray, err := SqlOperate.Query()
		if err != nil {
			fmt.Println("查询操作日志记录异常！ " + time.Now().Format("2006-01-02 15:04:05 PM"))
			fmt.Println(err)
			//continue
		} else {
			if len(logArray) > 0 {
				fmt.Println("成功查询" + strconv.Itoa(len(logArray)) + "操作日志！ " + time.Now().Format("2006-01-02 15:04:05 PM"))
			}
		}
		if len(logArray) > 0 {
			result, err := HttpOperate.PostOperationLog(logArray,zeusIP,zeusPort)
			if err != nil {
				fmt.Println("请求zeus异常！ " + time.Now().Format("2006-01-02 15:04:05 PM"))
				fmt.Println(err)
				//continue
			} else {
				fmt.Println("请求zeus成功！ " + time.Now().Format("2006-01-02 15:04:05 PM"))
			}
			if strings.Compare(result, "Success") == 0 {
				err := SqlOperate.Update(logArray)
				if err != nil {
					fmt.Println("更新操作日志记录异常！ " + time.Now().Format("2006-01-02 15:04:05 PM"))
					fmt.Println(err)
					//continue
				}
			}
		}
	//}
}

func getEnv(){
	ServerIP := os.Getenv("ZUES_IP")
	if ServerIP != "" {
		zeusIP = ServerIP
	}
	ServerPort:= os.Getenv("ZUES_PORT")
	if ServerPort != "" {
		zeusPort = ServerPort
	}
	interval_time, _ := strconv.Atoi(os.Getenv("INTERVAL_TIME"))
	if interval_time != 0 {
		intervalTime = interval_time
	}
}