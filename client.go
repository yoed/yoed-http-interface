package client

import (
	"net/http"
	"net/url"
	"log"
	"fmt"
)

type YoedClient interface {
	Handle(username string)
	GetConfig() *BaseYoedClientConfig
}

type BaseYoedClientConfig struct {
	Listen   string `json:"listen"`
	ServerUrl string `json:"serverUrl"`
	Handle string `json:"handle"`
}

type BaseYoedClient struct {
	Config *BaseYoedClientConfig
}
func (c *BaseYoedClient) GetConfig() (*BaseYoedClientConfig) {
	return c.Config
}

func Run(c YoedClient) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		c.Handle(username)
	})

	config := c.GetConfig()

	server := http.Server{
		Addr:    config.Listen,
		Handler: mux,
	}

	log.Printf("Send server Yo message...")
	resp, err := http.PostForm(config.ServerUrl+"/yo", url.Values{"handle":{config.Handle}, "callback_url":{"http://"+config.Listen}})

	if err != nil {
		panic(fmt.Sprintf("failed contacting server : %s", err))
	}
	log.Printf("Yoed server answer... %s", resp)
	log.Printf("Listening...")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}