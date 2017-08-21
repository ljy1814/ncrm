package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type DomainHandler struct {
}

func NewDomainHandler() *DomainHandler {
	return &DomainHandler{}
}

func (this *DomainHandler) handleListDomain(w http.ResponseWriter, r *http.Request) {
	oUrl := r.URL
	reqBody, err := r.GetBody()
	if err != nil {

	}
	defer func() {
		r.Body = reqBody
		r.URL = oUrl
	}()

}

func getMajorDomain(r *http.Request) (majorDomain string, majorDomainId, developerId int64, err error) {
	if r.Method != "POST" {
		return
	}
	var url string
	var client http.Client
	var resp http.ResponseWriter
	for k, host := range Hosts {
		url = "http://" + host + "/zc-platform/v1/getMajorDomain"
		resp, err = client.Do(r)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		var respBody []byte
		resyBody, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		var res map[string]interface{}
		err = json.Unmarshal(respBody, &res)
		if err != nil {
			return res["majorDomain"], res["majorDomainId"], res["developerId"], nil
		}
	}
	return
}

func getDeveloper(r *http.Request) (majorDomain string, majorDomainId, developerId int64, err error) {
	if r.Method != "POST" {
		return
	}
	var url string
	var client http.Client
	var resp http.ResponseWriter
	for k, host := range Hosts {
		url = "http://" + host + "/zc-platform/v1/getMajorDomain"
		req, err = http.NewRequest("POST", url, nil)
		resp, err = client.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		var respBody []byte
		resyBody, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		var res map[string]interface{}
		err = json.Unmarshal(respBody, &res)
		if err != nil {
			return res["majorDomain"], res["majorDomainId"], res["developerId"], nil
		}
	}
	return
}
