package danmu

import (
	"github.com/sirupsen/logrus"
	"github.com/taydy/danmu/douyu"
	"testing"
	"time"
)

func TestGetClient(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	jobChan := make(chan string, 1000)
	var client *douyu.Client
	var err error
	jobChan <- "1300804"
	for {
		select {
		case roomId := <- jobChan:
			client, err = GetClient(roomId)
			if err != nil {
				t.Fatal(err)
			}
			client.HeartBeatErrHandler = func(roomInfo *douyu.RoomInfo) {
				_ = client.Close()
				jobChan <- roomId
			}
			go client.Serve()

			go func(client *douyu.Client) {
				tick := time.Tick(50000 * time.Second)
				for {
					select {
					case <- tick:
						_ = client.Close()
					}
				}
			}(client)
		}

	}

}

