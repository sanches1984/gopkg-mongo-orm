// +build !ci

package mongodb

import (
	"context"
	"errors"
	"github.com/Kamva/mgm"
	"github.com/Kamva/mgm/operator"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
	"time"
)

type testItem struct {
	collection       struct{} `bson:"books"`
	mgm.DefaultModel `bson:",inline"`
	Name             string `bson:"name"`
	Value            int    `bson:"value"`
}

func TestCRUD(t *testing.T) {
	ctx := context.Background()
	client, err := Connect("test", Config{
		Host:     "localhost:8113",
		DbName:   "test",
		User:     "root",
		Password: "password",
		Timeout:  5 * time.Second,
	})
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
	err = client.GetByID(ctx, v)
	require.NoError(t, err)
	require.Equal(t, item.GetID(), v.GetID())
	require.Equal(t, item.Name, v.Name)
	require.Equal(t, item.Value, v.Value)

	arr := []*testItem{}
	err = client.Find(ctx, &arr, NewSearchFilter(bson.M{"name": bson.M{operator.Eq: "hello"}}, 1, 10, "created_at desc"))
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(arr), 1)

	err = client.Delete(ctx, item)
	require.NoError(t, err)
}

func TestTransaction(t *testing.T) {
	ctx := context.Background()
	client, err := Connect("test", Config{
		Host:       "localhost:8113",
		DbName:     "test",
		ReplicaSet: "rs",
		User:       "root",
		Password:   "password",
		Timeout:    5 * time.Second,
	})
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
