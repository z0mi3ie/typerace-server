package session

import (
	"strings"

	"github.com/google/uuid"
	"github.com/z0mi3ie/typerace-server/server/dictionary"
	"nhooyr.io/websocket"
)

type Player struct {
	ID         string
	Name       string
	Score      int
	Connection *websocket.Conn
}

type Session struct {
	ID         string
	Dictionary *dictionary.Dictionary
	Players    map[string]*Player
}

// client: player name entered
// client: session-id
func New() *Session {
	// temporarily load the dictionary here and use the words list directly
	// this will be changed with its own events but isn't necessary right now
	dict := dictionary.New()
	dict.Load()

	s := &Session{
		ID:         strings.Split(uuid.NewString(), "-")[0],
		Dictionary: dict,
	}

	return s
}

func (s *Session) AddPlayer(n string, c *websocket.Conn) {
	p := &Player{
		Name:       n,
		Connection: c,
		Score:      0,
		ID:         strings.Split(uuid.NewString(), "-")[0],
	}
	s.Players[n] = p
}
