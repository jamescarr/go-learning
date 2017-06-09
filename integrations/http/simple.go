package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Hook struct {
	ID       int64  `json:"id"`
	State    string `json:"state"`
	Accepted bool   `json:"accepted"`
}

func main() {

	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		hook := Hook{
			Accepted: true,
			State:    "Prcoessing",
			ID:       time.Now().Unix(),
		}
		resp, _ := json.Marshal(hook)
		fmt.Fprintf(w, string(resp))
	})

	s := &http.Server{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go s.ListenAndServe()

	ticker := time.NewTicker(time.Second)
	go func() {
		for _ = range ticker.C {
			resp, _ := http.Get("http://localhost:8080/bar")
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			val := Hook{}
			json.Unmarshal(body, &val)
			log.Println(val)
		}
	}()
	time.Sleep(time.Second * 6)
}
