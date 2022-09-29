package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Model interface {
	PrepareID(id interface{}) (interface{}, error)
	IsNew() bool
	GetID() interface{}
	SetID(id interface{})
	Creating()
	Updating()
}

// DefaultModel struct contain model's default fields.
type DefaultModel struct {
	IDField    `bson:",inline"`
	DateFields `bson:",inline"`
}

// IDField struct contain model's ID field.
type IDField struct {
	ID primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
}

// DateFields struct contain `created_at` and `updated_at`
// fields that autofill on insert/update model.
type DateFields struct {
	CreatedAt time.Time `json:"created_at" bson:"createdAt"`
	UpdatedAt time.Time `json:"updated_at" bson:"updatedAt"`
}

// PrepareID method prepare id value to using it as id in filtering,...
// e.g convert hex-string id value to bson.ObjectId
func (f *IDField) PrepareID(id interface{}) (interface{}, error) {
	if idStr, ok := id.(string); ok {
		return primitive.ObjectIDFromHex(idStr)
	}

	// Otherwise id must be ObjectId
	return id, nil
}

// IsNew method check and say that model is new or not.
//
// Deprecated: this method is deprecated and remove in version 2.
func (f *IDField) IsNew() bool {
	return f.GetID() == primitive.ObjectID{}
}

// GetID method return model's id
func (f *IDField) GetID() interface{} {
	return f.ID
}

// SetID set id value of model's id field.
func (f *IDField) SetID(id interface{}) {
	f.ID = id.(primitive.ObjectID)
}

//--------------------------------
// DateField methods
//--------------------------------

// Creating hook
func (f *DateFields) Creating() {
	f.CreatedAt = time.Now().UTC()
	f.UpdatedAt = time.Now().UTC()
}

// Updating hook
func (f *DateFields) Updating() {
	f.UpdatedAt = time.Now().UTC()
}
