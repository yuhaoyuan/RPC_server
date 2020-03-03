package corn

import (
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
)

type Server struct {
	addr string
	fMap map[string]reflect.Value
}

// Register a new method
func (t *Server) Register(name string, f interface{}) {
	if _, ok := t.fMap[name]; ok {
		fmt.Println("the func already exist")
		return
	}
	t.fMap[name] = reflect.ValueOf(f)
}

func (t *Server) Run() {
	ls, err := net.Listen("tcp", t.addr)
	if err != nil {
		log.Printf("listen failed!")
		return
	}
	connCount := 0
	for {
		conn, err2 := ls.Accept() //会等待下一个呼叫，并返回一个该呼叫的Conn接口。

		if err2 != nil {
			log.Printf("accept failed!")
			continue
		}
		// 在这里监控当前的连接数量
		connCount ++
		log.Println("now conCount = ", connCount)
		go func() {   //  goroutine的数量 等于 长连接的数量，每一个长链接
			transporter := NewCustomAgreement(conn)
			for {
				req, err := transporter.Receive()
				if err != nil {
					if err != io.EOF {
						log.Printf(fmt.Sprintf("Receive failed! err=%s", err))
						return
					}
				}
				// 获得client调用的方法
				f, ok := t.fMap[req.Name]
				if !ok {
					// 没有此方法
					logerr := fmt.Sprintf("the func %s does not exist", req.Name)
					log.Println(logerr)
					err := transporter.Send(ProtoData{
						Name: req.Name,
						Err:  logerr,
					})
					if err != nil {
						log.Printf("transporter---Send failed!")
					}
					continue
				}
				log.Printf("rpc-api-called, name=%s", req.Name)   // 输出日志，接口名字

				// 获得函数需要的参数
				fArgs := make([]reflect.Value, len(req.Args))
				for i := range req.Args {
					fArgs[i] = reflect.ValueOf(req.Args[i])
				}

				// 默认f必然是函数
				funcRsp := f.Call(fArgs)
				// package rsp
				RspInfo := make([]interface{}, len(funcRsp)-1)
				for i := 0; i < len(funcRsp)-1; i++ {
					RspInfo[i] = funcRsp[i].Interface()
				}
				var RspErr string
				rErr, ok := funcRsp[len(funcRsp)-1].Interface().(error)
				if !ok {
					RspErr = ""
				} else {
					RspErr = rErr.Error()
				}

				log.Println("RspInfo= -----------------")
				log.Println(RspInfo)
				log.Println(RspErr)
				log.Println("but req-args=", fArgs)
				log.Println("----------------------")

				// send rsp to client
				err = transporter.Send(ProtoData{
					Name: req.Name,
					Args: RspInfo,
					Err: RspErr,
				})
				if err != nil {
					log.Println(fmt.Sprintf("transporter---Send failed, err = %s", err))
				}
			}
		}()
	}

}

func NewServer(addr string) *Server {
	return &Server{addr, make(map[string]reflect.Value)}
}
