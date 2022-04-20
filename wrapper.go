package mongodb

import (
	"context"
	"errors"
	"github.com/Kamva/mgm"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
)

const fieldCollection = "collection"

var errIncorrectModelInterface = errors.New("incorrect model interface")
var errNoReplicaSet = errors.New("use replica set for transactions")

type dbWrapper struct {
	db    *mongo.Database
	hasRS bool
}

func (w *dbWrapper) DB() *mongo.Database {
	return w.db
}

func (w *dbWrapper) Ping(ctx context.Context) error {
	return w.db.Client().Ping(ctx, nil)
}

func (w *dbWrapper) Close() error {
	return w.db.Client().Disconnect(context.Background())
}

// WithTX run in transaction (need replica set!)
func (w *dbWrapper) WithTX(ctx context.Context, fn func(context.Context) error) error {
	if !w.hasRS {
		return errNoReplicaSet
	}
	return mgm.TransactionWithCtx(ctx, func(session mongo.Session, sc mongo.SessionContext) error {
		if err := fn(sc); err != nil {
			rollbackErr := session.AbortTransaction(sc)
			if rollbackErr != nil {
				// todo get logger from context
				log.Error().Err(rollbackErr).Msg("failed to rollback transaction")
			}
			return err
		}

		return session.CommitTransaction(sc)
	})
}

func (w *dbWrapper) Create(ctx context.Context, rec interface{}) error {
	coll, elem, err := getCollectionAndModel(rec)
	if err != nil {
		return err
	}
	return coll.CreateWithCtx(ctx, elem)
}

func (w *dbWrapper) Update(ctx context.Context, rec interface{}) error {
	coll, elem, err := getCollectionAndModel(rec)
	if err != nil {
		return err
	}
	return coll.UpdateWithCtx(ctx, elem)
}

func (w *dbWrapper) Upsert(ctx context.Context, rec interface{}) error {
	coll, elem, err := getCollectionAndModel(rec)
	if err != nil {
		return err
	}

	return coll.UpdateWithCtx(ctx, elem, options.Update().SetUpsert(true))
}

func (w *dbWrapper) UpdateWhere(ctx context.Context, rec interface{}, filter Filter) (int64, error) {
	coll, err := getCollection(rec)
	if err != nil {
		return 0, err
	}

	res, err := coll.UpdateMany(ctx, filter.Conditions(), nil, options.Update().SetUpsert(true))
	if err != nil {
		return 0, err
	}

	return res.ModifiedCount, nil
}

func (w *dbWrapper) Delete(ctx context.Context, rec interface{}) error {
	coll, elem, err := getCollectionAndModel(rec)
	if err != nil {
		return err
	}
	return coll.DeleteWithCtx(ctx, elem)
}

func (w *dbWrapper) DeleteWhere(ctx context.Context, rec interface{}, filter Filter) (int64, error) {
	coll, err := getCollection(rec)
	if err != nil {
		return 0, err
	}
	res, err := coll.DeleteMany(ctx, filter.Conditions())
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

func (w *dbWrapper) GetByID(ctx context.Context, rec interface{}) error {
	coll, elem, err := getCollectionAndModel(rec)
	if err != nil {
		return err
	}

	return coll.FindByIDWithCtx(ctx, elem.GetID(), elem)
}

func (w *dbWrapper) Find(ctx context.Context, rec interface{}, filter SearchFilter) error {
	coll, err := getCollectionFromSlice(rec)
	if err != nil {
		return err
	}

	return coll.SimpleFindWithCtx(ctx, rec, filter.Conditions(),
		options.Find().SetSkip(filter.Skip()),
		options.Find().SetLimit(filter.Limit()),
		options.Find().SetSort(filter.Order()),
	)
}

func getCollectionFromSlice(arr interface{}) (*mgm.Collection, error) {
	v := reflect.ValueOf(arr).Elem()
	if v.Kind() != reflect.Slice {
		return nil, errIncorrectModelInterface
	}

	obj := reflect.New(v.Type().Elem()).Elem().Interface()
	return getCollection(obj)
}

func getCollection(item interface{}) (*mgm.Collection, error) {
	_, ok := item.(mgm.Model)
	if !ok {
		return nil, errIncorrectModelInterface
	}

	if field, ok := reflect.TypeOf(item).Elem().FieldByName(fieldCollection); ok {
		collName := field.Tag.Get("bson")
		if collName != "" {
			return mgm.CollectionByName(collName), nil
		}
	}

	return nil, errIncorrectModelInterface
}

func getCollectionAndModel(item interface{}) (*mgm.Collection, mgm.Model, error) {
	v, ok := item.(mgm.Model)
	if !ok {
		return nil, nil, errIncorrectModelInterface
	}

	if field, ok := reflect.TypeOf(item).Elem().FieldByName(fieldCollection); ok {
		collName := field.Tag.Get("bson")
		if collName != "" {
			return mgm.CollectionByName(collName), v, nil
		}
	}

	return nil, nil, errIncorrectModelInterface
}
