package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// serveWS upgrades the HTTP connection to a WebSocket, seats the player in the hub,
// and runs the per-client read loop until the connection closes.
func serveWS(hub *Hub, upgrader *websocket.Upgrader, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	slot, slotIndex, err := hub.TryJoin(conn)
	if err != nil {
		conn.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "game is full"),
		)
		conn.Close()
		return
	}

	defer hub.HandleDisconnect(slotIndex)

	// Tell the client which symbol they are.
	waitMsg := "You are Player " + string(slot.Symbol) + ". Waiting for opponent..."
	if hub.BothConnected() {
		waitMsg = "You are Player " + string(slot.Symbol) + ". Game starting!"
	}
	hub.SendTo(slotIndex, MsgPlayerAssigned{
		Type:    "player_assigned",
		Player:  string(slot.Symbol),
		Message: waitMsg,
	})

	// If both players are now connected, notify the waiting player and broadcast the
	// actual game state from the hub (avoids stale state if slot 0 reconnects).
	if hub.BothConnected() {
		hub.SendTo(0, MsgOpponentJoined{
			Type:    "opponent_joined",
			Message: "Opponent connected. Game starting!",
		})
		hub.BroadcastGameState()
	}

	// Read loop — blocks until the connection closes or an error occurs.
	for {
		var msg InboundMsg
		if err := conn.ReadJSON(&msg); err != nil {
			return
		}

		switch msg.Type {
		case "move":
			hub.HandleMove(slotIndex, msg.Row, msg.Col)
		case "restart":
			hub.HandleRestart()
		case "quit":
			hub.HandleQuit(slotIndex)
			return
		default:
			hub.SendTo(slotIndex, MsgError{Type: "error", Message: "unknown message type: " + msg.Type})
		}
	}
}
