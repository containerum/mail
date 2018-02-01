package misc

import (
	"database/sql/driver"
	"time"
)

// NullTime is more convenient replacement for pq.NullTime
type NullTime struct {
	time.Time
	Valid bool
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

// UnmarshalJSON extends a standard Time.UnmarshalJSON functionality. If received data is null, set IsNull to true
func (nt *NullTime) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		nt.Valid = false
		return nil
	}
	nt.Valid = true
	return nt.Time.UnmarshalJSON(data)
}

// MarshalJSON extends a standard Time.Marshal functionality. If IsNull set to true, return nil data
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return nt.Time.MarshalJSON()
}

func WrapTime(value time.Time) NullTime {
	return NullTime{Valid: true, Time: value}
}
