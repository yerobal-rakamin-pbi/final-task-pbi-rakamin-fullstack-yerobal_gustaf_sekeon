package config

type Application struct {
	Server  Server  `json:"server"`
	SQL     SQL     `json:"sql"`
	Storage Storage `json:"storage"`
}

type Server struct {
	Port                 string `json:"port"`
	Host                 string `json:"host"`
	RequestTimeoutSecond int    `json:"requestTimeoutSecond"`
	JWT                  JWT    `json:"jwt"`
}

type JWT struct {
	Secret string `json:"secret"`
	ExpSec int64  `json:"expSec"`
}

type SQL struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type Storage struct {
	BucketName string `json:"bucketName"`
}
