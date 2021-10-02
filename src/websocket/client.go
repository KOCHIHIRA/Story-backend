package websocket

import (
	"fmt"
	"log"
	"server/src/sqls"

	"github.com/gorilla/websocket"
)

type Client struct {
	Name            string //クライアントの名前
	Conn            *websocket.Conn
	WritePermission bool //ユーザーの書き込み権限
}

type Message struct {
	Type      string       `json:"type"`
	Users     []string     `json:"users"`
	RoomName  string       `json:"name"`
	RoomTitle string       `json:"title"`
	Owner     string       `json:"owner"`
	Vote      int          `json:"vote"`
	Storys    []sqls.Story `json:"storys"`
}

func (c *Client) Read(room *Room) {
	defer func() {
		room.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if c.WritePermission {
			c.WritePermission = false
			message := Message{Type: "ADD_STORY", Storys: []sqls.Story{{UserName: c.Name, Sentence: string(p)}}}
			room.Broadcast <- message
			fmt.Printf("Message Received: %+v\n", message)
		}
	}
}
