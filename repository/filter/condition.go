package filter

import (
	"go.mongodb.org/mongo-driver/bson"
)

// Condition common facade for filter conditions
type Condition interface {
	Condition() bson.M
}

// Eq field name equal to value
type Eq map[string]interface{}

// Ne field name not equal to value
type Ne map[string]interface{}

// Lt field name less than value
type Lt map[string]interface{}

// Le field name less than value or equal
type Le map[string]interface{}

// Gt field name greater than value
type Gt map[string]interface{}

// Ge field name greater than value or equal
type Ge map[string]interface{}

// In field name contains in list of values
type In map[string][]interface{}

// NotIn field name not contains in list of values
type NotIn map[string][]interface{}

// Match filter
type Match map[string]string

// Contains filter
type Contains map[string]string

// IsNull field name equal to NULL
type IsNull string

// NotNull field name not equal to NULL
type NotNull string

// Or filter
type Or []Condition

// And filter
type And []Condition

// Not filter
type Not []Condition

func (c Eq) Condition() bson.M {
	for key, val := range c {
		return bson.M{key: bson.M{"$eq": val}}
	}
	return nil
}

func (c Ne) Condition() bson.M {
	for key, val := range c {
		return bson.M{key: bson.M{"$ne": val}}
	}
	return nil
}

func (c Lt) Condition() bson.M {
	for key, val := range c {
		return bson.M{key: bson.M{"$lt": val}}
	}
	return nil
}

func (c Le) Condition() bson.M {
	for key, val := range c {
		return bson.M{key: bson.M{"$lte": val}}
	}
	return nil
}

func (c Gt) Condition() bson.M {
	for key, val := range c {
		return bson.M{key: bson.M{"$gt": val}}
	}
	return nil
}

func (c Ge) Condition() bson.M {
	for key, val := range c {
		return bson.M{key: bson.M{"$gte": val}}
	}
	return nil
}

func (c In) Condition() bson.M {
	for key, val := range c {
		return bson.M{key: bson.M{"$in": val}}
	}
	return nil
}

func (c NotIn) Condition() bson.M {
	for key, val := range c {
		return bson.M{key: bson.M{"$nin": val}}
	}
	return nil
}

func (c Contains) Condition() bson.M {
	for key, val := range c {
		return bson.M{key: bson.M{"$regex": val, "$options": "i"}}
	}
	return nil
}

func (c IsNull) Condition() bson.M {
	return bson.M{string(c): bson.M{"$eq": nil}}
}

func (c NotNull) Condition() bson.M {
	return bson.M{string(c): bson.M{"$ne": nil}}
}

func (c Or) Condition() bson.M {
	var result []interface{}
	for _, cond := range c {
		result = append(result, cond.Condition())
	}
	return bson.M{"$or": result}
}

func (c And) Condition() bson.M {
	var result []interface{}
	for _, cond := range c {
		result = append(result, cond.Condition())
	}
	return bson.M{"$and": result}
}

func (c Not) Condition() bson.M {
	var result []interface{}
	for _, cond := range c {
		result = append(result, cond.Condition())
	}
	return bson.M{"$not": result}
}
