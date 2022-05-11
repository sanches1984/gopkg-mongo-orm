package mongodb

import (
	"errors"
	"fmt"
	"github.com/Kamva/mgm"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"time"
)

const defaultTimeout = 60 * time.Second

type Config struct {
	Host       string
	Database   string
	ReplicaSet string
	Username   string
	Password   string
	Timeout    time.Duration
}

func Connect(appName string, config *Config) (IClient, error) {
	if config == nil {
		return nil, errors.New("config not set")
	}

	opts := make([]*options.ClientOptions, 0, 5)
	opts = append(opts,
		options.Client().SetAppName(appName),
		options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", config.Host)),
		options.Client().SetConnectTimeout(config.Timeout),
		options.Client().SetReplicaSet(config.ReplicaSet))
	if config.Username != "" && config.Password != "" {
		opts = append(opts,
			options.Client().SetAuth(options.Credential{
				Username: config.Username,
				Password: config.Password,
			}))
	}

	err := mgm.SetDefaultConfig(&mgm.Config{CtxTimeout: config.Timeout}, config.Database, opts...)
	if err != nil {
		return nil, err
	}

	_, _, db, err := mgm.DefaultConfigs()
	if err != nil {
		return nil, err
	}

	return &dbWrapper{db: db, hasRS: config.ReplicaSet != ""}, nil
}

func ParseURL(dsn string) (*Config, error) {
	uri, err := connstring.Parse(dsn)
	if err != nil {
		return nil, err
	}

	host := "localhost"
	if len(uri.Hosts) != 0 {
		host = uri.Hosts[0]
	}

	return &Config{
		Host:       host,
		Database:   uri.Database,
		ReplicaSet: uri.ReplicaSet,
		Username:   uri.Username,
		Password:   uri.Password,
		Timeout:    defaultTimeout,
	}, nil
}
