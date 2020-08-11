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
	log.Println("[命令执行]将要执行:" + cmdstr)
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
	log.Printf("[命令执行]	结果: %s", bytes)
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
	//考虑到要传递密码，又不想做会话管理，那么全部使用post方法好了
	if !auth(w, r) {
		return
	}
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
	if !auth(w, r) {
		return
	}
	stop := executeCommand("/etc/init.d/brook-pf stop")
	log.Printf("[关闭Brook]")
	js, _ := json.Marshal(MessageResponse{Code: 200, Msg: string(stop)})
	w.Write(js)
	return
}

//重启brook
func restartBrook(w http.ResponseWriter, r *http.Request) {
	if !auth(w, r) {
		return
	}
	stop := executeCommand("/etc/init.d/brook-pf stop")
	start := executeCommand("/etc/init.d/brook-pf start")
	log.Printf("[重启Brook]")
	js, _ := json.Marshal(MessageResponse{Code: 200, Msg: string(stop) + "\n" + string(start)})
	w.Write(js)
	return
}

//启动brook
func startBrook(w http.ResponseWriter, r *http.Request) {
	if !auth(w, r) {
		return
	}
	start := executeCommand("/etc/init.d/brook-pf start")
	log.Printf("[启动Brook]")
	js, _ := json.Marshal(MessageResponse{Code: 200, Msg: string(start)})
	w.Write(js)
	return
}

/**
*添加端口转发
*先检查端口是否已经有了记录
*使用post方法
**/

