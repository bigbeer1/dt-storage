package main

import (
	"database/sql"
	"fmt"
	_ "github.com/taosdata/driver-go/v3/taosRestful"
	"log"
	"time"
)

func main() {

	var taosUri = "root:taosdata@http(127.0.0.1:6041)/"
	taos, err := sql.Open("taosRestful", taosUri)
	if err != nil {
		fmt.Println("failed to connect TDengine, err:", err)
		return
	}

	defer taos.Close()
	createMonitor(taos, 99999)

	time.Sleep(time.Second * 1)
	a := fmt.Sprintf("CREATE STABLE chint.monitor_point (ts TIMESTAMP, data FLOAT) TAGS (tenant_id BINARY(36))")

	CreateStable(taos, a)

	time.Sleep(time.Second * 1)

	fmt.Println("初始化成功")

}

func createMonitor(taos *sql.DB, number int) {
	dataString := fmt.Sprintf("CREATE DATABASE chint BUFFER 50 KEEP %vd  VGROUPS 5  ", number)
	_, err := taos.Exec(dataString)
	if err != nil {
		log.Fatalln("failed to create database, err:", err)
	}
}

func CreateStable(taos *sql.DB, dataString string) {
	_, err := taos.Exec(dataString)
	if err != nil {
		log.Fatalln("failed to create stable, err:", err)
	}
}
