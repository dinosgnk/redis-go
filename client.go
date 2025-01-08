package main

import "net"

type Client struct {
	conn      net.Conn
	cmdArgsCh chan [][]byte
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		conn:      conn,
		cmdArgsCh: make(chan [][]byte),
	}
}

func (c *Client) Send(resp []byte) {
	c.conn.Write(resp)
}
