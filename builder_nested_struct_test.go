package cfgbuild

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNestedStructs(t *testing.T) {

	os.Setenv("MY_INT", "42")
	os.Setenv("NESTED_JSON_CHILD", `{"b":true,"i":123,"s":"ahoy"}`)
	os.Setenv("DEFAULT_JSON_CHILD", `{"b":true,"s":"ahoy"}`)

	b := Builder[*TestNestedParentConfig]{debug: false, throwPanics: true}

	cfg, err := b.Build()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, 42, cfg.MyInt)

	assert.Equal(t, 123, cfg.NestedJSONChild.MyInt)
	assert.Equal(t, "ahoy", cfg.NestedJSONChild.MyString)
	assert.True(t, cfg.NestedJSONChild.MyBool)

	assert.Equal(t, 123, cfg.PointerJSONChild.MyInt)
	assert.Equal(t, "ahoy", cfg.PointerJSONChild.MyString)
	assert.True(t, cfg.PointerJSONChild.MyBool)

	assert.Equal(t, 3, cfg.DefaultJSONChild.MyInt)
	assert.Equal(t, "ahoy", cfg.DefaultJSONChild.MyString)
	assert.True(t, cfg.DefaultJSONChild.MyBool)

}

type TestNestedParentConfig struct {
	MyInt            int            `envvar:"MY_INT"`
	NestedJSONChild  TestJSONChild  `envvar:"NESTED_JSON_CHILD,unmarshalJSON"`
	PointerJSONChild *TestJSONChild `envvar:"NESTED_JSON_CHILD,unmarshalJSON,required"`
	DefaultJSONChild TestJSONChild  `envvar:"DEFAULT_JSON_CHILD,unmarshalJSON,default={\"i\":3}"`
}

type TestJSONChild struct {
	MyInt    int    `json:"i"`
	MyString string `json:"s"`
	MyBool   bool   `json:"b"`
}
