package mongodb

import (
	"go.mongodb.org/mongo-driver/bson"
	"strings"
)

const defaultPageSize = 100

type Filter interface {
	Conditions() bson.M
}

type SearchFilter interface {
	Conditions() bson.M
	Skip() int64
	Limit() int64
	Order() bson.D
}

type filter struct {
	conditions bson.M
	page       int64
	pageSize   int64
	order      bson.D
}

func NewFilter(cond bson.M) Filter {
	return &filter{conditions: cond}
}

func NewSearchFilter(cond bson.M, page, pageSize int64, order string) SearchFilter {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = defaultPageSize
	}

	return &filter{
		conditions: cond,
		page:       page,
		pageSize:   pageSize,
		order:      parseOrder(order),
	}
}

func (f *filter) Conditions() bson.M {
	return f.conditions
}

func (f *filter) Skip() int64 {
	return (f.page - 1) * f.pageSize
}

func (f *filter) Limit() int64 {
	return f.pageSize
}

func (f *filter) Order() bson.D {
	return bson.D{{"name", 1}}
}

func parseOrder(order string) bson.D {
	var sort bson.D
	if len(order) != 0 {
		vals := strings.Split(order, " ")
		if len(vals) == 1 {
			sort = bson.D{{vals[0], 1}}
		} else {
			if strings.ToLower(vals[1]) == "desc" {
				sort = bson.D{{vals[0], -1}}
			} else {
				sort = bson.D{{vals[0], 1}}
			}
		}
	}
	return sort
}
