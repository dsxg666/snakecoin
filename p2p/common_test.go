package p2p

import (
	"fmt"
	"testing"
)

func TestGetHost(t *testing.T) {
	fmt.Println("主机IPv4地址为：", GetHost())
}
