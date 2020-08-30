package alphavantage

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

const nullJSONString string = `"null"`

// Money just money type
type Money int64

func normalizeExp(v string) string {
	switch len(v) {
	case 4:
		return v
	case 3:
		return v + "0"
	case 2:
		return v + "00"
	case 1:
		return v + "000"
	case 0:
		return "0000"
	default:
		return ""
	}
}

// panicParseMoney parses alphavantage money types
func panicParseMoney(value string) Money {
	if value == "" || value == "None" || value == "nil" {
		return 0
	}
	parts := strings.Split(value, ".")
	if len(parts) != 1 && len(parts) != 2 {
		panic(errors.Errorf("Cannot parse money type '%s'", value))
	}

	exp := ""
	if len(parts) == 2 {
		exp = parts[1]
	}
	exp = normalizeExp(exp)
	if exp == "" {
		panic(errors.Errorf("Cannot parse money type '%s'", value))
	}
	return Money(panicParseInt64ish(parts[0] + exp))
}

// UnmarshalJSON decodes DateKey
func (m *Money) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == nullJSONString {
		*m = 0
		return
	}
	*m = panicParseMoney(s)
	return
}

// MarshalJSON converts DailyRawResponse to json string
func (m Money) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, m.String())), nil
}

// String converts Money to string
func (m Money) String() string {
	s := fmt.Sprintf("%+06d", int64(m))
	first := s[0 : len(s)-4]
	last4 := s[len(s)-4:]
	return fmt.Sprintf(`%s.%s`, strings.ReplaceAll(first, "+", ""), last4)
}
