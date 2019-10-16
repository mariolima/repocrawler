package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mariolima/repocrawler/pkg/crawler"
	log "github.com/sirupsen/logrus"
)

// MatchServer Options passed during creation
type MatchServer struct {
	Port     int
	Hostname string
	CertFile string
}

// Setup Sets up the MatchServer with it's listeners
func (ms MatchServer) Setup() {
	go Start()
	http.HandleFunc("/ws", WsPage)
	// panic(http.ListenAndServeTLS(fmt.Sprintf("%s:%d", ms.Hostname, ms.Port), "../../configs/certs/server.crt", "../../configs/certs/server.key", nil))
	// staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("../../web/build/static")))
	// router.PathPrefix("/static/").Handler(staticHandler)
	// http.Handle("/static", staticHandler)
	http.Handle("/", http.FileServer(http.Dir("../../web/build")))
	log.Warnf("Started UI client in http://%s:%d/\n", ms.Hostname, ms.Port)
	panic(http.ListenAndServe(fmt.Sprintf("%s:%d", ms.Hostname, ms.Port), nil))
	// TODO error handling
}

type matchData struct {
	Time  int64         `json:"time"`
	Match crawler.Match `json:"match,omitempty"`
}

// PushState Broadcasts given CrawlerState to all websocket clients
func (ms MatchServer) PushState(state []crawler.TaskState) error {
	mg := Message{
		Event:   STATE,
		Sender:  ms.Hostname,
		Content: state,
	}
	val, err := json.Marshal(mg)
	BroadcastData(val)
	return err
}

// PushMatch Broadcasts given Match to all websocket clients
func (ms MatchServer) PushMatch(match crawler.Match) error {
	mg := Message{
		Event:  MATCH,
		Sender: ms.Hostname,
		Content: matchData{
			Time:  time.Now().Unix(),
			Match: match,
		},
	}
	val, _ := json.Marshal(mg)
	BroadcastData(val)
	return nil
}

type logData struct {
	Time  int64  `json:"time"`
	Level string `json:"level,omitempty"`
	Msg   string `json:"msg,omitempty"`
}

// PushLogEntry Broadcasts given Logrus Entry to all websocket clients
func (ms MatchServer) PushLogEntry(entry log.Entry) error {
	// DebugMsg(entry.Message)
	mg := Message{
		Event:  DEBUG,
		Sender: ms.Hostname,
		Content: logData{
			Time:  entry.Time.Unix(),
			Level: entry.Level.String(),
			Msg:   entry.Message,
		},
	}
	val, _ := json.Marshal(mg)
	BroadcastData(val)
	return nil
}
