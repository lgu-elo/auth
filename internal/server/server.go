package server

import (
	"github.com/lgu-elo/auth/internal/auth"
)

type (
	AuthHandler auth.IHandler

	API struct {
		AuthHandler
	}
)

func NewAPI(auth auth.IHandler) *API {
	return &API{auth}
}
