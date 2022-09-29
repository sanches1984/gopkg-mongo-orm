package mongodb

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseURL(t *testing.T) {
	cases := []struct {
		url    string
		config *Config
	}{
		{
			url: "mongodb://admin:password@myhost:5555/test",
			config: &Config{
				DSN:      "mongodb://admin:password@myhost:5555/test",
				Schema:   "mongodb",
				Host:     "myhost:5555",
				Database: "test",
				Username: "admin",
				Password: "password",
				Timeout:  defaultTimeout,
			},
		},
	}

	for _, c := range cases {
		cfg, err := ParseURL(c.url)
		require.NoError(t, err)
		require.Equal(t, c.config, cfg)
	}
}
