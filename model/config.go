/**
 *	Author Fanxu(746439274@qq.com)
 */
package model

type Config struct{
		Session ServerConfig  `json:"session"`
		Http ServerConfig  `json:"http"`
		Rpc ServerConfig  `json:"rpc"`
		WebSocket ServerConfig  `json:"webSocket"`
		RpcServers []ServerConfig  `json:"rpcServers"`
		Mysql Dsn `json:"mysql"`
		MysqlPoolCount int32 `json:"mysqlPoolCount"`
		UploadPath string `json:"uploadPath"`
		SecretKey string `json:"secretKey"`
		WhiteIps []string `json:"whiteIps"`
}

type ServerConfig struct {
	Ip	string `json:"ip"`
	Port int32 `json:"port"`
}

type Dsn struct {
	Conn string `json:"conn"`
}
