package miniRPC

import (
	"bytes"
	"encoding/gob"
)

// Data是对 server 与 client 之间传递的数据的抽象
type Data struct {
	name string        // 函数的名称
	args []interface{} // request body 或 response body
	err  string        // 记录在服务端执行函数时发生的错误，如果没有错误，则err为""
}

// 序列化
func encode(data Data) ([]byte, error) {
	var buf bytes.Buffer
	// gob包管理gob流————在编码器（发送器）和解码器（接受器）之间交换的binary值。
	// 一般用于传递远端程序调用（RPC）的参数和结果
	encoder := gob.NewEncoder(&buf)
	// 发生错误
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	// 返回正常结果
	return buf.Bytes(), nil
}

// 反序列化
func decode([] byte) (Data, error) {
	buf := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buf)
	var data Data

	// func (dec *Decoder) Decode(e interface{}) error
	// Decode从输入流读取下一个之并将该值存入e。
	// 如果e是nil，将丢弃该值；否则e必须是可接收该值的类型的指针。
	if err := decoder.Decode(&data); err != nil{
		return Data{}, err
	}
	return data, nil
}
