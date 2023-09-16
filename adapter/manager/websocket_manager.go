package manager

import (
	"context"
	"fmt"
	"github.com/olahol/melody"
	"log"
	"regexp"
	"reimagined-chainsaw/gateway"
)

const RegexPatternToIdentifyStockCodeCommand = `\/stock=([a-zA-Z0-9_.]+)`

type WebSocketManager struct {
	ConnID string
	*melody.Melody
	gateway.RabbitMQClient
}

func (m *WebSocketManager) SetHandleMessage() {
	m.HandleMessage(func(session *melody.Session, msg []byte) {
		err := m.BroadcastFilter(msg, func(q *melody.Session) bool {
			return q.Request.URL.Path == session.Request.URL.Path
		})
		if err != nil {
			log.Println(err)
		}

		regex := regexp.MustCompile(RegexPatternToIdentifyStockCodeCommand)
		match := regex.FindStringSubmatch(string(msg))
		hasMatch := len(match) > 1
		if hasMatch {
			err := m.PublishMessage(
				context.Background(),
				fmt.Sprintf("%s:%s", m.ConnID, match[1]))
			if err != nil {
				log.Println(err)
			}
		}
	})
}
