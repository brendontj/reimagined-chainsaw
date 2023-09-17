package server

import (
	"errors"
	"fmt"
	"github.com/olahol/melody"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"reimagined-chainsaw/adapter/manager"
	"reimagined-chainsaw/gateway"
	"reimagined-chainsaw/infrastructure/rabbitmq"
	"reimagined-chainsaw/infrastructure/storage"
	"strings"
	"time"
)

var ErrUserUnauthorized = errors.New("user unauthorized")

type Server struct {
	mux            *http.ServeMux
	melodySessions map[string]*manager.WebSocketManager
	messageCh      chan string
	rabbitmqClient gateway.RabbitMQClient
	storage        gateway.Db
}

func NewServer() *Server {
	rc := rabbitmq.Client{
		URL:                rabbitmq.LocalInstanceRabbitURL,
		PublisherQueueName: rabbitmq.FindStockPriceQueueName,
		ConsumerQueueName:  rabbitmq.StockPriceResultQueueName,
	}

	if err := rc.Connect(); err != nil {
		panic(err)
	}

	return &Server{
		http.NewServeMux(),
		make(map[string]*manager.WebSocketManager),
		make(chan string),
		&rc,
		storage.NewInMemoryStorage(),
	}
}

func (s *Server) WithHandlers() {
	s.mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		urlQueryTokens := strings.Split(r.URL.RawQuery, "=")
		channelName := urlQueryTokens[len(urlQueryTokens)-1]

		s.melodySessions[channelName].HandleRequest(w, r)
	})

	s.mux.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			http.ServeFile(w, r, "html/signin.html")
			return

		case "POST":
			if err := r.ParseForm(); err != nil {
				log.Println(err)
				return
			}

			if err := s.authorizeUser(r.FormValue("user"), r.FormValue("password")); err != nil {
				log.Println(err)
				http.ServeFile(w, r, "html/signin.html")
				return
			}

			http.ServeFile(w, r, "html/index.html")
			return
		}
	})

	s.mux.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			http.ServeFile(w, r, "html/signup.html")
			return

		case "POST":
			if err := r.ParseForm(); err != nil {
				log.Println(err)
				return
			}

			if err := s.storage.AddNewUser(r.FormValue("user"), r.FormValue("password")); err != nil {
				log.Println(err)
				http.ServeFile(w, r, "html/signin.html")
				return
			}

			http.ServeFile(w, r, "html/index.html")
			return
		}
	})

	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "html/signup.html")
			return
		}

		urlPathTokens := strings.Split(r.URL.Path, "/")
		channelName := urlPathTokens[len(urlPathTokens)-1]

		m, exist := s.melodySessions[channelName]
		if !exist {
			s.melodySessions[channelName] = &manager.WebSocketManager{
				ConnID:         channelName,
				Melody:         melody.New(),
				RabbitMQClient: s.rabbitmqClient,
			}
			m = s.melodySessions[channelName]
		}
		m.SetHandleMessage()

		http.ServeFile(w, r, "html/channel.html")
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
				tokens := strings.Split(string(c.Body), ":")
				channelID := tokens[0]
				messageTemplate := "%s quote is $%s per share"
				_ = s.melodySessions[channelID].BroadcastFilter([]byte(fmt.Sprintf(messageTemplate, tokens[1], tokens[2])), func(q *melody.Session) bool {
					return q.Request.URL.Path == q.Request.URL.Path
				})
			default:
				time.Sleep(10 * time.Second)
			}
		}
	}(msgCh)
}

func (s *Server) Start() {
	err := http.ListenAndServe(":8000", s.mux)
	if err != nil {
		defer s.Close()
		panic(err)
	}
}

func (s *Server) Close() {
	s.rabbitmqClient.Close()
}

func (s *Server) authorizeUser(username, pw string) error {
	p, err := s.storage.FindUserPassword(username)
	if err != nil {
		return err
	}

	if p != pw {
		return ErrUserUnauthorized
	}

	return nil
}
