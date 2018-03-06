package misc

import "database/sql"

// NullString is extended replacement for sql.NullString
type NullString struct {
	sql.NullString
}

// MarshalJSON marshals string directly to json. If it is null (valid == false), marshals as "null"
func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return []byte("\"" + ns.String + "\""), nil
}

// UnmarshalJSON unmarshals string from json. If null was received, marks as not valid
func (ns *NullString) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		ns.Valid = false
		return nil
	}
	ns.String = string(data)
	ns.Valid = true
	return nil
}

func WrapString(value string) (ret NullString) {
	ret.Valid = true
	ret.String = value
	return
}
