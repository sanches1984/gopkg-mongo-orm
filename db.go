package mongodb

import "C"
import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"time"
)

const defaultTimeout = 20 * time.Second

type Config struct {
	DSN        string
	Schema     string
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

	db, err := connect(context.TODO(), config)
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
		DSN:        dsn,
		Schema:     uri.Scheme,
		Host:       host,
		Database:   uri.Database,
		ReplicaSet: uri.ReplicaSet,
		Username:   uri.Username,
		Password:   uri.Password,
		Timeout:    defaultTimeout,
	}, nil
}
func connect(ctx context.Context, config *Config) (*mongo.Database, error) {
	opts := options.Client().ApplyURI(config.DSN)
	client, err := mongo.NewClient(opts)
	if err != nil {
		return nil, err
	}
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("can't ping database: %v", err)
	}

	return client.Database(config.Database), nil
}
