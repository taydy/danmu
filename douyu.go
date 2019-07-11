package danmu

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"taydy/danmu/douyu"
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
	client.Serve()
}

func chatmsg(msg *douyu.Message) {
	rid := msg.GetStringField("rid")
	uid := msg.GetStringField("uid")
	level := msg.GetStringField("level")
	nn := msg.GetStringField("nn")
	txt := msg.GetStringField("txt")
	logrus.Info(fmt.Sprintf("danmu -------> rid(%s) uid(%s) - level(%s) - nickname(%s) >>> content(%s)", rid, uid, level, nn, txt))
}

func liwu(msg *douyu.Message) {
	rid := msg.GetStringField("rid")
	uid := msg.GetStringField("uid")
	level := msg.GetStringField("level")
	str := msg.GetStringField("str")
	nn := msg.GetStringField("nn")
	gfid := msg.GetStringField("gfid")
	gfcnt := msg.GetStringField("gfcnt")
	logrus.Info(fmt.Sprintf("liwu --------> rid(%s) uid(%s) - level(%s) - nickname(%s) - str(%s) -  >>> gfid(%s) - gfcnt(%s)", rid, uid, level, nn, str, gfid, gfcnt))
}

func userenter(msg *douyu.Message) {
	rid := msg.GetStringField("rid")
	uid := msg.GetStringField("uid")
	level := msg.GetStringField("level")
	str := msg.GetStringField("str")
	nn := msg.GetStringField("nn")
	logrus.Info(fmt.Sprintf("uenter --------> rid(%s) uid(%s) - level(%s) - nickname(%s) - str(%s)", rid, uid, level, nn, str))
}
