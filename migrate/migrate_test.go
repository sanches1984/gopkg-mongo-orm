//+build !ci

package migrate

import (
	"context"
	"fmt"
	"github.com/Kamva/mgm"
	"github.com/joho/godotenv"
	mongodb "github.com/sanches1984/gopkg-mongo-orm"
	"github.com/sanches1984/gopkg-mongo-orm/repository/opt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

type testItem struct {
	collection       struct{} `bson:"books"`
	mgm.DefaultModel `bson:",inline"`
	Name             string `bson:"name"`
	Value            int    `bson:"value"`
}

func TestMigrate_Run(t *testing.T) {
	ctx := context.Background()
	cfg := getConfig()
	migrator := NewMigrator("test/migrations", os.Getenv("DSN"), WithClean("test"))
	err := migrator.Run(context.Background())
	require.NoError(t, err)

	client, err := mongodb.Connect("test", cfg)
	require.NoError(t, err)

	arr := []*testItem{}
	err = client.Find(ctx, &arr, opt.List(opt.Eq("name", "some book")))

	require.NoError(t, err)
	require.Len(t, arr, 1)
	require.Equal(t, arr[0].Name, "some book")
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
