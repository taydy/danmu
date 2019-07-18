package danmu

import (
	"github.com/sirupsen/logrus"
	"github.com/taydy/danmu/douyu"
	"testing"
	"time"
)

func TestGetClient(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	roomId := 99999
	client, err := GetClient(roomId)
	if err != nil {
		t.Fatal(err)
	}
	go func(client *douyu.Client) {
		tick := time.Tick(500000 * time.Second)
		for {
			select {
			case <- tick:
				_ = client.Close()
			}
		}
	}(client)
	client.Serve()
}

