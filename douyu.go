package danmu

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/taydy/danmu/douyu"
	"time"
)

func GetClient(roomId int) (*douyu.Client, error) {
	client, err := douyu.Connect(douyu.DouYuBarrageAddress, nil)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	_ = client.HandlerRegister.Add(douyu.MsgTypeChatMsg, douyu.Handler(chatMsg), douyu.MsgTypeChatMsg)
	_ = client.HandlerRegister.Add(douyu.MsgTypeDGB, douyu.Handler(gift), douyu.MsgTypeDGB)
	_ = client.HandlerRegister.Add(douyu.MsgTypeUserEnter, douyu.Handler(userEnter), douyu.MsgTypeUserEnter)
	_ = client.HandlerRegister.Add(douyu.MsgTypeNoble, douyu.Handler(noble), douyu.MsgTypeNoble)
	_ = client.HandlerRegister.Add(douyu.MsgTypeFrank, douyu.Handler(frank), douyu.MsgTypeFrank)
	_ = client.HandlerRegister.Add(douyu.MsgTypeRSS, douyu.Handler(rss), douyu.MsgTypeRSS)
	_ = client.HandlerRegister.Add(douyu.MsgTypeRnewbc, douyu.Handler(rnewbc), douyu.MsgTypeRnewbc)
	_ = client.HandlerRegister.Add(douyu.MsgTypeAnbc, douyu.Handler(anbc), douyu.MsgTypeAnbc)
	if err := client.JoinRoom(roomId); err != nil {
		logrus.Error(fmt.Sprintf("Join room fail, %s", err.Error()))
		return nil, err
	}

	return client, nil
}

func chatMsg(roomInfo *douyu.RoomInfo, msg *douyu.Message) {
	rid := msg.GetStringField("rid")
	uid := msg.GetStringField("uid")
	level := msg.GetIntField("level")
	nn := msg.GetStringField("nn")
	txt := msg.GetStringField("txt")
	cst := int64(msg.GetIntField("cst"))
	sendTime := time.Unix(cst/1000, cst%1000*1000000)
	nl := msg.GetIntField("nl")
	logrus.Info(fmt.Sprintf("danmu -------> rid(%s) uid(%s) - level(%d) - nl(%d) - nickname(%s) - sendTime(%s) >>> content(%s)", rid, uid, level, nl, nn, sendTime, txt))
}

func gift(roomInfo *douyu.RoomInfo, msg *douyu.Message) {
	rid := msg.GetStringField("rid")
	uid := msg.GetStringField("uid")
	level := msg.GetIntField("level")
	nn := msg.GetStringField("nn")
	gfId := msg.GetStringField("gfid")
	gfCnt := msg.GetStringField("gfcnt")
	logrus.Info(fmt.Sprintf("liwu --------> rid(%s) uid(%s) - level(%d) - nickname(%s) - >>> gfid(%s) - gfcnt(%s)", rid, uid, level, nn, gfId, gfCnt))
}

func userEnter(roomInfo *douyu.RoomInfo, msg *douyu.Message) {
	rid := msg.GetStringField("rid")
	uid := msg.GetStringField("uid")
	level := msg.GetIntField("level")
	nn := msg.GetStringField("nn")
	logrus.Info(fmt.Sprintf("uenter --------> rid(%s) uid(%s) - level(%d) - nickname(%s)", rid, uid, level, nn))
}

func noble(roomInfo *douyu.RoomInfo, msg *douyu.Message) {
	rid := msg.GetStringField("rid")
	vn := msg.GetIntField("vn")
	logrus.Infof(fmt.Sprintf("noble ---------> rid(%s) vn(%d)", rid, vn))
}

func frank(roomInfo *douyu.RoomInfo, msg *douyu.Message) {
	rid := msg.GetStringField("rid")
	fc := msg.GetIntField("fc")
	logrus.Info("---------------------------------------")
	logrus.Info(msg.BodyString())
	logrus.Infof(fmt.Sprintf("frank ---------> rid(%s) fc(%d)", rid, fc))
}

func rss(roomInfo *douyu.RoomInfo, msg *douyu.Message) {
	rid := msg.GetStringField("rid")
	ss := msg.GetIntField("ss")
	rt := msg.GetStringField("rt")
	endTime := msg.GetIntField("endtime")
	logrus.Infof(fmt.Sprintf("rss ---------> rid(%s) ss(%d) rt(%s) endtime(%d)", rid, ss, rt, endTime))
}

func rnewbc(roomInfo *douyu.RoomInfo, msg *douyu.Message) {
	uid := msg.GetStringField("uid")
	unk := msg.GetStringField("unk")
	drid := msg.GetStringField("drid")
	donk := msg.GetStringField("donk")
	rid := msg.GetStringField("rid")
	nl := msg.GetIntField("nl")

	// 不是本直播间的贵族消息，不执行后续逻辑
	//if roomInfo.RoomId != drid && roomInfo.RoomId != rid {
	//	return
	//}
	logrus.Info(msg.BodyString())
	logrus.Infof(fmt.Sprintf("rnewbc ---------> uid(%s) unk(%s) drid(%s) rid(%s) donk(%s) nl(%d)", uid, unk, drid, rid, donk, nl))
}

func anbc(roomInfo *douyu.RoomInfo, msg *douyu.Message) {
	uid := msg.GetStringField("uid")
	unk := msg.GetStringField("unk")
	drid := msg.GetStringField("drid")
	donk := msg.GetStringField("donk")
	rid := msg.GetStringField("rid")
	nl := msg.GetIntField("nl")

	// 不是本直播间的贵族消息，不执行后续逻辑
	//if roomInfo.RoomId != drid && roomInfo.RoomId != rid {
	//	return
	//}
	logrus.Info(msg.BodyString())
	logrus.Infof(fmt.Sprintf("anbc ---------> uid(%s) unk(%s) drid(%s) rid(%s) donk(%s) nl(%d)", uid, unk, drid, rid, donk, nl))
}
