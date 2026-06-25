package ctrader

import (
	"log/slog"
	"net"
)

type Client struct {
	ApplicationId     string
	ApplicationSecret string
	AccountId         int64
	Logger            *slog.Logger

	Conn *net.Conn
}

func (c *Client) Start() error {
	conn, err := net.Dial("tcp", "")
	if err != nil {
		return err
	}
	defer conn.Close()
	c.Conn = &conn
	c.Logger.Info("Connected to ctrader successfully")

	return nil

}

func (c *Client) Stop() error {
	return nil
}

// func Command[T proto.Message, E proto.Message]() (E, error) {
// 	return E, nil
// }
