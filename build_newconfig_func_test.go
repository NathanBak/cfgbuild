package cfgbuild

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigFuncEnvVars(t *testing.T) {

	os.Setenv("MY_INT", "42")
	os.Setenv("MY_UINT", "142")
	os.Setenv("MY_FLOAT", "2.718")
	os.Setenv("MY_TIME", "2022-10-10T21:01:16+00:00")
	os.Setenv("MY_DURATION", "3s")
	os.Setenv("MY_BYTES", "secretPassword")
	os.Setenv("MY_STRING", "Nobody expects the Spanish Inquisition!")
	os.Setenv("MY_BOOL", "tRuE")

	cfg, err := NewConfig[*TestConfig]()

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

	assert.True(t, cfg.InitCalled)
	assert.True(t, cfg.ValidateCalled)
}
