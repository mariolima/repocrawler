package webserver

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"net/http"
	"time"
)

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	id     string
	socket *websocket.Conn
	send   chan []byte
}

type Event string

const (
	DEBUG    Event = "debug"
	MATCH    Event = "match"
	WARNING  Event = "warning"
	ANNOUNCE Event = "announcement"
)

type ContentData struct {
	Time int64  `json:"time"`
	Data string `json:"msg"`
}

type Message struct {
	Sender    string      `json:"sender,omitempty"`
	Event     Event       `json:"event"`
	Recipient string      `json:"recipient,omitempty"`
	Content   ContentData `json:"data,omitempty"`
	Time      int64       `json:"time"`
}

var manager = ClientManager{
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}

func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.register:
			manager.clients[conn] = true
			//jsonMessage, _ := json.Marshal(&Message{Event: DEBUG,Content: "A new socket has connected."})
			//manager.send(jsonMessage, conn)
			announce_msg("A new socket has connected", conn)
		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				delete(manager.clients, conn)
				//jsonMessage, _ := json.Marshal(&Message{Event: DEBUG, Content: "A socket has disconnected."})
				//manager.send(jsonMessage, conn)
				announce_msg("A socket has disconnected", conn)
			}
		case message := <-manager.broadcast:
			for conn := range manager.clients {
				select {
				case conn.send <- message:
				default:
					close(conn.send)
					delete(manager.clients, conn)
				}
			}
		}
	}
}

func (manager *ClientManager) Send(message []byte, ignore *Client) {
	for conn := range manager.clients {
		if conn != ignore {
			conn.send <- message
		}
	}
}

func (c *Client) read() {
	defer func() {
		manager.unregister <- c
		c.socket.Close()
	}()

	for {
		_, _, err := c.socket.ReadMessage()
		if err != nil {
			manager.unregister <- c
			c.socket.Close()
			break
		}
		//jsonMessage, _ := json.Marshal(&Message{Event: DEBUG, Sender: c.id, Content: "Message was read successfully"})
		//jsonMessage, _ := json.Marshal(&hit{domain: "asd"})
		jsonMessage, _ := json.Marshal(&Message{
			Event:  DEBUG,
			Sender: c.id,
			Time:   time.Now().Unix(),
			Content: ContentData{
				Time: time.Now().Unix(),
				Data: "Message was read successfully",
			},
		})

		manager.broadcast <- jsonMessage
	}
}

func (c *Client) write() {
	defer func() {
		c.socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func WsPage(res http.ResponseWriter, req *http.Request) {
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if error != nil {
		http.NotFound(res, req)
		return
	}
	uu := uuid.NewV4()
	client := &Client{id: uu.String(), socket: conn, send: make(chan []byte)}

	manager.register <- client

	go client.read()
	go client.write()
}

const server string = "RepoCrawl"

func Start() {
	manager.start()
}

func DebugMsg(msg string) {
	fmt.Sprintf("[DEBUG] %s", msg)
	mg := Message{
		Event:  DEBUG,
		Sender: server,
		Time:   time.Now().Unix(),
		Content: ContentData{
			Time: time.Now().Unix(),
			Data: msg,
		},
	}
	val, _ := json.Marshal(mg)
	//logs = append(logs, mg)
	manager.Send(val, nil)
}

func BroadcastData(message []byte) {
	manager.Send(message, nil)
}

func announce_msg(msg string, ignore *Client) {
	mg := Message{
		Event:  ANNOUNCE,
		Sender: server,
		Time:   time.Now().Unix(),
		Content: ContentData{
			Time: time.Now().Unix(),
			Data: msg,
		},
	}
	val, _ := json.Marshal(mg)
	manager.Send(val, ignore)
}
