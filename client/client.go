package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/z0mi3ie/typerace-server/comms"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type Client struct {
	Connection *websocket.Conn
	SessionID  string
	Words      []string
	Ready      bool
	Score      int
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

func (c *Client) SendEvent(e comms.EventRequest) {
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
	var resp any
	wsjson.Read(ctx, c.Connection, &resp)
	return resp
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
	case comms.EventNewSession:
		log.Printf("handling %s", comms.EventNewSession)
		var er comms.NewSessionResponse
		err = json.Unmarshal(jsonString, &er)
		if err != nil {
			log.Println("error unmarshalling event response")
		}
		c.SessionID = er.SessionID
	case comms.EventGetDictionary:
		log.Printf("handling %s", comms.EventGetDictionary)
		var er comms.GetDictionaryResponse
		err = json.Unmarshal(jsonString, &er)
		if err != nil {
			log.Println("error unmarshalling event response")
		}
		log.Println("words from event... " + string(er.Words[0]))
		c.Words = er.Words
	case comms.EventClientReady:
		log.Printf("handling %s", comms.EventClientReady)
		var er comms.ClientReadyResponse
		err = json.Unmarshal(jsonString, &er)
		if err != nil {
			log.Println("error unmarshalling event response")
		}
		c.Ready = true
	case comms.EventClientScore:
		log.Printf("handling %s", comms.EventClientScore)
		var er comms.ClientScoreResponse
		err = json.Unmarshal(jsonString, &er)
		if err != nil {
			log.Println("error unmarshalling event response")
		}
		c.Score = er.Score
	case comms.EventGameSummary:
		log.Printf("handling %s", comms.EventGameSummary)
	}
}

func main() {
	client := NewClient()
	client.Connect()
	defer client.Connection.Close(websocket.StatusInternalError, "client error, yikes")

	go func(c *Client) {
		for {
			resp := c.ReadMessage()
			log.Printf("%v", resp)
			c.HandleEvent(resp)
		}
	}(client)

	// Accept input from user for testing
	for {
		var message string
		fmt.Scanln(&message)
		if "message" == "quit" {
			break
		}

		log.Printf(">> client state >>")
		log.Printf(">> SessionID %s", client.SessionID)
		log.Printf(">> Words %s", client.Words)
		log.Printf(">> Ready %v", client.Ready)
		log.Printf(">> Score %d", client.Score)

		er := comms.EventRequest{
			SessionID: client.SessionID,
			Event:     message,
		}

		client.SendEvent(er)
	}

	// TODO: send client disconnecting message server
	client.Connection.Close(websocket.StatusNormalClosure, "")
}
