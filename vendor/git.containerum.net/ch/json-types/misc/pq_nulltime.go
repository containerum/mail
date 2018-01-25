package misc

import (
	"database/sql/driver"
	"time"
)

// PqNullTime is more convenient replacement for pq.NullTime
type PqNullTime struct {
	time.Time
	IsNull bool
}

// Scan implements the Scanner interface.
func (nt *PqNullTime) Scan(value interface{}) error {
	nt.Time, nt.IsNull = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt PqNullTime) Value() (driver.Value, error) {
	if !nt.IsNull {
		return nil, nil
	}
	return nt.Time, nil
}

// UnmarshalJSON extends a standard Time.UnmarshalJSON functionality. If received data is null, set IsNull to true
func (nt PqNullTime) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		nt.IsNull = true
		return nil
	}
	nt.IsNull = false
	return nt.Time.UnmarshalJSON(data)
}

// MarshalJSON extends a standard Time.Marshal functionality. If IsNull set to true, return nil data
func (nt PqNullTime) MarshalJSON() ([]byte, error) {
	if nt.IsNull {
		return nil, nil
	}
	return nt.Time.MarshalJSON()
}
