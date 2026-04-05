package socket

import "encoding/json"

type WSPayload interface {
	wsPayload()
}

type WSResponse struct {
	Type    string    `json:"type"`
	Payload WSPayload `json:"payload,omitempty"`
}

// all types of Payload that implements WSPayload
// =============================================
type InQueuePayload struct {
	Message string `json:"message"`
}

// type "game_started"
type GameStartedPayload struct {
	GameID string `json:"game_id"`
	Color  string `json:"color"`
}

// type "move"
type MovePayload struct {
	GameID string `json:"game_id"`
	From   string `json:"from"`
	To     string `json:"to"`
}

// type "leave_game"
type LeaveGamePayload struct {
	GameID string `json:"game_id"`
}

// type "win_leave"
type WinCauseLeavePayload struct {
	Message string `json:"message"`
}

// type "win_disconnect"
type WinCauseDisconnectPayload struct {
	Message string `json:"message"`
}

// type "draw"
type DrawPayload struct {
	Message string `json:"message"`
}

// =============================================

func (InQueuePayload) wsPayload()            {}
func (GameStartedPayload) wsPayload()        {}
func (MovePayload) wsPayload()               {}
func (LeaveGamePayload) wsPayload()          {}
func (WinCauseLeavePayload) wsPayload()      {}
func (WinCauseDisconnectPayload) wsPayload() {}
func (DrawPayload) wsPayload()               {}

// use RawMessage to UnMarshal to raw data and then consider of payload
type WSRequest struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
