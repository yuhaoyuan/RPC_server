package corn

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"io"
	"log"
	"net"
)

// 自定义数据格式
type ProtoData struct {
	Name string       //  函数名
	Args []interface{} //  存放函数所需要的参数
	Err  string
}

// 序列化
func encode(data ProtoData) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

// 解析
func decode(bs []byte) (ProtoData, error) {
	buf := bytes.NewBuffer(bs)
	decoder := gob.NewDecoder(buf)
	var data ProtoData
	err := decoder.Decode(&data)     // 非指针的话-报错：err=gob: attempt to decode into a non-pointer
	if err != nil {
		return ProtoData{}, err
	}
	return data, err
}

// 自定义协议  ：header(固定长度=4) + body(变长)
type CustomAgreement struct {
	conn net.Conn
}

func NewCustomAgreement(conn net.Conn) *CustomAgreement {
	return &CustomAgreement{conn}
}

func (t *CustomAgreement) Send(req ProtoData) error {
	b, err := encode(req)
	if err != nil {
		return err
	}
	buf := make([]byte, 4+len(b))
	binary.BigEndian.PutUint32(buf[:4], uint32(len(b))) // set header  默认使用大端序
	copy(buf[4:], b)                                      // set body
	_, err = t.conn.Write(buf)
	return err
}

func (t *CustomAgreement) Receive() (ProtoData, error) {
	header := make([]byte, 4)
	_, err := io.ReadFull(t.conn, header) // read精准的长度
	if err != nil {
		return ProtoData{}, err
	}
	dataLen := binary.BigEndian.Uint32(header)
	data := make([]byte, dataLen)
	_, err = io.ReadFull(t.conn, data)
	if err != nil {
		return ProtoData{}, err
	}
	rsp, err := decode(data)

	log.Printf("receive------rsp----------")
	log.Println(rsp)
	log.Printf("--------------\n\n")

	return rsp, err
}
