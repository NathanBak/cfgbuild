package cfgbuild

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNestedStructs(t *testing.T) {

	os.Setenv("MY_INT", "42")
	os.Setenv("MY_STRING", "Nobody expects the Spanish Inquisition!")
	os.Setenv("MY_BOOL", "tRuE")

	os.Setenv("MY_CHILD_INT", "17")
	os.Setenv("MY_CHILD_STRING", "I didn't expect the Spanish Inquisition.")
	os.Setenv("MY_CHILD_BOOL", "FaLsE")

	os.Setenv("POINTER_MY_CHILD_INT", "16")
	os.Setenv("POINTER_MY_CHILD_STRING", "Fetch the comfy chair.")
	os.Setenv("POINTER_MY_CHILD_BOOL", "true")

	b := Builder[*TestParentConfig]{debug: true}

	cfg, err := b.Build()
	assert.NoError(t, err)

	assert.NotNil(t, cfg)

	assert.Equal(t, 42, cfg.MyInt)
	assert.Equal(t, "Nobody expects the Spanish Inquisition!", cfg.MyString)
	assert.True(t, cfg.MyBool)

	assert.Equal(t, 17, cfg.MyChild.MyInt)
	assert.Equal(t, "I didn't expect the Spanish Inquisition.", cfg.MyChild.MyString)
	assert.False(t, cfg.MyChild.MyBool)

	assert.Equal(t, 16, cfg.MyPointerChild.MyInt)
	assert.Equal(t, "Fetch the comfy chair.", cfg.MyPointerChild.MyString)
	assert.True(t, cfg.MyPointerChild.MyBool)
}

type TestParentConfig struct {
	MyInt          int    `envvar:"MY_INT"`
	MyString       string `envvar:"MY_STRING"`
	MyBool         bool   `envvar:"MY_BOOL"`
	MyChild        TestChildConfig
	MyPointerChild *TestChildConfig `envvar:",prefix=POINTER_"`
}

type TestChildConfig struct {
	MyInt    int    `envvar:"MY_CHILD_INT"`
	MyString string `envvar:"MY_CHILD_STRING"`
	MyBool   bool   `envvar:"MY_CHILD_BOOL"`
}
