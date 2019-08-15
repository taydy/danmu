package douyu

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Gift mapping
type GiftMapping struct {
	Id    string  `json:"id"`   // gift id
	Name  string  `json:"name"` // gift name
	Type  string  `json:"type"` // gift type : 1 (鱼丸礼物); 2 (鱼翅礼物)
	Price float64 `json:"pc"`   // gift price : 鱼翅礼物(元) / 鱼丸礼物(鱼丸)
}

// room info
type RoomInfo struct {
	RoomId      string         `json:"room_id"`     // 房间 ID
	CateId      string         `json:"cate_id"`     // 房间所属分类呢
	CateName    string         `json:"cate_name"`   // 房间所属分类名称
	RoomName    string         `json:"room_name"`   // 房间名称
	RoomStatus  string         `json:"room_status"` // 房间开播状态
	StartTime   string         `json:"start_time"`  // 最近开播时间
	OwnerName   string         `json:"owner_name"`  // 房间所属主播昵称
	Avatar      string         `json:"avatar"`      // 房间所属主播头像地址
	Online      int            `json:"online"`      // 原人气字段，现在与热度值同步，后续很可能会依据情况废除该字段
	HN          int            `json:"hn"`          // 在线热度值
	OwnerWeight string         `json:"owner_name"`  // 直播间主播体重
	FansNum     string         `json:"fans_num"`    // 直播间关注数
	Gift        []*GiftMapping `json:"gift"`        // 直播间礼物信息列表
}

type Noble struct {
	Level           int     `json:"level"`             // id
	NobleName       string  `json:"noble_name"`        // 贵族名称
	IsOnSell        int     `json:"is_on_sell"`        // 是否在售
	FirstOpenPrice  float64 `json:"first_open_price"`  // 首次开通价格
	FirstRemandGold float64 `json:"first_remand_gold"` // 首次开通返还贵族鱼翅数量
	RenewPrice      float64 `json:"renew_price"`       // 续费价格
	RenewRemandGold float64 `json:"renew_remand_gold"` // 续费返还贵族鱼翅数量
}

// Get public gift mapping
func GetCommonGiftMapping() ([]*GiftMapping, error) {
	serviceURL := "https://webconf.douyucdn.cn/resource/common/prop_gift_list/prop_gift_config.json"
	resp, err := http.Get(serviceURL)
	if err != nil {
		logrus.Errorf("get gift mapping error, %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errBody, _ := ioutil.ReadAll(resp.Body)
		logrus.Errorf("get gift mapping error, %v", string(errBody))
		return nil, fmt.Errorf("get gift mapping error, %v", string(errBody))
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	respBody = respBody[17 : len(respBody)-2]
	mappings := make(map[string]interface{})
	if err := json.Unmarshal([]byte(respBody), &mappings); err != nil {
		logrus.Error(err)
		return nil, err
	}
	logrus.Info(mappings)

	giftMappings := make([]*GiftMapping, 0)
	for k, v := range mappings["data"].(map[string]interface{}) {
		mapping := v.(map[string]interface{})
		giftMappings = append(giftMappings, &GiftMapping{
			Id:    k,
			Type:  strconv.FormatFloat(mapping["type"].(float64), 'f', 0, 64),
			Name:  mapping["name"].(string),
			Price: mapping["pc"].(float64),
		})
	}
	return giftMappings, nil
}

// Get live room gift information
func GetRoomInfo(roomId string) (*RoomInfo, error) {
	serviceURL := fmt.Sprintf("http://open.douyucdn.cn/api/RoomApi/room/%s", roomId)
	resp, err := http.Get(serviceURL)
	if err != nil {
		logrus.Errorf("get gift mapping by room %s error, %v", roomId, err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errBody, _ := ioutil.ReadAll(resp.Body)
		logrus.Errorf("get gift mapping by room %s error, %v", roomId, string(errBody))
		return nil, fmt.Errorf("get gift mapping by room %s error, %v", roomId, string(errBody))
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	type RoomInfoResp struct {
		Error int       `json:"error"`
		Data  *RoomInfo `json:"data"`
	}

	roomInfoResp := &RoomInfoResp{}
	if err := json.Unmarshal([]byte(respBody), roomInfoResp); err != nil {
		logrus.Error(err)
		return nil, err
	}

	return roomInfoResp.Data, nil
}

func GetNobleInfo() ([]*Noble, error) {
	serviceURL := "https://www.douyu.com/noble/confignw"
	resp, err := http.Get(serviceURL)
	if err != nil {
		logrus.Errorf("get noble info by error, %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errBody, _ := ioutil.ReadAll(resp.Body)
		logrus.Errorf("get noble info by error, %v", string(errBody))
		return nil, fmt.Errorf("get noble info by error, %v", string(errBody))
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	type NobleInfoResp struct {
		Error int      `json:"error"`
		Data  []*Noble `json:"data"`
	}

	nobleInfoResp := &NobleInfoResp{}
	if err := json.Unmarshal([]byte(respBody), nobleInfoResp); err != nil {
		logrus.Error(err)
		return nil, err
	}

	return nobleInfoResp.Data, nil
}
