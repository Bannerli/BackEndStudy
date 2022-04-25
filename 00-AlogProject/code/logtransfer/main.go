package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"logtransfer/conf"
	"logtransfer/es"
	"logtransfer/kafka"
	"strings"
)
// log transfer
// 将日志数据从kafka取出来发往ES

func main() {
	// 0. 加载配置文件
	var cfg conf.LogTansfer
	err := ini.MapTo(&cfg, "./conf/config.ini")
	if err != nil {
		fmt.Println("init config, err:%v\n", err)
		return
	}
	fmt.Printf("cfg:%v\n", cfg)

	// 1. 初始化ES
	// 1.1 初始化一个ES连接的client
	err = es.Init(cfg.ES.Address, cfg.ES.ChanMaxSize, cfg.ES.Workers)
	if err != nil {
		fmt.Printf("init ES client failed,err:%v\n", err)
		return
	}
	fmt.Println("init ES client success.")
	// 2. 初始化kafka
	// 2.1 连接kafka，创建分区的消费者
	// 2.2 每个分区的消费者分别取出数据 通过SendToChan()将数据发往管道
	// 2.3 初始化时就开起协程去channel中取数据发往ES
	err = kafka.Init(strings.Split(cfg.Kafka.Address, ";"), cfg.Kafka.Topic)
	if err != nil {
		fmt.Printf("init kafka consumer failed,err:%v\n", err)
		return
	}
	fmt.Println("init kafka success.")

	// 3. 从kafka取日志数据并放入channel中
	kafka.Run()
}