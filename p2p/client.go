package p2p

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

// 解析地址函数，格式为（ip:port）
func parseAddr(addr string) net.UDPAddr {
	t := strings.Split(addr, ":")
	port, _ := strconv.Atoi(t[1])
	return net.UDPAddr{
		IP:   net.ParseIP(t[0]),
		Port: port,
	}
}

func TODO() {
	if len(os.Args) < 5 {
		fmt.Println("./client tag remoteIP remotePort port")
		return
	}
	port, _ := strconv.Atoi(os.Args[4])
	tag := os.Args[1]
	remoteIP := os.Args[2]
	remotePort, _ := strconv.Atoi(os.Args[3])

	//一定要绑定固定端口，否则介绍人不好介绍
	localAddr := net.UDPAddr{Port: port}

	//与服务器建立联系（严格意义上，UDP不能叫连接）
	conn, err := net.DialUDP("udp", &localAddr, &net.UDPAddr{IP: net.ParseIP(remoteIP), Port: remotePort})
	if err != nil {
		log.Panic("Failed ot DialUDP", err)
	}

	//自我介绍，亮明身份，但其实说啥都行
	conn.Write([]byte("我是:" + tag))

	buf := make([]byte, 256)
	//从服务器获得目标地址
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		log.Panic("Failed to ReadFromUDP", err)
	}
	conn.Close()
	toAddr := parseAddr(string(buf[:n]))
	fmt.Println("获得对象地址:", toAddr)
	//两个人建立P2P通信
	p2p(&localAddr, &toAddr)

}

func p2p(srcAddr *net.UDPAddr, dstAddr *net.UDPAddr) {
	//1. 请求与对方建立联系
	conn, _ := net.DialUDP("udp", srcAddr, dstAddr)
	//2.发送打洞消息
	conn.Write([]byte("打洞消息\n"))

	//启动一个goroutine监控标准输入
	go func() {
		buf := make([]byte, 256)
		for {
			//接收UDP消息并打印
			n, _, _ := conn.ReadFromUDP(buf)
			if n > 0 {
				fmt.Printf("收到消息:%sp2p>", buf[:n])
			}

		}
	}()
	//接下来监控标准输入，发送给对方
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("p2p>")
		//读取标准输入，以换行为读取标志
		data, _ := reader.ReadString('\n')
		conn.Write([]byte(data))

	}
}
