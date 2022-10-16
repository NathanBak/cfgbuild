package cfgbuild

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfigBuilderDefaults(t *testing.T) {
	os.Unsetenv("MY_INT")
	os.Unsetenv("MY_DURATION")
	os.Unsetenv("MY_TIME")
	os.Unsetenv("MY_BYTES")

	b := Builder[*TestConfig]{}

	cfg, err := b.Build()
	assert.NoError(t, err)

	assert.NotNil(t, cfg)

	assert.Equal(t, 8081, cfg.MyInt)
	assert.Equal(t, time.Second*3, cfg.MyDuration)
	assert.Equal(t, 17, cfg.MyTime.Day())
	assert.Nil(t, cfg.MyBytes)

}

func TestConfigBuilderEnvVars(t *testing.T) {

	os.Setenv("MY_INT", "42")
	os.Setenv("MY_UINT", "142")
	os.Setenv("MY_FLOAT", "2.718")
	os.Setenv("MY_TIME", "2022-10-10T21:01:16+00:00")
	os.Setenv("MY_DURATION", "3s")
	os.Setenv("MY_BYTES", "secretPassword")
	os.Setenv("MY_STRING", "Nobody expects the Spanish Inquisition!")
	os.Setenv("MY_BOOL", "tRuE")

	b := Builder[*TestConfig]{}

	cfg, err := b.Build()
	assert.NoError(t, err)

	assert.NotNil(t, cfg)

	assert.Equal(t, 42, cfg.MyInt)
	assert.Equal(t, uint(142), cfg.MyUInt)
	assert.EqualValues(t, 2.718, cfg.MyFloat)
	assert.Equal(t, 16, cfg.MyTime.Second())
	assert.Equal(t, 3*time.Second, cfg.MyDuration)
	assert.Equal(t, []byte("secretPassword"), cfg.MyBytes)
	assert.Equal(t, "Nobody expects the Spanish Inquisition!", cfg.MyString)
	assert.True(t, cfg.MyBool)
}

type TestConfig struct {
	MyInt      int           `envvar:"MY_INT"`
	MyUInt     uint          `envvar:"MY_UINT"`
	MyFloat    float32       `envvar:"MY_FLOAT"`
	MyDuration time.Duration `envvar:"MY_DURATION"`
	MyTime     time.Time     `envvar:"MY_TIME"`
	MyBytes    []byte        `envvar:"MY_BYTES"`
	MyString   string        `envvar:"MY_STRING"`
	MyBool     bool          `envvar:"MY_BOOL"`
}

func (cfg *TestConfig) Init() error {

	cfg.MyInt = 8081
	cfg.MyDuration = 3 * time.Second
	cfg.MyTime = time.Date(2000, time.March, 17, 0, 13, 37, 0, time.UTC)

	return nil
}

func TestParseDuration(t *testing.T) {
	tsts := map[string]int64{
		"3s":        3000000000,
		"1m":        60000000000,
		"100ms":     100000000,
		"100000000": 100000000,
	}

	for s, expected := range tsts {
		assert.NoError(t, os.Setenv("MY_DURATION", s))
		cfg, err := (&Builder[*TestConfig]{}).Build()
		assert.NoError(t, err)
		assert.Equal(t, time.Duration(expected), cfg.MyDuration)
	}
}
