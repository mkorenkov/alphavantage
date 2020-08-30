package alphavantage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeExp(t *testing.T) {
	testCases := map[string]string{
		"":           "0000",
		"0":          "0000",
		"9":          "9000",
		"01":         "0100",
		"99":         "9900",
		"012":        "0120",
		"999":        "9990",
		"0123":       "0123",
		"9999":       "9999",
		"01234":      "",
		"99999":      "",
		"fourty two": "",
	}

	for input, expected := range testCases {
		assert.Equal(t, expected, normalizeExp(input))
	}
}

func TestMoneyParse(t *testing.T) {
	testCases := map[string]int64{
		"None":     0,
		".42":      4200,
		"0.0":      0,
		"0.1234":   1234,
		"-0.12":    -1200,
		"-1.0":     -10000,
		"100.0":    1000000,
		"28282828": 282828280000,
	}

	for input, expectedResult := range testCases {
		actualResult := panicParseMoney(input)
		assert.Equal(t, expectedResult, int64(actualResult))
	}
}

func TestMoneyParsePanic(t *testing.T) {
	testCases := []string{
		"fourty two",
		"one third",
		"127.0.0.1",
	}

	for _, input := range testCases {
		assert.Panics(t, func() { panicParseMoney(input) })
	}
}

func TestMoneyToString(t *testing.T) {
	testCases := map[string]string{
		"None":     "0.0000",
		".42":      "0.4200",
		"0.0":      "0.0000",
		"0.1234":   "0.1234",
		"-0.12":    "-0.1200",
		"-1.0":     "-1.0000",
		"100.0":    "100.0000",
		"28282828": "28282828.0000",
	}

	for input, expectedResult := range testCases {
		actualResult := panicParseMoney(input)
		assert.Equal(t, expectedResult, actualResult.String())
	}
}

func TestMoneyMarshalJSON(t *testing.T) {
	testCases := map[string]string{
		"None":     `"0.0000"`,
		".42":      `"0.4200"`,
		"0.0":      `"0.0000"`,
		"0.1234":   `"0.1234"`,
		"-0.12":    `"-0.1200"`,
		"-1.0":     `"-1.0000"`,
		"100.0":    `"100.0000"`,
		"28282828": `"28282828.0000"`,
	}

	for input, expectedResult := range testCases {
		actualResult, err := panicParseMoney(input).MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, expectedResult, string(actualResult))
	}
}
