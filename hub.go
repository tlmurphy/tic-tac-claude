package main

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

// PlayerSlot holds a connected player's WebSocket connection and assigned symbol.
type PlayerSlot struct {
	Conn   *websocket.Conn
	Symbol CellValue
}

// Hub manages the two player connections and all game state.
// All exported methods are safe for concurrent use.
type Hub struct {
	mu      sync.Mutex
	players [2]*PlayerSlot
	game    *Game
}

func NewHub() *Hub {
	return &Hub{
		game: NewGame(),
	}
}

// TryJoin seats a new connection into the first available slot.
// Returns the assigned slot, its index (0=X, 1=O), and any error.
// If both slots are occupied, returns an error and the caller should close the connection.
func (h *Hub) TryJoin(conn *websocket.Conn) (*PlayerSlot, int, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for i, slot := range h.players {
		if slot == nil {
			symbol := X
			if i == 1 {
				symbol = O
			}
			h.players[i] = &PlayerSlot{Conn: conn, Symbol: symbol}

			// Start the game as soon as both slots are filled, regardless of which
			// slot just joined (slot 0 may reconnect while slot 1 is still present).
			if h.players[0] != nil && h.players[1] != nil {
				h.game.Status = StatusPlaying
			}

			return h.players[i], i, nil
		}
	}
	return nil, -1, fmt.Errorf("game is full")
}

// BothConnected reports whether both player slots are occupied.
func (h *Hub) BothConnected() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.players[0] != nil && h.players[1] != nil
}

// HandleDisconnect cleans up a player slot, resets game state, and notifies the other player.
func (h *Hub) HandleDisconnect(slotIndex int) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.players[slotIndex] == nil {
		return
	}
	h.players[slotIndex] = nil
	h.game = NewGame()

	other := 1 - slotIndex
	if h.players[other] != nil {
		h.writeJSON(other, MsgOpponentLeft{
			Type:    "opponent_left",
			Message: "Your opponent disconnected. Waiting for a new player...",
		})
	}
}

// HandleMove validates and applies a move from the given slot, then broadcasts the result.
func (h *Hub) HandleMove(slotIndex int, row, col int) {
	h.mu.Lock()
	defer h.mu.Unlock()

	slot := h.players[slotIndex]
	if slot == nil {
		return
	}

	if err := h.game.ApplyMove(slot.Symbol, row, col); err != nil {
		h.writeJSON(slotIndex, MsgError{Type: "error", Message: err.Error()})
		return
	}

	h.broadcastGameState()

	if h.game.Status == StatusWon || h.game.Status == StatusDraw {
		msg := h.buildGameOverMsg()
		h.broadcastLocked(msg)
	}
}

// HandleRestart resets the game and broadcasts the new state to both players.
func (h *Hub) HandleRestart() {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.game.Status == StatusPlaying || h.game.Status == StatusWaiting {
		h.writeJSON(0, MsgError{Type: "error", Message: "cannot restart while game is in progress"})
		return
	}
	if h.players[0] == nil || h.players[1] == nil {
		return
	}

	// Loser goes first next round. On a draw, alternate: whoever went first last goes second.
	nextFirst := O
	if h.game.Winner == O {
		nextFirst = X
	}
	h.game.Reset(nextFirst)
	h.broadcastGameState()
}

// HandleQuit treats a quit as a disconnect: notifies the other player and closes the
// quitting player's connection. The server keeps running for a new player to join.
// The deferred HandleDisconnect in the caller will find the slot already nil and no-op.
func (h *Hub) HandleQuit(slotIndex int) {
	h.mu.Lock()
	defer h.mu.Unlock()

	quitter := h.players[slotIndex]
	if quitter == nil {
		return
	}

	other := 1 - slotIndex
	if h.players[other] != nil {
		h.writeJSON(other, MsgOpponentLeft{
			Type:    "opponent_left",
			Message: "Your opponent quit the game. Waiting for a new player...",
		})
	}

	h.players[slotIndex] = nil
	h.game = NewGame()
	quitter.Conn.Close()
}

// SendTo sends a message to a specific slot index only.
func (h *Hub) SendTo(slotIndex int, msg any) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.writeJSON(slotIndex, msg)
}

// Broadcast sends a message to both connected players.
func (h *Hub) Broadcast(msg any) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.broadcastLocked(msg)
}

// BroadcastGameState sends the current game state to both connected players.
func (h *Hub) BroadcastGameState() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.broadcastGameState()
}

// broadcastLocked sends msg to both slots. Caller must hold h.mu.
func (h *Hub) broadcastLocked(msg any) {
	for i := range h.players {
		h.writeJSON(i, msg)
	}
}

// broadcastGameState builds and broadcasts the current game_state. Caller must hold h.mu.
func (h *Hub) broadcastGameState() {
	var boardStrings [3][3]string
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			boardStrings[r][c] = string(h.game.Board[r][c])
		}
	}
	msg := MsgGameState{
		Type:   "game_state",
		Board:  boardStrings,
		Turn:   string(h.game.Turn),
		Status: string(h.game.Status),
	}
	h.broadcastLocked(msg)
}

func (h *Hub) buildGameOverMsg() MsgGameOver {
	msg := MsgGameOver{
		Type:         "game_over",
		Winner:       string(h.game.Winner),
		WinningCells: h.game.WinningCells,
	}
	switch {
	case h.game.Status == StatusWon && h.game.Winner == X:
		msg.Message = "Player X wins!"
	case h.game.Status == StatusWon && h.game.Winner == O:
		msg.Message = "Player O wins!"
	default:
		msg.Message = "It's a draw!"
	}
	return msg
}

// writeJSON writes a JSON message to the given slot. Caller must hold h.mu.
// If the write fails, the connection is treated as dead.
func (h *Hub) writeJSON(slotIndex int, msg any) {
	if slotIndex < 0 || slotIndex > 1 {
		return
	}
	slot := h.players[slotIndex]
	if slot == nil {
		return
	}
	if err := slot.Conn.WriteJSON(msg); err != nil {
		slot.Conn.Close()
		h.players[slotIndex] = nil
	}
}
