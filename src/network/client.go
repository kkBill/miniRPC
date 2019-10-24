package network

import (
	"errors"
	"net"
	"reflect"
)

type Client struct {
	conn net.Conn
}

func CreateClient(conn net.Conn) *Client {
	return &Client{conn: conn}
}

// 客户端希望调用远程的 serviceMethod 方法
// 客户端提供 serviceMethod 名称，以及该方法的函数前面 fptr，服务器传回的结果会存放在 reply 中
func (c *Client) Call(serviceMethod string, fptr interface{} /*, reply interface{}*/) {
	//
	fn := reflect.ValueOf(fptr).Elem()

	// 这个函数是要传给 MakeFunc() 的
	f := func(req []reflect.Value) []reflect.Value {
		clientConn := CreateConnector(c.conn)

		// 异常处理函数（这个函数看不懂!）
		errorHandler := func(err error) []reflect.Value {
			// NumOut：返回func类型的返回值个数，如果不是函数，将会panic
			outArgs := make([]reflect.Value, fn.Type().NumOut())
			for i := 0; i < len(outArgs)-1; i++ {
				// Zero()：Zero返回一个持有类型typ的零值的Value
				// Out()：返回func类型的第i个返回值的类型，如非函数或者i不在[0, NumOut())内将会panic
				outArgs[i] = reflect.Zero(fn.Type().Out(i)) //??
			}
			outArgs[len(outArgs)-1] = reflect.ValueOf(&err).Elem()
			return outArgs
		}

		// package request arguments
		inArgs := make([]interface{}, 0, len(req))
		for i := range req {
			inArgs = append(inArgs, req[i].Interface()) // appends elements to the end of a slice.
		}
		// send request to server
		err := clientConn.Send(Data{Name: serviceMethod, Args: inArgs})
		if err != nil { // local network error or encode error
			return errorHandler(err)
		}

		// receive response from server
		response, err := clientConn.Receive()
		if err != nil {
			return errorHandler(err)
		}
		// remote server error
		if response.Err != "" {
			return errorHandler(errors.New(response.Err))
		}

		// ??
		if len(response.Args) == 0 {
			response.Args = make([]interface{}, fn.Type().NumOut())
		}

		// unpackage response arguements
		numOut := fn.Type().NumOut()
		outArgs := make([]reflect.Value, numOut)
		for i := 0; i < numOut; i++ {
			if i != numOut-1 { //
				// if argument is nil (gob will ignore "Zero" in transmission), set "Zero" value
				// ??
				if response.Args[i] == nil {
					outArgs[i] = reflect.Zero(fn.Type().Out(i))
				} else {
					outArgs[i] = reflect.ValueOf(response.Args[i])
				}
			} else { // 处理 error 参数
				outArgs[i] = reflect.Zero(fn.Type().Out(i))
			}
		}
		return outArgs
	}

	fn.Set(reflect.MakeFunc(fn.Type(), f))
}

// 该方法调用serviceMethod服务对应的函数，等待其完成，并返回error状态码
// args 是调用服务时传入的参数，reply 存放服务执行的结果
//func (c *Client) Call2(serviceMethod string, args interface{}, reply interface{}) error {
//
//}