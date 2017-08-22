package main

import (
	"ac-common-go/developer"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/fatih/color"
)

type DomainHandler struct {
}

func NewDomainHandler() *DomainHandler {
	return &DomainHandler{}
}

func (this *DomainHandler) handleListDomain(w http.ResponseWriter, r *http.Request) {
	oUrl := r.URL
	reqBody := r.Body
	reqURI := r.RequestURI
	oldDevSign := r.Header.Get("X-Zc-Developer-Signature")
	defer func() {
		r.Body = reqBody
		r.URL = oUrl
		r.RequestURI = reqURI
		r.Header.Set("X-Zc-Developer-Signature", oldDevSign)
	}()
	r.Header.Set("X-Zc-Developer-Id", fmt.Sprintf("%d", SignDeveloperId))
	r.Header.Set("X-Zc-Access-Key", AccessKey)
	r.Header.Set("X-Zc-Major-Domain", "test_ac")
	var err error
	domain, domainId, developerId, err := getMajorDomain(r)
	w.Write([]byte(fmt.Sprintf("domain: %v , domainId: %v, developerId: %v\r\n", domain, domainId, developerId)))

	company, err := getDeveloper(r, fmt.Sprintf("%d", developerId))
	if err != nil {
		w.Write([]byte(fmt.Sprintf("internal error: %v\r\n", err)))
	}
	w.Write([]byte(fmt.Sprintf("company : %v\n", company)))

	/*
		allInfo, err := getAllProjects(r, fmt.Sprintf("%d", developerId))
		if err != nil {
			w.Write([]byte(fmt.Sprintf("%v\n", allInfo)))
		}
	*/
}

//func getMajorDomain(r *http.Request) (majorDomain interface{}, majorDomainId, developerId interface{}, err error) {
func getMajorDomain(r *http.Request) (majorDomain string, majorDomainId, developerId int64, err error) {
	if r.Method != "POST" {
		return
	}
	var client http.Client
	var resp *http.Response
	for _, host := range Hosts {
		sUrl := "http://" + host + "/zc-platform/v1/getMajorDomain"
		method := "getMajorDomain"
		var req *http.Request
		jbody := []byte(`{"majorDomainId":` + r.Header.Get("X-Zc-Major-Domain-Id") + `}`)
		req, err = http.NewRequest("POST", sUrl, bytes.NewBuffer(jbody))
		req.Header.Set("X-Zc-Developer-Id", fmt.Sprintf("%d", SignDeveloperId))
		req.Header.Set("X-Zc-Access-Key", AccessKey)
		req.Header.Set("X-Zc-Major-Domain", SignDomain)
		var signature string
		req.Header.Set("X-Zc-Timeout", r.Header.Get("X-Zc-Timeout"))
		req.Header.Set("X-Zc-Timestamp", fmt.Sprintf("%d", time.Now().Unix()))
		req.Header.Set("X-Zc-Nonce", r.Header.Get("X-Zc-Nonce"))
		req.Header.Set("Content-Type", "application/x-zc-object")

		signature, err = genAccessSignature(req, method)
		if err != nil || signature == "" {
			color.Red("gen signature error: %v\n", err)
			return
		}
		req.Header.Set("X-Zc-Developer-Signature", signature)

		color.Blue("major domain req : %v\n", req)
		resp, err = client.Do(req)
		color.Yellow("resp: %v, err: %v\n", resp, err)
		if err != nil {
			//TODO 检查error
			color.Green("[%s] send request failed: err[%v]\n", host, err)
			break
		}
		msgName := resp.Header.Get("X-Zc-Msg-Name")
		if msgName == "X-Zc-Err" {
			//TODO 请求错误
			color.Green("[%s] get result failed: err[%v]\n", host, err)
			continue
		}
		defer resp.Body.Close()
		var respBody []byte
		respBody, err = ioutil.ReadAll(resp.Body)
		color.Red("major domain respBody : %v\n", string(respBody))
		if err != nil {
			return
		}
		var res map[string]interface{}
		err = json.Unmarshal(respBody, &res)
		if err == nil {
			majorDomain = GetString(res["majorDomain"])
			majorDomainId = GetInt(res["majorDomainId"])
			developerId = GetInt(res["developerId"])
			return
		}
	}
	return
}

