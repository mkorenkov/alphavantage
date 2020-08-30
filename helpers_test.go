package alphavantage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInt64Parse(t *testing.T) {
	testCases := map[string]int64{
		"None":          0,
		"0":             0,
		"-1":            -1,
		"1":             1,
		"152186000000":  152186000000,
		"-152186000000": -152186000000,
	}

	for input, expectedResult := range testCases {
		actualResult := panicParseInt64ish(input)
		assert.Equal(t, actualResult, expectedResult)
	}
}

func TestInt64ParsePanic(t *testing.T) {
	testCases := []string{
		"fourty two",
		"wow!",
		"-",
	}

	for _, input := range testCases {
		assert.Panics(t, func() { panicParseInt64ish(input) })
	}
}

func TestDateParse(t *testing.T) {
	testCases := map[string]Date{
		"2019-12-31": Date(time.Date(2019, 12, 31, 0, 0, 0, 0, time.UTC)),
		"2015-09-30": Date(time.Date(2015, 9, 30, 0, 0, 0, 0, time.UTC)),
	}

	for input, expectedResult := range testCases {
		actualResult := panicParseDate(input)
		assert.Equal(t, actualResult, expectedResult)
	}
}

func TestDateParsePanic(t *testing.T) {
	testCases := []string{
		"Jan 1, 2020",
		"12/31/2019",
		"31/01/2018",
		"08-25-2017",
	}

	for _, input := range testCases {
		assert.Panics(t, func() { panicParseDate(input) })
	}
}
