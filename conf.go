package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type NCrmConf struct {
	Filename      string
	ServiceName   string
	ServiceAddr   string
	ServicePort   string
	SandboxMajor  string
	TestMajor     string
	RouterMajor   string
	EurouterMajor string
	UsrouterMajor string
	EarouterMajor string
}

var CrmConf map[string]interface{}

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
