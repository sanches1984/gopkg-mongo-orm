package filter

import (
	"go.mongodb.org/mongo-driver/bson"
)

type Filter []Condition

func (f Filter) Apply() bson.M {
	var result []interface{}
	for _, cond := range f {
		result = append(result, cond.Condition())
	}
	return bson.M{"$and": result}
}
