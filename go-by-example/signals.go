package main

import (
  "log"
  "os"
  "os/signal"
  "syscall"
)

func main() {
  sigs := make(chan os.Signal, 1)
  done := make(chan bool, 1)

  signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

  go func() {
    sig := <-sigs
    log.Printf("Caught signal")
    log.Println(sig)
    done <- true
  }()

  log.Printf("Awaiting signal")
  <-done
  log.Printf("Exiting!")

}
