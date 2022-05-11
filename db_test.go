package mongodb

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseURL(t *testing.T) {
	cfg, err := ParseURL("mongodb://admin:password@myhost:5555/test")

	require.NoError(t, err)
	require.Equal(t, &Config{
		Host:     "myhost:5555",
		Database: "test",
		Username: "admin",
		Password: "password",
		Timeout:  defaultTimeout,
	}, cfg)
}
