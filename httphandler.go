package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "hello world")
}

func executeCommand(cmdstr string) []byte {
	if cmdstr == "" {
		return nil
	}
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
	fmt.Println("%s", bytes)
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

func demoCommandHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	cmd := exec.Command("/bin/bash", "-c", "ifconfig")

	//创建获取命令输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
		return
	}
	//执行命令
	if err := cmd.Start(); err != nil {
		fmt.Println("Error:The command is err,", err)
		return
	}
	//读取所有输出
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println("ReadAll Stdout:", err.Error())
		return
	}
	//阻塞等待到命令执行完毕
	if err := cmd.Wait(); err != nil {
		fmt.Println("wait:", err.Error())
		return
	}
	fmt.Fprintf(w, "%s", bytes)

}
