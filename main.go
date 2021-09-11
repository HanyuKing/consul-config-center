package main

import (
	"fmt"
	"net/http"
	"unite.config.center/config"

	"github.com/ngaut/log"
)

func main() {
	Init()

	http.HandleFunc("/keys", myHandler) //	测试. 打印出所有的key, value
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	for key, val := range config.KVMap {
		s := fmt.Sprintf("key: %s, val: %s", key, val)
		fmt.Println(s)
	}
}

func Init() {
	initConfigCenter()
}

func initConfigCenter() {
	path := "my/service_config" // 初始化监听的路径
	config.Init(path)
}
