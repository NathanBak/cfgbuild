package cfgbuild

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigBuilderErrors(t *testing.T) {

	os.Setenv("MY_INT", "forty-two")
	os.Setenv("MY_UINT", "-42")
	os.Setenv("MY_FLOAT", "pi")
	os.Setenv("MY_TIME", "1999")
	os.Setenv("MY_DURATION", "3ly")
	os.Setenv("MY_BOOL", "supposition")

	tsts := []struct{ varName, varVal, expected string }{
		{"MY_INT", "forty-two", `error reading "MY_INT" (strconv.ParseInt: parsing "forty-two": invalid syntax)`},
		{"MY_UINT", "-42", `error reading "MY_UINT" (strconv.ParseUint: parsing "-42": invalid syntax)`},
		{"MY_FLOAT", "pi", `error reading "MY_FLOAT" (strconv.ParseFloat: parsing "pi": invalid syntax)`},
		{"MY_TIME", "1999", `error reading "MY_TIME" (parsing time "1999" as "2006-01-02T15:04:05Z07:00": cannot parse "" as "-")`},
		{"MY_DURATION", "3ly", `error reading "MY_DURATION" (time: unknown unit "ly" in duration "3ly")`},
		{"MY_BOOL", "supposition", `error reading "MY_BOOL" (string "supposition" is not a valid boolean value)`},
	}

	for i, tst := range tsts {
		os.Clearenv()
		os.Setenv(tst.varName, tst.varVal)

		b := Builder[*TestConfig]{}

		_, err := b.Build()
		assert.Error(t, err, i)
		assert.Equal(t, tst.expected, err.Error())
	}
}
