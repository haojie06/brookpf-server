package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	config_path = "./"
	config_name = "brookpf"
	config_type = "yaml"
	release     = ""
	username    = ""
	password    = ""
)

func main() {
	//配置文件相关设置
	//先检测配置文件是否存在
	viper.SetConfigType(config_type)
	viper.SetConfigName(config_name)
	viper.AddConfigPath(config_path)
	//默认配置文件，如果指定路径没有配置文件则使用该配置来创建
	viper.SetDefault("info", "")
	viper.SetDefault("password", "admin")
	viper.SetDefault("username", "admin")
	viper.SetDefault("port", 8000)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("未找到配置文件，自动创建\n%s\n", err)
		//注意 必须是write as 才会创建新文件（如果文件不存在的话）
		viper.WriteConfigAs(config_path + config_name + "." + config_type)
	}
	viper.WatchConfig() //监听配置变化
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置发生变更,请重启服务端应用变化")
	})
	port := viper.GetString("port")
	username = viper.GetString("username")
	password = viper.GetString("password")

	//检查发行版
	fmt.Println("检查发行版")
	ret := executeCommand(CheckReleaseScript)
	release = strings.TrimSpace((string(ret)))
	fmt.Println("发行版为:" + release)
	//http.HandleFunc("/", indexHandler)
	//http.HandleFunc("/command", commandHandler)
	http.HandleFunc("/api/band", queryBandwidth)
	http.HandleFunc("/api/getstatus", getStatus)
	http.HandleFunc("/api/stopbrook", stopBrook)
	http.HandleFunc("/api/addpf", addPortForward)
	http.HandleFunc("/api/delpf", deletePortForward)
	http.HandleFunc("/api/listpf", listPortForward)
	http.HandleFunc("/api/editpf", editPortForward)
	http.HandleFunc("/api/stop", stopBrook)
	http.HandleFunc("/api/start", startBrook)
	http.HandleFunc("/api/restart", restartBrook)
	http.HandleFunc("/api/switchpf", changePortForward)

	fmt.Println("Brook面板受控端开始运行")
	fmt.Println("配置文件定义 端口:", port, " 用户名:", username, "密码:", password)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("服务器启动错误:\n" + err.Error())
	}
}
