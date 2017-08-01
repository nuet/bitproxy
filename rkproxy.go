/*
	By moli molisoft@qq.com
	Blog huoxr.com

*/
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"

	"rkproxy/libs"
	"rkproxy/manager"
	"rkproxy/manager/api"
	"rkproxy/utils"
)

var (
	config_path = "config.json"
	pid_path    = "rkproxy.pid"
)

var man *manager.Manager
var log *utils.Logger = utils.NewLogger("Main")

func listenSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		select {
		case sig := <-c:
			fmt.Printf("ctrl+c (%v)\n", sig)
			if man != nil {
				man.StopAll()
			}
			os.Remove(pid_path)
			os.Exit(0)
		}
	}()
}

func initFlag() {
	flag.StringVar(&config_path, "c", "config.json", "配置文件")
	flag.StringVar(&pid_path, "p", "rkproxy.pid", "进程id路径")
	flag.Parse()
}

func initPid() {
	if _, err := os.Stat(pid_path); os.IsExist(err) {
		os.Remove(pid_path)
	}
	file, err := os.OpenFile(pid_path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("写pid文件错误：", err)
		return
	}
	defer file.Close()
	file.WriteString(fmt.Sprintf("%d", os.Getpid()))
}

func initApi() error {
	if manager.Man.Config().Api != nil {
		config := manager.Man.Config().Api
		return api.Start(config.Password, config.Port)
	}
	return nil
}

func initRedis() {
	if manager.Man.Config().Redis != nil {
		config := manager.Man.Config().Redis
		libs.InitRedis(config.Host, config.Port)
	}
}

func initStats() {
	libs.StartStats()
}

func initBlackList() {
	libs.StartBlackList()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	manager.New(config_path)

	err := manager.Man.ParseConfig()
	if err != nil {
		log.Info(err.Error())
		return
	}

	initFlag()
	initPid()
	listenSignal()

	initRedis()
	initStats()
	initBlackList()

	manager.Man.RunAll()

	initApi()
	select {}
}
