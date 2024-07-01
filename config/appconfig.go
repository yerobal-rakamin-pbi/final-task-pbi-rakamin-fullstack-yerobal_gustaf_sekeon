package config

type Application struct {
	Server  Server  `json:"server"`
	SQL     SQL     `json:"sql"`
	Storage Storage `json:"storage"`
}

type Server struct {
	Port              string   `json:"port"`
	Host              string   `json:"host"`
	RequestTimeoutSec int64    `json:"requestTimeoutSec"`
	JWT               JWT      `json:"jwt"`
	Password          Password `json:"password"`
}

type JWT struct {
	Secret string `json:"secret"`
	ExpSec int64  `json:"expSec"`
}

type Password struct {
	SaltRound int64 `json:"saltRound"`
}

type SQL struct {
	Host       string     `json:"host"`
	Port       string     `json:"port"`
	Username   string     `json:"username"`
	Password   string     `json:"password"`
	Database   string     `json:"database"`
	PoolConfig PoolConfig `json:"poolConfig"`
}

type PoolConfig struct {
	MaxIdle            int64 `json:"maxIdle"`
	MaxOpen            int64 `json:"maxOpen"`
	ConnIdleSec        int64 `json:"connIdleSec"`
	ConnMaxLifetimeSec int64 `json:"connMaxLifetimeSec"`
}

type Storage struct {
	BucketName string `json:"bucketName"`
}