func getDeveloper(r *http.Request, developerIdStr string) (company string, err error) {
	if r.Method != "POST" {
		return
	}
	var client http.Client
	var resp *http.Response
	for _, host := range Hosts {
		sUrl := "http://" + host + "/zc-platform/v1/getDeveloper"
		method := "getDeveloper"
		jbody := []byte(`{"developerId":` + developerIdStr + `}`)
		//rbody := bytes.NewBuffer(jbody)
		var req *http.Request
		req, err = http.NewRequest("POST", sUrl, bytes.NewBuffer(jbody))
		req.Header.Set("X-Zc-Developer-Id", fmt.Sprintf("%d", SignDeveloperId))
		req.Header.Set("X-Zc-Access-Key", AccessKey)
		req.Header.Set("X-Zc-Major-Domain", SignDomain)
		var signature string
		req.Header.Set("X-Zc-Timeout", r.Header.Get("X-Zc-Timeout"))
		req.Header.Set("X-Zc-Timestamp", fmt.Sprintf("%d", time.Now().Unix()))
		req.Header.Set("X-Zc-Nonce", r.Header.Get("X-Zc-Nonce"))
		req.Header.Set("Content-Type", "application/x-zc-object")

		signature, err = genAccessSignature(req, method)
		if err != nil || signature == "" {
			color.Red("gen signature error: %v\n", err)
			return
		}
		req.Header.Set("X-Zc-Developer-Signature", signature)

		color.Blue("major domain req : %v\n", req)
		resp, err = client.Do(req)
		color.Yellow("resp: %v, err: %v\n", resp, err)
		if err != nil {
			//TODO 检查error
			color.Green("[%s] send request failed: err[%v]\n", host, err)
			break
		}
		msgName := resp.Header.Get("Z-Zc-Msg-Name")
		if msgName == "X-Zc-Err" {
			color.Green("[%s] get result failed: err[%v]\n", host, err)
			continue
		}
		defer resp.Body.Close()
		var respBody []byte
		respBody, err = ioutil.ReadAll(resp.Body)
		color.Red("developer respBody : %v\n", string(respBody))
		if err != nil {
			return
		}
		var res map[string]interface{}
		err = json.Unmarshal(respBody, &res)
		if err == nil {
			company = GetString(res["company"])
			return
		}
	}
	return
}

func getAllProjects(r *http.Request, developerIdStr string) (allInfo map[int64][]ProjectEntry, err error) {
	if r.Method != "POST" {
		return
	}
	var client http.Client
	var resp *http.Response
	allInfo = make(map[int64][]ProjectEntry)
	for k, host := range Hosts {
		sUrl := "http://" + host + "/zc-platform/v1/getAllProjects"
		var req *http.Request
		method := "getAllProjects"
		msgName := resp.Header.Get("Z-Zc-Msg-Name")
		if msgName == "X-Zc-Err" {
			color.Green("[%s] get result failed: err[%v]\n", host, err)
			continue
		}
		jbody := []byte(`{"developerId":` + developerIdStr + `}`)
		//rbody := bytes.NewBuffer(jbody)
		var req *http.Request
		req, err = http.NewRequest("POST", sUrl, bytes.NewBuffer(jbody))
		req.Header.Set("X-Zc-Developer-Id", fmt.Sprintf("%d", SignDeveloperId))
		req.Header.Set("X-Zc-Access-Key", AccessKey)
		req.Header.Set("X-Zc-Major-Domain", SignDomain)
		var signature string
		req.Header.Set("X-Zc-Timeout", r.Header.Get("X-Zc-Timeout"))
		req.Header.Set("X-Zc-Timestamp", fmt.Sprintf("%d", time.Now().Unix()))
		req.Header.Set("X-Zc-Nonce", r.Header.Get("X-Zc-Nonce"))
		req.Header.Set("Content-Type", "application/x-zc-object")

		signature, err = genAccessSignature(req, method)
		if err != nil || signature == "" {
			color.Red("gen signature error: %v\n", err)
			return
		}
		req.Header.Set("X-Zc-Developer-Signature", signature)

		color.Blue("major domain req : %v\n", req)
		resp, err = client.Do(req)
		color.Yellow("resp: %v, err: %v\n", resp, err)
		if err != nil {
			//TODO 检查error
			color.Green("[%s] send request failed: err[%v]\n", host, err)
			break
		}
		msgName := resp.Header.Get("Z-Zc-Msg-Name")
		if msgName == "X-Zc-Err" {
			color.Green("[%s] get result failed: err[%v]\n", host, err)
			continue
		}
		defer resp.Body.Close()
		var respBody []byte
		respBody, err = ioutil.ReadAll(resp.Body)
		color.Red("respBody : %v\n", string(respBody))
		if err != nil {
			return
		}
		info := make([]ProjectEntry, 0)
		err = json.Unmarshal(respBody, &info)
		if err != nil {

			return
		}
		allInfo[k] = info
	}
	return
}

func genAccessSignature(r *http.Request, method string) (signature string, err error) {

	signer := developer.Signer{}
	signer.DeveloperId = int64(SignDeveloperId)
	if err != nil {
		return
	}
	signer.MajorDomain = "test_ac"
	//	signer.MajorDomain = r.Header.Get("X-Zc-Major-Domain")
	//	signer.SubDomain = r.Header.Get("X-Zc-Sub-Domain")
	signer.Timestamp, err = strconv.ParseInt(r.Header.Get("X-Zc-Timestamp"), 10, 64)
	if err != nil {
		return
	}
	signer.Timeout, err = strconv.ParseInt(r.Header.Get("X-Zc-timeout"), 10, 64)
	signer.Nonce = r.Header.Get("X-Zc-Nonce")
	signer.Method = method

	return signer.Sign(SecretKey)

}
