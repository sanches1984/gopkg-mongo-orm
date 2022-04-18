package mongodb

import (
	"fmt"
	"github.com/Kamva/mgm"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Config struct {
	Host       string
	DbName     string
	ReplicaSet string
	User       string
	Password   string
	Timeout    time.Duration
}

func Connect(appName string, config Config) (IClient, error) {
	err := mgm.SetDefaultConfig(
		&mgm.Config{CtxTimeout: config.Timeout},
		config.DbName,
		options.Client().SetAppName(appName),
		options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", config.Host)),
		options.Client().SetConnectTimeout(config.Timeout),
		options.Client().SetReplicaSet(config.ReplicaSet),
		options.Client().SetAuth(options.Credential{
			Username: config.User,
			Password: config.Password,
		}),
	)
	if err != nil {
		return nil, err
	}

	_, _, db, err := mgm.DefaultConfigs()
	if err != nil {
		return nil, err
	}

	return &dbWrapper{db: db}, nil
}
