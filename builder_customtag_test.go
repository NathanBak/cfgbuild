package cfgbuild

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigBuilderCustomTag(t *testing.T) {

	os.Setenv("myInt", "42")
	os.Setenv("myString", "Nobody expects the Spanish Inquisition!")
	os.Setenv("myBool", "tRuE")

	b := Builder[*TestCustomTagConfig]{TagName: "voodoo"}

	cfg, err := b.Build()
	assert.NoError(t, err)

	assert.NotNil(t, cfg)

	assert.Equal(t, 42, cfg.MyInt)
	assert.Equal(t, "Nobody expects the Spanish Inquisition!", cfg.MyString)
	assert.True(t, cfg.MyBool)
	assert.Equal(t, 1234, cfg.MyDefaultInt)
}

type TestCustomTagConfig struct {
	MyInt        int    `voodoo:"myInt"`
	MyString     string `voodoo:"myString"`
	MyBool       bool   `voodoo:"myBool"`
	MyDefaultInt int    `voodoo:"myDefaultInt,default=1234"`
}
