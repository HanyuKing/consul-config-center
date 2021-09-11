package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/ngaut/log"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	scheme = "http"
)

func GetKV(key string) (string, error) {
	dc, err := getDataCenter()
	if err != nil {
		return "", err
	}
	if strings.Index(key, "/") == 0 {
		key = key[1:]
	}

	var requestURI = "/v1/kv/"
	url := fmt.Sprintf("%s://%s%s%s?dc=%s", scheme, GetConsulAddr(), requestURI, getPath()+key, dc)
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("config.center: GetKV occur an error. %s", err)
		return "", err
	}

	defer resp.Body.Close()

	var entries []*api.KVPair
	if err := decodeBody(resp, &entries); err != nil {
		log.Errorf("config.center: GetKV decodeBody occur an error: %s", err)
		return "", err
	}
	if len(entries) > 0 {
		return string(entries[0].Value), nil
	}
	return "", errors.New("no value")
}

// decodeBody is used to JSON decode a body
func decodeBody(resp *http.Response, out interface{}) error {
	dec := json.NewDecoder(resp.Body)
	return dec.Decode(out)
}

func GetPathKeys(path string) ([]string, error) {
	dc, err := getDataCenter()
	if err != nil {
		return []string{}, err
	}

	var requestURI = "/v1/kv/"
	url := fmt.Sprintf("%s://%s%s%s?keys&dc=%s", scheme, GetConsulAddr(), requestURI, path, dc)

	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("config.center: GetPathKeys occur an error. %s", err)
		return []string{}, err
	}

	defer resp.Body.Close()

	bodyB, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("config.center: GetPathKeys read resp occur an error: %s", err)
		return []string{}, err
	}

	if len(bodyB) == 0 {
		return []string{}, nil
	}

	var keys []string
	err = json.Unmarshal(bodyB, &keys)
	if err != nil {
		log.Errorf("config.center: GetPathKeys Unmarshal resp occur an error: %s", err)
		return []string{}, err
	}
	return keys, nil
}

func getDataCenter() (string, error) {

	var requestURI = "/v1/agent/self"

	resp, err := http.Get(fmt.Sprintf("%s://%s%s", scheme, GetConsulAddr(), requestURI))
	if err != nil {
		log.Errorf("config.center: getDataCenter occur an error. %s", err)
		return "", err
	}

	defer resp.Body.Close()

	bodyB, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("config.center: getDataCenter read resp occur an error: %s", err)
		return "", err
	}

	var out map[string]map[string]interface{}
	err = json.Unmarshal(bodyB, &out)
	if err != nil {
		log.Errorf("config.center: getDataCenter Unmarshal resp occur an error: %s", err)
		return "", errors.New("config.center: getDataCenter Unmarshal resp occur an error")
	}

	cfg, ok := out["Config"]

	if !ok {
		return "", errors.New("config.center:: self.Config not found")
	}

	dc, ok := cfg["Datacenter"].(string)
	if !ok {
		return "", errors.New("config.center:: self.DataCenter not found")
	}

	return dc, nil
}
