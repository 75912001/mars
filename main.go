// mars project main.go
package main

import (
	"math/rand"
	"runtime"
	"sms"
	"time"

	"github.com/75912001/goz/zhttp"
	"github.com/75912001/goz/ztcp"
	"github.com/75912001/goz/zutility"
)

func main() {
	rand.Seed(time.Now().Unix())
	////////////////////////////////////////////////////////////////////////////////
	//加载配置文件bench.ini
	if zutility.IsWindows() {
		ini.Load("./bench.ini.bak")
	} else if zutility.IsDarwin() {
		ini.Load("/Users/mlc/go/src/github.com/75912001/mars/bench.ini.bak")
	} else {
		return
	}

	goProcessMax := ini.GetInt("common", "go_process_max", runtime.NumCPU())
	runtime.GOMAXPROCS(goProcessMax)
	////////////////////////////////////////////////////////////////////////////////
	//初始化log
	log_level := ini.GetInt("log", "log_level", 8)
	log_path := ini.GetString("log", "path", "default.log.")
	log.Init(log_path, 1000)
	log.SetLevel(log_level)
	defer log.DeInit()

	////////////////////////////////////////////////////////////////////////////////
	//redis
	{
		const benchFileSection string = "redis_server"
		ip := ini.GetString(benchFileSection, "ip", "")
		port := ini.GetUint16(benchFileSection, "port", 0)
		redisDatabases := ini.GetInt(benchFileSection, "databases", 0)

		//链接redis
		err := redis.server.Connect(ip, port, redisDatabases)
		if nil != err {
			log.Crit("redis connect err:", err)
			return
		}
	}

	////////////////////////////////////////////////////////////////////////////////
	//初始化服务器tcp
	{
		ztcp.SetLog(&log)
		//设置回调函数
		tcpServer.OnInit = onInit
		tcpServer.OnFini = onFini
		tcpServer.OnPeerConnClosed = onCliConnClosed
		tcpServer.OnPeerConn = onCliConn
		tcpServer.OnPeerPacket = onCliPacket
		tcpServer.OnParseProtoHead = onParseProtoHead

		//运行
		delay := true

		tcpServer.PacketLengthMax = ini.GetUint32("server", "packet_length_max", 81920)

		serverIP := ini.GetString("server", "ip", "")
		serverPort := ini.GetUint16("server", "port", 0)
		recvChanMaxCnt := ini.GetUint32("server", "recv_chan_max_cnt", 1000)
		if 0 != serverPort {
			log.Trace(serverIP, serverPort, delay)
			go tcpServer.Run(serverIP, serverPort, delay, recvChanMaxCnt)
		}
	}
	////////////////////////////////////////////////////////////////////////////////
	//HTTP SERVER
	{
		ip := ini.GetString("http_server", "ip", "")
		port := ini.GetUint16("http_server", "port", 0)
		zhttp.SetLog(&log)

		{ //启动手机注册功能
			err := sms.GSmsRegister.Init(ini)
			if nil != err {
				log.Crit("sms.GSmsRegister.Init err:", err)
				return
			}
			httpServer.AddHandler(sms.Pattern, sms.SmsRegisterHttpHandler)
		}
		go httpServer.Run(ip, port)
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
