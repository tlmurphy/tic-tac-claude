package main

// Outbound message types (server → client).

type MsgPlayerAssigned struct {
	Type    string `json:"type"`
	Player  string `json:"player"`
	Message string `json:"message"`
}

type MsgOpponentJoined struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type MsgGameState struct {
	Type   string       `json:"type"`
	Board  [3][3]string `json:"board"`
	Turn   string       `json:"turn"`
	Status string       `json:"status"`
}

type MsgGameOver struct {
	Type         string   `json:"type"`
	Winner       string   `json:"winner"`
	WinningCells [][2]int `json:"winning_cells"`
	Message      string   `json:"message"`
}

type MsgOpponentLeft struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type MsgError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// Inbound message type (client → server).
// A single envelope covers all client messages since the field set is small.

type InboundMsg struct {
	Type string `json:"type"`
	Row  int    `json:"row"`
	Col  int    `json:"col"`
}
