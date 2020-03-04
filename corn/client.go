package corn

import (
	"errors"
	"github.com/yuhaoyuan/RPC_server/config"
	"io"
	"log"
	"net"
	"reflect"
	"time"
)

// Client 客户端
type Client struct {
	conn net.Conn
}

// NewClient 创建客户端
func NewClient(conn net.Conn) *Client { // 并发读写会乱序.......
	return &Client{conn}
}

// Close 关闭
func (t *Client) Close() {
	if t.conn != nil {
		_ = t.conn.Close()
	}
}

// CheckConn 检查连接是否有问题，如果有问题则重新创建连接。
func (t *Client) CheckConn() {
	_ = t.conn.SetReadDeadline(time.Now())
	var one = []byte{}
	_, err := t.conn.Read(one)
	if err != nil {
		log.Println("Net err: ", err)
	}
	if err == io.EOF {
		return
	}
	d := net.Dialer{}
	t.conn, err = d.Dial("tcp", config.BaseConf.Addr)
	if err != nil {
		log.Println("client-dial failed!, err ", err)
	}
	log.Println("make RpcClient-conn")
}

// Call 客户端发送请求的逻辑
func (t *Client) Call(name string, funcPointer interface{}) {
	container := reflect.ValueOf(funcPointer).Elem()

	f := func(req []reflect.Value) []reflect.Value {
		log.Println("rpc-client-Call-----in")
		clientTransport := NewCustomAgreement(t.conn)

		handleError := func(err error) []reflect.Value {
			log.Println("Call-------handleError-------err=", err)
			outArgs := make([]reflect.Value, container.Type().NumOut())
			for i := 0; i < len(outArgs)-1; i++ {
				outArgs[i] = reflect.Zero(container.Type().Out(i))
			}
			outArgs[len(outArgs)-1] = reflect.ValueOf(&err).Elem()
			return outArgs
		}
		// package
		fArgs := make([]interface{}, 0, len(req))
		for i := range req {
			fArgs = append(fArgs, req[i].Interface())
		}
		log.Println("rpc-client-Call-----ready------reqArgs=", fArgs)
		// send
		err := clientTransport.Send(ProtoData{
			Name: name,
			Args: fArgs,
		})
		if err != nil {
			// 处理err
			return handleError(err)
		}
		// recieve
		rsp, err := clientTransport.Receive()
		if err != nil {
			return handleError(err)
		}
		if rsp.Err != "" {
			return handleError(errors.New(rsp.Err))
		}
		log.Println("------client-------data---------check")
		log.Println("send Args = ", fArgs)
		log.Println("Receive data = ", rsp.Args)
		log.Println("------client-------data---------check------end")
		if len(rsp.Args) == 0 {
			rsp.Args = make([]interface{}, container.Type().NumOut())
		}
		argsCount := container.Type().NumOut()
		formatArgs := make([]reflect.Value, argsCount)
		for i := 0; i < argsCount-1; i++ {
			if rsp.Args[i] == nil { // 避免被自动干掉，填充0
				formatArgs[i] = reflect.Zero(container.Type().Out(i))
			} else {
				formatArgs[i] = reflect.ValueOf(rsp.Args[i])
			}
		}
		// handle error
		formatArgs[argsCount-1] = reflect.Zero(container.Type().Out(argsCount - 1))
		return formatArgs
	}
	container.Set(reflect.MakeFunc(container.Type(), f))
}
