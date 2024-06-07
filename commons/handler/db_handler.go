package handler

import (
	"context"
	"github.com/Jackzode/painting/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"time"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
	"xorm.io/xorm/names"
)

var (
	Engine      *xorm.Engine
	RedisClient *redis.Client
)

func InitRedisCache(conf *config.RedisConf) error {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Auth,     // 密码
		PoolSize: conf.PoolSize, // 连接池大小
	})
	err := redisClient.Ping(context.TODO()).Err()
	if err != nil {
		return err
	}
	RedisClient = redisClient
	return nil
}

func InitDataBase(debug bool, dataConf *config.Database) error {

	engine, err := xorm.NewEngine(dataConf.Driver, dataConf.Connection)
	if err != nil {
		return err
	}

	if debug {
		engine.ShowSQL(true)
	} else {
		engine.SetLogLevel(log.LOG_ERR)
	}

	if err = engine.Ping(); err != nil {
		return err
	}

	if dataConf.MaxIdleConn > 0 {
		engine.SetMaxIdleConns(dataConf.MaxIdleConn)
	}
	if dataConf.MaxOpenConn > 0 {
		engine.SetMaxOpenConns(dataConf.MaxOpenConn)
	}
	if dataConf.ConnMaxLifeTime > 0 {
		engine.SetConnMaxLifetime(time.Duration(dataConf.ConnMaxLifeTime) * time.Second)
	}
	engine.SetColumnMapper(names.GonicMapper{})
	Engine = engine
	return nil
}
