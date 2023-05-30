package main

const (
	EventClientReady   = "CLIENT_READY"
	EventGetDictionary = "GET_DICTIONARY"
	EventClientScore   = "CLIENT_SCORE"
	EventGameSummary   = "GAME_SUMMARY"
	EventNewSession    = "NEW_SESSION"
)

type EventRequest struct {
	SessionID string `json:"sessionId"`
	Event     string `json:"event"`
}
