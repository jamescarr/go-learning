# Kafka Producer / Consumer in Go
This demonstrates using sarama to do consumer and producer in go.

## Running the Example
This uses docker compose and [glide](https://github.com/Masterminds/glide) for dependency management.

* `docker-compose up -d` to bring the broker up
* `glide install` to install dependencies
* in separate terminals, run `go run producer.go` and `go run consumer.go`

## To Do
* [ ] Retry when broker is unavailable.
* [ ] Replay items in the stream
* [ ] Explore different operations

