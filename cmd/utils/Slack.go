package utils

import (
	"net/http"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"fmt"
	"strings"
	log "github.com/sirupsen/logrus"
)

func SlackNotify(line , webhook string) error {
	jsonData := map[string]string{"text": line, "thread_ts":"1567432492.005400"}
    jsonValue, _ := json.Marshal(jsonData)
    response, err := http.Post(webhook, "application/json", bytes.NewBuffer(jsonValue))
    if err != nil {
		log.Error("Couldn't send Slack msg: ", err)
		return err
    } else {
        data, _ := ioutil.ReadAll(response.Body)
		log.Info("Slack response ",data)
    }
	return nil
}

func SlackHighlightWords(line string, words []string) (res string) {
	if words == nil {
		return line
	}
	for _, word := range words {
		res=strings.ReplaceAll(line,word,fmt.Sprintf("`%s`",word))
	}
	return res
}



