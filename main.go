package main

import (
	"flag"
	"fmt"
	"github.com/open-falcon/agent/cron"
	"github.com/open-falcon/agent/funcs"
	"github.com/open-falcon/agent/g"
	"github.com/open-falcon/agent/http"
	"os"
)

func main() {

	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	check := flag.Bool("check", false, "check collector")

	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	if *check {
		funcs.CheckCollector()
		os.Exit(0)
	}

	g.ParseConfig(*cfg)

	g.InitRootDir()
	g.InitLocalIps()
	g.InitRpcClients()

	funcs.BuildMappers()

	go cron.InitDataHistory()
	// 上报本机状态
	cron.ReportAgentStatus()
	// 同步插件
	cron.SyncMinePlugins()
	// 同步监控端口、路径、进程和URL
	cron.SyncBuiltinMetrics()
	// 后门调试agent,允许执行shell指令的ip列表
	cron.SyncTrustableIps()
	// 开始数据次采集
	cron.Collect()
	// 启动dashboard server
	go http.Start()

	select {}

}
