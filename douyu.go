package danmu

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/taydy/danmu/douyu"
	"time"
)

func Work(roomId int) {
	client, err := douyu.Connect("openbarrage.douyutv.com:8601", nil)
	if err != nil {
		logrus.Error(err)
		return
	}

	_ = client.HandlerRegister.Add(douyu.MsgTypeChatMsg, douyu.Handler(chatmsg), douyu.MsgTypeChatMsg)
	_ = client.HandlerRegister.Add(douyu.MsgTypeDGB, douyu.Handler(liwu), douyu.MsgTypeDGB)
	_ = client.HandlerRegister.Add(douyu.MsgTypeUserEnter, douyu.Handler(userenter), douyu.MsgTypeUserEnter)
	if err := client.JoinRoom(roomId); err != nil {
		logrus.Error(fmt.Sprintf("Join room fail, %s", err.Error()))
		return
	}
	go func(client *douyu.Client) {
		tick := time.Tick(10 * time.Second)
		for {
			select {
			case <- tick:
				client.Close()
			}
		}
	}(client)
	client.Serve()
}

func chatmsg(msg *douyu.Message) {
	rid := msg.GetStringField("rid")
	uid := msg.GetStringField("uid")
	level := msg.GetIntField("level")
	nn := msg.GetStringField("nn")
	txt := msg.GetStringField("txt")
	logrus.Info(fmt.Sprintf("danmu -------> rid(%s) uid(%s) - level(%d) - nickname(%s) >>> content(%s)", rid, uid, level, nn, txt))
}

func liwu(msg *douyu.Message) {
	rid := msg.GetStringField("rid")
	uid := msg.GetStringField("uid")
	level := msg.GetIntField("level")
	nn := msg.GetStringField("nn")
	gfid := msg.GetStringField("gfid")
	gfcnt := msg.GetStringField("gfcnt")
	logrus.Info(fmt.Sprintf("liwu --------> rid(%s) uid(%s) - level(%d) - nickname(%s) - >>> gfid(%s) - gfcnt(%s)", rid, uid, level, nn, gfid, gfcnt))
}

func userenter(msg *douyu.Message) {
	rid := msg.GetStringField("rid")
	uid := msg.GetStringField("uid")
	level := msg.GetIntField("level")
	nn := msg.GetStringField("nn")
	logrus.Info(fmt.Sprintf("uenter --------> rid(%s) uid(%s) - level(%d) - nickname(%s)", rid, uid, level, nn))
}
