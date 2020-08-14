package main

import "log"

var (
	CheckReleaseScript = `if [[ -f /etc/redhat-release ]]; then
	release="centos"
elif cat /etc/issue | grep -q -E -i "debian"; then
	release="debian"
elif cat /etc/issue | grep -q -E -i "ubuntu"; then
	release="ubuntu"
elif cat /etc/issue | grep -q -E -i "centos|red hat|redhat"; then
	release="centos"
elif cat /proc/version | grep -q -E -i "debian"; then
	release="debian"
elif cat /proc/version | grep -q -E -i "ubuntu"; then
	release="ubuntu"
elif cat /proc/version | grep -q -E -i "centos|red hat|redhat"; then
	release="centos"
fi
echo ${release}`
)

//修改系统iptables true为增加，false为删除 同时还要增加流量统计
func changeIptables(add bool, port string) {
	//修改iptables
	log.Println("[iptables修改]修改iptables")
	if add {
		log.Printf("[iptables修改]%s\n", string(executeCommand("iptables -D INPUT -p tcp --dport "+port)))
		log.Printf("[iptables修改]%s\n", string(executeCommand("iptables -D INPUT -p udp --dport "+port)))
		log.Printf("[iptables修改]%s\n", string(executeCommand("iptables -D OUTPUT -p tcp --sport "+port)))
		log.Printf("[iptables修改]%s\n", string(executeCommand("iptables -D OUTPUT -p udp --sport "+port)))
		log.Printf("[iptables修改]添加TCP:\n%s\n", string(executeCommand("iptables -I INPUT -m state --state NEW -m tcp -p tcp --dport "+port+" -j ACCEPT")))
		log.Printf("[iptables修改]添加UDP:\n%s\n", string(executeCommand("iptables -I INPUT -m state --state NEW -m udp -p udp --dport "+port+" -j ACCEPT")))
		log.Printf("[iptables修改]%s\n", string(executeCommand("iptables -I INPUT -p tcp --dport "+port)))
		log.Printf("[iptables修改]%s\n", string(executeCommand("iptables -I INPUT -p udp --dport "+port)))
		log.Printf("[iptables修改]%s\n", string(executeCommand("iptables -I OUTPUT -p tcp --sport "+port)))
		log.Printf("[iptables修改]%s\n", string(executeCommand("iptables -I OUTPUT -p udp --sport "+port)))
		// 							iptables -I INPUT -p tcp --dport $port
		// iptables -I INPUT -p udp --dport $port
		// iptables -I OUTPUT -p tcp --sport $port
		// iptables -I OUTPUT -p udp --sport $port
	} else {
		log.Printf("[iptables修改]%s\n", string(executeCommand("iptables -D INPUT -m state --state NEW -m tcp -p tcp --dport "+port+" -j ACCEPT")))
		log.Printf("[iptables修改]%s\n", string(executeCommand("iptables -D INPUT -m state --state NEW -m udp -p udp --dport "+port+" -j ACCEPT")))
		log.Printf("[iptables修改]%s\n", string(executeCommand("iptables -D INPUT -p tcp --dport "+port)))
		log.Printf("[iptables修改]%s\n", string(executeCommand("iptables -D INPUT -p udp --dport "+port)))
		log.Printf("[iptables修改]%s\n", string(executeCommand("iptables -D OUTPUT -p tcp --sport "+port)))
		log.Printf("[iptables修改]%s\n", string(executeCommand("iptables -D OUTPUT -p udp --sport "+port)))
	}
	//保存对iptables的修改
	if release == "centos" {
		log.Printf("%s\n", string(executeCommand("service iptables save")))
	} else if release == "ubuntu" {
		log.Printf("%s\n", string(executeCommand("iptables-save > /etc/iptables.up.rules")))
	} else {
		log.Println("脚本不支持当前发行版", release)
	}
}
