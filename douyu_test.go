package danmu

import (
	"github.com/taydy/danmu/douyu"
	"testing"
	"time"
)

func TestGetClient(t *testing.T) {
	roomId := 24422
	client, err := GetClient(roomId)
	if err != nil {
		t.Fatal(err)
	}
	go func(client *douyu.Client) {
		tick := time.Tick(10 * time.Second)
		for {
			select {
			case <- tick:
				_ = client.Close()
			}
		}
	}(client)
	client.Serve()
}

