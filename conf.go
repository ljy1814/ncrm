package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type NCrmConf struct {
	Filename     string
	ServiceName  string
	ServiceAddr  string
	ServicePort  string
	MajorDomains map[string]string `json:"major_domains"`
	//	SandboxMajor  string
	//	TestMajor     string
	//	RouterMajor   string
	//	EurouterMajor string
	//	UsrouterMajor string
	//	EarouterMajor string
	Servers map[int64]string `json:"servers"`
}

//var CrmConf map[string]interface{}
var CrmConf NCrmConf

func ReadConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	fstream, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(fstream, &CrmConf)
}
