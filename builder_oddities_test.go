package cfgbuild

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"testing"
	"time"

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

	// net.IP implements the encoding.TextUnmarshaler interface
	MyIP net.IP `envvar:"-,default=192.168.0.42"`

	// time.Time implements the encoding.TextUnmarshaler interface
	MyTime time.Time `envvar:"-,default=2000-03-17T13:37:00Z"`

	MyURL        url.URL  `envvar:"MY_URL"`
	MyURLPointer *url.URL `envvar:"MY_URL_POINTER"`
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

	fmt.Println(time.Now().String())

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
	assert.Equal(t, "192.168.0.42", cfg.MyIP.String())
	assert.Equal(t, 17, cfg.MyTime.Local().Day())
	assert.Equal(t, time.March, cfg.MyTime.Local().Month())
	assert.Equal(t, 2000, cfg.MyTime.Local().Year())
}

func TestURL(t *testing.T) {
	os.Setenv("MY_URL", "https://www.nathanbak.com/?p=744")
	os.Setenv("MY_URL_POINTER", "https://www.nathanbak.com/?p=711")

	b := Builder[*TestOddConfig]{debug: true}
	cfg, err := b.Build()
	assert.NoError(t, err)

	expected1, err := url.Parse("https://www.nathanbak.com/?p=744")
	assert.NoError(t, err)

	expected2, err := url.Parse("https://www.nathanbak.com/?p=711")
	assert.NoError(t, err)

	assert.Equal(t, *expected1, cfg.MyURL)
	assert.Equal(t, expected2, cfg.MyURLPointer)

}
