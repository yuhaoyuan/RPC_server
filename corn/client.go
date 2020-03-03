package corn

import (
	"log"
	"net"
	"reflect"
)

type Client struct {
	conn net.Conn
}

func NewClient(conn net.Conn) *Client {
	return &Client{conn}
}
func (t *Client) Close() {
	if t.conn != nil {
		_ = t.conn.Close()
	}
}

func (t *Client) Call(name string, funcPointer interface{}) {
	container := reflect.ValueOf(funcPointer).Elem()

	f := func(req []reflect.Value) []reflect.Value {
		clientTransport := NewCustomAgreement(t.conn)

		handleError := func(err error) []reflect.Value {
			outArgs := make([]reflect.Value, container.Type().NumOut())
			for i := 0; i < len(outArgs)-1; i++ {
				outArgs[i] = reflect.Zero(container.Type().Out(i))
			}
			outArgs[len(outArgs)-1] = reflect.ValueOf(&err).Elem()
			return outArgs
		}

		// package
		fArgs := make([]interface{}, len(req))
		for i := range req {
			fArgs[i] = req[i].Interface()
		}

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
			return handleError(err)
		}
		log.Println("------client-------data---------check")
		log.Println("send Args = ", fArgs)
		log.Println("Receive data = ", rsp.Args)
		log.Println("------client-------data---------check------end")


		// if rsp.Args == []
		// handle rsp-Args
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
