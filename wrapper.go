package mongodb

import (
	"context"
	"errors"
	"github.com/sanches1984/gopkg-mongo-orm/model"
	"github.com/sanches1984/gopkg-mongo-orm/repository/opt"
	"go.mongodb.org/mongo-driver/bson"
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
	return errors.New("method not implemented")
	//if !w.hasRS {
	//	return errNoReplicaSet
	//}
	//return mgm.TransactionWithCtx(ctx, func(session mongo.Session, sc mongo.SessionContext) error {
	//	if err := fn(sc); err != nil {
	//		rollbackErr := session.AbortTransaction(sc)
	//		if rollbackErr != nil {
	//			// todo get logger from context
	//			log.Error().Err(rollbackErr).Msg("failed to rollback transaction")
	//		}
	//		return err
	//	}
	//
	//	return session.CommitTransaction(sc)
	//})
}

func (w *dbWrapper) Create(ctx context.Context, rec interface{}) error {
	coll, elem, err := getCollectionAndModel(rec)
	if err != nil {
		return err
	}

	elem.Creating()
	res, err := w.db.Collection(coll).InsertOne(ctx, rec)
	if err != nil {
		return err
	}

	elem.SetID(res.InsertedID)
	return nil
}

func (w *dbWrapper) Update(ctx context.Context, rec interface{}) error {
	coll, elem, err := getCollectionAndModel(rec)
	if err != nil {
		return err
	}

	elem.Updating()
	_, err = w.db.Collection(coll).UpdateOne(ctx, bson.M{"_id": elem.GetID()}, bson.M{"$set": elem})
	return err
}

func (w *dbWrapper) Upsert(ctx context.Context, rec interface{}) error {
	coll, elem, err := getCollectionAndModel(rec)
	if err != nil {
		return err
	}

	if elem.IsNew() {
		elem.Creating()
	} else {
		elem.Updating()
	}

	res, err := w.db.Collection(coll).UpdateOne(ctx, bson.M{"_id": elem.GetID()}, bson.M{"$set": elem}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}

	elem.SetID(res.UpsertedID)
	return nil
}

func (w *dbWrapper) UpdateWhere(ctx context.Context, rec interface{}, opts []opt.FnOpt) (int64, error) {
	return 0, errors.New("method not implemented")
	//coll, err := getCollection(rec)
	//if err != nil {
	//	return 0, err
	//}
	//
	//res, err := coll.UpdateMany(ctx, opt.GetFilter(opts...), nil, options.Update().SetUpsert(true))
	//if err != nil {
	//	return 0, err
	//}
	//
	//return res.ModifiedCount, nil
}

func (w *dbWrapper) Delete(ctx context.Context, rec interface{}) error {
	coll, elem, err := getCollectionAndModel(rec)
	if err != nil {
		return err
	}

	_, err = w.db.Collection(coll).DeleteOne(ctx, bson.M{"_id": elem.GetID()})
	return err
}

func (w *dbWrapper) DeleteWhere(ctx context.Context, rec interface{}, opts []opt.FnOpt) (int64, error) {
	return 0, errors.New("method not implemented")
	//coll, err := getCollection(rec)
	//if err != nil {
	//	return 0, err
	//}
	//res, err := coll.DeleteMany(ctx, opt.GetFilter(opts...))
	//if err != nil {
	//	return 0, err
	//}
	//
	//return res.DeletedCount, nil
}

func (w *dbWrapper) FindByID(ctx context.Context, rec interface{}) error {
	coll, elem, err := getCollectionAndModel(rec)
	if err != nil {
		return err
	}

	res := w.db.Collection(coll).FindOne(ctx, bson.M{"_id": elem.GetID()})
	if err != nil {
		return err
	}

	return res.Decode(rec)
}

func (w *dbWrapper) Find(ctx context.Context, rec interface{}, opts []opt.FnOpt) error {
	coll, err := getCollectionFromSlice(rec)
	if err != nil {
		return err
	}

	res, err := w.db.Collection(coll).Find(ctx, opt.GetFilter(opts...), opt.GetOptions(opts...))
	if err != nil {
		return err
	}
	return res.All(ctx, rec)
}

func getCollectionFromSlice(arr interface{}) (string, error) {
	v := reflect.ValueOf(arr).Elem()
	if v.Kind() != reflect.Slice {
		return "", errIncorrectModelInterface
	}

	obj := reflect.New(v.Type().Elem()).Elem().Interface()
	return getCollection(obj)
}

func getCollection(item interface{}) (string, error) {
	_, ok := item.(model.Model)
	if !ok {
		return "", errIncorrectModelInterface
	}

	if field, ok := reflect.TypeOf(item).Elem().FieldByName(fieldCollection); ok {
		collName := field.Tag.Get("bson")
		if collName != "" {
			return collName, nil
		}
	}

	return "", errIncorrectModelInterface
}

func getCollectionAndModel(item interface{}) (string, model.Model, error) {
	v, ok := item.(model.Model)
	if !ok {
		return "", nil, errIncorrectModelInterface
	}

	if field, ok := reflect.TypeOf(item).Elem().FieldByName(fieldCollection); ok {
		collName := field.Tag.Get("bson")
		if collName != "" {
			return collName, v, nil
		}
	}

	return "", nil, errIncorrectModelInterface
}
