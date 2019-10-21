package network

import (
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
)

// 注意结构体成员变量首个字母的大小写
// 如果首字母是小写的，则表示该字段不允许被 外部包 访问
// 如果首字母是大写的，则表示该字段允许被 外部包 访问
type Server struct {
	ipAddr   string                   //网络地址
	function map[string]reflect.Value //函数名与函数实体之间的映射
}

// 函数名首字母大小写问题，同结构体成员变量一样
func CreateServer(ipAddr string) *Server {
	return &Server{ipAddr: ipAddr, function: make(map[string]reflect.Value)}
}

func (s *Server) Run() {
	listener, err := net.Listen("tcp", s.ipAddr)
	if err != nil {
		log.Printf("listen on %s error: %v\n", s.ipAddr, err)
		return
	}
	for {

		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept error: %v\n", err)
			continue
		}

		go func() {
			connector := CreateConnector(conn)
			for {
				// 从客户端接收数据
				// 客户端传输过来的函数名和参数列表均存放在 data 结构体中
				data, err := connector.Receive()
				if err != nil {
					log.Printf("receive error: %v\n", err)
					if err != io.EOF {
						log.Printf("read error: %v\n", err)
					}
					return
				}

				// 根据函数名获取函数
				f, ok := s.function[data.Name]
				if !ok { // 客户端请求的函数不存在
					e := fmt.Sprintf("function %s does not exist!", data.Name)
					log.Println(e)
					if err = connector.Send(Data{Name: data.Name, Err: e}); err != nil {
						log.Printf("transmit error-info from server to client error: %v\n", err)
					}
					continue
				}

				log.Printf("function %s is called\n", data.Name)

				// 提取参数
				inArgs := make([]reflect.Value, len(data.Args))
				for i := range data.Args {
					inArgs[i] = reflect.ValueOf(data.Args[i]) // ??
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
				err = connector.Send(Data{Name: data.Name, Args: outArgs, Err: e})
				if err != nil {
					log.Printf("transmit result from server to client error: %v\n", err)
				}
			}
		}()
	}
}

// 注册服务
func (s *Server) Register(name string, f interface{}) {
	// 如果方法已经注册过了，就直接返回；否则就需要注册
	if _, ok := s.function[name]; ok {
		return
	}
	s.function[name] = reflect.ValueOf(f)
}