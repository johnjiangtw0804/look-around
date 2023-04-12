package envconfig

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_env(t *testing.T) {
	// required env
	os.Setenv("DATABASE_URL", "http://example.com")
	_, err := New()
	require.NoError(t, err)

	os.Setenv("PORT", "8080")
	_, err = New()
	require.NoError(t, err)

	// optional env
	require.NoError(t, os.Setenv("DEBUG", "true"))
	_, err = New()
	require.NoError(t, err)

	require.Equal(t, os.Getenv("DEBUG"), "true")
	require.Equal(t, os.Getenv("PORT"), "8080")
	require.Equal(t, os.Getenv("DATABASE_URL"), "http://example.com")
}
