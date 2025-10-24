package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Array []string

func (a *Array) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		return json.Unmarshal(s, a)
	case string:
		return json.Unmarshal([]byte(s), a)
	}
	return errors.New("Array type not support")
}

func (a Array) Value() (driver.Value, error) {
	return json.Marshal(a)
}
