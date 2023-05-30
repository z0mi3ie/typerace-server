package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type Client struct {
	Connection *websocket.Conn
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Connect() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	conn, _, err := websocket.Dial(ctx, "ws://localhost:8080", nil)
	if err != nil {
		log.Println(err)
	}
	c.Connection = conn
}

func (c *Client) SendEvent(e EventRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err := wsjson.Write(ctx, c.Connection, e)
	if err != nil {
		log.Println(err)
	}
}

func (c *Client) ReadMessage() any {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	var resp any
	wsjson.Read(ctx, c.Connection, &resp)
	return resp
}

type EventRequest struct {
	SessionID string `json:"sessionId"`
	Event     string `json:"event"`
}

func main() {
	client := NewClient()
	client.Connect()
	defer client.Connection.Close(websocket.StatusInternalError, "client error, yikes")

	// Read messages from server and display them as they appear
	go func() {
		for {
			resp := client.ReadMessage()
			log.Printf("%v", resp)
		}
	}()

	// Accept input from user
	for {
		var message string
		fmt.Scanln(&message)
		if "message" == "quit" {
			break
		}

		er := EventRequest{
			SessionID: "fakesession",
			Event:     message,
		}

		client.SendEvent(er)
	}

	// TODO: send client disconnecting message server

	client.Connection.Close(websocket.StatusNormalClosure, "")
}
