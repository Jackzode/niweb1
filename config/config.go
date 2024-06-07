package config

import (
	"fmt"
	"github.com/Jackzode/painting/commons/constants"
	"github.com/spf13/viper"
	"path/filepath"
)

func GetConfigFilePath() string {
	return filepath.Join(constants.ConfigFileDir, constants.DefaultConfigFileName)
}

func ReadConfig(configFilePath string) (c *MainConfig, err error) {
	if configFilePath == "" {
		configFilePath = GetConfigFilePath()
	}
	c = &MainConfig{}
	v := viper.New()
	v.SetConfigFile(configFilePath)
	err = v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	if err = v.Unmarshal(&c); err != nil {
		return nil, err
	}
	fmt.Println("get config : ", *c.Data, "  ||  ", *c.Cache)
	return c, nil
}

type MainConfig struct {
	Debug       bool         `json:"debug" mapstructure:"debug" yaml:"debug"`
	Data        *Database    `json:"data" mapstructure:"data" yaml:"data"`
	Cache       *RedisConf   `json:"redis" mapstructure:"redis" yaml:"redis"`
	Addr        string       `json:"addr" mapstructure:"addr"`
	EmailConfig *EmailConfig `json:"email" mapstructure:"email" yaml:"email"`
}

// EmailConfig email config
type EmailConfig struct {
	FromEmail          string `json:"from_email" mapstructure:"from_email" yaml:"from_email"`
	FromName           string `json:"from_name" mapstructure:"from_name" yaml:"from_name"`
	SMTPHost           string `json:"smtp_host" mapstructure:"smtp_host" yaml:"smtp_host"`
	SMTPPort           int    `json:"smtp_port" mapstructure:"smtp_port" yaml:"smtp_port"`
	Encryption         string `json:"encryption" mapstructure:"encryption" yaml:"encryption"`
	SMTPUsername       string `json:"smtp_username" mapstructure:"smtp_username" yaml:"smtp_username"`
	SMTPPassword       string `json:"smtp_password" mapstructure:"smtp_password" yaml:"smtp_password"`
	SMTPAuthentication bool   `json:"smtp_authentication" mapstructure:"smtp_authentication" yaml:"smtp_authentication"`
}

type Database struct {
	Driver          string `json:"driver" mapstructure:"driver" yaml:"driver"`
	Connection      string `json:"connection" mapstructure:"connection" yaml:"connection"`
	ConnMaxLifeTime int    `json:"conn_max_life_time" mapstructure:"conn_max_life_time" yaml:"conn_max_life_time,omitempty"`
	MaxOpenConn     int    `json:"max_open_conn" mapstructure:"max_open_conn" yaml:"max_open_conn,omitempty"`
	MaxIdleConn     int    `json:"max_idle_conn" mapstructure:"max_idle_conn" yaml:"max_idle_conn,omitempty"`
}

// RedisConf cache
type RedisConf struct {
	Addr       string `json:"addr" yaml:"addr" mapstructure:"addr"`
	MaxOpen    int    `json:"maxOpen" yaml:"maxOpen" mapstructure:"maxOpen"`
	MaxIdle    int    `json:"maxIdle" yaml:"maxIdle" mapstructure:"maxIdle"`
	MaxConnect int    `json:"maxConnect" yaml:"maxConnect" mapstructure:"maxConnect"`
	Timeout    int    `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	PoolSize   int    `json:"poolSize" yaml:"poolSize" mapstructure:"poolSize"`
	Auth       string `json:"auth" yaml:"auth" mapstructure:"auth"`
}
