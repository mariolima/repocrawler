package webserver

import (
	"encoding/json"
	"github.com/mariolima/repocrawl/pkg/crawler"
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
	panic(http.ListenAndServeTLS("gobh:8090", "/home/msclima/go/src/github.com/mariolima/repocrawl/configs/certs/server.crt", "/home/msclima/go/src/github.com/mariolima/repocrawl/configs/certs/server.key", nil))
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
		Sender: "RepoCrawl",
		Content: MatchData{
			Time: time.Now().Unix(),
			// Data: match.Line,
			Match: match,
		},
	}
	val, _ := json.Marshal(mg)
	BroadcastData(val)

	return nil
	// DebugMsg(match.Line)
}
