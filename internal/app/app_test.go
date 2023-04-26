package app_test

import (
	"testing"

	"github.com/lgu-elo/auth/internal/app"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func TestValidateApp(t *testing.T) {
	err := fx.ValidateApp(app.CreateApp())
	require.NoError(t, err)
}
