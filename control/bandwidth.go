package main

import (
	"regexp"
	"strings"
)

//流量查询
func getBandwidth(port string) string {
	//流量字节切片
	bandWidth := []string{"0", "0", "0", "0"}
	//先查询入
	outStr := string(executeCommand("iptables -n -v -x -t filter -L OUTPUT |grep " + port))
	outStr = strings.TrimSpace(delete_extra_space(outStr))
	// fmt.Printf("流量统计%s端口出端口流量\n%s", port, outStr)
	outRecords := strings.Split(outStr, "\n")
	// fmt.Println("出端口流量")
	if len(outRecords) != 1 {
		for _, r := range outRecords {
			items := strings.Split(strings.TrimSpace(delete_extra_space(r)), " ")
			if items[2] == "tcp" {
				// fmt.Printf("tcp出端口:%s\n", items)
				bandWidth[0] = items[1]
			} else if items[2] == "udp" {
				// fmt.Printf("udp出端口:%s\n", items)
				bandWidth[2] = items[1]
			}
		}
	}
	//查询出端口流量
	inStr := string(executeCommand("iptables -n -v -x -t filter -L OUTPUT |grep " + port))
	inStr = strings.TrimSpace(inStr)
	inRecords := strings.Split(inStr, "\n")
	// fmt.Println("入端口流量")
	if len(inRecords) != 1 {
		for _, r := range inRecords {
			// fmt.Printf("%s\n", strings.TrimSpace(delete_extra_space(r)))
			items := strings.Split(strings.TrimSpace(delete_extra_space(r)), " ")
			if items[2] == "tcp" {
				// fmt.Printf("tcp入端口:%s\n", items)
				bandWidth[1] = items[1]
			} else if items[2] == "udp" {
				// fmt.Printf("udp入端口:%s\n", items)
				bandWidth[3] = items[1]
			}
		}
	}
	//tcp出 tcp入 udp出 udp入
	return bandWidth[0] + " " + bandWidth[1] + " " + bandWidth[2] + " " + bandWidth[3]
}

func delete_extra_space(s string) string {
	//删除字符串中的多余空格，有多个空格时，仅保留一个空格
	s1 := strings.Replace(s, "	", " ", -1)       //替换tab为空格
	regstr := "\\s{2,}"                          //两个及两个以上空格的正则表达式
	reg, _ := regexp.Compile(regstr)             //编译正则表达式
	s2 := make([]byte, len(s1))                  //定义字符数组切片
	copy(s2, s1)                                 //将字符串复制到切片
	spc_index := reg.FindStringIndex(string(s2)) //在字符串中搜索
	for len(spc_index) > 0 {                     //找到适配项
		s2 = append(s2[:spc_index[0]+1], s2[spc_index[1]:]...) //删除多余空格
		spc_index = reg.FindStringIndex(string(s2))            //继续在字符串中搜索
	}
	return string(s2)
}
