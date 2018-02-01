package misc

import (
	"database/sql"
	"strconv"
)

// NullBool is extended replacement for sql.NullBool
type NullBool struct {
	sql.NullBool
}

// UnmarshalJSON implements json.Unmarshaler interface. If received data is "null" it will mark it as not valid
func (nb *NullBool) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 0 || string(data) == "null" {
		nb.Valid = false
		return
	}
	nb.Valid = true
	nb.Bool, err = strconv.ParseBool(string(data))
	return
}

// MarshalJSON implements json.Marshaler interface. If type marked as not valid, send null
func (nb NullBool) MarshalJSON() ([]byte, error) {
	if !nb.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatBool(nb.Bool)), nil
}

func WrapBool(value bool) (ret NullBool) {
	ret.Valid = true
	ret.Bool = value
	return
}
