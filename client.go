package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/cenkalti/backoff"
)

type Handler interface {
	Handle(username, handle string)
}

type Config struct {
	Listen    string   `json:"listen"`
	ServerUrl string   `json:"server_url"`
	Handles   []string `json:"handles"`
}

type Client struct {
	config  *Config
	handler Handler
}

func LoadConfig(configPath string, v interface{}) error {

	configFile, err := os.Open(configPath)

	if err != nil {
		return err
	}

	configJson, err := ioutil.ReadAll(configFile)

	if err != nil {
		return err
	}

	if err := json.Unmarshal(configJson, v); err != nil {
		return err
	}

	return nil
}

func New(handler Handler, config *Config) *Client {
	return &Client{
		config:  config,
		handler: handler,
	}
}

func (c *Client) Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		handle := r.FormValue("handle")
		log.Printf("Yo'ed by %s on handle %s", username, handle)
		c.handler.Handle(username, handle)
	})

	server := http.Server{
		Addr:    c.config.Listen,
		Handler: mux,
	}

	b := backoff.NewExponentialBackOff()
	ticker := backoff.NewTicker(b)

	var err error
	for _ = range ticker.C {
		if err = c.connect(); err != nil {
			log.Println(err, "will retry...")
			continue
		}

		break
	}

	if err != nil {
		panic(fmt.Sprintf("Failed contacting server : %s", err))
	}

	log.Printf("Listening...")

	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}

func (c *Client) connect() error {

	log.Printf("Send server Yo message...")

	resp, err := http.PostForm(c.config.ServerUrl+"/yo", url.Values{
		"handles":      {strings.Join(c.config.Handles, ",")},
		"callback_url": {"http://" + c.config.Listen},
	})

	log.Printf("Yoed server answer... %s", resp)

	return err
}
