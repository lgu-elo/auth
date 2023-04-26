package auth

import (
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type storage struct {
	*pgx.Conn
	log  *logrus.Logger
	lock *sync.Mutex
}

func NewStorage(db *pgx.Conn, logger *logrus.Logger, lock *sync.Mutex) Repository {
	return &storage{db, logger, lock}
}
