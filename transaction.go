package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type transactionFunc func(session mongo.Session, sc mongo.SessionContext) error

func transaction(ctx context.Context, client *mongo.Client, f transactionFunc) error {
	session, err := client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	if err = session.StartTransaction(); err != nil {
		return err
	}

	wrapperFn := func(sc mongo.SessionContext) error {
		return f(session, sc)
	}

	return mongo.WithSession(ctx, session, wrapperFn)
}
