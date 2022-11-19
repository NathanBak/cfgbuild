/*
BSD 2-Clause License

Copyright (c) 2022, Nathan Bak
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

 1. Redistributions of source code must retain the above copyright notice, this
    list of conditions and the following disclaimer.

 2. Redistributions in binary form must reproduce the above copyright notice,
    this list of conditions and the following disclaimer in the documentation
    and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/
package cfgbuild

import (
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"math/bits"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// NewConfig will create and initialize a Config of the provided type.
func NewConfig[T any]() (T, error) {
	b := Builder[T]{}
	return b.Build()
}

// The InitConfig function accepts an existing Config and will perform the initialization steps.
func InitConfig(cfg interface{}) error {
	b := Builder[interface{}]{cfg: cfg, instantiated: true}
	_, err := b.Build()
	return err
}

// Default values for builder configuration fields
const (
	DefaultTagKey            = "envvar"
	DefaultListSeparator     = ","
	DefaultKeyValueSeparator = ":"
)

// A Builder is able to create and initialize a Config.  After creating a Builder, run the Build()
// method.
type Builder[T interface{}] struct {
	cfg          T
	instantiated bool
	setProps     map[string]bool
	debug        bool
	throwPanics  bool
	indent       string
	prefix       string
	// ListSeparator splits items in a list (slice).  Default is comma (,).
	ListSeparator string
	// TagKey used to identify the field tag value to be used.  Default is "envvar".
	TagKey string
	// KeyValueSeparator splits keys and values for maps.  Default is colon (:)
	KeyValueSeparator string
	// Uint8Lists designates that []uint8 and []byte should be treated as a list (ie 1,2,3,4).  The
	// default is false meaning that value will be treated as a series of bytes.
	Uint8Lists bool
}

type initInterface interface {
	CfgBuildInit() error
}

type validateInterface interface {
	CfgBuildValidate() error
}

func (b *Builder[T]) Build() (cfg T, err error) {
	b.printDebugFunctionStart()
	defer b.printDebugFunctionFinish()

	if !b.throwPanics {
		// Don't Panic!
		defer func() {
			if panicErr := recover(); panicErr != nil {
				err = fmt.Errorf("builder panic:  %v", panicErr)
			}
		}()
	}

	err = b.instantiateCfg()
	if err != nil {
		return b.cfg, err
	}
	b.printDebugf("building type %T", b.cfg)

	err = b.validateCfgTags()
	if err != nil {
		return b.cfg, err
	}

	// If config has CfgBuildInit() function, run it.
	initter, ok := any(b.cfg).(initInterface)
	if ok {
		err = initter.CfgBuildInit()
		if err != nil {
			return b.cfg, err
		}
	}

	err = b.setDefaults()
	if err != nil {
		return b.cfg, err
	}

	b.setProps = make(map[string]bool)

	err = b.readEnvVars()
	if err != nil {
		return b.cfg, err
	}

	err = b.checkRequired()
	if err != nil {
		return b.cfg, err
	}

	// If config has a CfgBuildValidate() function, run it.
	validator, ok := any(b.cfg).(validateInterface)
	if ok {
		err = validator.CfgBuildValidate()
	}
	return b.cfg, err
}

func (b *Builder[T]) validateCfgTags() error {
	b.printDebugFunctionStart()
	defer b.printDebugFunctionFinish()

	typ := reflect.TypeOf(b.cfg).Elem()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Name
		tagValue, ok := field.Tag.Lookup(b.getTagKey())

		// Ignore field if tag value isn't set
		if !ok {
			continue
		}

		// Tags may not be set on non-public fields
		if !isPublicField(field) {
			msg := "non-public fields may not have the tag set"
			return &TagSyntaxError{
				FieldName: fieldName,
				TagKey:    b.getTagKey(),
				TagValue:  tagValue,
				msg:       msg,
			}
		}

		envVarName := getTagEnvVarName(tagValue)

		if envVarName == "" {
			msg := "tag does not have the name attribute set"
			return &TagSyntaxError{
				FieldName: fieldName,
				TagKey:    b.getTagKey(),
				TagValue:  tagValue,
				msg:       msg,
			}
		}

		_, defaultSet := getTagAttribute(tagValue, tagAttrDefault)
		if envVarName == ">" && defaultSet {
			msg := "the \"default\" attribute is not allowed on \">\" nested config fields"
			return &TagSyntaxError{
				FieldName: fieldName,
				TagKey:    b.getTagKey(),
				TagValue:  tagValue,
				msg:       msg,
			}
		}

		_, requiredSet := getTagAttribute(tagValue, tagAttrRequired)
		if envVarName == "-" && requiredSet {
			msg := "the \"required\" attribute is not allowed on \"-\" fields"
			return &TagSyntaxError{
				FieldName: fieldName,
				TagKey:    b.getTagKey(),
				TagValue:  tagValue,
				msg:       msg,
			}
		}

		_, marshalJSONSet := getTagAttribute(tagValue, tagAttrUnmarshalJSON)
		if marshalJSONSet {
			value := reflect.ValueOf(b.cfg).Elem()
			fieldVal := value.Field(i)
			fieldInterface := fieldVal.Addr().Interface()
			err := json.Unmarshal([]byte("{}"), fieldInterface)
			if err != nil {
				msg := "field type does not support \"unmarshalJSON\" tag attribute"
				return &TagSyntaxError{
					FieldName: fieldName,
					TagKey:    b.getTagKey(),
					TagValue:  tagValue,
					msg:       msg,
				}
			}
		}

		_, prefixSet := getTagAttribute(tagValue, tagAttrPrefix)
		if envVarName != ">" && prefixSet {
			msg := `the "prefix" attribute is only allowed on ">" nested config fields`
			return &TagSyntaxError{
				FieldName: fieldName,
				TagKey:    b.getTagKey(),
				TagValue:  tagValue,
				msg:       msg,
			}
		}

		attrNames := getTagAttributeNames(tagValue)
		for _, attrName := range attrNames {
			found := false
			for _, attr := range allTagAttr {
				if attrName == string(attr) {
					found = true
					break
				}
			}
			if !found && attrName != "" {
				msg := fmt.Sprintf(`tag value contains non-existent attribute %q`, attrName)
				return &TagSyntaxError{
					FieldName: fieldName,
					TagKey:    b.getTagKey(),
					TagValue:  tagValue,
					msg:       msg,
				}
			}
		}

		for _, attr := range allTagAttr {
			if _, found := getTagAttribute(tagValue, attr); found {
				if attr.hasValue() && !strings.Contains(tagValue, string(attr)+"=") {
					msg := fmt.Sprintf(`the %q attribute requires a value`, attr)
					return &TagSyntaxError{
						FieldName: fieldName,
						TagKey:    b.getTagKey(),
						TagValue:  tagValue,
						msg:       msg,
					}
				}

				if !attr.hasValue() && strings.Contains(tagValue, string(attr)+"=") {
					msg := fmt.Sprintf(`the %q attribute may not have a value`, attr)
					return &TagSyntaxError{
						FieldName: fieldName,
						TagKey:    b.getTagKey(),
						TagValue:  tagValue,
						msg:       msg,
					}
				}
			}
		}
	}
	return nil
}

func (b *Builder[T]) setDefaults() error {
	b.printDebugFunctionStart()
	defer b.printDebugFunctionFinish()
	return b.fieldLoop(true)
}

func (b *Builder[T]) readEnvVars() error {
	b.printDebugFunctionStart()
	defer b.printDebugFunctionFinish()
	return b.fieldLoop(false)
}

func (b *Builder[T]) fieldLoop(setDefault bool) error {
	b.printDebugFunctionStart()
	defer b.printDebugFunctionFinish()

	typ := reflect.TypeOf(b.cfg).Elem()
	value := reflect.ValueOf(b.cfg).Elem()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Name

		tagValue, ok := field.Tag.Lookup(b.getTagKey())
		if !ok {
			b.printDebugf("skipping %q because it does not have the %q tag set", fieldName,
				b.getTagKey())
			continue
		}

		envVarName := getTagEnvVarName(tagValue)

		if !setDefault && envVarName == "-" {
			b.printDebugf("skipping field %q because env var name is set to \"-\"", fieldName)
			continue
		}

		defaultVal, defaultAttributeSet := getTagAttribute(tagValue, tagAttrDefault)

		if setDefault && !defaultAttributeSet {
			continue
		}

		if envVarName == ">" {
			myTyp := value.Field(i).Type()
			myNew := reflect.New(myTyp)
			myVal := myNew.Interface()

			if strings.HasPrefix(myTyp.String(), "*") {
				myVal = myNew.Elem().Interface()
			}

			cb := Builder[interface{}]{
				cfg:               myVal,
				debug:             b.debug,
				indent:            b.indent,
				ListSeparator:     b.ListSeparator,
				KeyValueSeparator: b.KeyValueSeparator,
				TagKey:            b.TagKey,
				Uint8Lists:        b.Uint8Lists,
			}

			cb.prefix, _ = getTagAttribute(tagValue, tagAttrPrefix)

			ccfg, err := cb.Build()
			if err != nil {
				return err
			}

			rvo := reflect.ValueOf(ccfg)

			if len(cb.setProps) > 0 {
				if strings.HasPrefix(myTyp.String(), "*") {
					value.Field(i).Set(rvo)
				} else {
					ele := rvo.Elem()
					value.Field(i).Set(ele)
				}
				b.setProps[fieldName] = true
			} else {
				b.printDebugf("no properties set for field %q", fieldName)
			}
		} else {
			var valStr string
			if setDefault {
				valStr = defaultVal
			} else {
				if envVarVal, ok := os.LookupEnv(b.prefix + envVarName); ok {
					valStr = envVarVal
				} else {
					continue
				}
			}

			if _, tagFound := getTagAttribute(tagValue, tagAttrUnmarshalJSON); tagFound {
				fieldVal := value.Field(i)
				fieldInterface := fieldVal.Addr().Interface()
				err := json.Unmarshal([]byte(valStr), fieldInterface)
				if err != nil {
					return err
				}
				b.printDebugf("unmarshaled value for field %q", field.Name)

				if !setDefault {
					b.setProps[fieldName] = true
				}
			} else {

				err := b.setFieldValue(fieldName, value.Field(i), valStr)
				if err != nil {
					if setDefault {
						return fmt.Errorf("error setting default value for %q (%s)", b.prefix+envVarName, err.Error())
					}
					return fmt.Errorf("error reading %q (%s)", b.prefix+envVarName, err.Error())
				}
				b.printDebugf("set value for field %q", fieldName)
				if !setDefault {
					b.setProps[fieldName] = true
				}
			}
		}
	}

	return nil
}

func (b *Builder[T]) instantiateCfg() error {
	b.printDebugFunctionStart()
	defer b.printDebugFunctionFinish()
	if !b.instantiated {
		typ := reflect.TypeOf(b.cfg)
		val := reflect.New(typ.Elem()).Interface().(T)
		b.cfg = val
		b.instantiated = true
	}
	return nil
}

func (b *Builder[T]) setFieldValue(fieldName string, v reflect.Value, s string) error {
	b.printDebugFunctionStart()
	defer b.printDebugFunctionFinish()

	b.printDebugf("fieldName: %q\n%stype:      %q\n%skind:      %q\n%sstringval: %q\n%s",
		fieldName, b.indent,
		v.Type().String(), b.indent,
		v.Kind().String(), b.indent,
		s, b.indent,
	)

	if !v.CanAddr() {
		return errors.New("unable to obtain field address")
	}

	if !v.CanSet() {
		return errors.New("unable to set field value")
	}

	sep := b.ListSeparator

	switch v.Type() {

	case reflect.TypeOf(time.Now()): // Time
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(t))

	case reflect.TypeOf(time.Duration(3)): // Duration
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			d, err := time.ParseDuration(s)
			if err != nil {
				return err
			}
			i = int64(d)
		}
		v.SetInt(int64(i))

	case reflect.TypeOf([]string{}):
		vals := split(s, sep)
		v.Set(reflect.ValueOf(vals))

	case reflect.TypeOf([]int{}):
		vals, err := parseIntegers[int](s, sep, true, bits.UintSize)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(vals))

	case reflect.TypeOf([]int8{}):
		vals, err := parseIntegers[int8](s, sep, true, 8)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(vals))

	case reflect.TypeOf([]int16{}):
		vals, err := parseIntegers[int16](s, sep, true, 16)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(vals))

	case reflect.TypeOf([]int32{}):
		vals, err := parseIntegers[int32](s, sep, true, 32)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(vals))

	case reflect.TypeOf([]int64{}):
		vals, err := parseIntegers[int64](s, sep, true, 64)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(vals))

	case reflect.TypeOf([]uint{}):
		vals, err := parseIntegers[uint](s, sep, false, bits.UintSize)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(vals))

	case reflect.TypeOf([]uint8{}):
		if b.Uint8Lists {
			vals, err := parseIntegers[uint8](s, sep, false, 8)
			if err != nil {
				return err
			}
			v.Set(reflect.ValueOf(vals))
		} else {
			// be default we assume []uint8 to actually be []byte
			v.Set(reflect.ValueOf([]uint8(s)))
		}

	case reflect.TypeOf([]uint16{}):
		vals, err := parseIntegers[uint16](s, sep, false, 16)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(vals))

	case reflect.TypeOf([]uint32{}):
		vals, err := parseIntegers[uint32](s, sep, false, 32)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(vals))

	case reflect.TypeOf([]uint64{}):
		vals, err := parseIntegers[uint64](s, sep, false, 64)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(vals))

	case reflect.TypeOf([]float32{}):
		vals, err := parseFloats[float32](s, sep, 32)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(vals))

	case reflect.TypeOf([]float64{}):
		vals, err := parseFloats[float64](s, sep, 64)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(vals))

	case reflect.TypeOf(map[string]string{}):
		kvsep := b.KeyValueSeparator
		if kvsep == "" {
			kvsep = ":"
		}

		mp := make(map[string]string)
		pairs := split(s, sep)
		for _, pair := range pairs {
			kv := split(pair, kvsep)
			if len(kv) != 2 {
				return fmt.Errorf("key/value pair must contain exactly one %q separator", kvsep)
			}
			mp[kv[0]] = kv[1]
		}
		v.Set(reflect.ValueOf(mp))

	default:

		if v.CanInterface() {
			vi := v.Interface()
			textUnmarshaler, ok := vi.(encoding.TextUnmarshaler)
			if !ok {
				textUnmarshaler, ok = v.Addr().Interface().(encoding.TextUnmarshaler)
			}

			if ok {
				return textUnmarshaler.UnmarshalText([]byte(s))
			}
		}

		switch v.Kind() {

		case reflect.Bool:
			switch strings.ToUpper(s) {
			case "TRUE":
				v.SetBool(true)
			case "FALSE":
				v.SetBool(false)
			default:
				return fmt.Errorf("string %q is not a valid boolean value", s)
			}

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return err
			}
			if v.OverflowInt(i) {
				return errors.New("overflow error")
			}
			v.SetInt(i)

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			u, err := strconv.ParseUint(s, 10, 64)
			if err != nil {
				return err
			}
			if v.OverflowUint(u) {
				return errors.New("overflow error")
			}
			v.SetUint(u)

		case reflect.Float32, reflect.Float64:
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return err
			}
			if v.OverflowFloat(f) {
				return errors.New("overflow error")
			}
			v.SetFloat(f)

		case reflect.Pointer:
			switch v.Type().String() {
			case "*int":
				i64, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					return err
				}
				i := int(i64)
				v.Set(reflect.ValueOf(&i))

			case "*int8":
				i64, err := strconv.ParseInt(s, 10, 8)
				if err != nil {
					return err
				}
				i8 := int8(i64)
				v.Set(reflect.ValueOf(&i8))

			case "*int16":
				i64, err := strconv.ParseInt(s, 10, 16)
				if err != nil {
					return err
				}
				i16 := int16(i64)
				v.Set(reflect.ValueOf(&i16))

			case "*int32":
				i64, err := strconv.ParseInt(s, 10, 32)
				if err != nil {
					return err
				}
				i32 := int32(i64)
				v.Set(reflect.ValueOf(&i32))

			case "*int64":
				i64, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					return err
				}
				v.Set(reflect.ValueOf(&i64))

			case "*uint":
				ui64, err := strconv.ParseUint(s, 10, 64)
				if err != nil {
					return err
				}
				ui := uint(ui64)
				v.Set(reflect.ValueOf(&ui))

			case "*uint8":
				ui64, err := strconv.ParseUint(s, 10, 8)
				if err != nil {
					return err
				}
				ui8 := uint8(ui64)
				v.Set(reflect.ValueOf(&ui8))

			case "*uint16":
				ui64, err := strconv.ParseUint(s, 10, 16)
				if err != nil {
					return err
				}
				ui16 := uint16(ui64)
				v.Set(reflect.ValueOf(&ui16))

			case "*uint32":
				ui64, err := strconv.ParseUint(s, 10, 32)
				if err != nil {
					return err
				}
				ui32 := uint32(ui64)
				v.Set(reflect.ValueOf(&ui32))

			case "*uint64":
				ui64, err := strconv.ParseUint(s, 10, 64)
				if err != nil {
					return err
				}
				v.Set(reflect.ValueOf(&ui64))

			case "*float32":
				f64, err := strconv.ParseFloat(s, 32)
				if err != nil {
					return err
				}
				f32 := float32(f64)
				v.Set(reflect.ValueOf(&f32))

			case "*float64":
				f64, err := strconv.ParseFloat(s, 64)
				if err != nil {
					return err
				}
				v.Set(reflect.ValueOf(&f64))

			case "*time.Duration":
				i, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					d, err := time.ParseDuration(s)
					if err != nil {
						return err
					}
					i = int64(d)
				}
				d := time.Duration(i)
				v.Set(reflect.ValueOf(&d))

			case "*string":
				str := s
				v.Set(reflect.ValueOf(&str))

			case "*bool":
				var b bool
				switch strings.ToUpper(s) {
				case "TRUE":
					b = true
				case "FALSE":
					b = false
				default:
					return fmt.Errorf("string %q is not a valid boolean value", s)
				}
				v.Set(reflect.ValueOf(&b))
			}

		case reflect.String:
			v.SetString(s)

		default:
			return fmt.Errorf("unsupported type/kind \"%s/%s\"",
				v.Type().String(), v.Kind().String())
		}
	}
	return nil
}

// checkRequired looks at each field and ensures that each field with a "required" tag was
// previously set from an env var.  An error is returned if any required fields were not set.
func (b *Builder[T]) checkRequired() error {
	b.printDebugFunctionStart()
	defer b.printDebugFunctionFinish()
	typ := reflect.TypeOf(b.cfg).Elem()
	missingRequired := []string{}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Name
		tagValue, ok := field.Tag.Lookup(b.getTagKey())
		if !ok {
			continue
		}

		envVarName := getTagEnvVarName(tagValue)
		_, required := getTagAttribute(tagValue, tagAttrRequired)

		if envVarName == "-" {
			continue
		}
		if required && !b.setProps[fieldName] {
			missingRequired = append(missingRequired, fieldName)
		}
	}

	switch len(missingRequired) {
	case 0:
		return nil
	case 1:
		return fmt.Errorf("missing required var %q", missingRequired[0])
	default:
		return fmt.Errorf("missing required vars: %s", strings.Join(missingRequired, ","))
	}
}

// getTagKey returns the user-specified tag name or defaults to "envvar" if none is specified.
func (b *Builder[T]) getTagKey() string {
	if b.TagKey == "" {
		return DefaultTagKey
	}
	return b.TagKey
}

func (b *Builder[T]) printDebugFunctionStart() {
	if b.debug {
		pc, _, line, _ := runtime.Caller(1)
		fmt.Printf("%sRunning function %s [line %d]\n", b.indent, funcName(pc), line)
		b.indent += "> "
	}
}

func (b *Builder[T]) printDebugFunctionFinish() {
	if b.debug {
		pc, _, line, _ := runtime.Caller(1)
		b.indent = b.indent[2:]
		fmt.Printf("%sFinished running function %s [line %d]\n", b.indent, funcName(pc), line)
	}
}

func (b *Builder[T]) printDebugf(msg string, a ...any) {
	if b.debug {
		_, _, line, _ := runtime.Caller(1)
		fmt.Printf("%s%s [line %d]\n", b.indent, fmt.Sprintf(msg, a...), line)
	}
}

func funcName(pc uintptr) string {
	fn := runtime.FuncForPC(pc).Name()
	split := strings.Split(fn, ".")
	return split[len(split)-1] + "()"
}

func split(s, sep string) []string {
	if sep == "" {
		sep = DefaultListSeparator
	}
	vals := strings.Split(s, sep)
	out := []string{}
	for _, v := range vals {
		out = append(out, strings.TrimSpace(v))
	}
	return out
}

type integers interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

func parseIntegers[T integers](s, sep string, signed bool, bitsize int) ([]T, error) {
	vals := split(s, sep)
	integers := []T{}
	for _, v := range vals {
		if signed {
			i64, err := strconv.ParseInt(v, 10, bitsize)
			if err != nil {
				return integers, err
			}
			integers = append(integers, T(i64))
		} else {
			u64, err := strconv.ParseUint(v, 10, bitsize)
			if err != nil {
				return integers, err
			}
			integers = append(integers, T(u64))
		}
	}
	return integers, nil
}

type floats interface {
	float32 | float64
}

func parseFloats[T floats](s, sep string, bitsize int) ([]T, error) {
	vals := split(s, sep)
	floats := []T{}
	for _, v := range vals {

		f64, err := strconv.ParseFloat(v, bitsize)
		if err != nil {
			return floats, err
		}
		floats = append(floats, T(f64))

	}
	return floats, nil
}

func getTagEnvVarName(tagVal string) string {
	return strings.Split(tagVal, ",")[0]
}

// getTagAttribute looks at the tag value and returns the attribute value for the specified
// attribute name and a bool indicator as to whether or not the attribute exists in the tag value.
func getTagAttribute(tagVal string, attributeName tagAttr) (string, bool) {
	prefix := string(attributeName) + "="
	for _, a := range strings.Split(tagVal, ",") {
		if a == string(attributeName) {
			return "", true
		}
		if strings.HasPrefix(a, prefix) {
			return strings.TrimPrefix(a, prefix), true
		}
	}
	return "", false
}

type TagSyntaxError struct {
	FieldName string
	TagKey    string
	TagValue  string
	msg       string
}

func (e TagSyntaxError) Error() string {
	return e.msg
}

func isPublicField(f reflect.StructField) bool {
	var first rune
	for _, c := range f.Name {
		first = c
		break
	}
	return !unicode.IsLower(first)
}

type tagAttr string

const (
	tagAttrDefault       tagAttr = "default"
	tagAttrPrefix        tagAttr = "prefix"
	tagAttrRequired      tagAttr = "required"
	tagAttrUnmarshalJSON tagAttr = "unmarshalJSON"
)

var allTagAttr = []tagAttr{
	tagAttrDefault,
	tagAttrPrefix,
	tagAttrRequired,
	tagAttrUnmarshalJSON,
}

func (a tagAttr) hasValue() bool {
	switch a {
	case tagAttrDefault, tagAttrPrefix:
		return true
	default:
		return false
	}
}

func getTagAttributeNames(tagValue string) []string {
	attrs := []string{}
	first := true
	for _, a := range strings.Split(tagValue, ",") {
		if first {
			first = false
			continue
		}
		kv := strings.Split(a, "=")
		attrs = append(attrs, kv[0])
	}
	return attrs
}
