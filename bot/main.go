package main

import (
	"context"
	"fmt"
	"github.com/levigross/grequests"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"reimagined-chainsaw/infrastructure/rabbitmq"
	"strings"
	"time"
)

func main() {
	c := rabbitmq.Client{
		URL:                "xpto",
		PublisherQueueName: "stock_price",
		ConsumerQueueName:  "find_stock_price",
	}
	err := c.Connect()
	if err != nil {
		panic(err)
	}
	defer c.Close()

	msgCh, err := c.ConsumeMessages()
	if err != nil {
		panic(err)
	}

	go func(ch <-chan amqp091.Delivery) {
		for {
			select {
			case stockCode := <-ch:
				fmt.Println(stockCode.Body)
				resp, err := grequests.Get(fmt.Sprintf("https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv", stockCode.Body), nil)
				if err != nil {
					log.Fatalln("Unable to make request: ", err)
				}
				lines := strings.Split(resp.String(), "\n")
				values := strings.Split(lines[1], ",")

				fmt.Println(values[1])
				err = c.PublishMessage(context.Background(), values[1])
				if err != nil {
					panic(err)
				}
			default:
				time.Sleep(5 * time.Second)
			}
		}
	}(msgCh)
	time.Sleep(time.Minute * 60)
}
