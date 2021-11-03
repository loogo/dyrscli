package cli

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func Recreate(server string, connect_name string) {

	// paramters verification

	if connect_name == "sample-sink" {
		log.Fatalln("请输入一个有效的连接器名称")
		return
	}

	httpClient := &http.Client{}

	// 2. HTTP GET 获取连接器config
	log.Println("Get connector...")

	reqOne, err := http.NewRequest(http.MethodGet, server+"/connectors/" +connect_name+ "/config", nil)

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

	reqTwo, err := http.NewRequest(http.MethodDelete, server+"/connectors/" +connect_name, nil)

	if err != nil {
		log.Println("删除连接器失败:", err)
	}

	respTwo, err := httpClient.Do(reqTwo)
	defer respTwo.Body.Close()

	// 4. HTTP PUT 重建这个连接器
	log.Println("Recreate connector...")


	bodyThree := strings.NewReader(connectConfigString)

	reqThree, err := http.NewRequest(http.MethodPut, server+"/connectors/" +connect_name+ "/config", bodyThree)
	reqThree.Header.Set("Content-Type", "application/json")

	if err != nil {
		log.Println("重建连接器失败:", err)
	}

	respThree, err := httpClient.Do(reqThree)
	defer respThree.Body.Close()

	log.Println("Recreate finished...")
}