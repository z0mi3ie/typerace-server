package main

type GetDictionaryResponse struct {
	SessionID string   `json:"sessionId"`
	Words     []string `json:"words"`
}

type ClientReadyResponse struct {
	SessionID string `json:"sessionId"`
	Name      string `json:"name"`
}

type GameSummaryResponse struct {
	SessionID string `json:"sessionId"`
	Score     string `json:"score"`
}

type ClientScoreResponse struct {
	SessionID string `json:"sessionId"`
	Score     string `json:"score"`
}

type NewSessionResponse struct {
	SessionID string `json:"sessionId"`
}

type ErrorResponse struct {
	SessionID string `json:"sessionId"`
	Error     string `json:"error"`
}
