package corn

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"io"
	"log"
	"net"
)

// ProtoData 消息传输格式
type ProtoData struct {
	Name string        //  函数名
	Args []interface{} //  存放函数所需要的参数
	Err  string
}

// encode 消息序列化
func encode(data ProtoData) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

// decode 消息反序列化
func decode(bs []byte) (ProtoData, error) {
	buf := bytes.NewBuffer(bs)
	decoder := gob.NewDecoder(buf)
	var data ProtoData
	err := decoder.Decode(&data) // 非指针的话-报错：err=gob: attempt to decode into a non-pointer
	if err != nil {
		return ProtoData{}, err
	}
	return data, err
}

// CustomAgreement 简单的协议 header(固定长度=4) + body(变长)
type CustomAgreement struct {
	conn net.Conn
}

// NewCustomAgreement create transport
func NewCustomAgreement(conn net.Conn) *CustomAgreement {
	return &CustomAgreement{conn}
}

// Send 发送消息(客户端和服务端通用)
func (t *CustomAgreement) Send(req ProtoData) error {
	log.Println("SSSSSend-in,req=", req)
	b, err := encode(req)
	log.Println("SSSSSend-encode. b=", b)
	if err != nil {
		log.Println("Send-encode error=", err)
		return err
	}
	buf := make([]byte, 4+len(b))
	binary.BigEndian.PutUint32(buf[:4], uint32(len(b))) // set header  默认使用大端序
	copy(buf[4:], b)                                    // set body
	log.Println("SSSSSend-copy.  buf=", buf)
	_, err = t.conn.Write(buf) // 这个并发写操作看起来并不安全呢    那意味着客户端并发需要控制
	log.Println("SSSSSend-to-other Done!,req=", req)
	if err != nil {
		log.Println("Send-Write error=", err)
	}
	return err
}

// Receive 接收消息(客户端和服务端通用)
func (t *CustomAgreement) Receive() (ProtoData, error) {
	log.Println("RRRReceive-")
	header := make([]byte, 4)
	_, err := io.ReadFull(t.conn, header) // read精准的长度, 读不到的话会阻塞？

	log.Println("RRRReceive-io-read, header=", header)
	if err != nil {
		log.Println("receive error=", err)
		return ProtoData{}, err
	}
	dataLen := binary.BigEndian.Uint32(header)
	data := make([]byte, dataLen)
	_, err = io.ReadFull(t.conn, data) // 应控制并发读
	log.Println("RRRReceive-read data, data=", data)
	if err != nil {
		log.Println("receive error=", err)
		return ProtoData{}, err
	}
	rsp, err := decode(data)
	log.Println("RRRReceive-decode data done!, rsp=", rsp)
	return rsp, err
}
