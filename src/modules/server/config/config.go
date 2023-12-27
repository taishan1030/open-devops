package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type MySQLConf struct {
	DB   string `yaml:"db"`
	Max  int    `yaml:"max"`
	Idle int    `yaml:"idle"`
}

type Config struct {
	MysqlS *MySQLConf `yaml:"mysql_s"`
}

// 根据io read读取配置文件后的字符串解析yaml
func Load(s []byte) (*Config, error) {
	cfg := &Config{}

	err := yaml.Unmarshal(s, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// 根据conf路径读取内容
func LoadFile(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg, err := Load(content)
	if err != nil {
		fmt.Printf("[parsing Yaml file err...][err:%v]\n", err)
		return nil, err
	}
	return cfg, nil
}
