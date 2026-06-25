package ctrader

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"

	"github.com/lifeashansen/ctrader"
	"github.com/lifeashansen/ctrader/openapi"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	ApplicationClientId     string
	ApplicationClientSecret string
	AccountId               int64
	AccessToken             string
	RefreshToken            string
	Logger                  *slog.Logger
	Live                    bool
	Conn                    *net.Conn
	HandlerFunc             func(proto.Message)
}

func (c *Client) Start() error {
	var host string

	if c.Live {
		host = "live.ctraderapi.com:5035"
	} else {
		host = "demo.ctraderapi.com:5035"
	}

	conn, err := net.Dial("tcp", host)
	if err != nil {
		return err
	}
	defer conn.Close()

	c.Conn = &conn
	c.Logger.Info("Connected to ctrader successfully")

	var len [4]byte

	for {
		_, err := io.ReadFull(*c.Conn, len[:])

		if err != nil {
			if err == io.EOF {
				c.Logger.Error("Error EOF")

				return nil

			}
			return err
		}
		c.Logger.Info(fmt.Sprintf("Message: %s\n", string(len[:])))

	}

}

func (c *Client) Stop() error {
	return nil
}

func (c *Client) keepAlive() {
	req := &openapi.ProtoHeartbeatEvent{}

}

func (c *ctrader.Client) Send(ctx context.Context, client *ctrader.Client, msg proto.Message) error {
	client.

	return nil
}
