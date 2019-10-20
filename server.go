package miniRPC

import (
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
)

type Server struct {
	ipAddr   string                   //网络地址
	function map[string]reflect.Value //函数名与函数实体之间的映射
}

func createServer(ipAddr string) *Server {
	return &Server{ipAddr: ipAddr, function: make(map[string]reflect.Value)}
}

func (s *Server) Run() {
	//返回在一个本地网络地址laddr上监听的Listener。
	// 网络类型参数net必须是面向流的网络："tcp"、"tcp4"、"tcp6"、"unix"或"unixpacket"。
	listener, err := net.Listen("tcp", s.ipAddr)
	if err != nil {
		log.Printf("listen on %s error: %v\n", s.ipAddr, err)
		return
	}
	for {
		// Accept用于实现Listener接口的Accept方法；
		// 他会等待下一个呼叫，并返回一个该呼叫的Conn接口。
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept error: %v\n", err)
			continue
		}

		go func() {
			connector := createConnector(conn)
			for {
				// 从客户端接收数据
				data, err := connector.Receive()
				if err != nil {
					if err != io.EOF {
						log.Printf("read error: %v\n", err)
					}
					return
				}

				// 根据函数名获取函数
				f, ok := s.function[data.name]
				if !ok { // 客户端请求的函数不存在
					e := fmt.Sprintf("function %s does not exist!", data.name)
					log.Println(e)
					if err = connector.Send(Data{name: data.name, err: e}); err != nil {
						log.Printf("transmit error-info from server to client error: %v\n", err)
					}
					continue
				}

				log.Printf("function %s is called\n", data.name)

				// 提取参数
				inArgs := make([]reflect.Value, len(data.args))
				for i := range data.args {
					inArgs[i] = reflect.ValueOf(data.args[i]) // ??
				}

				// 调用相应的函数，并返回相应的结果
				out := f.Call(inArgs) // Call calls the function v with the input arguments inArgs.
				//
				outArgs := make([]interface{}, len(out)-1)
				for i := 0; i < len(out)-1; i++ {
					outArgs[i] = out[i].Interface()
				}

				// ??
				var e string
				if _, ok := out[len(out)-1].Interface().(error); !ok {
					e = ""
				} else {
					e = out[len(out)-1].Interface().(error).Error()
				}

				// 把结果返回给 客户端
				err = connector.Send(Data{name: data.name, args: outArgs, err: e})
				if err != nil {
					log.Printf("transmit result from server to client error: %v\n", err)
				}
			}
		}()
	}
}

// 根据名字注册方法
// Register 是什么哪里定义的接口？？
func (s *Server) Register(name string, f interface{}) {
	// 如果方法已经注册过了，就直接返回；否则就需要注册
	if _, ok := s.function[name]; ok {
		return
	}
	s.function[name] = reflect.ValueOf(f)
}
