package webserver

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

// ClientManager Server Manager used to relay Websocket messages/chan/wait
type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

// Client User that successfully connects to Ws
type Client struct {
	id     string
	socket *websocket.Conn
	send   chan []byte
}

// Event sent to Clients signaling message type
type Event string

const (
	DEBUG    Event = "debug" // DEBUG state
	MATCH    Event = "match"
	STATE    Event = "state"
	WARNING  Event = "warning"
	ANNOUNCE Event = "announcement"
)

// ContentData Message raw data with Timestamp
type ContentData struct {
	Time int64  `json:"time"`
	Data string `json:"msg"`
}

// Message Telemetry sent to Client with arbitrary json `data`
type Message struct {
	Sender    string      `json:"sender,omitempty"`
	Event     Event       `json:"event"`
	Recipient string      `json:"recipient,omitempty"`
	Content   interface{} `json:"data,omitempty"`
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
			announceMsg("A new socket has connected", conn)
		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				delete(manager.clients, conn)
				//jsonMessage, _ := json.Marshal(&Message{Event: DEBUG, Content: "A socket has disconnected."})
				//manager.send(jsonMessage, conn)
				announceMsg("A socket has disconnected", conn)
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

// Send sends Raw data to all clients except the one given
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

// WsPage Websocket page handler
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

// Start starts Ws server
func Start() {
	manager.start()
}

// DebugMsg Broadcast msg of the type DEBUG
func DebugMsg(msg string) {
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

// BroadcastData Broadcasts array of bytes to all Ws clients
func BroadcastData(message []byte) {
	manager.Send(message, nil)
}

// announceMsg Broadcasts Message struct  to all clients except the given Ignored one
func announceMsg(msg string, ignore *Client) {
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
