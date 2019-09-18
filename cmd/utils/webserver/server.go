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

type MatchData struct {
	Time  int64         `json:"time"`
	Match crawler.Match `json:"match,omitempty"`
}

func (ms MatchServer) PushMatch(match crawler.Match) error {
	mg := Message{
		Event:  MATCH,
		Sender: ms.Hostname,
		Content: MatchData{
			Time:  time.Now().Unix(),
			Match: match,
		},
	}
	val, _ := json.Marshal(mg)
	BroadcastData(val)
	return nil
}

type LogData struct {
	Time  int64  `json:"time"`
	Level string `json:"level,omitempty"`
	Msg   string `json:"msg,omitempty"`
}

func (ms MatchServer) PushLogEntry(entry log.Entry) error {
	// DebugMsg(entry.Message)
	mg := Message{
		Event:  DEBUG,
		Sender: ms.Hostname,
		Content: LogData{
			Time:  entry.Time.Unix(),
			Level: entry.Level.String(),
			Msg:   entry.Message,
		},
	}
	val, _ := json.Marshal(mg)
	BroadcastData(val)
	return nil
}
