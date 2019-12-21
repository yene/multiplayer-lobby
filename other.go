package main

import (
	"sync"
	"time"
)

type GameState int

const (
	GameError GameState = iota
	WaitingInLobby
	Waiting
	West
)

type Games struct {
	sync.RWMutex
	lobbies map[string]*Lobby
}

type Lobby struct {
	LobbyID      int
	Hashid       string
	CreationDate time.Time
	Players      []int
}

type JoinAnswer struct {
	*Lobby
	YourPlayerID int
}
