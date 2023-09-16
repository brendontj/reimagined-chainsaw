package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"regexp"
	"reimagined-chainsaw/infrastructure/rabbitmq"
	"time"
)

type Server struct {
	ginEngine      *gin.Engine
	melody         *melody.Melody
	messageCh      chan string
	rabbitmqClient rabbitmq.Client
}

func NewServer() *Server {
	rc := rabbitmq.Client{
		URL:                "xpto",
		PublisherQueueName: "find_stock_price",
		ConsumerQueueName:  "stock_price",
	}

	if err := rc.Connect(); err != nil {
		panic(err)
	}

	return &Server{
		gin.Default(),
		melody.New(),
		make(chan string),
		rc,
	}
}

func (s *Server) WithHandlers() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "index.html")
			return
		}

		http.ServeFile(w, r, "chan.html")
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		s.melody.HandleRequest(w, r)
	})

	//s.ginEngine.GET("/", func(c *gin.Context) {
	//	http.ServeFile(c.Writer, c.Request, "index.html")
	//	return
	//})
	//
	//s.ginEngine.GET("/channel", func(c *gin.Context) {
	//	http.ServeFile(c.Writer, c.Request, "chan.html")
	//	return
	//})
	//
	//s.ginEngine.GET("/ws", func(c *gin.Context) {
	//	err := s.melody.HandleRequest(c.Writer, c.Request)
	//	if err != nil {
	//		panic(err)
	//	}
	//})
}

func (s *Server) SetHandleMessage() {
	s.melody.HandleMessage(func(session *melody.Session, msg []byte) {
		_ = s.melody.BroadcastFilter(msg, func(q *melody.Session) bool {
			return q.Request.URL.Path == session.Request.URL.Path
		})

		msgAsString := string(msg)

		pattern := `\/stock=([a-zA-Z0-9_.]+)`
		regex := regexp.MustCompile(pattern)
		match := regex.FindStringSubmatch(msgAsString)
		if len(match) > 1 {
			err := s.rabbitmqClient.PublishMessage(context.Background(), match[1])
			if err != nil {
				log.Println(err)
			}
		}
	})
}

func (s *Server) RunCallbackFn() {
	msgCh, err := s.rabbitmqClient.ConsumeMessages()
	if err != nil {
		panic(err)
	}

	go func(ch <-chan amqp091.Delivery) {
		for {
			select {
			case c := <-ch:
				_ = s.melody.BroadcastFilter(c.Body, func(q *melody.Session) bool {
					return q.Request.URL.Path == q.Request.URL.Path
				})
			default:
				time.Sleep(5 * time.Second)
			}
		}
	}(msgCh)
}

func (s *Server) Start() {
	//err := s.ginEngine.Run(":5000")
	//if err != nil {
	//	panic(err)
	//}
	http.ListenAndServe(":5000", nil)
}
