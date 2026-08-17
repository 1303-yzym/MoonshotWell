package service

import "github.com/1303-yzym/MoonshotWell/internal/infrastructure/state"

type Service struct {
}

func New(appState *state.AppState) *Service {
	return &Service{}
}
