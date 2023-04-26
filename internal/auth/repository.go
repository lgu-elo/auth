//go:generate mockgen -source repository.go -destination ./../../mocks/auth/repository.go -package auth_mock
package auth

import (
	"context"

	"github.com/lgu-elo/auth/internal/auth/model"
	"github.com/pkg/errors"
)

type Repository interface {
	GetUserByUsername(username string) (*model.User, error)
	CreateUser(user *model.User) error
}

func (db *storage) GetUserByUsername(username string) (*model.User, error) {
	db.lock.Lock()
	defer db.lock.Unlock()

	var user model.User

	err := db.QueryRow(
		context.Background(),
		"SELECT id,login,password FROM users WHERE login=$1",
		username,
	).Scan(&user.ID, &user.Login, &user.Password)

	db.log.Println(user.Login)
	if err != nil {
		return nil, errors.Wrap(err, "user not found")
	}

	return &user, err
}

func (db *storage) CreateUser(user *model.User) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	row := db.QueryRow(
		context.Background(),
		"INSERT INTO users (login, password, name, role) VALUES ($1,$2,$3,$4) RETURNING id, login, name, role",
		user.Login,
		user.Password,
		user.Name,
		user.Role,
	)
	err := row.Scan(&user.ID, &user.Login, &user.Name, &user.Role)
	if err != nil {
		return errors.Wrap(err, "cannot create user")
	}

	return nil
}
