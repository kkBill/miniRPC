package network

import (
	"encoding/binary" // binary包实现了简单的数字与字节序列的转换以及变长值的编解码。
	"io"
	"log"
	"net"
)

type Connector struct {
	conn net.Conn
}

func CreateConnector(conn net.Conn) *Connector {
	return &Connector{conn}
}

func (c *Connector) Send(data Data) error {
	bytes, err := encode(data) // encode data into bytes
	if err != nil {
		log.Printf("Send data error: %v\n", err)
		return err
	}
	buf := make([]byte, len(bytes)+4)
	binary.BigEndian.PutUint32(buf[:4], uint32(len(bytes))) // set header(the length of the piece of data)
	copy(buf[4:], bytes)                                  // real data
	_, err = c.conn.Write(buf)                            // writes data to the connection.
	return err
}

func (c *Connector) Receive() (Data, error) {
	// 先读取一条数据的头部（即本条数据的长度）
	header := make([]byte, 4)
	_, err := io.ReadFull(c.conn, header) //??
	if err != nil {
		return Data{}, err
	}
	// 解析头部
	dataLen := binary.BigEndian.Uint32(header)
	dataBuf := make([]byte, dataLen)
	_, err = io.ReadFull(c.conn, dataBuf)
	if err != nil {
		return Data{}, nil
	}
	data, err := decode(dataBuf)
	return data, err
}
