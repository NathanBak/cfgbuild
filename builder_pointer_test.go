package cfgbuild

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPointerConfigBuilderEnvVars(t *testing.T) {

	os.Setenv("MY_INT", "42")
	os.Setenv("MY_INT8", "8")
	os.Setenv("MY_INT16", "160")
	os.Setenv("MY_INT32", "3200")
	os.Setenv("MY_INT64", "64000")

	os.Setenv("MY_UINT", "142")
	os.Setenv("MY_UINT8", "18")
	os.Setenv("MY_UINT16", "1160")
	os.Setenv("MY_UINT32", "13200")
	os.Setenv("MY_UINT64", "164000")

	os.Setenv("MY_FLOAT32", "2.718")
	os.Setenv("MY_FLOAT64", "3.142")

	os.Setenv("MY_DURATION", "3s")
	os.Setenv("MY_STRING", "Nobody expects the Spanish Inquisition!")
	os.Setenv("MY_BOOL", "tRuE")

	b := Builder[*TestPointerConfig]{}

	cfg, err := b.Build()
	assert.NoError(t, err)

	assert.NotNil(t, cfg)

	assert.Equal(t, 42, *cfg.MyInt)
	assert.Equal(t, int8(8), *cfg.MyInt8)
	assert.Equal(t, int16(160), *cfg.MyInt16)
	assert.Equal(t, int32(3200), *cfg.MyInt32)
	assert.Equal(t, int64(64000), *cfg.MyInt64)

	assert.Equal(t, uint(142), *cfg.MyUInt)
	assert.Equal(t, uint8(18), *cfg.MyUInt8)
	assert.Equal(t, uint16(1160), *cfg.MyUInt16)
	assert.Equal(t, uint32(13200), *cfg.MyUInt32)
	assert.Equal(t, uint64(164000), *cfg.MyUInt64)

	assert.EqualValues(t, 2.718, *cfg.MyFloat32)
	assert.EqualValues(t, 3.142, *cfg.MyFloat64)

	assert.Equal(t, 3*time.Second, *cfg.MyDuration)
	assert.Equal(t, "Nobody expects the Spanish Inquisition!", *cfg.MyString)
	assert.True(t, *cfg.MyBool)
}

type TestPointerConfig struct {
	MyInt   *int   `envvar:"MY_INT"`
	MyInt8  *int8  `envvar:"MY_INT8"`
	MyInt16 *int16 `envvar:"MY_INT16"`
	MyInt32 *int32 `envvar:"MY_INT32"`
	MyInt64 *int64 `envvar:"MY_INT64"`

	MyUInt   *uint   `envvar:"MY_UINT"`
	MyUInt8  *uint8  `envvar:"MY_UINT8"`
	MyUInt16 *uint16 `envvar:"MY_UINT16"`
	MyUInt32 *uint32 `envvar:"MY_UINT32"`
	MyUInt64 *uint64 `envvar:"MY_UINT64"`

	MyFloat32 *float32 `envvar:"MY_FLOAT32"`
	MyFloat64 *float64 `envvar:"MY_FLOAT64"`

	MyDuration *time.Duration `envvar:"MY_DURATION"`
	MyString   *string        `envvar:"MY_STRING"`
	MyBool     *bool          `envvar:"MY_BOOL"`
}

func (cfg *TestPointerConfig) Init() error {
	return nil
}
