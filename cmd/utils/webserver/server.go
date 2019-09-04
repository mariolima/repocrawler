package webserver

import(
	"github.com/mariolima/repocrawl/pkg/crawler"
	"net/http"
	"time"
	"encoding/json"
)

func Serve(matchesChan chan crawler.Match) {
	go Start()
	http.HandleFunc("/ws", WsPage)
	// panic(http.ListenAndServe("me:8090", nil))
	// go func(matchesChan chan crawler.Match){
	// 	for{
	// 		PushMatch(<-matchesChan)
	// 	}
	// }(matchesChan)
	panic(http.ListenAndServeTLS("gobh:8090","/home/msclima/go/src/github.com/mariolima/repocrawl/configs/certs/server.crt", "/home/msclima/go/src/github.com/mariolima/repocrawl/configs/certs/server.key", nil))
}

type MatchMessage struct {
    Sender    string										`json:"sender,omitempty"`
    Event     Event											`json:"event"`
    Recipient string										`json:"recipient,omitempty"`
    Content   MatchData										`json:"data,omitempty"`
	Time      int64											`json:"time"`
}

type MatchData struct{
	Time int64												`json:"time"`
	// Data string												`json:"msg"`
	Match crawler.Match										`json:"match,omitempty"`
}

func PushMatch(match crawler.Match) {
	mg:= MatchMessage{
		Event: MATCH,
		Sender: "RepoCrawl",
		Content: MatchData{
			Time: time.Now().Unix(),
			// Data: match.Line,
			Match: match,
		},
	}
	val, _ :=json.Marshal(mg)
	BroadcastData(val)

	// DebugMsg(match.Line)
}
