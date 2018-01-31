package misc

import "database/sql"

// SqlNullString is extended replacement for sql.NullString
type SqlNullString struct {
	sql.NullString
}

// MarshalJSON marshals string directly to json. If it is null (valid == false), marshals as "null"
func (ns SqlNullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return []byte(ns.String), nil
}

func (ns *SqlNullString) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		ns.Valid = false
		return nil
	}
	ns.String = string(data)
	ns.Valid = true
	return nil
}
