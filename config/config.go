package config

import (
	"fmt"
	"strconv"

	"github.com/go-micro/plugins/v4/config/source/consul"
	"go-micro.dev/v4/config"
)

const (
	Host   = "127.0.0.1"
	Port   = 8500
	Prefix = "/config"
)

func GetConsulConfig() (config.Config, error) {
	consulSource := consul.NewSource(
		// optionally specify consul address; default to localhost:8500
		consul.WithAddress(Host+":"+strconv.FormatInt(Port, 10)),
		// optionally specify prefix; defaults to /micro/config
		consul.WithPrefix(Prefix),
		// optionally strip the provided prefix from the keys, defaults to false
		consul.StripPrefix(true),
	)

	// Create new config
	conf, err := config.NewConfig()
	if err != nil {
		fmt.Println("load config:", err)
		return conf, err
	}

	// Load consul source
	conf.Load(consulSource)
	return conf, nil
}

type MysqlConfig struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Pwd      string `json:"pwd"`
	Database string `json:"database"`
	Port     int64  `json:"port"`
}

// GetMysqlFromConsul 获取mysql的配置
func GetMysqlFromConsul(config config.Config, path ...string) (*MysqlConfig, error) {
	mysqlConfig := &MysqlConfig{}
	//获取配置
	err := config.Get(path...).Scan(mysqlConfig)
	if err != nil {
		return nil, err
	}
	return mysqlConfig, nil
}

type KongGatewayConfig struct {
	Url string `json:"url"`
}

// TODO写个通用的 get config 用泛型

func GetKongGatewayFromConsul(config config.Config, path ...string) (*KongGatewayConfig, error) {
	kongConfig := &KongGatewayConfig{}
	//获取配置
	err := config.Get(path...).Scan(kongConfig)
	if err != nil {
		return nil, err
	}
	return kongConfig, nil
}
