package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
)

func createConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms
	return config
}

func main() {
	done := make(chan bool, 1)
	signals := make(chan os.Signal, 1)
	brokerList := []string{"localhost:9092"}
	producer, err := sarama.NewAsyncProducer(brokerList, createConfig())

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}
	go func() {
		for err := range producer.Errors() {
			log.Println("Failed to write access log entry:", err)
		}
	}()

	go func() {
		for {
			time.Sleep(time.Second)
			msg := &sarama.ProducerMessage{
				Topic: "important",
				Value: sarama.StringEncoder("This is my test"),
			}
			select {
			case producer.Input() <- msg:
				log.Printf("message published")
			case err := <-producer.Errors():
				log.Printf("Error thrown", err)
			case <-signals:
				done <- true
			}
		}
	}()

	<-done

}
