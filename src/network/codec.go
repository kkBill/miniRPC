package network

import (
	"bytes"
	"encoding/gob"
	"log"
)

// Data是对 server 与 client 之间传递的数据的抽象
type Data struct {
	Name string        // 函数的名称
	Args []interface{} // request body 或 response body
	Err  string        // 记录在服务端执行函数时发生的错误；如果没有错误，则err为""
}

// 序列化
func encode(data Data) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		log.Printf("encode error: %v\n", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

// 反序列化
func decode(b [] byte) (Data, error) {
	buf := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buf)
	var data Data

	// func (dec *Decoder) Decode(e interface{}) error
	// Decode从输入流读取下一个之并将该值存入e。
	// 如果e是nil，将丢弃该值；否则e必须是可接收该值的类型的指针。
	if err := decoder.Decode(&data); err != nil{
		log.Printf("decode error: %v\n", err)
		return Data{}, err
	}
	return data, nil
}
