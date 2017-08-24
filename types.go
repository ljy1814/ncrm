package main

import (
	"encoding/json"
	"reflect"
	"time"
)

type ProjectEntry struct {
	Id            int64     `json:"projectId"`
	Name          string    `json:"name,omitempty"`
	DeveloperId   int64     `json:"developerId"`
	MajorDomainId int64     `json:"majorDomainId"`
	SubDomainId   int64     `json:"subDomainId"`
	Description   string    `json:"description"`
	CreateTime    time.Time `json:"createTime"`
	ModifyTime    time.Time `json:"modifyTime"`
	Status        int8      `json:"status"`
}

type MajorInfo struct {
	MajorDomainId       int64  `json:"majorDomainId"`
	MajorDomain         string `json:"majorDomain"`
	CompanyName         string `json:"companyName"`
	AllLicenseCount     int64  `json:"allLicenseCount"`
	AllLicenseAllocated int64  `json:"allLicenseAllocated"`
	AllDeviceImport     int64  `json:"allDeviceImport"`
	AllDeviceActived    int64  `json:"allDeviceActived"`
	AccountCount        int64  `json:"accountCount"`
}

type ProductInfo struct {
	ProductName      string `json:"productName"`
	SubDomainId      int64  `json:"subDomainId"`
	LicenseMode      int8   `json:"licenseMode"`
	LicenseAllocated int64  `json:"licenseAllocated"`
	DeviceImport     int64  `json:"deviceImport"`
	DeviceActived    int64  `json:"deviceActived"`
	Environment      string `json:"environment"`
}

type AllMajorInfo struct {
	MajorInfo `json:"major"`
	Products  []*ProductInfo `json:"products"`
}

type Product struct {
	Domain            string `json:"domain"`                                  // id of major domain
	SubDomain         string `json:"sub_domain"`                              // id of sub domain
	DomainName        string `json:"domain_name"`                             // name of major domain
	SubDomainName     string `json:"sub_domain_name"`                         // name of sub domain
	Name              string `json:"name"`                                    // product name
	ProductType       string `json:"type"`                                    // product type 独立设备/网关设备/安卓设备/子设备
	Model             string `json:"model"`                                   // product model
	Description       string `json:"description"`                             // product description
	Category          int8   `json:"category"`                                // product category
	ProductImageUrl   string `json:"productImageUrl"`                         // 产品图片url
	CreateTime        string `json:"createTime"`                              // 产品的创建时间
	OS                string `json:"os" confidential:"level1"`                // product operating system
	Protocol          string `json:"protocol" confidential:"level1"`          // message json/klv/binary
	TransportProtocol string `json:"transportProtocol" confidential:"level1"` // transport protocol tcp/simple tcp/mqtt/http, default tcp
	Communication     string `json:"communication" confidential:"level1"`     // communication protocol /wifi/bluetooth/ethernet/cellular
	ThirdCloud        string `json:"thirdCloud" confidential:"level1"`        // product third cloud, default cloud is ablecloud
	SecType           string `json:"secType" confidential:"level1"`           // product sec type, default RSA
	LicenseMode       int8   `json:"licenseMode" confidential:"level2"`       // product's license mode @refer to LicenseModeType
	DeviceMode        int8   `json:"deviceMode" confidential:"level2"`        // device manager mode; 0: 非绑定模式 1: 管理员绑定模式 2: 普通绑定模式
	TaskMode          int8   `json:"taskMode" confidential:"level2"`          // 设备定时模式 0: 无定时 1: 云端定时(默认) 2: 设备定时(支持云端定时)
	TaskUpdatePolicy  int8   `json:"taskUpdatePolicy" confidential:"level2"`  // 定时任务删除机制；1: 普通用户解绑不删除定时任务，管理员解绑删除所有定时任务(默认) 2: 用户解绑删除定时任务
	MaxDeviceNum      int64  `json:"maxDeviceNum" confidential:"level2"`      // 支持最大设备连接数，测试环境默认10
	DeviceUdsService  string `json:"deviceUdsService" confidential:"level2"`  // 产品映射的用于设备上报数据处理的UDS服务，格式 `3.6/DemoService` `3./DemoService`
}

type SubDomainEntry struct {
	DeveloperId int64  `json:"developerId"`
	SubDomainId int64  `json:"subDomainId"`
	MajorDomain string `json:"majorDomain"`
	SubDomain   string `json:"subDomain"`
}

type LicenseQuotaInfo struct {
	QuotaTotal   int64 `json:"quotaTotal"`
	QuotaRemains int64 `json:"quotaRemains"`
	QuotaUsed    int64 `json:"quotaUsed"`
}

func GetInt(v interface{}) int64 {
	if v == nil {
		return 0
	}
	kind := reflect.TypeOf(v).Kind()
	switch kind {
	case reflect.Int:
		return int64(v.(int))
	case reflect.Int64:
		return int64(v.(int64))
	case reflect.Int8:
		return int64(v.(int8))
	case reflect.Int16:
		return int64(v.(int16))
	case reflect.Int32:
		return int64(v.(int32))
	case reflect.Uint:
		return int64(v.(uint))
	case reflect.Uint8:
		return int64(v.(uint8))
	case reflect.Uint16:
		return int64(v.(uint16))
	case reflect.Uint32:
		return int64(v.(uint32))
	case reflect.Uint64:
		return int64(v.(uint64))
	case reflect.Float32:
		return int64(v.(float32))
	case reflect.Float64:
		return int64(v.(float64))
	}
	return 0
}
func GetFloat(v interface{}) float64 {
	if v == nil {
		return 0
	}
	switch reflect.TypeOf(v).Kind() {
	case reflect.Float64, reflect.Float32:
		return v.(float64)
	default:
		return 0
	}
}
func GetString(v interface{}) string {
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
}

var (
	envs = map[int64]string{
		1: "dev",
		2: "test",
		3: "生产",
		4: "华北环境",
		5: "北美环境",
		6: "欧洲环境",
		7: "东南亚环境",
	}
)

func ParseError(oriErr []byte) (errMsg string, code int64, err error) {
	if len(oriErr) <= 0 {
		return
	}
	var errMap map[string]interface{}
	err = json.Unmarshal(oriErr, &errMap)
	if err != nil {
		return
	}
	errMsg = GetString(errMap["error"])
	code = GetInt(errMap["errorCode"])
	return
}
