// mybeeimClient project main.go
package main

import (
	//"bytes"
	"github.com/wyq756543431/client/client"
	//"fmt"
	"log"
	"net"
	"runtime"
//	"sync"
	"fmt"
)

const (
	MAX_Client = 20
)

var (
	//datatype   = []byte{0x00, 0x00, 0x00, 0x02}
	datatype   = uint32(2)
	datalength = []byte{0x00, 0x00, 0x00, 0x08}
	TestPacketHead = []byte{/*hat*/0xde, 0xad, 0xbe, 0xef}
	data       = []byte{'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o','r'}
	//packket = []byte{ 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x08, 0x12, 0x23, 0x34, 0x45, 0x56, 0x67, 0x78, 0x89 }
)

func main() {
	fmt.Println(string(data))
	packet := client.NewPacket()
	packet.SetType(datatype)
	rpacket := client.NewPacket()
	runtime.GOMAXPROCS(runtime.NumCPU())

	var err error
	conn, err := net.Dial("tcp", "127.0.0.1:1114")
	if err != nil {
		log.Println("error connetion:%s", err)
	}
	packetR := client.NewPacketReader(conn)
	packetW := client.NewPacketWriter(conn)
	for{
		var data1=make([]byte,2048)
		fmt.Print("请输入要发送的消息:")
		n, err := fmt.Scan(&data1)
		if err != nil {
			fmt.Println("数据输入异常:", err.Error())
		}
		log.Printf("%d--->%s",n,string(data1))
		packet.SetData(data1)
		if err := packetW.WritePacket(packet); err != nil {
			log.Println("data transfer fail")
			//conn.Close()
		}
		log.Printf("@@@@@@@@@@@\n")
		if err = packetW.Flush(); err != nil {
			log.Panicln(err)
		}
		log.Printf("###############\n")
		rpacket, err = packetR.ReadPacket()
		if err != nil {
			log.Println("read errors: %s", err)
			//conn.Close()
		}
		log.Println("Read SuccessPacket successfully and SuccessPacket:=", rpacket)
	}

}
