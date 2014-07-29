// mybeeimClient project main.go
package main

import (
	//"bytes"
	"client/client"
	//"fmt"
	"log"
	"net"
	"time"

//	"sync"
)

const (
	MAX_Client = 20
)

var (
	//datatype   = []byte{0x00, 0x00, 0x00, 0x02}
	datatype   = uint32(2)
	datalength = []byte{0x00, 0x00, 0x00, 0x08}
	data       = []byte{0x12, 0x23, 0x34, 0x45, 0x56, 0x67, 0x78, 0x89}
	//packket = []byte{ 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x08, 0x12, 0x23, 0x34, 0x45, 0x56, 0x67, 0x78, 0x89 }
)

func main() {
	packet := client.NewPacket()
	packet.SetType(datatype)
	packet.SetData(data)
	for i := int(0); i < MAX_Client; i++ {
		go func() {
			rpacket := client.NewPacket()
			var err error
			conn, err := net.Dial("tcp", "192.168.0.104:1114")
			if err != nil {
				log.Println("error connetion:%s", err)
			}
			packetR := client.NewPacketReader(conn)
			packetW := client.NewPacketWriter(conn)
			for {
				if err := packetW.WritePacket(packet); err != nil {
					log.Println("data transfer fail")
					conn.Close()
					return
				}
				//log.Println("write packet success", packet)
				if err = packetW.Flush(); err != nil {
				}
				time.Sleep(time.Second)
				rpacket, err = packetR.ReadPacket()
				if err != nil {
					log.Println("read errors: %s", err)
					conn.Close()
					return
				}
				log.Println("Read SuccessPacket successfully and SuccessPacket:=", rpacket)
				/*if rpacket.GetType() != datatype && bytes.Compare(rpacket.GetData(), data) != 0 {
					log.Println("rPacket != packet")
				}*/
				time.Sleep(time.Second * 10)
			}
		}()
	}
	time.Sleep(time.Hour)
	return
}
