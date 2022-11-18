package cfgbuild

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestOddConfig struct {
	// No associated environment variable, just set to default
	// Use - to designate fields with no environment variable
	DefaultOnly int `envvar:"-,default=7"`

	// If a nested config doesn't have anything special, set tag value to ">"
	// Nested configs can be regular values
	Nested TestNestedConfig `envvar:">"`

	// If a nested config needs a prefix, use tag with name ">" followed by prefix attribute
	// Nested configs can also be pointers
	AltNested *TestNestedConfig `envvar:">,prefix=ALT_"`

	// If a nested config shouldn't be initialized, do not tag
	IgnoredNested TestNestedConfig

	// Nested config (like any field) can also be ignored with "-"
	IgnoredNested2 TestNestedConfig `envvar:"-"`

	// Two different fields can be associated with the same environment variable
	MyInt     int `envvar:"MY_INT"`
	MySameInt int `envvar:"MY_INT"`

	// If not a config, field will not be initialized
	NotConfig *TestNotConfig

	// Non-public fields (starting with lowercase letter) are not set
	notPublic int `envvar:"MY_INT"`
}

type TestNestedConfig struct {
	MyVal string `envvar:"MY_VAL"`
}

type TestNotConfig struct {
	MyInt8 int8
}

func TestOddities(t *testing.T) {
	os.Setenv("MY_VAL", "my val")
	os.Setenv("ALT_MY_VAL", "alt my val")
	os.Setenv("MY_INT", "42")

	b := Builder[*TestOddConfig]{debug: true}
	cfg, err := b.Build()
	assert.NoError(t, err)

	assert.Equal(t, 7, cfg.DefaultOnly)
	assert.Equal(t, "my val", cfg.Nested.MyVal)
	assert.Equal(t, "alt my val", cfg.AltNested.MyVal)
	assert.Equal(t, "", cfg.IgnoredNested.MyVal)
	assert.Equal(t, "", cfg.IgnoredNested.MyVal)
	assert.Equal(t, 42, cfg.MyInt)
	assert.Equal(t, 42, cfg.MySameInt)
	assert.Nil(t, cfg.NotConfig)
	assert.Zero(t, cfg.notPublic)

}
