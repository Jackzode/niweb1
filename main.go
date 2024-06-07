package main

import (
	"fmt"
	"github.com/Jackzode/painting/commons/handler"
	glog "github.com/Jackzode/painting/commons/logger"
	"github.com/Jackzode/painting/config"
	"github.com/Jackzode/painting/routers"
)

func main() {
	glog.InitLogger("./log/test.log", "debug")
	defer func() {
		err := glog.Klog.Sync()
		fmt.Println(err)
	}()
	fmt.Println("server is starting...")
	filename := `./conf/config.yaml`
	c, err := config.ReadConfig(filename)
	if err != nil {
		panic(err)
	}
	//init db
	err = handler.InitDataBase(c.Debug, c.Data)
	if err != nil {
		panic(err)
	}
	err = handler.InitRedisCache(c.Cache)
	if err != nil {
		panic(err)
	}
	//init email
	handler.InitEmailService(c.EmailConfig)
	//listen
	server := routers.NewHTTPServer(c.Debug)
	err = server.Run(c.Addr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
