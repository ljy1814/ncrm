package main

var (
	/*		1: "dev",
			2: "test",
			3: "生产",
			4: "华北",
			5: "北美",
			6: "欧洲",
			7: "东南亚",
	*/
	/*
		Hosts = map[int64]string{
			5: "usrouter.ablecloud.cn:5000",
			4: "router.ablecloud.cn:5000",
			6: "eurouter.ablecloud.cn:5000",
			7: "earouter.ablecloud.cn:5000",
			2: "test.ablecloud.cn:5000",
					1: "sandbox.ablecloud.cn:5000",
		}
	*/

	//	AccessKey = "9e8f02b7403c1c79804b843a8998fb81"
	//	SecretKey = "542c1fad40b4628f80fc95580f2b8afc"
	//	SignDeveloperId = 182

	AccessKey       = "96d3197e4061db2780f518984e731995"
	SecretKey       = "ffb33b61406e582c80008ab4c6314cc1"
	SignDeveloperId = 1
	SignDomain      = "test_ac"

	ZC_MSG_NAME_ACK = "X-Zc-Ack"
	ZC_MSG_NAME_ERR = "X-Zc-Err"
)

var Hosts map[int64]string
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

const (
	ENV_TEST    = 2
	ENV_PROD    = 1
	ENV_SANDBOX = 3
)
