package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Room struct {
	RoomID      int64     `json:"room_id"`
	AdminUserID int64     `json:"admin_user_id"`
	Description string    `json:"description"`
	CategoryID  int64     `json:"category_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type GroupchatClient struct {
	host    string
	timeout time.Duration
}

type GroupchatClientItf interface {
	GetGroupchatRoom(roomID int64) *Room
}

func NewClient(host string, timeout time.Duration) GroupchatClientItf {
	return &GroupchatClient{
		host:    host,
		timeout: timeout,
	}
}

func (gc *GroupchatClient) GetGroupchatRoom(roomID int64) *Room {
	client := &http.Client{
		Timeout: gc.timeout,
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%v/groupchat/%d", gc.host, roomID), nil)
	if err != nil {
		return nil
	}

	respRaw, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer respRaw.Body.Close()

	if err != nil {
		return nil
	}

	if respRaw.StatusCode != 200 {
		return nil
	}

	respByte, err := ioutil.ReadAll(respRaw.Body)
	if err != nil {
		log.Print(err)
		return nil
	}

	var resp struct {
		Err  string `json:"err"`
		Data Room   `json:"data"`
	}

	err = json.Unmarshal(respByte, &resp)
	if err != nil {
		return nil
	}

	return &resp.Data
}
