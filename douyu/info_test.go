package douyu

import "testing"

func TestGetCommonGiftMapping(t *testing.T) {
	giftMappings, err := GetCommonGiftMapping()
	if err != nil {
		t.FailNow()
	}
	for _, v := range giftMappings {
		t.Log(v)
	}
}

func TestGetRoomInfo(t *testing.T) {
	roomId := 9999
	roomInfo, err := GetRoomInfo(roomId)
	if err != nil {
		t.FailNow()
	}
	t.Log(roomInfo)
}
