package corn

import (
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
	"runtime/debug"
)

// Server 端
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

// Run Server
func (t *Server) Run() {
	ls, err := net.Listen("tcp", t.addr)
	if err != nil {
		log.Printf("listen failed!")
		return
	}
	connCount := 0
	for {
		conn, err2 := ls.Accept() //  一个rpcClient会创建一个conn

		if err2 != nil {
			log.Printf("accept failed!")
			continue
		}
		// 在这里监控当前的连接数量
		connCount++
		log.Println("now conCount = ", connCount)
		/*
			当rpcClient建立tcp连接之后，
			每一个下游客户端请求连接都会开一个goroutine 去执行请求
		*/
		go func() {
			defer func() {
				log.Println("go routine done!")
				if err := recover(); err != nil {
					log.Println("goroutine error = ", err)
					log.Println("stack-info = ", string(debug.Stack()))
				}
				log.Println("go routine defer done.")
			}()
			transporter := NewCustomAgreement(conn)
			for {
				log.Println("\n\n收到客户端请求，服务端处理起点") // 打点证明conn没有断
				req, err := transporter.Receive()  // 当客户端建立连接后send的时候，这边接收其请求
				if err != nil {
					if err != io.EOF {
						log.Printf(fmt.Sprintf("Receive failed! err=%s", err))
						return
					}
				}
				log.Println("rpc-api-Receive , req=", req) // 当客户端的请求过来时，打点日志
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
				err = transporter.Send(ProtoData{ // 处理完成，将数据发送给客户端
					Name: req.Name,
					Args: RspInfo,
					Err:  RspErr,
				})
				if err != nil {
					log.Println(fmt.Sprintf("transporter---Send failed, err = %s", err))
				}
			}
		}()
	}

}

// NewServer create a server
func NewServer(addr string) *Server {
	return &Server{addr, make(map[string]reflect.Value)}
}
