package ctrader

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"sync"
	"time"

	"github.com/havvvsen/ctrader/openapi"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	AppClientId     string
	AppClientSecret string
	AccountId       int64
	AccessToken     string
	RefreshToken    string
	Conn            net.Conn
	Logger          *slog.Logger
	HandlerFunc     func(proto.Message) error
	Live            bool
	Wg              *sync.WaitGroup
}

func (c *Client) Start() error {
	c.Wg = &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	var addr string
	if c.Live {
		addr = "live.ctraderapi.com:5035"
	} else {
		addr = "demo.ctraderapi.com:5035"
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	c.Conn = conn
	c.Logger.Info("Connected to ctrader successfully")

	c.Wg.Go(func() {
		err := c.keepAlive()

		if err != nil {
			c.Logger.Error(fmt.Sprintf("Failed to send protoHeartBeartEvent: %s\n", err.Error()))
		}
	})

	c.Wg.Go(func() {
		c.listen()
	})

	return nil
}

func (c *Client) listen() {
	for {
		// lenBuf := make([]uint8, 4)
		var len uint32
		err := binary.Read(c.Conn, binary.BigEndian, &len)

		if err != nil {
			if err == io.EOF {
				c.Logger.Error(fmt.Sprintf("Server closed connection: %s\n", err))
				os.Exit(1)
			} else {
				c.Logger.Error(fmt.Sprintf("An error occured while listening: %s\n", err))
			}

			continue
		}

		// len := binary.BigEndian.Uint32(lenBuf)

		msgRawBuf := make([]uint8, len)
		msg := &openapi.ProtoMessage{}

		_, err = io.ReadFull(c.Conn, msgRawBuf)
		if err != nil {
			c.Logger.Error(fmt.Sprintf("An error occured while listening: %s\n", err))

		}

		err = proto.Unmarshal(msgRawBuf, msg)

		if err != nil {
			c.Logger.Error(fmt.Sprintf("An error occured while unmarshling message: %s\n", err))

		}

	}

}

func (c *Client) Stop() error {

	err := c.Conn.Close()

	if err != nil {
		c.Logger.Error("Failed to stop TCP client")

	}
	return nil
}

func (c *Client) Send(ctx context.Context, protoMsg proto.Message) error {
	fmt.Println("Sending: ", protoMsg)
	msg, err := proto.Marshal(protoMsg)

	if err != nil {
		// c.Logger.Error(fmt.Sprintf("Failed to marshal proto message: %s\n", err))
		return err
	}

	msgLen := uint32(len(msg))
	err = binary.Write(c.Conn, binary.BigEndian, msgLen)

	if err != nil {
		// c.Logger.Error(fmt.Sprintf("Failed to write length prefix to connection: %s\n", err))
		return err
	}

	_, err = c.Conn.Write(msg)

	if err != nil {
		// c.Logger.Error(fmt.Sprintf("Failed to write message to connection: %s\n", err))
		return err
	}

	fmt.Println("Successfully sent")

	return nil

}

func (c *Client) keepAlive(ctx context.Context) error {
	protoHeartBeatEvent := &openapi.ProtoHeartbeatEvent{}
	ticker := time.NewTicker(time.Second * 5)

	for _ = range ticker.C {
		fmt.Println("New tick")

		err := c.Send(context.Background(), protoHeartBeatEvent)

		if err != nil {
			return err
		}
	}
	return nil
}

// func Command[T proto.Message, E proto.Message]() (E, error) {
// 	return E, nil
// }
