package douyu

import "sync"

/**
	聚合数据。
 */

type Aggregate struct {
	gift     *Gift
	msg      *Msg
	audience *Audience

	mu sync.RWMutex
}

// 礼物
type Gift struct {
	GiftAllWorth float64 // 所有礼物价值总和，包括付费礼物、免费礼物和道具等
	GiftWorth    float64 // 付费礼物，主播有收益的礼物，还没去掉分成前的礼物价值
	GiftUserNum  int     // 给主播送礼(包含免费礼物)的人数，当天多次送礼的用户不重复计算
}

// 弹幕
type Msg struct {
	MsgNum     int // 弹幕总数
	MsgUserNum int // 弹幕人数，当天多次发送弹幕的用户不重复计算
}

// 观众
type Audience struct {
	InteractNum int // 活跃观众人数
	NobleNum    int // 贵族人数
}

func (a *Aggregate) AddMsg() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.msg.MsgNum++
}

func (a *Aggregate) AddGift(price float64, isPay bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.gift.GiftAllWorth += price
	if isPay {
		a.gift.GiftWorth += price
	}
}
