package SqlOperate

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"encoding/json"
	"os"
	"strconv"
	"time"
)

var mysqlIP string = "10.158.224.15"
var mysqlPort string = "32759"
var mysqlUsername string = "root"
var mysqlPassword string = "123456"
var databaseName string = "cloud"
var tableName string = "operatelog"
var maxReadNum string = "100"

func Query() ([]string, error) {
	var str []string
	//获取环境变量中mysql参数
	getMysqlEnv()
	//dataSourceName := "root:123456@tcp(10.158.224.15:32759)/cloud?charset=utf8"
	dataSourceName :=  mysqlUsername + ":" + mysqlPassword + "@tcp(" + mysqlIP + ":" + mysqlPort + ")" + "/" + databaseName +"?charset=utf8"
	//连接数据库
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		fmt.Println("连接数据库错误！ " + time.Now().Format("2006-01-02 15:04:05 PM"))
		return str, err
	}
	//结束连接
	defer db.Close();
	//查询
	querySql := "SELECT * FROM " + tableName + " WHERE readFlag=false LIMIT " + maxReadNum
	rows, err := db.Query(querySql)
	if err != nil {
		fmt.Println("查询数据错误！ " + time.Now().Format("2006-01-02 15:04:05 PM"))
		return str, err
	}

	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	var logArray []string
	//將每条数据转为map，并json序列化后放到logArray中
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		record := make(map[string]string)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}
		log, _ := json.Marshal(record)
		logArray = append(logArray,string(log[:]))
	}

	return logArray, err
}



func Update(logArray []string) error {
	dataSourceName :=  mysqlUsername + ":" + mysqlPassword + "@tcp(" + mysqlIP + ":" + mysqlPort + ")" + "/" + databaseName + "?charset=utf8"
	//连接数据
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		fmt.Println("连接数据库错误！ " + time.Now().Format("2006-01-02 15:04:05 PM"))
		return err
	}
	//结束连接
	defer db.Close();

	len := len(logArray)
	tx,_ := db.Begin()
	for i := 0; i < len; i++ {
		log := logArray[i]
		var logMap map[string]interface{}
		if err := json.Unmarshal([]byte(log), &logMap); err == nil {
			logId,_ := strconv.Atoi(logMap["id"].(string))
			updateSql := "UPDATE " + tableName + " SET readFlag=true WHERE id=?"
			//执行更新
			tx.Exec(updateSql, logId)
		} else {
			if err != nil {
				fmt.Println("log转logMap错误！" + time.Now().Format("2006-01-02 15:04:05 PM"))
				return err
			}
		}
	}
	tx.Commit()
	fmt.Println("已更新 " + strconv.Itoa(len) + " 条操作日志! " + time.Now().Format("2006-01-02 15:04:05 PM"))

	return err
}

func getMysqlEnv(){
	serverIP := os.Getenv("MYSQL_IP")
	if serverIP != "" {
		mysqlIP = serverIP
	}
	serverPort := os.Getenv("MYSQL_PORT")
	if serverPort != "" {
		mysqlPort = serverPort
	}
	username := os.Getenv("MYSQL_USERNAME")
	if username != "" {
		mysqlUsername = username
	}
	password := os.Getenv("MYSQL_PASSWORD")
	if password != "" {
		mysqlPassword = password
	}
	database := os.Getenv("DATABASE_NAME")
	if database != "" {
		databaseName = database
	}
	table := os.Getenv("TABLE_NAME")
	if table != "" {
		tableName = table
	}
	max_read_num := os.Getenv("MAX_READ_NUM")
	if max_read_num != "" {
		maxReadNum= max_read_num
	}
}
