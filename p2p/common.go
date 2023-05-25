package p2p

import (
	"fmt"
	"net"
)

// GetHost 获取本地主机的有效 IPv4 地址
func GetHost() string {
	// 获取主机的网络接口列表
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("无法获取网络接口列表：", err)
		return ""
	}
	// 遍历网络接口列表，查找非回环接口的有效 IPv4 地址
	for _, iface := range interfaces {
		// 排除回环接口和无效接口
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			addresses, err := iface.Addrs()
			if err != nil {
				fmt.Println("无法获取接口地址：", err)
				continue
			}
			// 遍历接口地址，查找有效 IPv4 地址
			for _, addr := range addresses {
				ipNet, ok := addr.(*net.IPNet)
				if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
					return ipNet.IP.String()
				}
			}
		}
	}
	return ""
}
