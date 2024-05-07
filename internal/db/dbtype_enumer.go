// Code generated by "enumer -type=DBType"; DO NOT EDIT.

package db

import (
	"fmt"
	"strings"
)

const _DBTypeName = "PostgresSqlite"

var _DBTypeIndex = [...]uint8{0, 8, 14}

const _DBTypeLowerName = "postgressqlite"

func (i DBType) String() string {
	if i < 0 || i >= DBType(len(_DBTypeIndex)-1) {
		return fmt.Sprintf("DBType(%d)", i)
	}
	return _DBTypeName[_DBTypeIndex[i]:_DBTypeIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _DBTypeNoOp() {
	var x [1]struct{}
	_ = x[Postgres-(0)]
	_ = x[Sqlite-(1)]
}

var _DBTypeValues = []DBType{Postgres, Sqlite}

var _DBTypeNameToValueMap = map[string]DBType{
	_DBTypeName[0:8]:       Postgres,
	_DBTypeLowerName[0:8]:  Postgres,
	_DBTypeName[8:14]:      Sqlite,
	_DBTypeLowerName[8:14]: Sqlite,
}

var _DBTypeNames = []string{
	_DBTypeName[0:8],
	_DBTypeName[8:14],
}

// DBTypeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func DBTypeString(s string) (DBType, error) {
	if val, ok := _DBTypeNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _DBTypeNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to DBType values", s)
}

// DBTypeValues returns all values of the enum
func DBTypeValues() []DBType {
	return _DBTypeValues
}

// DBTypeStrings returns a slice of all String values of the enum
func DBTypeStrings() []string {
	strs := make([]string, len(_DBTypeNames))
	copy(strs, _DBTypeNames)
	return strs
}

// IsADBType returns "true" if the value is listed in the enum definition. "false" otherwise
func (i DBType) IsADBType() bool {
	for _, v := range _DBTypeValues {
		if i == v {
			return true
		}
	}
	return false
}
