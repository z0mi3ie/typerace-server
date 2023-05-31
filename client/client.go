package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type Client struct {
	Connection *websocket.Conn
	SessionID  string
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
	log.Printf("SendEvent: %v", e)
	err := wsjson.Write(ctx, c.Connection, e)
	if err != nil {
		log.Println(err)
	}
}

func (c *Client) ReadMessage() any {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	// Got an invalid frame payload when directly reading to EventResponse
	var resp any
	wsjson.Read(ctx, c.Connection, &resp)
	return resp
}

type NewSessionResponse struct {
	EventName string `json:"eventName"`
	SessionID string `json:"sessionId"`
}

type EventRequest struct {
	SessionID string `json:"sessionId"`
	Event     string `json:"event"`
}

func (c *Client) HandleEvent(e any) {

	v, ok := e.(map[string]any)
	if !ok {
		fmt.Println("received response not map[string]any")
		return
	}

	jsonString, err := json.Marshal(e)
	if err != nil {
		log.Println("error marshalling response")
		return
	}

	switch v["eventName"] {
	case "NEW_SESSION":
		log.Println("handling NEW_SESSION_RESPONSE")
		var er NewSessionResponse
		err = json.Unmarshal(jsonString, &er)
		if err != nil {
			log.Println("error unmarshalling event response")
		}
		log.Println("er.EventName: " + er.EventName)
		log.Println("er.SessionID: " + er.SessionID)
		c.SessionID = er.SessionID
	}
}

func main() {
	client := NewClient()
	client.Connect()
	defer client.Connection.Close(websocket.StatusInternalError, "client error, yikes")

	// Read messages from server and display them as they appear
	go func(c *Client) {
		for {
			resp := c.ReadMessage()
			log.Printf("%v", resp)
			c.HandleEvent(resp)
		}
	}(client)

	// Accept input from user
	for {
		var message string
		fmt.Scanln(&message)
		if "message" == "quit" {
			break
		}

		er := EventRequest{
			SessionID: client.SessionID,
			Event:     message,
		}

		client.SendEvent(er)
	}

	// TODO: send client disconnecting message server

	client.Connection.Close(websocket.StatusNormalClosure, "")
}
