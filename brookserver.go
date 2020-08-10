package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
	fmt.Println("接收到请求")
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
	if err := cmd.Wait(); err != nil {
		fmt.Println("wait:", err.Error())
		return
	}

	fmt.Printf("命令输出stdout:\n\n %s", bytes)
	fmt.Fprintf(w, "%s", bytes)
}
func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/democ", demoCommandHandler)
	http.ListenAndServe(":8000", nil)
}
