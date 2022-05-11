package opt

import (
	"github.com/sanches1984/gopkg-mongo-orm/repository/filter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
)

// Opt is options for database requests
type Opt struct {
	Skip      int64
	Limit     int64
	SortBy    string
	SortOrder int
	Filter    filter.Filter
}

// FnOpt is a function that modifies options
type FnOpt func(*Opt)

// New creates new Opt
func New(optFn ...FnOpt) *Opt {
	o := &Opt{}
	for _, fn := range optFn {
		if fn != nil {
			fn(o)
		}
	}

	return o
}

// GetFilter returns filter
func (o *Opt) GetFilter() bson.M {
	if o.IsFilter() {
		return o.Filter.Apply()
	}

	return nil
}

// GetOptions returns find options
func (o *Opt) GetOptions() *options.FindOptions {
	opts := options.Find()
	if o.IsPaging() {
		opts.SetSkip(o.Skip).SetLimit(o.Limit)
	}

	if o.IsSorting() {
		opts.SetSort(bson.D{{o.SortBy, o.SortOrder}})
	}

	return opts
}

// GetFilter returns filter
func GetFilter(optFn ...FnOpt) bson.M {
	return New(optFn...).GetFilter()
}

// GetOptions returns find options
func GetOptions(optFn ...FnOpt) *options.FindOptions {
	return New(optFn...).GetOptions()
}

// IsPaging responds whether pagination options set
func (o *Opt) IsPaging() bool {
	return o.Limit > 0 || o.Skip > 0
}

// IsSorting responds whether sorting options set
func (o *Opt) IsSorting() bool {
	return o.SortOrder != 0 && o.SortBy != ""
}

// IsFilter responds whether filter options set
func (o *Opt) IsFilter() bool {
	return len(o.Filter) > 0
}

// List converts periodic opts args into slice
func List(optFn ...FnOpt) []FnOpt {
	return optFn
}

// Paging sets both page and page size options
func Paging(page, size int) FnOpt {
	return func(opt *Opt) {
		opt.Limit = int64(size)
		opt.Skip = int64((page - 1) * size)
	}
}

// Asc sets ascending order options
func Asc(column string) FnOpt {
	return func(opt *Opt) {
		opt.SortBy = column
		opt.SortOrder = 1
	}
}

// Desc sets descending order options
func Desc(column string) FnOpt {
	return func(opt *Opt) {
		opt.SortBy = column
		opt.SortOrder = -1
	}
}

// Eq adds to filter equal condition
func Eq(column string, val interface{}) FnOpt {
	return func(opt *Opt) {
		opt.Filter = append(opt.Filter, filter.Eq{column: val})
	}
}

// Gt adds to filter great than condition
func Gt(column string, val interface{}) FnOpt {
	return func(opt *Opt) {
		opt.Filter = append(opt.Filter, filter.Gt{column: val})
	}
}

// Ge adds to filter great and equal condition
func Ge(column string, val interface{}) FnOpt {
	return func(opt *Opt) {
		opt.Filter = append(opt.Filter, filter.Ge{column: val})
	}
}

// Lt adds to filter less than condition
func Lt(column string, val interface{}) FnOpt {
	return func(opt *Opt) {
		opt.Filter = append(opt.Filter, filter.Lt{column: val})
	}
}

// Le adds to filter less and equal condition
func Le(column string, val interface{}) FnOpt {
	return func(opt *Opt) {
		opt.Filter = append(opt.Filter, filter.Le{column: val})
	}
}

// Neq adds to filter not-equal condition
func Neq(column string, val interface{}) FnOpt {
	return func(opt *Opt) {
		opt.Filter = append(opt.Filter, filter.Ne{column: val})
	}
}

// In sets condition for IN operation
func In(column string, vals interface{}) FnOpt {
	return func(opt *Opt) {
		if reflect.TypeOf(vals).Kind() != reflect.Slice {
			vals = []interface{}{vals}
		}

		in := []interface{}{}
		v := reflect.ValueOf(vals)
		for i := 0; i < v.Len(); i++ {
			in = append(in, v.Index(i).Interface())
		}

		opt.Filter = append(opt.Filter, filter.In{column: in})
	}
}

// Contains builds a condition with contains statement
func Contains(column string, val string) FnOpt {
	return func(opt *Opt) {
		opt.Filter = append(opt.Filter, filter.Contains{column: val})
	}
}

// Match builds a condition with match statement
func Match(column string, val string) FnOpt {
	return func(opt *Opt) {
		opt.Filter = append(opt.Filter, filter.Match{column: val})
	}
}

// Or adds set of conditions joined with OR statement
func Or(optFn ...FnOpt) FnOpt {
	return func(opt *Opt) {
		o := New(optFn...)
		opt.Filter = append(opt.Filter, filter.Or(o.Filter))
	}
}

// And adds set of conditions joined with AND statement
func And(optFn ...FnOpt) FnOpt {
	return func(opt *Opt) {
		o := New(optFn...)
		opt.Filter = append(opt.Filter, filter.And(o.Filter))
	}
}

// NotNull adds `IS NOT NULL` condition
func NotNull(column string) FnOpt {
	return func(opt *Opt) {
		opt.Filter = append(opt.Filter, filter.NotNull(column))
	}
}

// IsNull adds `IS NULL` condition
func IsNull(column string) FnOpt {
	return func(opt *Opt) {
		opt.Filter = append(opt.Filter, filter.IsNull(column))
	}
}

// Not adds `NOT` condition
func Not(optFn ...FnOpt) FnOpt {
	return func(opt *Opt) {
		o := New(optFn...)
		opt.Filter = append(opt.Filter, filter.Not(o.Filter))
	}
}