func addPortForward(w http.ResponseWriter, r *http.Request) {
	if !auth(w, r) {
		return
	}
	var messageResponse MessageResponse
	if r.Method != http.MethodPost {
		messageResponse.Code = 405
		messageResponse.Msg = "错误的请求方法"
		//w.WriteHeader(405)
		log.Println("[添加端口转发]	错误的请求方法")
		//return
	} else {
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
		request["Name"] = r.PostFormValue("Name")
		request["Description"] = r.PostFormValue("Description")
		//处理一下非法输入
		log.Println("[添加端口转发]	请求记录", request)
		for key, value := range request {
			fmt.Println("KEY:" + key + "---" + "VALUE:" + value)
		}
		//首先检查端口是否重复

		completed := true
		for index, data := range request {
			if data == "" {
				//request post参数不完整
				completed = false
				messageResponse.Code = 400
				messageResponse.Msg = "request参数不完整 缺少" + index
				log.Printf("[添加中转端口]	request参数不完整，缺少%s\n", index)
			}
		}
		if completed {
			var ifDuplicated bool = false
			if _, err := os.Stat(brook_conf); err == nil {
				log.Println("[添加端口转发]Brook配置文件存在:")
				if dat, err := ioutil.ReadFile(brook_conf); err == nil {
					//去除一下首尾的字符
					datas := strings.Split(strings.TrimSpace(string(dat)), "\n")
					for index, data := range datas {
						log.Printf("%d.:%s\n", index, data)
						lp := strings.Split(data, " ")[0]
						if request["LocalPort"] == lp {
							log.Println("[添加端口转发]端口冲突，无法添加转发:" + lp)
							ifDuplicated = true
						}
					}
					if ifDuplicated {
						//端口冲突 无法添加
						messageResponse.Code = 400
						messageResponse.Msg = "端口冲突无法添加转发"

					} else {
						log.Println("可以添加转发端口" + request["LocalPort"])
						//写入文件
						f, err := os.OpenFile(brook_conf, os.O_APPEND|os.O_WRONLY, 0600)
						if err != nil {
							log.Println("[添加端口转发]打开配置文件出错", err.Error())

						}
						if _, err := f.Write([]byte(request["LocalPort"] + " " + request["Host"] + " " + request["RemotePort"] + " " + request["Enable"] + " " + request["Name"] + " " + request["Description"] + "\n")); err != nil {
							log.Fatal(err)
						}
						if err := f.Close(); err != nil {
							log.Fatal(err)
						} else {
							log.Println("[添加端口转发]写入配置文件")
						}
						//判断是否成功添加（看上去不判断也没问题，上一步只要没执行出错就不会有问题。）
						//修改iptables
						log.Println("[添加端口转发]修改iptables")
						log.Printf("%s\n", string(executeCommand("iptables -I INPUT -m state --state NEW -m tcp -p tcp --dport "+request["LocalPort"]+" -j ACCEPT")))
						log.Printf("%s\n", string(executeCommand("iptables -I INPUT -m state --state NEW -m udp -p udp --dport "+request["LocalPort"]+" -j ACCEPT")))
						if release == "centos" {
							log.Printf("%s\n", string(executeCommand("service iptables save")))
						} else if release == "ubuntu" {
							log.Printf("%s\n", string(executeCommand("iptables-save > /etc/iptables.up.rules")))
						} else {
							log.Println("脚本不支持当前发行版", release)
						}
						//重启brook
						executeCommand("/etc/init.d/brook-pf stop")
						executeCommand("/etc/init.d/brook-pf start")
						log.Println("成功添加端口转发")
						messageResponse.Code = 200
						messageResponse.Msg = "成功添加端口转发"
					}
				} else {
					log.Println("[添加端口转发]打开配置文件失败" + err.Error())
					messageResponse.Code = 400
					messageResponse.Msg = "打开brook配置文件失败"
				}
			} else {
				log.Println("[添加端口转发]rook配置文件不存在")
				messageResponse.Code = 400
				messageResponse.Msg = "[添加端口转发]brook配置文件不存在"
			}
		}

		w.Header().Set("Content-Type", "application/json")
		js, err := json.Marshal(messageResponse)
		if err != nil {
			log.Println("JSON转换失败" + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(js)
		return
	}

}

//删除端口转发
//检查本地是否有该端口
//直接调用shell删除就好了
func deletePortForward(w http.ResponseWriter, r *http.Request) {
	if !auth(w, r) {
		return
	}
	var mr MessageResponse
	r.ParseForm()
	port := r.Form["port"]
	log.Println("[删除端口转发] 开始删除, 本地端口:", port)
	if dat, err := ioutil.ReadFile(brook_conf); err == nil {
		//去除一下首尾的字符
		datas := strings.Split(strings.TrimSpace(string(dat)), "\n")
		found := false
		for index, data := range datas {
			lport := strings.Split(data, " ")[0]
			if port[0] == lport {
				log.Printf("[删除端口转发]	找到对应端口记录: %d.:%s", index, data)
				found = true
				//进行删除
				ret := executeCommand(`sed -i "/^` + lport + `/d" ` + brook_conf)
				log.Println("[删除端口转发]	" + string(ret))
				//查看是否成功删除 暂时不做这个
				//下面这个方法有可能出现bug，因为如果域名，ip中带有和端口一样的数字，则无法判断是否成功删除
				//ret = executeCommand(`cat ` + brook_conf + `| grep ` + lport)
				//重启一下 brook
				executeCommand("/etc/init.d/brook-pf stop")
				executeCommand("/etc/init.d/brook-pf start")
				mr.Code = 200
				mr.Msg = "已删除端口转发记录"

				break
			}
		}
		if !found {
			log.Println("[删除端口转发]	配置文件中未找到对应端口")
			mr.Code = 400
			mr.Msg = "配置文件中未找到对应端口"
		}
	} else {
		log.Println("[删除端口转发 ]打开配置文件失败" + err.Error())
		mr.Code = 400
		mr.Msg = "打开配置文件失败"
	}
	js, _ := json.Marshal(mr)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return
}

//查找端口转发列表
func listPortForward(w http.ResponseWriter, r *http.Request) {
	if !auth(w, r) {
		return
	}
	var dr DataResponse
	if dat, err := ioutil.ReadFile(brook_conf); err == nil {
		//去除一下首尾的字符
		datas := strings.Split(strings.TrimSpace(string(dat)), "\n")
		dr.Code = 200
		dataMap := make(map[string]interface{})
		dataMap["records"] = datas
		dr.Data = dataMap
	} else {
		log.Println("[查询端口转发 ]打开配置文件失败" + err.Error())
		dr.Code = 400
		dr.Data = nil
	}
	js, err := json.Marshal(dr)
	if err != nil {
		log.Println("[查询端口转发]	JSON转换失败")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

//修改端口转发 也是使用post方法
func editPortForward(w http.ResponseWriter, r *http.Request) {
	if !auth(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req AddPortForwardRequest
	req.LocalPort = r.PostFormValue("LocalPort")
	req.RemotePort = r.PostFormValue("RemotePort")
	req.NewPort = r.PostFormValue("NewPort")
	req.Host = r.PostFormValue("Host")
	req.Name = r.PostFormValue("Name")
	req.Description = r.PostFormValue("Description")
	//检查本地端口，有的话删除那一行再加一行新的。
	var dr MessageResponse
	if dat, err := ioutil.ReadFile(brook_conf); err == nil {
		//去除一下首尾的字符
		datas := strings.Split(strings.TrimSpace(string(dat)), "\n")

		found := false
		for index, data := range datas {
			lport := strings.Split(data, " ")[0]
			if req.LocalPort == lport {
				log.Printf("[编辑端口转发]找到对应端口记录: %d.:%s %s==%s", index, data, req.LocalPort, lport)
				found = true
				//进行删除
				ret := executeCommand(`sed -i "/^` + lport + `/d" ` + brook_conf)
				log.Println("[编辑端口转发]	删除记录" + string(ret))
				//删除后添加新记录
				addLine := req.NewPort + " " + req.Host + " " + req.RemotePort + " " + req.Enable + " " + req.Name + " " + req.Description
				ret = executeCommand(`echo "` + addLine + `" >> ` + brook_conf)
				log.Println("[编辑端口转发]添加记录" + string(ret))
				//最后重启一下
				log.Println("[编辑端口转发]重启使编辑生效")
				executeCommand("/etc/init.d/brook-pf stop")
				executeCommand("/etc/init.d/brook-pf start")
				dr.Code = 200
				dr.Msg = "成功编辑"
				break
			}
		}
		if !found {
			log.Println("[编辑端口转发]	配置文件中未找到对应端口")
			dr.Code = 400
			dr.Msg = "配置文件中未找到对应端口"
		}
	} else {
		log.Println("[删除端口转发 ]打开配置文件失败" + err.Error())
		dr.Code = 400
		dr.Msg = "打开配置文件失败"
	}
	w.Header().Set("Content-Type", "application/json")
	js, _ := json.Marshal(dr)
	w.Write(js)
}

//授权
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
