package webserver

import (
	"encoding/json"
	"fmt"
	"github.com/mariolima/repocrawl/pkg/crawler"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type MatchServer struct {
	Port     int
	Hostname string
	CertFile string
}

func (ms MatchServer) Setup() {
	go Start()
	http.HandleFunc("/ws", WsPage)
	panic(http.ListenAndServeTLS(fmt.Sprintf("%s:%d", ms.Hostname, ms.Port), "../../configs/certs/server.crt", "../../configs/certs/server.key", nil))
	// TODO error handling
}

type MatchMessage struct {
	Sender    string    `json:"sender,omitempty"`
	Event     Event     `json:"event"`
	Recipient string    `json:"recipient,omitempty"`
	Content   MatchData `json:"data,omitempty"`
	Time      int64     `json:"time"`
}

type MatchData struct {
	Time  int64         `json:"time"`
	Match crawler.Match `json:"match,omitempty"`
}

func (ms MatchServer) PushMatch(match crawler.Match) error {
	mg := MatchMessage{
		Event:  MATCH,
		Sender: ms.Hostname,
		Content: MatchData{
			Time: time.Now().Unix(),
			// Data: match.Line,
			Match: match,
		},
	}
	val, _ := json.Marshal(mg)
	BroadcastData(val)
	return nil
}

func (ms MatchServer) PushLogEntry(entry log.Entry) error {
	DebugMsg(entry.Message)
	// val, _ := json.Marshal(entry)
	// BroadcastData(val)
	return nil
}
