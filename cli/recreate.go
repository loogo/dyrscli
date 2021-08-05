package cli

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func Recreate() {

	// 1. cmd接收输入参数
	// param host: Debezium 连接器服务端地址, 示例: http://127.0.0.1:8083
	// param connect_name : 指定的并且已经存在的连接器Source或Sink名称

	DebeziumServerPtr := flag.String("server", "http://172.16.10.246:8083", "Debezium server address, sample: http://127.0.0.1:8083")
	DebeziumConnectNamePtr := flag.String("connect_name", "sample-sink", "Source or Sink name")

	flag.Parse()

	// paramters verification

	if *DebeziumConnectNamePtr == "sample-sink" {
		log.Fatalln("请输入一个有效的连接器名称")
		return
	}

	httpClient := &http.Client{}

	// 2. HTTP GET 获取连接器config
	log.Println("Get connector...")

	reqOne, err := http.NewRequest(http.MethodGet, *DebeziumServerPtr+"/connectors/" +*DebeziumConnectNamePtr+ "/config", nil)

	if err != nil {
		log.Println("获取连接器失败:", err)
	}

	respOne, err := httpClient.Do(reqOne)
	defer respOne.Body.Close()

	robots, err := ioutil.ReadAll(respOne.Body)

	connectConfigString := string(robots)

	fmt.Println(connectConfigString)

	// 3. HTTP DELETE 删除已经存在的连接器
	log.Println("Delete connector...")

	reqTwo, err := http.NewRequest(http.MethodDelete, *DebeziumServerPtr+"/connectors/" +*DebeziumConnectNamePtr, nil)

	if err != nil {
		log.Println("删除连接器失败:", err)
	}

	respTwo, err := httpClient.Do(reqTwo)
	defer respTwo.Body.Close()

	// 4. HTTP PUT 重建这个连接器
	log.Println("Recreate connector...")


	bodyThree := strings.NewReader(connectConfigString)

	reqThree, err := http.NewRequest(http.MethodPut, *DebeziumServerPtr+"/connectors/" +*DebeziumConnectNamePtr+ "/config", bodyThree)
	reqThree.Header.Set("Content-Type", "application/json")

	if err != nil {
		log.Println("重建连接器失败:", err)
	}

	respThree, err := httpClient.Do(reqThree)
	defer respThree.Body.Close()

	log.Println("Recreate finished...")
}