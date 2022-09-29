//go:build !ci
// +build !ci

package migrate

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	mongodb "github.com/sanches1984/gopkg-mongo-orm"
	"github.com/sanches1984/gopkg-mongo-orm/model"
	"github.com/sanches1984/gopkg-mongo-orm/repository/opt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

type testItem struct {
	model.DefaultModel `bson:",inline"`
	collection         struct{} `bson:"books"`
	Name               string   `bson:"name"`
	Value              int      `bson:"value"`
}

func TestMigrate_Run(t *testing.T) {
	ctx := context.Background()
	cfg := getConfig()
	migrator := NewMigrator("test/migrations", os.Getenv("DSN"), WithClean("test"))
	err := migrator.Run(context.Background())
	require.NoError(t, err)

	client, err := mongodb.Connect("test", cfg)
	require.NoError(t, err)

	item := &testItem{
		Name:  "some new test name",
		Value: 567,
	}

	err = client.Create(ctx, item)
	require.NoError(t, err)
	require.NotEmpty(t, item.ID)

	item.Value = 777
	err = client.Update(ctx, item)
	require.NoError(t, err)

	arr := []*testItem{}
	err = client.Find(ctx, &arr, opt.List(opt.Eq("name", "some new test name")))
	require.NoError(t, err)
	require.Len(t, arr, 1)
	require.Equal(t, arr[0].Name, "some new test name")
	require.Equal(t, arr[0].Value, 777)

	err = client.Delete(ctx, item)
	require.NoError(t, err)
}

func getConfig() *mongodb.Config {
	if err := godotenv.Load(); err != nil {
		fmt.Println(".env file not found")
		return nil
	}

	cfg, err := mongodb.ParseURL(os.Getenv("DSN"))
	if err != nil {
		fmt.Println("can't parse dsn")
		return nil
	}

	return cfg
}
