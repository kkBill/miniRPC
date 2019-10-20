package miniRPC

import (
	"errors"
	"net"
	"reflect"
)

type Client struct {
	conn net.Conn
}

func createClient(conn net.Conn) *Client {
	return &Client{conn: conn}
}

// Call() 接口到底定在哪里？
func (c *Client) Call(name string, fptr interface{}) {
	//
	fn := reflect.ValueOf(fptr).Elem()

	//
	f := func(req []reflect.Value) []reflect.Value {
		clientConn := createConnector(c.conn)

		//
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
		err := clientConn.Send(Data{name: name, args: inArgs})
		if err != nil { // local network error or encode error
			return errorHandler(err)
		}

		// receive response from server
		response, err := clientConn.Receive()
		if err != nil {
			return errorHandler(err)
		}
		// remote server error
		if response.err != "" {
			return errorHandler(errors.New(response.err))
		}

		// ??
		if len(response.args) == 0 {
			response.args = make([]interface{}, fn.Type().NumOut())
		}

		// unpackage response arguements
		numOut := fn.Type().NumOut()
		outArgs := make([]reflect.Value, numOut)
		for i := 0; i < numOut; i++ {
			if i != numOut-1 { //
				// if argument is nil (gob will ignore "Zero" in transmission), set "Zero" value
				// ??
				if response.args[i] == nil{
					outArgs[i] = reflect.Zero(fn.Type().Out(i))
				}else{
					outArgs[i] = reflect.ValueOf(response.args[i])
				}
			} else { // 处理 error 参数
				outArgs[i] = reflect.Zero(fn.Type().Out(i))
			}
		}
		return outArgs
	}

	fn.Set(reflect.MakeFunc(fn.Type(), f))
}
