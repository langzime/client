package client

var (
	INQUEUE chan Packet
	OUTQEUE chan Packet
)

func init(){
	INQUEUE=make(chan Packet,100)
	OUTQEUE=make(chan Packet,100)
}