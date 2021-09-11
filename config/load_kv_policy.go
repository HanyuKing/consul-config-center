package config

import (
	"encoding/json"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/ngaut/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var KVMap = make(map[string]interface{})
var path = "" // start with /, end with /

func Init(pathPrefix string) {
	initPath(pathPrefix)
	initKV()
	listening()
}

func initPath(pathPrefix string) {
	path = pathPrefix
	if strings.LastIndex(path, "/") != (len(path) - 1) {
		path = path + "/"
	}
	if strings.Index(path, "/") != 0 {
		path = "/" + path
	}
}

func getPath() string {
	return path
}

func initKV() {
	// load from local snapshots
	KVMap = readFromSnapshots()

	// read from remote
	remoteKVMap := make(map[string]interface{})
	keys, err := GetPathKeys(getPath())
	if err != nil {
		return
	}

	var pathPrefix = getPath()

	for _, key := range keys {
		key = strings.TrimPrefix(key, pathPrefix[1:])
		if key == "" || key == "/" {
			continue
		}
		value, err := GetKV(key)
		if err == nil {
			remoteKVMap[key] = value
		}
	}

	if len(keys) == len(remoteKVMap) {
		KVMap = remoteKVMap
	}

	// update snapshots
	writeSnapshots(KVMap)
}

func listening() {
	params := make(map[string]interface{})
	params["type"] = "keyprefix"
	params["prefix"] = getPath()
	plan, err := watch.Parse(params)
	if err != nil {
		panic(err)
	}
	plan.Handler = func(index uint64, result interface{}) {
		if kvs, ok := result.(api.KVPairs); ok {
			remoteKVMap := make(map[string]interface{})
			for _, kv := range kvs {
				key := strings.TrimPrefix(kv.Key, getPath()[1:])
				if key == "" || key == "/" {
					continue
				}
				remoteKVMap[key] = string(kv.Value)
			}

			KVMap = remoteKVMap

			go writeSnapshots(KVMap)
		}
	}
	go func() {
		// your consul agent addr
		if err = plan.Run(GetConsulAddr()); err != nil {
			panic(err)
		}
	}()
}

func writeSnapshots(kvMap map[string]interface{}) {
	f, _ := os.OpenFile(getFileName(), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)

	defer f.Close()

	kvJSONBytes, err := json.Marshal(kvMap)

	if err != nil {
		log.Errorf("writeSnapshots occur an error. kvMap: %+v, error: %s", kvMap, err)
		return
	}

	f.Write(kvJSONBytes)
}

func getFileName() string {
	dir, filePath := filepath.Split("data/config_center")
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0777); err != nil {
			println(err)
		}
	}
	var fileName = filepath.Join(dir, filePath+".snapshots")
	return fileName
}

func readFromSnapshots() map[string]interface{} {
	configData := make(map[string]interface{})

	fileName := getFileName()
	_, err := os.Stat(fileName)
	if err != nil && os.IsNotExist(err) {
		log.Warnf("file [%s] not exists, so create it later ...", fileName)
		return configData
	}

	fileDataB, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Errorf("readFromSnapshots read file occur an error. error: %s", err)
		return configData
	}

	err = json.Unmarshal(fileDataB, &configData)
	if err != nil {
		log.Errorf("readFromSnapshots Unmarshal occur an error. error: %s", err)
		return configData
	}
	return configData
}
