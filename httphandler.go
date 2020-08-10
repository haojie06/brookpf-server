package main

import (
	"fmt"
	"io/ioutil"
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
		fmt.Println("命令为空")
		return nil
	}
	fmt.Println("将要执行命令" + cmdstr)
	cmd := exec.Command("/bin/bash", "-c", cmdstr)
	//打开命令的标准输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		//打开输出管道失败
		fmt.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
		return nil
	}
	//执行命令 !注意，这样写的时候err为局部参数，只在if的作用域中有效
	if err := cmd.Start(); err != nil {
		fmt.Println("Error:The command is err,", err)
		return nil
	}
	//读取命令输出
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println("ReadAll Stdout:", err.Error())
		return nil
	}
	//阻塞等待到命令执行完毕，获取输出
	if err := cmd.Wait(); err != nil {
		fmt.Println("wait:", err.Error())
		return nil
	}
	//执行到这一步，命令已经执行完毕，也获得了命令的输出
	fmt.Printf("%s", bytes)
	return bytes
}

func commandHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//获取命令
	fmt.Println(r.Form["cmd"][0])
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	output := executeCommand(r.Form["cmd"][0])
	fmt.Fprintf(w, "%s", output)
}

//获取服务器状态与上面的brook状态 /api/getstatus
func getStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Println("查询服务器状态")
	//是否在线不用专门做，只要能返回信息就是在线
	//查询brook是否安装
	if _, err := os.Stat(brook_file); err == nil {
		fmt.Printf("Brook已经安装\n")
	} else {
		fmt.Printf("Brook未安装\n")
	}
	//查询brook是否启动
	pid := executeCommand(`ps -ef| grep "brook relays"| grep -v grep| grep -v ".sh"| grep -v "init.d"| grep -v "service"| awk '{print $2}'`)
	if spid := string(pid); spid == "" {
		fmt.Println("Brook未启动")
	} else {
		fmt.Println("Brook已启动 PID:" + spid)
	}
	//返回端口列表
	//先查看配置文件是否存在
	if _, err := os.Stat(brook_conf); err == nil {
		fmt.Println("Brook配置文件存在:")
		if dat, err := ioutil.ReadFile(brook_conf); err == nil {
			//去除一下首尾的字符
			datas := strings.Split(strings.TrimSpace(string(dat)), "\n")
			for index, data := range datas {
				fmt.Printf("%d.:%s\n", index, data)
			}
		} else {
			fmt.Println("打开配置文件失败" + err.Error())
		}
	} else {
		fmt.Println("Brook配置文件不存在")
	}
}

//重启brook
//添加端口转发
//删除端口转发
//修改端口转发
//编辑端口转发
//func addPf(w http.ResponseWriter, r *http.Request)
