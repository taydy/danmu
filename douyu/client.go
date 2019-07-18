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
	conn   net.Conn
	// Turn off heartbeat and barrage receiver
	closed chan struct{}

	// Message processor handler
	HandlerRegister *HandlerRegister

	rLock sync.Mutex
	wLock sync.Mutex
}

// Connect to douyu barrage server
// @Param connStr default "openbarrage.douyutv.com:8601"
func Connect(connStr string, handlerRegister *HandlerRegister) (*Client, error) {
	conn, err := net.Dial("tcp", connStr)
	if err != nil {
		return nil, err
	}

	logrus.Info(fmt.Sprintf("%s connected.", connStr))

	// init client
	client := &Client{
		conn:   conn,
		closed: make(chan struct{}),
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
func (c *Client) Receive() ([]byte, int16, error) {
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
		return buf, int16(code), err
	}
	// exclude ENDING
	return buf[:cl-1], int16(code), nil
}

// Close connnection
func (c *Client) Close() error {
	c.closed <- struct{}{} // heartbeat
	c.closed <- struct{}{} // receive
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
	// Verify that the code is correct
	if code != MsgFromServer {
		logrus.Errorf("Msg code is abnormal, except %d, actual %d", MsgFromServer, code)
		return fmt.Errorf("msg code is abnormal, except %d, actual %d", MsgFromServer, code)
	}
	logrus.Info(fmt.Sprintf("room %d joined", roomId))
	logrus.Info(string(b))
	loginRes := NewMessage(nil, MsgFromServer).Decode(b, code)

	// The field live stat doesn't seem to work at the moment.
	// Whether the anchor is on or off, it is 0.
	logrus.Info(fmt.Sprintf("room %d live status %s", roomId, loginRes.GetStringField("live_stat")))

	joinMessage := NewMessage(nil, MsgToServer).
		SetField("type", "joingroup").
		SetField("rid", roomId).
		SetField("gid", "-9999") // -9999 代表接收所有弹幕消息

	logrus.Info(fmt.Sprintf("joining group %d...", -9999))
	_, err = c.Send(joinMessage.Encode())
	if err != nil {
		return err
	}
	logrus.Info(fmt.Sprintf("group %d joined", -9999))
	return nil
}

// start to get
func (c *Client) Serve() {
loop:
	for {
		select {
		case <-c.closed:
			logrus.Infof("crawler close!")
			break loop
		default:
			b, code, err := c.Receive()
			if err != nil {
				logrus.Error(err)
				break loop
			}

			// decode message
			msg := NewMessage(nil, MsgFromServer).Decode(b, code)

			err, handlers := c.HandlerRegister.Get(msg.GetStringField("type"))
			if err != nil {
				logrus.Debugf(msg.BodyString())
				continue
			}
			for _, v := range handlers {
				go v.Run(msg)
			}
		}
	}
}

// Betta server heartbeat monitoring time is 45 seconds，For insurance use for 40 seconds
func (c *Client) heartbeat() {
	tick := time.Tick(40 * time.Second)
loop:
	for {
		select {
		case <-c.closed:
			logrus.Infof("heart beat close!")
			break loop
		case <-tick:
			heartbeatMsg := NewMessage(nil, MsgToServer).
				SetField("type", MsgTypeKeepAlive).
				SetField("tick", time.Now().Unix())

			_, err := c.Send(heartbeatMsg.Encode())
			if err != nil {
				logrus.Error("heartbeat failed, " + err.Error())
			}
		}
	}
}
