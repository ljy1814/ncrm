package main

import (
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

type Products []*ProductInfo
type AllMajorInfo struct {
	MajorInfo `json:"major"`
	//Products  []*ProductInfo `json:"products"`
	Products `json:"products"`
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

type StockStatisticsInfo struct {
	QuotaCount  int64 `json:"quotaCount"`
	StockCount  int64 `json:"stockCount"`
	ActiveCount int64 `json:"activeCount"`
}

func (ps Products) Len() int {
	return len(ps)
}

func (ps Products) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func (ps Products) Less(i, j int) bool {
	return ps[i].SubDomainId < ps[j].SubDomainId
}
