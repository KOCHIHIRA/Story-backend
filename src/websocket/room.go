package websocket

import (
	"fmt"
	"server/src/sqls"
)

type Room struct {
	Name       string
	Title      string
	Register   chan *Client
	Unregister chan *Client
	Clients    map[string]*Client
	Broadcast  chan Message
}

func NewRoom(roomName, roomTitle string) *Room {
	return &Room{
		Name:       roomName,
		Title:      roomTitle,
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[string]*Client),
		Broadcast:  make(chan Message),
	}
}

func (room *Room) GetUserList() []string {
	var users []string
	for _, clients := range room.Clients {
		users = append(users, clients.Name)
	}
	return users
}

func (room *Room) Start() {
	for {
		select {
		case client := <-room.Register:
			fmt.Println("Size of Connection Pool: ", len(room.Clients))
			if _, exist := room.Clients[client.Name]; !exist {
				room.Clients[client.Name] = client
				for _, clients := range room.Clients {
					clients.Conn.WriteJSON(Message{Type: "JOIN_USER", Users: []string{client.Name}})
				}
			}

			break
		case client := <-room.Unregister:
			delete(room.Clients, client.Name)
			fmt.Println("Size of Connection Pool: ", len(room.Clients), ": room_out_message")
			for _, clients := range room.Clients {
				clients.Conn.WriteJSON(Message{Type: "OUT_USER", Users: []string{client.Name}})
			}
			break
		case message := <-room.Broadcast:
			fmt.Println("Sending message to all clients in Pool")
			sqls.SetSentence("STORY", room.Name, message.Storys[0].UserName, message.Storys[0].Sentence)
			for _, client := range room.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
