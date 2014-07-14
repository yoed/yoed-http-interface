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
	"github.com/cenkalti/backoff"
	"time"
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

func connect(config *BaseYoedClientConfig) error {
	log.Printf("Send server Yo message...")
	resp, err := http.PostForm(config.ServerUrl+"/yo", url.Values{"handles":{strings.Join(config.Handles, ",")}, "callback_url":{"http://"+config.Listen}})

	log.Printf("Yoed server answer... %s", resp)
	return err
}

func ping(config *BaseYoedClientConfig) {
	ticker := time.NewTicker(time.Minute*5)

	doPing := func() error {
    	_, err := http.Get(config.ServerUrl)
    	return err
	}
	mustReconnect := false
	for t := range ticker.C {
        fmt.Println("Ping server at", t)

        b := backoff.NewExponentialBackOff()
    	pingBackoffTicker := backoff.NewTicker(b)

    	var err error
    	for _ = range pingBackoffTicker.C {
    	    if err = doPing(); err != nil {
    	        log.Println(err, "will retry...")
    	    	mustReconnect = true
    	        continue
    	    }

    	    break
    	}

    	if err != nil {
    		panic(fmt.Sprintf("Failed contacting server : %s", err))
    	}

    	if mustReconnect {
    		connect(config)
    	}
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

	b := backoff.NewExponentialBackOff()
	ticker := backoff.NewTicker(b)

	var err error
	for _ = range ticker.C {
	    if err = connect(config); err != nil {
	        log.Println(err, "will retry...")
	        continue
	    }

	    break
	}

	if err != nil {
		panic(fmt.Sprintf("Failed contacting server : %s", err))
	}
	
	log.Printf("Listening...")
	go ping(config)
	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}