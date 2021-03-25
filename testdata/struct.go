package testdata

import (
	"time"
)

type EnumField string

type TestStruct struct {
	Username string
	unexported string
	Skip string `json:"-"`
	Something EnumField
	Time time.Time
	Inner OtherStruct
	FieldPtr *string
	//FieldSlice []string
	//FieldSlicePtr []*string
	//FieldPtrSlice *[]string
	//FieldPtrSlicePtr *[]*string
	//
	//// This will be omitted
	//EmptyField interface{}
	////
	////MapField map[string]interface{}
	//NoPtr NotPtrStruct

	//S *OtherStruct
}

type OtherStruct struct {
	Field int
	FieldB string
}

type NotPtrStruct struct {
	FieldB int
	FieldA bool
}