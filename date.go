package alphavantage

import (
	"strings"
	"time"

	"github.com/pkg/errors"
)

const dateLayout = "2006-01-02"

// Date just money type
type Date time.Time

// panicParseDate parses alphavantage date types
func panicParseDate(v string) Date {
	res, err := time.Parse(dateLayout, v)
	if err != nil {
		panic(errors.Wrapf(err, "Cannot parse '%s'", v))
	}
	return Date(res.UTC())
}

// UnmarshalJSON decodes DateKey
func (d *Date) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	*d = panicParseDate(s)
	return
}

// MarshalJSON converts DailyRawResponse to json string
func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(d.String()), nil
}

// String converts Money to string
func (d Date) String() string {
	return (time.Time(d)).Format(dateLayout)
}
