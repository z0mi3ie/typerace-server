package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/z0mi3ie/typerace-server/comms"
	"github.com/z0mi3ie/typerace-server/server/dictionary"
	"github.com/z0mi3ie/typerace-server/server/session"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var allMessages string
var connections []*websocket.Conn

var words []string
var dict *dictionary.Dictionary

var sessions map[string]*session.Session

func init() {
	sessions = make(map[string]*session.Session)
}

func GetSession() string {
	return "fakesession"
}

func HandleEvent(e comms.EventRequest) any {
	switch e.Event {
	case comms.EventNewSession:
		log.Println("NewSessionEvent received.")
		s := session.New()
		sessions[s.ID] = s
		r := comms.NewSessionResponse{
			EventName: comms.NewSessionResponseEvent,
			SessionID: sessions[s.ID].ID,
		}
		return r
	case comms.EventGetDictionary:
		log.Println("EventGetDictionary received. Session: " + e.SessionID)
		session := sessions[e.SessionID]
		log.Println("Got session: " + session.ID)
		w := session.Dictionary.Random()
		r := comms.GetDictionaryResponse{
			EventName: comms.GetDictionaryResponseEvent,
			SessionID: GetSession(),
			Words:     []string{w},
		}
		return r
	case comms.EventClientReady:
		log.Println("EventClientReady received.")
		r := comms.ClientReadyResponse{
			EventName: comms.ClientReadyResponseEvent,
			SessionID: GetSession(),
			Ready:     true,
		}
		return r
	case comms.EventClientScore:
		log.Println("EventClientScore received.")
		r := comms.ClientScoreResponse{
			EventName: comms.ClientScoreResponseEvent,
			SessionID: GetSession(),
			Score:     1,
		}
		return r
	case comms.EventGameSummary:
		log.Println("EventGameSummary received.")
		r := comms.GameSummaryResponse{
			SessionID: GetSession(),
		}
		return r
	}

	log.Println("[ERROR] invalid event recieved")
	return comms.ErrorResponse{
		SessionID: GetSession(),
		Error:     "invalid event",
	}
}

func MessageHandler(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "internal server error, yikes")

	connections = append(connections, c)

	// Uncomment this when we are ready to support all the clients :)
	/*
		for _, cc := range connections {
			err = wsjson.Write(context.Background(), cc, "new client connected")
			if err != nil {
				log.Println(err)
			}
		}
	*/

	for {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*60)
		defer cancel()

		//var v any
		var v comms.EventRequest
		err = wsjson.Read(ctx, c, &v)
		if err != nil {
			log.Println(err)
			return
		}

		// parse and handle event
		r := HandleEvent(v)

		log.Println("after handle event")
		// respond with events back to clients
		for _, cc := range connections {
			err = wsjson.Write(ctx, cc, r)
			if err != nil {
				log.Println(err)
			}
		}

		// This close the specific connection that connects
		//if v.(string) == "quit" {
		if v.Event == "quit" {
			c.Close(websocket.StatusNormalClosure, "closed normally")
		}
	}
}

func main() {
	http.ListenAndServe(":8080", http.HandlerFunc(MessageHandler))
}
