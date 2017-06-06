package main

import (
  "log"
  "os"
  "os/signal"
  "syscall"

  "github.com/Shopify/sarama"
)

func main() {
  topic := "important"
  done := make(chan bool, 1)
  sigs := make(chan os.Signal, 1)
  msgCount := 0
  brokerList := []string{"localhost:9092"}
  config := sarama.NewConfig()
  config.Consumer.Return.Errors = true

  master, err := sarama.NewConsumer(brokerList, config)
  signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
  if err != nil {
    panic(err)
  }

  defer func() {
    if err := master.Close(); err != nil {
      panic(err)
    }
  }()
  consumer, err := master.ConsumePartition(topic, 0, sarama.OffsetOldest)
  if err != nil {
    panic(err)
  }
  go func() {
    for {
      select {
      case err := <-consumer.Errors():
        log.Println(err)
      case msg := <-consumer.Messages():
        msgCount++
        log.Println("Received messages", string(msg.Key), string(msg.Value))
      case <-sigs:
        log.Println("Interrupt is detected")
        done <- true
      }
    }
  }()

  <-done
  log.Println("Processed", msgCount, "messages")
}
