package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

//登录验证
func login(w http.ResponseWriter, r *http.Request) {
	if !auth(w, r) {
		log.Println("[登录验证]登录失败 用户ip:", r.RemoteAddr)
		return
	}
	mr := MessageResponse{Code: 200, Msg: "登陆成功"}
	js, _ := json.Marshal(mr)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	log.Println("[登录验证]登录成功 用户ip:", r.RemoteAddr)
	return
}

//获取服务器列表
func getServers(w http.ResponseWriter, r *http.Request) {
	if !auth(w, r) {
		return
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Println(err.Error())
	}
	dataMap := make(map[string]interface{})
	dataMap["servers"] = config.Servers
	mr := DataResponse{Code: 200, Msg: "查询成功", Data: dataMap}
	js, _ := json.Marshal(mr)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	log.Println("[服务器列表查询] 查询服务器列表")
	return
}

//添加服务器
func addServer(w http.ResponseWriter, r *http.Request) {
	if !auth(w, r) {
		return
	}
	var newServer Server
	newServer.Name = r.PostFormValue("Name")
	newServer.IP = r.PostFormValue("IP")
	newServer.UserName = r.PostFormValue("Username")
	newServer.Password = r.PostFormValue("Password")
	newServer.Port = r.PostFormValue("Port")
	newServer.Desc = r.PostFormValue("Desc")
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Println(err.Error())
	}
	servers := config.Servers
	newServer.ID = len(servers) + 1
	servers = append(servers, newServer)
	vip.Set("servers", servers)
	var mr MessageResponse
	if err := vip.WriteConfig(); err != nil {
		log.Println("[添加配置服务器]错误:", err.Error())
		mr = MessageResponse{Code: 400, Msg: err.Error()}
	} else {
		mr = MessageResponse{Code: 200, Msg: "成功添加"}
	}
	js, _ := json.Marshal(mr)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return
}

//删除服务器
func delServer(w http.ResponseWriter, r *http.Request) {
	if !auth(w, r) {
		return
	}
}

//修改服务器

//登录验证以及授权
func auth(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return false
	}
	uname := r.PostFormValue("Username")
	pword := r.PostFormValue("Password")
	if uname != username || pword != password {
		//授权失败
		log.Printf("授权失败 name:%s password:%s \n", uname, pword)
		mr := MessageResponse{Code: 401, Msg: "授权失败"}
		js, _ := json.Marshal(mr)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return false
	} else {
		//授权成功
		return true
	}
}
