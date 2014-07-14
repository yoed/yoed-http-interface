package client

import (
	"net/http"
	"net/url"
	"log"
	"fmt"
	"os"
	"encoding/json"
	"io/ioutil"
	"strings"
)

type YoedClient interface {
	Handle(username string)
	GetConfig() *BaseYoedClientConfig
}

type BaseYoedClientConfig struct {
	Listen   string `json:"listen"`
	ServerUrl string `json:"serverUrl"`
	Handles []string `json:"handles"`
}

type BaseYoedClient struct {
	config *BaseYoedClientConfig
}
func (c *BaseYoedClient) GetConfig() (*BaseYoedClientConfig) {
	return c.config
}
func (c *BaseYoedClient) loadConfig(configPath string) (*BaseYoedClientConfig, error) {
	configJson, err := ReadConfig(configPath)

	if err != nil {
		return nil, err
	}

	config := &BaseYoedClientConfig{}

	if err := json.Unmarshal(configJson, config); err != nil {
		return nil, err
	}

	return config, nil
}

func NewBaseYoedClient() (*BaseYoedClient, error) {
	c := &BaseYoedClient{}
	config, err := c.loadConfig("./config.json")

	if err != nil {
		panic(fmt.Sprintf("failed loading config: %s", err))
	}

	c.config = config

	return c, nil
}

func ReadConfig(configPath string) ([]byte, error) {
	configFile, err := os.Open(configPath)

	if err != nil {
		return nil, err
	}

	configJson, err := ioutil.ReadAll(configFile)

	if err != nil {
		return nil, err
	} else {
		return configJson, nil
	}
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
	resp, err := http.PostForm(config.ServerUrl+"/yo", url.Values{"handles":{strings.Join(config.Handles, ",")}, "callback_url":{"http://"+config.Listen}})

	if err != nil {
		panic(fmt.Sprintf("Failed contacting server : %s", err))
	}
	log.Printf("Yoed server answer... %s", resp)
	log.Printf("Listening...")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}