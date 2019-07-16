package danmu

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/taydy/danmu/douyu"
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
	if err := client.JoinRoom(roomId); err != nil {
		logrus.Error(fmt.Sprintf("Join room fail, %s", err.Error()))
		return nil, err
	}

	// 获取直播间详情
	roomInfo, err := douyu.GetRoomInfo(roomId)
	if err != nil {
		logrus.Errorf("get room %d info error, %v", roomId, err)
	}
	logrus.Infof("room %d info : %+v", roomId, roomInfo)

	return client, nil
}

func chatMsg(msg *douyu.Message) {
	rid := msg.GetStringField("rid")
	uid := msg.GetStringField("uid")
	level := msg.GetIntField("level")
	nn := msg.GetStringField("nn")
	txt := msg.GetStringField("txt")
	logrus.Info(fmt.Sprintf("danmu -------> rid(%s) uid(%s) - level(%d) - nickname(%s) >>> content(%s)", rid, uid, level, nn, txt))
}

func gift(msg *douyu.Message) {
	rid := msg.GetStringField("rid")
	uid := msg.GetStringField("uid")
	level := msg.GetIntField("level")
	nn := msg.GetStringField("nn")
	gfId := msg.GetStringField("gfid")
	gfCnt := msg.GetStringField("gfcnt")
	logrus.Info(fmt.Sprintf("liwu --------> rid(%s) uid(%s) - level(%d) - nickname(%s) - >>> gfid(%s) - gfcnt(%s)", rid, uid, level, nn, gfId, gfCnt))
}

func userEnter(msg *douyu.Message) {
	rid := msg.GetStringField("rid")
	uid := msg.GetStringField("uid")
	level := msg.GetIntField("level")
	nn := msg.GetStringField("nn")
	logrus.Info(fmt.Sprintf("uenter --------> rid(%s) uid(%s) - level(%d) - nickname(%s)", rid, uid, level, nn))
}

func noble(msg *douyu.Message)  {
	rid := msg.GetStringField("rid")
	vn := msg.GetIntField("vn")
	logrus.Info("---------------------------------------")
	logrus.Info(msg.BodyString())
	logrus.Infof(fmt.Sprintf("noble ---------> rid(%s) vn(%d)", rid, vn))
}
