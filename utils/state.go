package utils

import (
	"log/slog"
	"sync"
)

type State string

const (
	STARTED     State = "Started storage data transfer"
	IN_PROGRESS State = "Storage data transfer in progress"
	DONE        State = "Storage data transfer done"
)

type GlobalState struct {
	state State
	mu    sync.RWMutex
}

func NewGlobalState() *GlobalState {
	return &GlobalState{}
}

func (gsm *GlobalState) GetState() State {
	gsm.mu.RLock()
	defer gsm.mu.RUnlock()
	return gsm.state
}

func (gsm *GlobalState) SetState(newState State) {
	gsm.mu.Lock()
	defer gsm.mu.Unlock()
	gsm.state = newState
	slog.Info(
		"Change data transfer state",
		slog.Any("State", gsm.state),
	)
}