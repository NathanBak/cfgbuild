package cfgbuild

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringArrays(t *testing.T) {

	os.Setenv("MY_STRINGS", "one,two,three,four,five")

	cfg, err := NewConfig[*TestArrayConfig]()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, 5, len(cfg.MyStrings))
	assert.Equal(t, "one", cfg.MyStrings[0])
	assert.Equal(t, "two", cfg.MyStrings[1])
	assert.Equal(t, "three", cfg.MyStrings[2])
	assert.Equal(t, "four", cfg.MyStrings[3])
	assert.Equal(t, "five", cfg.MyStrings[4])
}

type TestArrayConfig struct {
	MyStrings []string `envvar:"MY_STRINGS"`
	MyInts    []int    `envvar:"MY_INTS"`
}
