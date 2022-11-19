package cfgbuild

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateConfigTags(t *testing.T) {

	tst := func(cfg interface{}, expectedTagVal, expectedFieldName, expectedMsg string) {
		err := InitConfig(cfg)
		assert.Error(t, err)
		e, ok := err.(*TagSyntaxError)
		assert.True(t, ok, "error should be a TagSyntaxError")
		assert.Equal(t, "envvar", e.TagKey)
		assert.Equal(t, expectedTagVal, e.TagValue)
		assert.Equal(t, expectedFieldName, e.FieldName)
		assert.Equal(t, expectedMsg, e.msg)
	}

	tst(&struct {
		notPublic int `envvar:"-"`
	}{}, "-", "notPublic", "non-public fields may not have the tag set")

	tst(&struct {
		DashEnvVarName int `envvar:"-,required"`
	}{}, "-,required", "DashEnvVarName", `the "required" attribute is not allowed on "-" fields`)

	tst(&struct {
		NoEnvVarName int `envvar:",required"`
	}{}, ",required", "NoEnvVarName", `tag does not have the name attribute set`)

	tst(&struct {
		NestedConfig TestChildConfig `envvar:">,default=foo"`
	}{}, ">,default=foo", "NestedConfig",
		`the "default" attribute is not allowed on ">" nested config fields`)

	tst(&struct {
		NotMarshalJSON int `envvar:"-,unmarshalJSON"`
	}{}, "-,unmarshalJSON", "NotMarshalJSON",
		`field type does not support "unmarshalJSON" tag attribute`)

	tst(&struct {
		MyInt int `envvar:"-,prefix=CHILD_"`
	}{}, "-,prefix=CHILD_", "MyInt",
		`the "prefix" attribute is only allowed on ">" nested config fields`)

	tst(&struct {
		MyInt int `envvar:"-,ninja"`
	}{}, "-,ninja", "MyInt",
		`tag value contains non-existent attribute "ninja"`)

	tst(&struct {
		MyInt int `envvar:"MY_INT,default"`
	}{}, "MY_INT,default", "MyInt",
		`the "default" attribute requires a value`)

	tst(&struct {
		MyInt int `envvar:"MY_INT,required=sure"`
	}{}, "MY_INT,required=sure", "MyInt",
		`the "required" attribute may not have a value`)

}
