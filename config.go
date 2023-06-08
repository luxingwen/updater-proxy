package updaterproxy

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	ProxyId    string   `json:"proxyId"`
	Servers    []string `json:"servers"`
	PkgServers []string `json:"pkgservers"`
	Port       string   `json:"port"`
}

var conf *Config

func init() {

	b, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	conf = &Config{}
	err = json.Unmarshal(b, conf)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func GetConfig() *Config {
	return conf
}
