package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lgu-elo/auth/internal/auth"
	auth_mock "github.com/lgu-elo/auth/mocks/auth"
	"github.com/lgu-elo/auth/pkg/pb"
)

func TestHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := auth_mock.NewMockIService(ctrl)

	tests := []struct {
		name      string
		withError bool
		want      *pb.Token
		payload   *pb.Credentials
		mock      func()
	}{
		{
			name:      "success login",
			withError: false,
			want: &pb.Token{
				Token: "t0k3n",
			},
			payload: &pb.Credentials{
				Username: "test",
				Password: "test",
			},
			mock: func() {
				svc.EXPECT().IsUserExist(gomock.Any()).Return(true)
				svc.EXPECT().ValidateUser(gomock.Any()).Return("t0k3n", nil)
			},
		},
		{
			name:      "user doesn't exists",
			withError: true,
			want:      nil,
			payload: &pb.Credentials{
				Username: "test",
				Password: "test",
			},
			mock: func() {
				svc.EXPECT().IsUserExist(gomock.Any()).Return(false)
			},
		},
		{
			name:      "wrong password",
			withError: true,
			want:      nil,
			payload: &pb.Credentials{
				Username: "test",
				Password: "test",
			},
			mock: func() {
				svc.EXPECT().IsUserExist(gomock.Any()).Return(true)
				svc.EXPECT().ValidateUser(gomock.Any()).Return("", errors.New("123"))
			},
		},
	}

	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {
			tcase.mock()

			h := auth.NewHandler(svc, nil, nil)
			result, err := h.Login(context.Background(), tcase.payload)
			if err != nil && !tcase.withError {
				t.Errorf("Login() got = %v, err = %v", err, tcase.withError)
			}

			if err != nil && tcase.withError {
				return
			}

			if result.Token != tcase.want.Token {
				t.Errorf("Login() got = %v, expected = %v", result, tcase.want)
			}
		})
	}
}

func TestHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := auth_mock.NewMockIService(ctrl)

	tests := []struct {
		name      string
		withError bool
		want      *pb.Token
		payload   *pb.BasicCredentials
		mock      func()
	}{
		{
			name:      "success register",
			withError: false,
			want: &pb.Token{
				Token: "t0k3n",
			},
			payload: &pb.BasicCredentials{
				Username: "username",
				Name:     "name",
				Password: "password",
				Role:     "role",
			},
			mock: func() {
				svc.EXPECT().IsUserExist(gomock.Any()).Return(false)
				svc.EXPECT().RegisterUser(gomock.Any()).Return("t0k3n", nil)
			},
		},
		{
			name:      "user already exists",
			withError: true,
			want:      nil,
			payload: &pb.BasicCredentials{
				Username: "username",
				Name:     "name",
				Password: "password",
				Role:     "role",
			},
			mock: func() {
				svc.EXPECT().IsUserExist(gomock.Any()).Return(true)
			},
		},
		{
			name:      "unsuccess register",
			withError: true,
			want:      nil,
			payload: &pb.BasicCredentials{
				Username: "username",
				Name:     "name",
				Password: "password",
				Role:     "role",
			},
			mock: func() {
				svc.EXPECT().IsUserExist(gomock.Any()).Return(false)
				svc.EXPECT().RegisterUser(gomock.Any()).Return("", errors.New("123"))
			},
		},
	}

	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {
			tcase.mock()

			h := auth.NewHandler(svc, nil, nil)
			result, err := h.Register(context.Background(), tcase.payload)
			if err != nil && !tcase.withError {
				t.Errorf("Register() got = %v, err = %v", err, tcase.withError)
			}

			if err != nil && tcase.withError {
				return
			}

			if result.Token != tcase.want.Token {
				t.Errorf("Register() got = %v, expected = %v", result, tcase.want)
			}
		})
	}
}
