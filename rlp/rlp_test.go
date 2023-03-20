package rlp

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"
)

type TestRlpStruct struct {
	A      uint
	B      string
	C      []byte
	BigInt *big.Int
}

func TestRlp(t *testing.T) {
	// 将一个整数数组序列化
	arrdata, err := EncodeToBytes([]uint{32, 28})
	fmt.Printf("unuse err:%v\n", err)
	//fmt.Sprintf("data=%s,err=%v", hex.EncodeToString(arrdata), err)
	// 将数组反序列化
	var intarray []uint
	err = DecodeBytes(arrdata, &intarray)
	fmt.Printf("intarray=%v\n", intarray)

	// 将一个布尔变量序列化到一个writer中
	writer := new(bytes.Buffer)
	err = Encode(writer, true)
	//fmt.Sprintf("data=%s,err=%v",hex.EncodeToString(writer.Bytes()),err)

	// 将一个布尔变量反序列化
	var b bool
	err = DecodeBytes(writer.Bytes(), &b)
	fmt.Printf("b=%v\n", b)

	// 将任意一个struct序列化
	//将一个struct序列化到reader中
	_, r, err := EncodeToReader(TestRlpStruct{3, "44", []byte{0x12, 0x32}, big.NewInt(32)})
	var teststruct TestRlpStruct
	err = Decode(r, &teststruct)
	//{A:0x3, B:"44", C:[]uint8{0x12, 0x32}, BigInt:32}
	fmt.Printf("teststruct=%#v\n", teststruct)

}
