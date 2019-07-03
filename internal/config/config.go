package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type HttpConf struct {
	Port        int      `yaml:"port"`
	Host        string   `yaml:"host"`
	AllowOrigin []string `yaml:"allow-origin"`
	DataDir     string   `yaml:"data-dir"`
}

type InfluxDBConf struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"dbName"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type ServerConfig struct {
	Http     *HttpConf     `yaml:"http"`
	InfluxDB *InfluxDBConf `yaml:"influxdb"`
}

func ErrHandler(op string, err error) {
	if err != nil {
		log.Fatalf("%s found error: %v", op, err)
	}
}

func NewServerConfig() *ServerConfig {
	var serverConf ServerConfig
	ymlFile, err := ioutil.ReadFile("configs/configs.yml")
	ErrHandler("opening file", err)

	err = yaml.Unmarshal(ymlFile, &serverConf)
	ErrHandler("unmarshal yaml", err)

	return &serverConf
}
