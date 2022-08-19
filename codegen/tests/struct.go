package tests

import "time"

type BaseStruct struct {
	IntValue  int
	BoolValue bool
	Embedded  EmbeddedStruct
}

type EmbeddedStruct struct {
	FieldA int
}

// ComplexStruct written descriptions on structs are show as descriptions in the table
type ComplexStruct struct {
	IntValue    int
	Int8Value   int8
	Int16Value  int16
	Int32Value  int32
	Int64Value  int64
	String      string
	StringArray []string
	BoolValue   bool
	// Comments written as descriptions on fields are show as descriptions
	Json map[string]interface{}
}

type RelationStruct struct {
	Int64Value  int64
	String      string
	StringArray []string
	SomeBases   []BaseStruct
	Inner       RelationalRelation
	Relations   []RelationalRelation
}

type RelationalRelation struct {
	EmbeddedString string
	SomeBases      []BaseStruct
}

type SimpleRelation struct {
	Column    int
	Relations []BaseStruct
}

type OtherStruct struct {
	OtherField int
}

type DeepRelationStruct struct {
	Relations []DeepRelationalRelation
}

type DeepRelationalRelation struct {
	Relations []DeepDeepRelationalRelation
}

type DeepDeepRelationalRelation struct {
	Relations []RelationalRelation
}

// DescriptionStruct describes itself as a struct.
// Also, a second sentence.
type DescriptionStruct struct {
	// Simple is a simple description
	Simple string

	// Value is a string. It has more properties, but we will
	// only tell you that it is a string. This is to make sure we have
	// multiple sentences.
	Complex string

	// A sentence with a full stop that should get removed.
	CreatedAt *time.Time
}
