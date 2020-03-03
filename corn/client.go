package corn

import (
	"errors"
	"log"
	"net"
	"reflect"
)

type Client struct {
	conn net.Conn
}

func NewClient(conn net.Conn) *Client {   // 并发读写会乱序.......
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
		log.Println("rpc-client-Call-----in")
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
		fArgs := make([]interface{}, 0, len(req))
		for i := range req {

			fArgs = append(fArgs, req[i].Interface())
			//fArgs[i] = req[i].Interface()   ???
		}

		log.Println("rpc-client-Call-----ready------")

		// send
		err := clientTransport.Send(ProtoData{
			Name: name,
			Args: fArgs,
		})
		if err != nil {
			// 处理err
			return handleError(err)
		}

		log.Println("rpc-client-Call-----send-done!")
		// recieve
		rsp, err := clientTransport.Receive()

		log.Println("rpc-client-Call-----Receive-done!")

		if err != nil {
			return handleError(err)
		}
		if rsp.Err != "" {
			return handleError(errors.New(rsp.Err))
		}

		log.Println("rpc-client-Call-----send-done!")

		log.Println("------client-------data---------check")
		log.Println("send Args = ", fArgs)
		log.Println("Receive data = ", rsp.Args)
		log.Println("------client-------data---------check------end")

		// 如果不做特殊处理的话.....
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
