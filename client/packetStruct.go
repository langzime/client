/*
  4byte(packagetype)+4byte(packagelength)+data+4byte(CRC32)
  包格式
*/
package client

import (
	"bufio"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
	//"time"
	"log"
)

const (
	PacketBufSize   = 2048
	MaxPacketLength = 0x00ffffff
)

var (
	ErrorDataTransfer = errors.New("packet transfer error!")
	ErrorDataLength   = errors.New("packet was too long!")
)

var (
	succPakcetType = uint32(1)
	succPacketData = []byte{0xff, 0xff, 0xff, 0xff, 0x01}
)

type Packet struct {
	packetType uint32
	packetData []byte
}

func NewPacket() (packet *Packet) {
	packet = &Packet{0, nil}
	return packet
}

func (p *Packet) SetType(t uint32) {
	p.packetType = t
}

func (p *Packet) GetType() (t uint32) {
	t = p.packetType
	return t
}

func (p *Packet) SetData(data []byte) {
	p.packetData = data
}

func (p *Packet) GetData() (data []byte) {
	data = p.packetData
	return data
}

////////////////////////服务端的响应包//////////////////////////
func SuccessPacket() *Packet {
	Succpacket := NewPacket()
	Succpacket.SetType(succPakcetType)
	Succpacket.SetData(succPacketData)
	return Succpacket
}

///////////////////////////////////////////////////////////////
type PacketReader struct {
	br *bufio.Reader
}

func NewPacketReader(rd io.Reader) *PacketReader {
	r := &PacketReader{}
	r.br = bufio.NewReaderSize(rd, PacketBufSize)
	return r
}

func (p *PacketReader) readHat() (packetType uint32, err error) {
	buf := make([]byte, 4)
	hasRead := int(0)
	for {
		n, err := p.br.Read(buf[hasRead:])
		if err != nil {
			return 0, err
		}
		hasRead += n
		if hasRead >= len(buf) {
			break
		}
	}
	log.Println("数据头",buf)
	packetType = binary.BigEndian.Uint32(buf)

	return packetType, nil

}

func (p *PacketReader) readType() (packetType uint32, err error) {
	buf := make([]byte, 4)
	hasRead := int(0)
	for {
		n, err := p.br.Read(buf[hasRead:])
		if err != nil {
			return 0, err
		}
		hasRead += n
		if hasRead >= len(buf) {
			break
		}
	}
	log.Println("数据类型",buf)
	packetType = binary.BigEndian.Uint32(buf)

	return packetType, nil

}

func (p *PacketReader) readSumCheck() (packetType uint32, err error) {
	buf := make([]byte, 4)
	hasRead := int(0)
	for {
		n, err := p.br.Read(buf[hasRead:])
		if err != nil {
			return 0, err
		}
		hasRead += n
		if hasRead >= len(buf) {
			break
		}
	}
	log.Println("数据校验",buf)
	packetType = binary.BigEndian.Uint32(buf)

	return packetType, nil

}

func (p *PacketReader) readLength() (packetSize uint32, err error) {
	bufLength := make([]byte, 4)
	hasRead := int(0)
	for {
		n, err := p.br.Read(bufLength[hasRead:])
		if err != nil {
			return 0, err
		}
		hasRead += n
		if hasRead >= len(bufLength) {
			break
		}
	}
	log.Println("数据长度",bufLength)
	packetSize = binary.BigEndian.Uint32(bufLength)
	return packetSize, nil

}

func (p *PacketReader) readData(data []byte) error {
	hasRead := uint32(0)
	for {
		n, err := p.br.Read(data[hasRead:])
		if err != nil {
			return err
		}
		hasRead += uint32(n)
		if hasRead >= uint32(len(data)) {
			break
		}
	}
	return nil

}

func (p *PacketReader) checkData(data []byte) error {
	bufCRC := make([]byte, 4)
	hasRead := int(0)
	for {
		n, err := p.br.Read(bufCRC[hasRead:])
		if err != nil {
			return err
		}
		hasRead += n
		if hasRead >= len(bufCRC) {
			break
		}
	}
	log.Println("数据检查",bufCRC)
	bufCRC32 := binary.BigEndian.Uint32(bufCRC)
	srcCRC := crc32.ChecksumIEEE(data)
	if bufCRC32 != srcCRC {
		return ErrorDataTransfer
	}

	return nil

}

//........................读取包...............................

func (p *PacketReader) ReadPacket() (packet *Packet, err error) {
	packet = NewPacket()

	//check packet hat
	_, err = p.readHat()
	if err != nil {
		return nil, err
	}
	//check packet type
	packetType, err := p.readType()
	if err != nil {
		return nil, err
	}
	packet.SetType(packetType)
	packetLength, err := p.readLength()
	if err != nil {
		return nil, err
	}
	log.Printf("%d>%d",packetLength,MaxPacketLength)
	if packetLength > MaxPacketLength {
		return nil, ErrorDataLength
	}
	packetData := make([]byte, packetLength)
	err = p.readData(packetData)
	if err != nil {
		return nil, err
	}
	log.Println("数据体",packetData)
	packet.SetData(packetData)
	//.............................................
	err = p.checkData(packetData)
	if err != nil {
		return nil, ErrorDataTransfer
	}
	//..............................................
	return packet, nil

}

/*newpacket := NewPacket()
//..........................................
log.Println("1111111111", newpacket)
packetType, err := p.readType()
if err != nil {
	return nil, err
}
log.Println("22222222222", newpacket)
newpacket.SetType(packetType)
log.Println("2222255555555552", newpacket)
//...........................................
packetLength, err := p.readLength()
if err != nil {
	return nil, err
}
if packetLength > MaxPacketLength {
	return nil, ErrorDataLength
}
//.............................................
packetData := make([]byte, packetLength)
err = p.readData(packetData)
if err != nil {
	return nil, err
}
newpacket.SetData(packetData)
//.............................................
err = p.checkData(packetData)
if err != nil {
	log.Println("3333333333", newpacket)
	return nil, ErrorDataTransfer
}
//..............................................
log.Println("Read a packet success:", newpacket)
return newpacket, nil*/

/////////////////////////////////////////////////////////////////////
type PacketWriter struct {
	bw *bufio.Writer
}

func NewPacketWriter(wr io.Writer) *PacketWriter {
	w := &PacketWriter{}
	w.bw = bufio.NewWriterSize(wr, PacketBufSize)
	return w
}

//........................写包............................

func (w *PacketWriter) WritePacket(packet *Packet) error {
	//..................写数据头............................
	buf := make([]byte, 4)

	_, err := w.bw.Write([]byte{0xde, 0xad, 0xbe, 0xef})
	if err != nil {
		return err
	}

	//..................写数据类型............................
	binary.BigEndian.PutUint32(buf, packet.GetType())
	_, err = w.bw.Write(buf)
	if err != nil {
		return err
	}

	//...................写数据长度............................
	binary.BigEndian.PutUint32(buf, uint32(len(packet.GetData())))
	_, err = w.bw.Write(buf)
	if err != nil {
		return err
	}

	//....................写数据...............................
	_, err = w.bw.Write(packet.GetData())
	if err != nil {
		return err
	}
	//.....................写CRC32校验..........................
	intCRC := crc32.ChecksumIEEE(packet.GetData())
	binary.BigEndian.PutUint32(buf, intCRC)
	_, err = w.bw.Write(buf)
	if err != nil {
		return err
	}
	log.Printf("剩余缓冲区：%d",w.bw.Available())
	return nil

}

//........................缓冲清除..........................

func (w *PacketWriter) Flush() error {
	return w.bw.Flush()

}

