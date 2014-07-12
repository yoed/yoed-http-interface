package client

import (
	"net/http"
	"net/url"
	"log"
	"fmt"
)

type YoedClient interface {
	Handle(username string)
	GetConfig() *YoedClientConfig
}

type YoedClientConfig struct {
	Listen   string `json:"listen"`
	ServerUrl string `json:"serverUrl"`
}

func Run(c YoedClient) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		c.Handle(username)
	})

	server := http.Server{
		Addr:    c.GetConfig().Listen,
		Handler: mux,
	}

	log.Printf("Send server Yo message...")
	resp, err := http.PostForm(c.GetConfig().ServerUrl, url.Values{"callback_url":{"http://"+c.GetConfig().Listen}})

	if err != nil {
		panic(fmt.Sprintf("failed contacting server : %s", err))
	}
	log.Printf("Yoed server answer... %s", resp)
	log.Printf("Listening...")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}