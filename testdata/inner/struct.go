package inner


type OtherPackageStruct struct {
	FieldSimple string
	OtherStruct *OtherPackageInnerStruct
}

type OtherPackageInnerStruct struct {
	Field int
}