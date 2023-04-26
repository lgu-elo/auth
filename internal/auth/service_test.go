package auth_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lgu-elo/auth/internal/auth"
	"github.com/lgu-elo/auth/internal/auth/model"
	auth_mock "github.com/lgu-elo/auth/mocks/auth"
	"github.com/lgu-elo/auth/pkg/pb"
	"golang.org/x/crypto/bcrypt"
)

func TestService_IsUserExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := auth_mock.NewMockRepository(ctrl)

	tests := []struct {
		name    string
		want    bool
		payload string
		mock    func()
	}{
		{
			name:    "user exists",
			want:    true,
			payload: "test",
			mock: func() {
				svc.EXPECT().GetUserByUsername(gomock.Any()).Return(nil, nil)
			},
		},
		{
			name:    "user doesn't exists",
			want:    false,
			payload: "test",
			mock: func() {
				svc.EXPECT().GetUserByUsername(gomock.Any()).Return(nil, errors.New("123"))
			},
		},
	}

	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {
			tcase.mock()

			h := auth.NewService(svc)
			result := h.IsUserExist(tcase.payload)

			if result != tcase.want {
				t.Errorf("IsUserExist() got = %v, want = %v", result, tcase.want)
			}
		})
	}
}

func TestService_RegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := auth_mock.NewMockRepository(ctrl)

	tests := []struct {
		name      string
		withError bool
		payload   *pb.BasicCredentials
		mock      func()
	}{
		{
			name:      "success",
			withError: false,
			payload: &pb.BasicCredentials{
				Name:     "name",
				Username: "username",
				Password: "password",
				Role:     "role",
			},
			mock: func() {
				svc.EXPECT().CreateUser(gomock.Any()).Return(nil)
			},
		},
		{
			name:      "failure",
			withError: true,
			payload: &pb.BasicCredentials{
				Name:     "name",
				Username: "username",
				Password: "password",
				Role:     "role",
			},
			mock: func() {
				svc.EXPECT().CreateUser(gomock.Any()).Return(errors.New("123"))
			},
		},
	}

	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {
			tcase.mock()

			h := auth.NewService(svc)
			res, err := h.RegisterUser(tcase.payload)

			if err != nil && !tcase.withError {
				t.Errorf("RegisterUser() got = %v, err = %v", err, tcase.withError)
			}

			if err != nil && tcase.withError {
				return
			}

			if res == "" {
				t.Error("RegisterUser() got empty string, expected jwt")
			}
		})
	}
}

func TestService_ValidateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	passB, _ := bcrypt.GenerateFromPassword([]byte("password"), 14)
	pass := string(passB)
	svc := auth_mock.NewMockRepository(ctrl)

	tests := []struct {
		name      string
		withError bool
		payload   *pb.Credentials
		mock      func()
	}{
		{
			name:      "success",
			withError: false,
			payload: &pb.Credentials{
				Username: "username",
				Password: "password",
			},
			mock: func() {
				svc.EXPECT().GetUserByUsername(gomock.Any()).Return(&model.User{
					ID:       1,
					Login:    "username",
					Password: pass,
					Name:     "name",
					Role:     "role",
				}, nil)
			},
		},
		{
			name:      "wrong password",
			withError: true,
			payload: &pb.Credentials{
				Username: "username",
				Password: "pass",
			},
			mock: func() {
				svc.EXPECT().GetUserByUsername(gomock.Any()).Return(&model.User{
					ID:       1,
					Login:    "username",
					Password: pass,
					Name:     "name",
					Role:     "role",
				}, nil)
			},
		},
		{
			name:      "invalid bcrypted password",
			withError: true,
			payload: &pb.Credentials{
				Username: "username",
				Password: "pass",
			},
			mock: func() {
				svc.EXPECT().GetUserByUsername(gomock.Any()).Return(&model.User{
					ID:       1,
					Login:    "username",
					Password: "pass",
					Name:     "name",
					Role:     "role",
				}, nil)
			},
		},
	}

	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {
			tcase.mock()

			h := auth.NewService(svc)
			res, err := h.ValidateUser(tcase.payload)

			if err != nil && !tcase.withError {
				t.Errorf("ValidateUser() got = %v, err = %v", err, tcase.withError)
			}

			if err != nil && tcase.withError {
				return
			}

			if res == "" {
				t.Error("ValidateUser() got empty string, expected jwt")
			}
		})
	}
}
