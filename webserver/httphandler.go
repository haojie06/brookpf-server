package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type MessageResponse struct {
	Code int
	Msg  string
}

//登录验证
func login(w http.ResponseWriter, r *http.Request) {
	if !auth(w, r) {
		return
	}
	mr := MessageResponse{Code: 200, Msg: "登陆成功"}
	js, _ := json.Marshal(mr)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return
}

//获取服务器列表
//添加服务器
//删除服务器
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
