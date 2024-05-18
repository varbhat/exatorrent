package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

type RFC3339TimeStamp struct {
	time.Time
}

func NewRFC3339TimeStamp(t time.Time) RFC3339TimeStamp {
	return RFC3339TimeStamp{
		Time: t,
	}
}

func RFC3339TimeStampNow() RFC3339TimeStamp {
	return NewRFC3339TimeStamp(time.Now())
}

var _ driver.Valuer = (*RFC3339TimeStamp)(nil)

// Value returns the time as string using timeFormat.
func (tm RFC3339TimeStamp) Value() (driver.Value, error) {
	return tm.Format(time.RFC3339), nil
}

var _ sql.Scanner = (*RFC3339TimeStamp)(nil)

// Scan scans the time parsing it if necessary using timeFormat.
func (tm *RFC3339TimeStamp) Scan(src interface{}) (err error) {
	switch src := src.(type) {
	case time.Time:
		*tm = NewRFC3339TimeStamp(src)
		return nil
	case string:
		tm.Time, err = time.Parse(time.RFC3339, src)
		return err
	case []byte:
		tm.Time, err = time.Parse(time.RFC3339, string(src))
		return err
	case nil:
		tm.Time = time.Time{}
		return nil
	default:
		return fmt.Errorf("unsupported data type: %T", src)
	}
}
