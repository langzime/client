// mybeeimClient project main.go
package main

import (
	//"bytes"
	"github.com/wyq756543431/client/client"
	"github.com/golang/protobuf/proto"
	"github.com/wyq756543431/client/client/protos"
	//"fmt"
	"log"
	"net"
	"runtime"
	//	"sync"
	"fmt"
	"time"
)

const (
	MAX_Client = 20
)

var (
	//datatype   = []byte{0x00, 0x00, 0x00, 0x02}
	datatype       = uint32(2)
	datalength     = []byte{0x00, 0x00, 0x00, 0x08}
	TestPacketHead = []byte{ /*hat*/ 0xde, 0xad, 0xbe, 0xef}
	data           = []byte{'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r'}
	//packket = []byte{ 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x08, 0x12, 0x23, 0x34, 0x45, 0x56, 0x67, 0x78, 0x89 }
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	PrintInfoList()
	var err error
	conn, err := net.Dial("tcp", "127.0.0.1:1114")
	if err != nil {
		log.Panicf("error connetion:%s", err)
	}

	go ReadRtn(conn)
	go OutPacketProcessor(conn)
	go InPacketProcessor(conn)
	packet := client.NewPacket()
	for {
		var data1 = make([]byte, 2048)
		fmt.Print("请输入你要选择的操作:")
		_, err := fmt.Scan(&data1)
		if err != nil {
			fmt.Println("数据输入异常:", err.Error())
		}
		if string(data1)=="1"{
			packet.SetType(client.PacketType_GetLoginToken)
			loginToken:=&protos.GetLoginToken{}
			loginToken.ClientType=proto.String("pc")
			bytes,err:=proto.Marshal(loginToken)
			if err!=nil{
				panic(err)
			}
			packet.SetData(bytes)
			client.OUTQUEUE<-*packet
		}
	}

}

func OutPacketProcessor(conn net.Conn){
	log.Printf("Daemon Thread for process out message start \n")
	packetW := client.NewPacketWriter(conn)
	for p:=range client.OUTQUEUE{
		if err := packetW.WritePacket(&p); err != nil {
			log.Panicf("data transfer fail\n")
			//conn.Close()
		}
		if err:= packetW.Flush(); err != nil {
			log.Panicln(err)
		}
	}
}

func PrintInfoList(){
	fmt.Println("操作列表")
	fmt.Println("1 获得登陆token")
	fmt.Println("2 登陆请求")
}

func ReadRtn(conn net.Conn) {
	log.Printf("Daemon Thread for read in message start \n")
	for {
		rpacket := client.NewPacket()
		packetR := client.NewPacketReader(conn)
		rpacket, err := packetR.ReadPacket()
		if err != nil {
			log.Println("read errors: %s", err)
			time.Sleep(time.Second*5)
		}
		client.INQUEUE<-*rpacket
		log.Println("Read SuccessPacket successfully and SuccessPacket:=", rpacket)
	}
}


func InPacketProcessor(conn net.Conn){
	log.Printf("Daemon Thread for process in message start \n")
	for p:=range client.INQUEUE{
		if p.GetType()==client.PacketType_GetLoginToken{
			rdata:=&protos.GetLoginTokenRtn{}
			err:=proto.Unmarshal(p.GetData(),rdata)
			if err!=nil{
				log.Panicln(err)
				continue
			}
			log.Printf("%s",*rdata.TokenId)
		}else if p.GetType()==client.PacketType_Login{
			rdata:=&protos.GetLoginTokenRtn{}
			err:=proto.Unmarshal(p.GetData(),rdata)
			if err!=nil{
				log.Panicln(err)
				continue
			}
			log.Printf("%s",*rdata.TokenId)
		}
	}
}
