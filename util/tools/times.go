package tools

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const (
	TimeFormat = "2006-01-02 15:04:05"
)

type FormatTime struct {
	time.Time
}

func (t *FormatTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+TimeFormat+`"`, string(data), time.Local)
	*t = FormatTime{
		now,
	}
	return
}

func (t FormatTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format(TimeFormat))
	return []byte(formatted), nil
}

func (t FormatTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *FormatTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = FormatTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
