package cfgbuild

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrefixFallback(t *testing.T) {

	os.Setenv("MY_INT", "42")
	os.Setenv("MY_STRING", "Nobody expects the Spanish Inquisition!")
	os.Setenv("MY_BOOL", "tRuE")

	b := Builder[*TestPrefixFallbackParentConfig]{debug: true, PrefixFallback: true}

	cfg, err := b.Build()
	assert.NoError(t, err)

	assert.NotNil(t, cfg)

	assert.Equal(t, 42, cfg.MyInt)
	assert.Equal(t, "Nobody expects the Spanish Inquisition!", cfg.MyString)
	assert.True(t, cfg.MyBool)

	assert.Equal(t, 42, cfg.MyChild.MyInt)
	assert.Equal(t, "Nobody expects the Spanish Inquisition!", cfg.MyChild.MyString)
	assert.True(t, cfg.MyChild.MyBool)
}

func TestNoPrefixFallbackWithoutFlag(t *testing.T) {

	os.Setenv("MY_INT", "42")
	os.Setenv("MY_STRING", "Nobody expects the Spanish Inquisition!")
	os.Setenv("MY_BOOL", "tRuE")

	b := Builder[*TestPrefixFallbackParentConfig]{debug: true}

	cfg, err := b.Build()
	assert.NoError(t, err)

	assert.NotNil(t, cfg)

	assert.Equal(t, 42, cfg.MyInt)
	assert.Equal(t, "Nobody expects the Spanish Inquisition!", cfg.MyString)
	assert.True(t, cfg.MyBool)

	assert.Equal(t, 0, cfg.MyChild.MyInt)
	assert.Equal(t, "", cfg.MyChild.MyString)
	assert.False(t, cfg.MyChild.MyBool)
}

func TestPartialPrefixFallback(t *testing.T) {

	os.Setenv("MY_INT", "42")
	os.Setenv("MY_STRING", "Nobody expects the Spanish Inquisition!")
	os.Setenv("MY_BOOL", "tRuE")

	os.Setenv("PREFIX_MY_STRING", "Fetch the comfy chair.")

	b := Builder[*TestPrefixFallbackParentConfig]{debug: true, PrefixFallback: true}

	cfg, err := b.Build()
	assert.NoError(t, err)

	assert.NotNil(t, cfg)

	assert.Equal(t, 42, cfg.MyInt)
	assert.Equal(t, "Nobody expects the Spanish Inquisition!", cfg.MyString)
	assert.True(t, cfg.MyBool)

	assert.Equal(t, 42, cfg.MyChild.MyInt)
	assert.Equal(t, "Fetch the comfy chair.", cfg.MyChild.MyString)
	assert.True(t, cfg.MyChild.MyBool)
}

type TestPrefixFallbackParentConfig struct {
	MyInt    int                           `envvar:"MY_INT"`
	MyString string                        `envvar:"MY_STRING"`
	MyBool   bool                          `envvar:"MY_BOOL"`
	MyChild  TestPrefixFallbackChildConfig `envvar:">,prefix=PREFIX_"`
}

type TestPrefixFallbackChildConfig struct {
	MyInt    int    `envvar:"MY_INT"`
	MyString string `envvar:"MY_STRING"`
	MyBool   bool   `envvar:"MY_BOOL"`
}
