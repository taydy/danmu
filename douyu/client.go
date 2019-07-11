package douyu

import (
	"encoding/binary"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"sync"
	"time"
)

type Client struct {
	conn            net.Conn
	HandlerRegister *HandlerRegister
	closed          chan struct{}

	rLock sync.Mutex
	wLock sync.Mutex
}

// Connect to douyu barrage server
func Connect(connStr string, handlerRegister *HandlerRegister) (*Client, error) {
	conn, err := net.Dial("tcp", connStr)
	if err != nil {
		return nil, err
	}

	logrus.Info(fmt.Sprintf("%s connected.", connStr))

	// server connected
	client := &Client{
		conn: conn,
	}

	if handlerRegister == nil {
		client.HandlerRegister = CreateHandlerRegister()
	} else {
		client.HandlerRegister = handlerRegister
	}

	go client.heartbeat()
	return client, nil
}

// Send message to server
func (c *Client) Send(b []byte) (int, error) {
	c.wLock.Lock()
	defer c.wLock.Unlock()
	return c.conn.Write(b)
}

// Receive message from server
func (c *Client) Receive() ([]byte, int, error) {
	c.rLock.Lock()
	defer c.rLock.Unlock()
	buf := make([]byte, 512)
	if _, err := io.ReadFull(c.conn, buf[:12]); err != nil {
		return buf, 0, err
	}

	// 12 bytes header
	// 4byte for packet length
	pl := binary.LittleEndian.Uint32(buf[:4])

	// ignore buf[4:8]

	// 2byte for message type
	code := binary.LittleEndian.Uint16(buf[8:10])

	// 1byte for secret
	// 1byte for reserved

	// body content length(include ENDING)
	cl := pl - 8

	if cl > 512 {
		// expand buffer
		buf = make([]byte, cl)
	}
	if _, err := io.ReadFull(c.conn, buf[:cl]); err != nil {
		return buf, int(code), err
	}
	// exclude ENDING
	return buf[:cl-1], int(code), nil
}

// Close connnection
func (c *Client) Close() error {
	c.closed <- struct{}{} // receive
	c.closed <- struct{}{} // heartbeat
	return c.conn.Close()
}

// JoinRoom
func (c *Client) JoinRoom(roomId int) error {
	loginMessage := NewMessage(nil, MsgToServer).
		SetField("type", MsgTypeLoginReq).
		SetField("roomid", roomId)

	logrus.Info(fmt.Sprintf("joining room %d...", roomId))
	if _, err := c.Send(loginMessage.Encode()); err != nil {
		return err
	}

	b, code, err := c.Receive()
	if err != nil {
		return err
	}

	// TODO assert(code == MESSAGE_FROM_SERVER)
	logrus.Info(fmt.Sprintf("room %d joined", roomId))
	loginRes := NewMessage(nil, MsgFromServer).Decode(b, code)
	logrus.Info(fmt.Sprintf("room %d live status %s", roomId, loginRes.GetStringField("live_stat")))

	joinMessage := NewMessage(nil, MsgToServer).
		SetField("type", "joingroup").
		SetField("rid", roomId).
		SetField("gid", "-9999")

	logrus.Info(fmt.Sprintf("joining group %d...", -9999))
	_, err = c.Send(joinMessage.Encode())
	if err != nil {
		return err
	}
	logrus.Info(fmt.Sprintf("group %d joined", -9999))
	return nil
}

func (c *Client) Serve() {
loop:
	for {
		select {
		case <-c.closed:
			break loop
		default:
			b, code, err := c.Receive()
			if err != nil {
				logrus.Error(err)
				break loop
			}

			// analize message
			msg := NewMessage(nil, MsgFromServer).Decode(b, code)
			err, handlers := c.HandlerRegister.Get(msg.GetStringField("type"))
			if err != nil {
				logrus.Debug(err)
				continue
			}
			for _, v := range handlers {
				go v.Run(msg)
			}
		}
	}
}

func (c *Client) heartbeat() {
	tick := time.Tick(45 * time.Second)
loop:
	for {
		select {
		case <-c.closed:
			break loop
		case <-tick:
			heartbeatMsg := NewMessage(nil, MsgToServer).
				SetField("type", "keeplive").
				SetField("tick", time.Now().Unix())

			_, err := c.Send(heartbeatMsg.Encode())
			if err != nil {
				logrus.Error("heartbeat failed, " + err.Error())
			}
		}
	}
}
