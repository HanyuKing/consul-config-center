package config

import (
	"encoding/json"
	"github.com/ngaut/log"
	"strconv"
)

func GetBoolOrElse(key string, defaultValue bool) bool {
	if v, ok := KVMap[key]; ok {
		if retV, err := strconv.ParseBool(v.(string)); err == nil {
			return retV
		}
	}
	return defaultValue
}

func GetStringOrElse(key string, defaultValue string) string {
	if v, ok := KVMap[key].(string); ok {
		return v
	}
	return defaultValue
}

func GetIntOrElse(key string, defaultValue int) int {
	v, ok := KVMap[key]
	if !ok {
		return defaultValue
	}

	intV, err := strconv.Atoi(v.(string))
	if err != nil {
		log.Errorf("config.center: GetInt value: %s, error: %s", v, err)
		return defaultValue
	}
	return intV
}

func GetInt64OrElse(key string, defaultValue int64) int64 {
	v, ok := KVMap[key]
	if !ok {
		return defaultValue
	}

	intV, err := strconv.ParseInt(v.(string), 10, 64)
	if err != nil {
		log.Errorf("config.center: GetInt64 value: %s, error: %s", v, err)
		return defaultValue
	}
	return intV
}

func GetUint64OrElse(key string, defaultValue uint64) uint64 {
	v, ok := KVMap[key]
	if !ok {
		return defaultValue
	}

	intV, err := strconv.ParseUint(v.(string), 10, 64)
	if err != nil {
		log.Errorf("config.center: GetUint64 value: %s, error: %s", v, err)
		return defaultValue
	}
	return intV
}

func GetStringMapOrElse(key string, defaultValue map[string]string) map[string]string {
	v, ok := KVMap[key]
	if !ok {
		return defaultValue
	}

	strV, ok := v.(string)
	if !ok {
		return defaultValue
	}

	mapV := make(map[string]string)

	err := json.Unmarshal([]byte(strV), &mapV)
	if err != nil {
		log.Errorf("config.center: GetStringMapOrElse value: %s, error: %s", v, err)
		return defaultValue
	}
	return mapV
}

func GetIntMapOrElse(key string, defaultValue map[string]int64) map[string]int64 {
	v, ok := KVMap[key]
	if !ok {
		return defaultValue
	}

	strV, ok := v.(string)
	if !ok {
		return defaultValue
	}

	mapV := make(map[string]int64)

	err := json.Unmarshal([]byte(strV), &mapV)
	if err != nil {
		log.Errorf("config.center: GetStringMapOrElse value: %s, error: %s", v, err)
		return defaultValue
	}
	return mapV
}
