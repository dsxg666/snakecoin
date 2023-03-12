package p2p

import (
	"fmt"
	"net"
	"time"
)

func P2PServer() {
	// 服务器启动侦听
	listener, _ := net.ListenUDP("udp", &net.UDPAddr{Port: 40404})
	defer listener.Close()
	fmt.Println("begin server at ", listener.LocalAddr().String())
	// 定义切片存放2个udp地址
	peers := make([]*net.UDPAddr, 2, 2)
	buf := make([]byte, 256)
	// 接下来从2个UDP消息中获得连接的地址A和B
	n, addr, _ := listener.ReadFromUDP(buf)
	fmt.Printf("read from<%s>:%s\n", addr.String(), buf[:n])
	peers[0] = addr
	n, addr, _ = listener.ReadFromUDP(buf)
	fmt.Printf("read from<%s>:%s\n", addr.String(), buf[:n])
	peers[1] = addr
	fmt.Println("begin nat \n")
	// 将A和B分别介绍给彼此
	listener.WriteToUDP([]byte(peers[0].String()), peers[1])
	listener.WriteToUDP([]byte(peers[1].String()), peers[0])
	// 睡眠10s确保消息发送完成，可以退出历史舞台
	time.Sleep(time.Second * 10)
}
