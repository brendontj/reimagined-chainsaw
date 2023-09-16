package main

import (
	"context"
	"fmt"
	"github.com/levigross/grequests"
	"log"
	"reimagined-chainsaw/infrastructure/rabbitmq"
	"strings"
	"time"
)

var StockPricesURL = "https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv"

func main() {
	log.Println("Starting stock prices bot execution...")
	c := rabbitmq.Client{
		URL:                rabbitmq.LocalInstanceRabbitURL,
		PublisherQueueName: rabbitmq.StockPriceResultQueueName,
		ConsumerQueueName:  rabbitmq.FindStockPriceQueueName,
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
	log.Println("Stock prices bot in execution...")
	for {
		select {
		case stockCode := <-msgCh:
			tokens := strings.Split(string(stockCode.Body), ":")
			resp, err := grequests.Get(fmt.Sprintf(StockPricesURL, tokens[1]), nil)
			if err != nil {
				log.Fatalln("Unable to make request: ", err)
			}
			lines := strings.Split(resp.String(), "\n")
			values := strings.Split(lines[1], ",")
			err = c.PublishMessage(context.Background(), fmt.Sprintf("%s:%s:%s", tokens[0], tokens[1], values[6]))
			if err != nil {
				panic(err)
			}
		default:
			time.Sleep(10 * time.Second)
		}
	}
}
