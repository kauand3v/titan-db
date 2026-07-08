package rpc

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

type Client struct {
	conn  net.Conn
	msgID uint32
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Connect(address string) error {
	var err error
	c.conn, err = net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("não foi possível conectar ao servidor: %w", err)
	}
	return nil
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Client) sendFrame(cmd uint8, payload []byte) error {
	c.msgID++
	header := make([]byte, 9)
	binary.BigEndian.PutUint32(header[0:4], uint32(len(payload)))
	binary.BigEndian.PutUint32(header[4:8], c.msgID)
	header[8] = cmd
	_, err := c.conn.Write(header)
	if err != nil {
		return err
	}
	if len(payload) > 0 {
		_, err = c.conn.Write(payload)
	}
	return err
}

func (c *Client) readResponse() (status uint8, data []byte, err error) {
	respHeader := make([]byte, 9)
	_, err = io.ReadFull(c.conn, respHeader)
	if err != nil {
		return 0, nil, err
	}
	respSize := binary.BigEndian.Uint32(respHeader[0:4])
	status = respHeader[8]
	data = make([]byte, respSize)
	if respSize > 0 {
		_, err = io.ReadFull(c.conn, data)
	}
	return status, data, err
}

func (c *Client) Put(key, value string) (uint8, error) {
	payload := EncodePutRequest(key, value)
	if err := c.sendFrame(CMD_PUT, payload); err != nil {
		return StatusError, err
	}
	status, _, err := c.readResponse()
	return status, err
}

func (c *Client) Get(key string) (string, uint8, error) {
	payload := EncodeGetRequest(key)
	if err := c.sendFrame(CMD_GET, payload); err != nil {
		return "", StatusError, err
	}
	status, data, err := c.readResponse()
	if err != nil {
		return "", status, err
	}
	return string(data), status, nil
}

func (c *Client) Delete(key string) (uint8, error) {
	payload := EncodeDeleteRequest(key)
	if err := c.sendFrame(CMD_DELETE, payload); err != nil {
		return StatusError, err
	}
	status, _, err := c.readResponse()
	return status, err
}

func (c *Client) ThreeSum() (string, error) {
	if err := c.sendFrame(CMD_3SUM, nil); err != nil {
		return "", err
	}
	status, data, err := c.readResponse()
	if err != nil {
		return "", err
	}
	if status != StatusOK {
		return "", fmt.Errorf("servidor retornou status %d: %s", status, string(data))
	}
	return string(data), nil
}