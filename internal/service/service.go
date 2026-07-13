package service

import (
	"time"

	"project-manager-go/internal/auth"
)

type Service struct {
	store  Store
	tokens auth.TokenManager
	now    func() time.Time
}

func New(store Store, tokens auth.TokenManager) *Service {
	return &Service{
		store:  store,
		tokens: tokens,
		now:    time.Now,
	}
}
