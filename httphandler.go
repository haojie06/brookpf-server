package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var (
	brook_file        = "/usr/local/brook-pf/brook"
	brook_conf        = "/usr/local/brook-pf/brook.conf"
	brook_log         = "/usr/local/brook-pf/brook.log"
	brook_server_conf = "/usr/local/brook-pf/bserver.conf"
	Crontab_file      = "/usr/bin/crontab"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "hello world")
}

func executeCommand(cmdstr string) []byte {
	if cmdstr == "" {
		log.Println("命令为空")
		return nil
	}
	log.Println("将要执行命令" + cmdstr)
	cmd := exec.Command("/bin/bash", "-c", cmdstr)
	//打开命令的标准输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		//打开输出管道失败
		log.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
		return nil
	}
	//执行命令 !注意，这样写的时候err为局部参数，只在if的作用域中有效
	if err := cmd.Start(); err != nil {
		log.Println("Error:The command is err,", err)
		return nil
	}
	//读取命令输出
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Println("ReadAll Stdout:", err.Error())
		return nil
	}
	//阻塞等待到命令执行完毕，获取输出
	if err := cmd.Wait(); err != nil {
		log.Println("wait:", err.Error())
		return nil
	}
	//执行到这一步，命令已经执行完毕，也获得了命令的输出
	log.Printf("%s", bytes)
	return bytes
}

func commandHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//获取命令
	log.Println(r.Form["cmd"][0])
	log.Println("path", r.URL.Path)
	log.Println("scheme", r.URL.Scheme)
	output := executeCommand(r.Form["cmd"][0])
	fmt.Fprintf(w, "%s", output)
}

//获取服务器状态与上面的brook状态 /api/getstatus
func getStatus(w http.ResponseWriter, r *http.Request) {
	log.Println("查询服务器状态")
	//是否在线不用专门做，只要能返回信息就是在线
	//查询brook是否安装
	var Response StatusResponse
	Response.Code = 200
	if _, err := os.Stat(brook_file); err == nil {
		log.Printf("Brook已经安装\n")
		Response.Installed = true
	} else {
		log.Printf("Brook未安装\n")
		Response.Installed = false
	}
	//查询brook是否启动
	pid := executeCommand(`ps -ef| grep "brook relays"| grep -v grep| grep -v ".sh"| grep -v "init.d"| grep -v "service"| awk '{print $2}'`)
	if spid := string(pid); spid == "" {
		log.Println("Brook未启动")
		Response.Enable = false
	} else {
		log.Println("Brook已启动 PID:" + spid)
		Response.Enable = true
	}
	//返回端口列表
	//先查看配置文件是否存在
	if _, err := os.Stat(brook_conf); err == nil {
		log.Println("Brook配置文件存在:")
		if dat, err := ioutil.ReadFile(brook_conf); err == nil {
			//去除一下首尾的字符
			datas := strings.Split(strings.TrimSpace(string(dat)), "\n")
			for index, data := range datas {
				log.Printf("%d.:%s\n", index, data)
			}
			Response.Records = datas
		} else {
			log.Println("打开配置文件失败" + err.Error())
			Response.Code = 201
		}
	} else {
		log.Println("Brook配置文件不存在")
		Response.Code = 201
	}
	js, err := json.Marshal(Response)
	log.Println(string(js))
	if err != nil {
		log.Println("JSON转换失败" + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

//关闭brook
func stopBrook(w http.ResponseWriter, r *http.Request) {
	stop := executeCommand("/etc/init.d/brook-pf stop")
	log.Printf("关闭:%s", stop)

}

//重启brook
func restartBrook(w http.ResponseWriter, r *http.Request) {
	stop := executeCommand("/etc/init.d/brook-pf stop")
	start := executeCommand("/etc/init.d/brook-pf start")
	log.Printf("%s\n%s", stop, start)
}

//启动brook
func startBrook(w http.ResponseWriter, r *http.Request) {
	start := executeCommand("/etc/init.d/brook-pf start")
	log.Printf("%s", start)
}

/**
*添加端口转发
*先检查端口是否已经有了记录
*使用post方法
**/

func addPortForward(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		log.Println("错误的请求方法")
		return
	}
	/*
		如果要用form-data而不是json来发送，那么就不要先解析
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(400)
			return
		}
	*/
	request := make(map[string]string)
	request["LocalPort"] = r.PostFormValue("LocalPort")
	request["RemotePort"] = r.PostFormValue("RemotePort")
	request["Host"] = r.PostFormValue("Host")
	request["Enable"] = r.PostFormValue("Enable")
	log.Println(request)
	for key, value := range request {
		fmt.Println("KEY:" + key + "---" + "VALUE:" + value)
	}
	//首先检查端口是否重复
	var ifDuplicated bool = false
	if _, err := os.Stat(brook_conf); err == nil {
		log.Println("[检查端口]Brook配置文件存在:")
		if dat, err := ioutil.ReadFile(brook_conf); err == nil {
			//去除一下首尾的字符
			datas := strings.Split(strings.TrimSpace(string(dat)), "\n")
			for index, data := range datas {
				log.Printf("%d.:%s\n", index, data)
				lp := strings.Split(data, " ")[0]
				if request["LocalPort"] == lp {
					log.Println("该端口已经添加过转发:" + lp)
					ifDuplicated = true
				}
			}
			if !ifDuplicated {
				log.Println("可以添加转发端口" + request["LocalPort"])
			} else {
				w.WriteHeader(401)
				//端口冲突 无法添加
			}
		} else {
			log.Println("[检查端口]打开配置文件失败" + err.Error())
			w.WriteHeader(http.StatusExpectationFailed)
			fmt.Fprintf(w, err.Error())
		}
	} else {
		log.Println("[检查端口]rook配置文件不存在")
		w.WriteHeader(http.StatusExpectationFailed)
		fmt.Fprintf(w, err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	return
	//log.Println(r.Form["cmd"][0])
	//log.Println("path", r.URL.Path)
	//log.Println("scheme", r.URL.Scheme)
}

//删除端口转发
//修改端口转发
//编辑端口转发
//func addPf(w http.ResponseWriter, r *http.Request)
