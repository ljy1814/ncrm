package main

import (
	log "ac-common-go/glog"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type DomainHandler struct {
}

func NewDomainHandler() *DomainHandler {
	return &DomainHandler{}
}

func getMajorDomain(r *http.Request) (majorDomain string, majorDomainId, developerId int64, acErr []byte, err error) {
	if r.Method != "POST" {
		return
	}
	var client http.Client
	var resp *http.Response
	var reqBody []byte
	reqBody, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	var rbody map[string]interface{}
	err = json.Unmarshal(reqBody, &rbody)
	if err != nil {
		return
	}
	majorDomainId = GetInt(rbody["majorDomainId"])
	for _, host := range Hosts {
		sUrl := "http://" + host + "/zc-platform/v1/getMajorDomain"
		var req *http.Request
		jbody := []byte(`{"majorDomainId":` + fmt.Sprintf("%d", majorDomainId) + `}`)
		req, err = http.NewRequest("POST", sUrl, bytes.NewBuffer(jbody))
		req.Header.Set("X-Zc-Developer-Id", "0")
		req.Header.Set("X-Zc-Access-Mode", "1")
		req.Header.Set("X-Zc-Major-Domain", GetString(CrmConf[host]))
		req.Header.Set("Content-Type", "application/x-zc-object")

		resp, err = client.Do(req)
		traceId := resp.Header.Get("X-Zc-Trace-Id")
		if err != nil {
			//TODO 检查error,请求发送失败
			log.Fatalf("TraceID[%s], [%s] send request failed: err[%v]\n", traceId, host, err)
			return
		}
		defer resp.Body.Close()
		var respBody []byte
		respBody, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		msgName := resp.Header.Get("X-Zc-Msg-Name")
		if msgName == "X-Zc-Err" {
			//TODO 请求错误
			log.Warningf("TraceId[%s], [%s] get result failed: err[%v]\n", traceId, host, string(respBody))
			errMsg, _, perr := ParseError(respBody)
			if perr != nil {
				break
			}
			err = errors.New(errMsg)
			acErr = respBody
			continue
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
func getDeveloper(r *http.Request, domain, developerIdStr string) (company string, acErr []byte, err error) {
	if r.Method != "POST" {
		return
	}
	var client http.Client
	var resp *http.Response
	for _, host := range Hosts {
		sUrl := "http://" + host + "/zc-platform/v1/getDeveloper"
		jbody := []byte(`{"developerId":` + developerIdStr + `}`)
		var req *http.Request
		req, err = http.NewRequest("POST", sUrl, bytes.NewBuffer(jbody))
		req.Header.Set("X-Zc-Developer-Id", developerIdStr)
		req.Header.Set("X-Zc-Access-Mode", "1")
		req.Header.Set("X-Zc-Major-Domain", domain)
		req.Header.Set("Content-Type", "application/x-zc-object")

		resp, err = client.Do(req)
		traceId := resp.Header.Get("X-Zc-Trace-Id")
		if err != nil {
			//TODO 检查error
			log.Fatalf("TraceId[%s], [%s] send request failed: err[%v]\n", traceId, host, err)
			break
		}
		defer resp.Body.Close()
		var respBody []byte
		respBody, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		msgName := resp.Header.Get("X-Zc-Msg-Name")
		if msgName == "X-Zc-Err" {
			//TODO 请求错误
			log.Warningf("TraceId[%s], [%s] get result failed: err[%v]\n", traceId, host, string(respBody))
			errMsg, _, perr := ParseError(respBody)
			if perr != nil {
				break
			}
			err = errors.New(errMsg)
			acErr = respBody
			continue
		}
		var res map[string]interface{}
		err = json.Unmarshal(respBody, &res)
		if err == nil {
			acErr = nil
			company = GetString(res["company"])
			return
		}
	}
	return
}

func getAllProjects(r *http.Request, domain, developerIdStr string) (allInfo map[int64][]ProjectEntry, acErr []byte, err error) {
	if r.Method != "POST" {
		return
	}
	var client http.Client
	var resp *http.Response
	allInfo = make(map[int64][]ProjectEntry)
	for k, host := range Hosts {
		sUrl := "http://" + host + "/zc-platform/v1/getAllProjects"
		var req *http.Request
		jbody := []byte(`{"developerId":` + developerIdStr + `}`)
		req, err = http.NewRequest("POST", sUrl, bytes.NewBuffer(jbody))
		req.Header.Set("X-Zc-Developer-Id", developerIdStr)
		req.Header.Set("X-Zc-Access-Mode", "1")
		req.Header.Set("X-Zc-Major-Domain", domain)
		req.Header.Set("Content-Type", "application/x-zc-object")

		resp, err = client.Do(req)
		traceId := resp.Header.Get("X-Zc-Trace-Id")
		if err != nil {
			//TODO 检查error
			log.Fatalf("TraceId[%s], [%s] send request failed: err[%v]\n", traceId, host, err)
			break
		}
		defer resp.Body.Close()
		var respBody []byte
		respBody, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		msgName := resp.Header.Get("X-Zc-Msg-Name")
		if msgName == "X-Zc-Err" {
			//TODO 请求错误
			log.Warningf("TraceId[%s], [%s] get result failed: err[%v]\n", traceId, host, string(respBody))
			errMsg, _, perr := ParseError(respBody)
			if perr != nil {
				break
			}
			err = errors.New(errMsg)
			acErr = respBody
			continue
		}
		info := make([]ProjectEntry, 0)
		err = json.Unmarshal(respBody, &info)
		if err != nil {

			return
		}
		acErr = nil
		allInfo[k] = info
	}
	return
}
func getAccountCount(r *http.Request, domain, developerIdStr string) (accountCount int64, acErr []byte, err error) {
	if r.Method != "POST" {
		return
	}
	var client http.Client
	var resp *http.Response
	for _, host := range Hosts {
		sUrl := "http://" + host + "/zc-account/v1/getAccountCount"
		var req *http.Request
		req, err = http.NewRequest("POST", sUrl, nil)
		req.Header.Set("X-Zc-Developer-Id", developerIdStr)
		req.Header.Set("X-Zc-Access-Mode", "1")
		req.Header.Set("X-Zc-Major-Domain", domain)
		req.Header.Set("Content-Type", "application/x-zc-object")

		resp, err = client.Do(req)
		traceId := resp.Header.Get("X-Zc-Trace-Id")
		if err != nil {
			//TODO 检查error
			log.Fatalf("TraceId[%s], [%s] send request failed: err[%v]\n", traceId, host, err)
			break
		}
		defer resp.Body.Close()
		var respBody []byte
		respBody, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		msgName := resp.Header.Get("X-Zc-Msg-Name")
		if msgName == "X-Zc-Err" {
			//TODO 请求错误
			log.Warningf("TraceId[%s], [%s] get result failed: err[%v]\n", traceId, host, string(respBody))
			errMsg, _, perr := ParseError(respBody)
			if perr != nil {
				break
			}
			err = errors.New(errMsg)
			acErr = respBody
			continue
		}
		var res map[string]interface{}
		err = json.Unmarshal(respBody, &res)
		if err == nil {
			acErr = nil
			accountCount += GetInt(res["count"])
		}
	}
	return
}
func getProduct(r *http.Request, host, domain, subDomain, developerIdStr string) (product Product, acErr []byte, err error) {
	if r.Method != "POST" {
		return
	}
	var client http.Client
	var resp *http.Response
	sUrl := "http://" + host + "/zc-product/v1/getProduct"
	var req *http.Request
	req, err = http.NewRequest("POST", sUrl, nil)
	req.Header.Set("X-Zc-Developer-Id", developerIdStr)
	req.Header.Set("X-Zc-User-Id", "0")
	req.Header.Set("X-Zc-Inner-Service", "")
	req.Header.Set("X-Zc-Access-Mode", "1")
	req.Header.Set("X-Zc-Major-Domain", domain)
	req.Header.Set("X-Zc-Sub-Domain", subDomain)
	req.Header.Set("Content-Type", "application/x-zc-object")

	resp, err = client.Do(req)
	traceId := resp.Header.Get("X-Zc-Trace-Id")
	if err != nil {
		//TODO 检查error
		log.Fatalf("TraceId[%s], [%s] send request failed: err[%v]\n", traceId, host, err)
		return
	}
	defer resp.Body.Close()
	var respBody []byte
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	msgName := resp.Header.Get("X-Zc-Msg-Name")
	if msgName == "X-Zc-Err" {
		//TODO 请求错误
		log.Warningf("TraceId[%s], [%s] get result failed: err[%v]\n", traceId, host, string(respBody))
		errMsg, _, perr := ParseError(respBody)
		if perr != nil {
			return
		}
		err = errors.New(errMsg)
		acErr = respBody
		return
	}
	var res map[string]Product
	err = json.Unmarshal(respBody, &res)
	if err == nil {
		product = res["product"]
		return
	}
	return
}
func getLicenseQuota(r *http.Request, host, domain, subDomain, developerIdStr string) (lq LicenseQuotaInfo, acErr []byte, err error) {
	if r.Method != "POST" {
		return
	}
	var client http.Client
	var resp *http.Response
	sUrl := "http://" + host + "/zc-warehouse/v1/getLicenseQuota"
	var req *http.Request
	req, err = http.NewRequest("POST", sUrl, nil)
	req.Header.Set("X-Zc-Developer-Id", developerIdStr)
	req.Header.Set("X-Zc-Access-Mode", "1")
	req.Header.Set("X-Zc-Major-Domain", domain)
	req.Header.Set("X-Zc-Sub-Domain", subDomain)
	req.Header.Set("Content-Type", "application/x-zc-object")

	resp, err = client.Do(req)
	traceId := resp.Header.Get("X-Zc-Trace-Id")
	if err != nil {
		//TODO 检查error
		log.Fatalf("TraceId[%s], [%s] send request failed: err[%v]\n", traceId, host, err)
		return
	}
	defer resp.Body.Close()
	var respBody []byte
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	msgName := resp.Header.Get("X-Zc-Msg-Name")
	if msgName == "X-Zc-Err" {
		//TODO 请求错误
		log.Warningf("TraceId[%s], [%s] get result failed: err[%v]\n", traceId, host, string(respBody))
		errMsg, _, perr := ParseError(respBody)
		if perr != nil {
			return
		}
		err = errors.New(errMsg)
		acErr = respBody
		return
	}
	err = json.Unmarshal(respBody, &lq)
	if err != nil {
		return
	}
	return
}
func getSubDomain(r *http.Request, host, domain, subDomain, developerIdStr string) (info SubDomainEntry, acErr []byte, err error) {
	sUrl := "http://" + host + "/zc-platform/v1/getSubDomain"
	jbody := []byte(`{"subDomainId":` + subDomain + `}`)
	var req *http.Request
	var resp *http.Response
	client := http.Client{}
	req, err = http.NewRequest("POST", sUrl, bytes.NewBuffer(jbody))
	req.Header.Set("X-Zc-Developer-Id", developerIdStr)
	req.Header.Set("X-Zc-Access-Mode", "1")
	req.Header.Set("X-Zc-Major-Domain", domain)
	req.Header.Set("X-Zc-Sub-Domain-Id", subDomain)
	req.Header.Set("Content-Type", "application/x-zc-object")

	resp, err = client.Do(req)
	traceId := resp.Header.Get("X-Zc-Trace-Id")
	if err != nil {
		//TODO 检查error
		log.Fatalf("TraceId[%s], [%s] send request failed: err[%v]\n", traceId, host, err)
		return
	}
	defer resp.Body.Close()
	var respBody []byte
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	msgName := resp.Header.Get("X-Zc-Msg-Name")
	if msgName == "X-Zc-Err" {
		//TODO 请求错误
		log.Warningf("TraceId[%s], [%s] get result failed: err[%v]\n", traceId, host, string(respBody))
		errMsg, _, perr := ParseError(respBody)
		if perr != nil {
			return
		}
		err = errors.New(errMsg)
		acErr = respBody
		return
	}
	err = json.Unmarshal(respBody, &info)
	if err != nil {
		return
	}

	return
}

func getDeviceCount(r *http.Request, host, domain, subDomain, developerIdStr string) (count, deviceCount int64, acErr []byte, err error) {
	if r.Method != "POST" {
		return
	}
	var client http.Client
	var resp *http.Response
	for _, host := range Hosts {
		sUrl := "http://" + host + "/zc-warehouse/v1/getDeviceCount"
		jbody := []byte(`{"developerId":` + developerIdStr + `}`)
		var req *http.Request
		req, err = http.NewRequest("POST", sUrl, bytes.NewBuffer(jbody))
		req.Header.Set("X-Zc-Developer-Id", fmt.Sprintf("%d", SignDeveloperId))
		req.Header.Set("X-Zc-Access-Mode", "1")
		req.Header.Set("X-Zc-Major-Domain", domain)
		req.Header.Set("X-Zc-Sub-Domain", subDomain)
		req.Header.Set("Content-Type", "application/x-zc-object")

		resp, err = client.Do(req)
		if err != nil {
			//TODO 检查error
			log.Fatalf("[%s] send request failed: err[%v]\n", host, err)
			return
		}
		traceId := resp.Header.Get("X-Zc-Trace-Id")
		defer resp.Body.Close()
		var respBody []byte
		respBody, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		msgName := resp.Header.Get("X-Zc-Msg-Name")
		if msgName == "X-Zc-Err" {
			//TODO 请求错误
			log.Warningf("TraceId[%s], [%s] get result failed: err[%v]\n", traceId, host, string(respBody))
			errMsg, _, perr := ParseError(respBody)
			if perr != nil {
				break
			}
			err = errors.New(errMsg)
			acErr = respBody
			return
		}
		var res map[string]interface{}
		err = json.Unmarshal(respBody, &res)
		if err != nil {
			return
		}
		count = GetInt(res["count"])
		deviceCount = GetInt(res["deviceCount"])
	}
	return
}

func (this *DomainHandler) handleListDomain(w http.ResponseWriter, r *http.Request) {
	var amp AllMajorInfo

	domain, domainId, developerId, acErr, err := getMajorDomain(r)
	if err != nil {
		writeError(acErr, w)
		return
	}
	if domainId <= 0 || domain == "" {
		writeError(nil, w)
		return
	}
	amp.MajorInfo.MajorDomainId = domainId
	amp.MajorInfo.MajorDomain = domain

	developerIdStr := fmt.Sprintf("%d", developerId)

	company, acErr, err := getDeveloper(r, domain, fmt.Sprintf("%d", developerId))
	if err != nil {
		writeError(acErr, w)
		return
	}
	amp.MajorInfo.CompanyName = company

	accounts, acErr, err := getAccountCount(r, domain, fmt.Sprintf("%d", developerId))
	if err != nil {
		writeError(acErr, w)
		return
	}
	amp.MajorInfo.AccountCount = accounts

	//var projects []lcommon.ProjectEntry
	allProjects, acErr, err := getAllProjects(r, domain, fmt.Sprintf("%d", developerId))
	if err != nil {
		writeError(acErr, w)
		return
	}
	//TODO 暂时还在测试获取所属环境的接口
	//　所有环境的产品
	for k, projects := range allProjects {
		//某个环境的所有产品
		for _, v := range projects {
			var pi ProductInfo
			pi.ProductName = v.Name
			pi.SubDomainId = v.SubDomainId
			subDomainIdStr := strconv.Itoa(int(v.SubDomainId))
			subDm, acErr, err := getSubDomain(r, Hosts[k], domain, subDomainIdStr, developerIdStr)
			if err != nil {
				//writeError(acErr, w)
				//return
				continue
			}
			//			plm, err := this.productClient.GetProduct(req, client.Hosts[k], domainIdStr, subDomainIdStr, developerId)
			plm, acErr, err := getProduct(r, Hosts[k], domain, subDm.SubDomain, developerIdStr)
			if err != nil {
				log.WarningfT(nil, "get product failed: domain[%s], subdomain[%d], developerId[%d] err[%v]", domain, pi.SubDomainId, developerId, err)
				writeError(acErr, w)
				return
			}
			pi.LicenseMode = plm.LicenseMode
			//TODO 错误处理,直接忽略?
			//	lq, err := this.warehouseClient.GetLicenseQuota(req, client.Hosts[k], domainIdStr, subDomainIdStr, developerId)
			lq, acErr, err := getLicenseQuota(r, Hosts[k], domain, subDm.SubDomain, developerIdStr)
			amp.MajorInfo.AllLicenseCount += lq.QuotaTotal
			amp.MajorInfo.AllLicenseAllocated += lq.QuotaUsed
			if err != nil {
				//				writeError(acErr, w)
				//		return
				pi.LicenseAllocated = -1
			} else {
				pi.LicenseAllocated = lq.QuotaUsed
			}

			//			count, activeCount, err := this.warehouseClient.GetDeviceCount(req, client.Hosts[k], domainIdStr, subDomainIdStr, developerId)
			count, activeCount, acErr, err := getDeviceCount(r, Hosts[k], domain, subDm.SubDomain, developerIdStr)
			amp.MajorInfo.AllDeviceImport += count
			amp.MajorInfo.AllDeviceActived += activeCount
			if err != nil {
				//				writeError(acErr, w)
				//				return

				pi.DeviceImport = -1
				pi.DeviceActived = -1
			} else {
				pi.DeviceImport = count
				pi.DeviceActived = activeCount
			}

			pi.Environment = envs[k]
			//TODO 环境
			amp.Products = append(amp.Products, &pi)
		}
	}

	jamp, err := json.Marshal(amp)
	if err != nil {
		writeError([]byte(err.Error()), w)
		return
	}
	//	resp.SetPayload(jamp, zc.ZC_MSG_PAYLOAD_JSON)
	w.Header().Set("Content-Type", "text/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(jamp)))
	w.Header().Set("X-Zc-Content-Length", fmt.Sprintf("%d", len(jamp)))
	w.Header().Set("X-Zc-Msg-Name", ZC_MSG_NAME_ACK)
	w.Write(jamp)
}
