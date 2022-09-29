//go:build !ci
// +build !ci

package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sanches1984/gopkg-mongo-orm/model"
	"github.com/sanches1984/gopkg-mongo-orm/repository/opt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

type testItem struct {
	collection         struct{} `bson:"books"`
	model.DefaultModel `bson:",inline"`
	Name               string `bson:"name"`
	Value              int    `bson:"value"`
}

func TestCRUD(t *testing.T) {
	ctx := context.Background()
	cfg := getConfig()
	client, err := Connect("test", cfg)
	require.NoError(t, err)

	defer client.Close()

	item := &testItem{Name: "hello", Value: 123}
	err = client.Create(ctx, item)
	require.NoError(t, err)
	require.NotEmpty(t, item.ID)
	require.NotEmpty(t, item.CreatedAt)

	item.Value = 456
	err = client.Update(ctx, item)
	require.NoError(t, err)
	require.NotEmpty(t, item.ID)
	require.NotEmpty(t, item.CreatedAt)

	v := &testItem{}
	v.SetID(item.GetID())
	err = client.FindByID(ctx, v)
	require.NoError(t, err)
	require.Equal(t, item.GetID(), v.GetID())
	require.Equal(t, item.Name, v.Name)
	require.Equal(t, item.Value, v.Value)

	arr := []*testItem{}
	err = client.Find(ctx, &arr, opt.List(
		opt.Eq("name", "hello"),
		opt.Desc("created_at"),
		opt.Paging(1, 10),
	))
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(arr), 1)

	err = client.Delete(ctx, item)
	require.NoError(t, err)
}

func TestTransaction(t *testing.T) {
	ctx := context.Background()
	cfg := getConfig()
	client, err := Connect("test", cfg)
	require.NoError(t, err)

	err = client.WithTX(ctx, func(ctx context.Context) error {
		item := &testItem{Name: "hello_bad", Value: 555}
		err = client.Create(ctx, item)
		require.NoError(t, err)

		return errors.New("rollback")
	})
	require.Error(t, err)

	err = client.WithTX(ctx, func(ctx context.Context) error {
		item := &testItem{Name: "hello_good", Value: 777}
		err = client.Create(ctx, item)
		require.NoError(t, err)

		return nil
	})
	require.NoError(t, err)
}

func getConfig() *Config {
	if err := godotenv.Load(); err != nil {
		fmt.Println(".env file not found")
		return nil
	}

	cfg, err := ParseURL(os.Getenv("DSN"))
	if err != nil {
		fmt.Println("can't parse dsn")
		return nil
	}

	return cfg
}
