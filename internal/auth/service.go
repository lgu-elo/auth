//go:generate mockgen -source service.go -destination ./../../mocks/auth/service.go -package auth_mock
package auth

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/lgu-elo/auth/internal/auth/model"
	"github.com/lgu-elo/auth/pkg/jwt"
	"github.com/lgu-elo/auth/pkg/pb"
	"golang.org/x/crypto/bcrypt"
)

type (
	IService interface {
		IsUserExist(username string) bool
		RegisterUser(creds *pb.BasicCredentials) (string, error)
		ValidateUser(creds *pb.Credentials) (string, error)
	}

	service struct {
		repo Repository
	}
)

const hashRound = 14

func NewService(repo Repository) IService {
	return &service{repo}
}

func (s *service) IsUserExist(username string) bool {
	if _, err := s.repo.GetUserByUsername(username); err != nil {
		return false
	}
	return true
}

func (s *service) RegisterUser(creds *pb.BasicCredentials) (string, error) {
	u := model.User{
		Login:    creds.Username,
		Password: creds.Password,
		Name:     creds.Name,
		Role:     creds.Role,
	}

	hashed, errHash := bcrypt.GenerateFromPassword([]byte(u.Password), hashRound)
	if errHash != nil {
		return "", errHash
	}

	u.Password = string(hashed)

	if err := s.repo.CreateUser(&u); err != nil {
		return "", errors.Wrap(err, "cannot create user")
	}

	token, err := jwt.GenJWT(fmt.Sprint(u.ID), u.Name)
	if err != nil {
		return "", errors.Wrap(err, "cannot generate JWT")
	}

	return token, nil
}

func (s *service) ValidateUser(creds *pb.Credentials) (string, error) {
	user, _ := s.repo.GetUserByUsername(creds.Username)

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))

	if err != nil {
		return "", errors.Wrap(err, "wrong password")
	}

	token, err := jwt.GenJWT(fmt.Sprint(user.ID), user.Login)
	if err != nil {
		return "", errors.Wrap(err, "cannot generate JWT")
	}

	return token, nil
}
