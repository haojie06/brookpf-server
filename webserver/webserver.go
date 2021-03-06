package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

var (
	config_path = "./"
	config_name = "webserver"
	config_type = "yaml"
	release     = ""
	username    = ""
	password    = ""
	vip         = viper.GetViper()
)

//这个为面板的后端，功能较为简单，从配置文件中读取面板用户和密码以及负责读写该用户已经添加的服务器 端口 密码 备注

func main() {
	viper.SetConfigType(config_type)
	viper.SetConfigName(config_name)
	viper.AddConfigPath(config_path)

	viper.SetDefault("desc", "描述")
	viper.SetDefault("broadcast", "公告")

	viper.SetDefault("port", 8001)
	viper.SetDefault("password", "admin")
	viper.SetDefault("username", "admin")
	//string map的切片
	defaultServers := make([]Server, 0)
	//默认服务器配置
	defaultServer := Server{}
	defaultServer.ID = 1
	defaultServer.Name = "tempserver"
	defaultServer.Desc = "temps desc"
	defaultServer.IP = "127.0.0.1"
	defaultServer.Port = "8000"
	defaultServer.UserName = "admin"
	defaultServer.Password = "admin"
	defaultServers = append(defaultServers, defaultServer)
	viper.SetDefault("servers", defaultServers)
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("未找到配置文件，自动创建\n%s\n", err)
		//注意 必须是write as 才会创建新文件（如果文件不存在的话）
		viper.WriteConfigAs(config_path + config_name + "." + config_type)
	}
	//此处为面板配置
	//下面读取默认服务器列表
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Println(err.Error())
	}
	port := config.Port
	username = config.UserName
	password = config.Password
	info := viper.GetString("info")

	//serverList := config.Servers
	//serverlist := viper.GetStringMap("servers")
	log.Printf("Brookweb后端程序开始运行\nusername:%s\npassword:%s\nport:%s\n", username, password, port)
	fmt.Println(`服务器信息` + info)
	//log.Println(serverList)
	//绑定监听方法
	http.HandleFunc("/web/login", login)
	http.HandleFunc("/web/getservers", getServers)
	http.HandleFunc("/web/addserver", addServer)
	http.HandleFunc("/web/delserver", delServer)
	http.HandleFunc("/web/editserver", editServer)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("服务器启动错误:\n" + err.Error())
	}

}
