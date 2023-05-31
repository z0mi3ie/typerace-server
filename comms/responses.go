package comms

const (
	NewSessionResponseEvent    = "NEW_SESSION"
	GetDictionaryResponseEvent = "GET_DICTIONARY"
	ClientReadyResponseEvent   = "CLIENT_READY"
	GameSummaryResponseEvent   = "GAME_SUMMARY"
	ClientScoreResponseEvent   = "CLIENT_SCORE"
	ErrorResponseEvent         = "ERROR"
)

type GetDictionaryResponse struct {
	EventName string   `json:"eventName"`
	SessionID string   `json:"sessionId"`
	Words     []string `json:"words"`
}

type ClientReadyResponse struct {
	EventName string `json:"eventName"`
	SessionID string `json:"sessionId"`
	Name      string `json:"name"`
}

type GameSummaryResponse struct {
	EventName string `json:"eventName"`
	SessionID string `json:"sessionId"`
	Score     string `json:"score"`
}

type ClientScoreResponse struct {
	EventName string `json:"eventName"`
	SessionID string `json:"sessionId"`
	Score     string `json:"score"`
}

type NewSessionResponse struct {
	EventName string `json:"eventName"`
	SessionID string `json:"sessionId"`
}

type ErrorResponse struct {
	EventName string `json:"eventName"`
	Event     string `json:"event"`
	SessionID string `json:"sessionId"`
	Error     string `json:"error"`
}
