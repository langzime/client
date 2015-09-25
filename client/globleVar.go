package client

var (
	INQUEUE chan Packet
	OUTQUEUE chan Packet
)

func init(){
	INQUEUE=make(chan Packet,100)
	OUTQUEUE=make(chan Packet,100)
}