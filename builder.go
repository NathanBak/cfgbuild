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
	DefaultTagName           = "envvar"
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
	indent       string
	prefix       string
	// ListSeparator splits items in a list (slice).  Default is comma (,).
	ListSeparator string
	// TagName used to identify the environment variable name for a field.  Default is "envvar".
	TagName string
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

	// Don't Panic!
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = fmt.Errorf("builder panic:  %v", panicErr)
		}
	}()

	err = b.instantiateCfg()
	if err != nil {
		return b.cfg, err
	}
	b.printDebugf("building type %T", b.cfg)

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

func (b *Builder[T]) readEnvVars() error {
	b.printDebugFunctionStart()
	defer b.printDebugFunctionFinish()

	typ := reflect.TypeOf(b.cfg).Elem()
	value := reflect.ValueOf(b.cfg).Elem()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		var first rune
		for _, c := range field.Name {
			first = c
			break
		}
		if unicode.IsLower(first) {
			b.printDebugf("skipping %q because it is not a public field", field.Name)
			continue
		}

		tag := field.Tag.Get(b.getTagName())
		key := getTagKey(tag)
		if key == "-" {
			b.printDebugf("skipping field %q because env var key is set to \"-\"", field.Name)
			continue
		}
		if key == "" {
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
				TagName:           b.TagName,
				Uint8Lists:        b.Uint8Lists,
			}

			cb.prefix, _ = getTagAttribute(tag, "prefix")

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
			} else {
				b.printDebugf("no properties set for field %q", field.Name)
			}
		} else if envVarVal, ok := os.LookupEnv(b.prefix + key); ok {
			err := b.setFieldValue(field.Name, value.Field(i), envVarVal)
			if err != nil {
				return fmt.Errorf("error reading %q (%s)", b.prefix+key, err.Error())
			}
			b.printDebugf("set value for field %q", field.Name)
			b.setProps[key] = true
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
			unmarshaller, ok := vi.(encoding.TextUnmarshaler)
			if !ok {
				if !ok {
					unmarshaller, ok = v.Addr().Interface().(encoding.TextUnmarshaler)
				}
			}

			if ok {
				return unmarshaller.UnmarshalText([]byte(s))
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
		structField := typ.Field(i)
		tag := structField.Tag.Get(b.getTagName())
		key := getTagKey(tag)
		_, required := getTagAttribute(tag, "required")

		if key == "-" {
			continue
		}
		if required && !b.setProps[key] {
			missingRequired = append(missingRequired, key)
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

func (b *Builder[T]) setDefaults() error {
	b.printDebugFunctionStart()
	defer b.printDebugFunctionFinish()

	typ := reflect.TypeOf(b.cfg).Elem()
	value := reflect.ValueOf(b.cfg).Elem()

	for i := 0; i < typ.NumField(); i++ {
		structField := typ.Field(i)
		tag := structField.Tag.Get(b.getTagName())
		key := getTagKey(tag)

		if key == "-" {
			continue
		}

		if defaultVal, ok := getTagAttribute(tag, "default"); ok {
			err := b.setFieldValue(structField.Name, value.Field(i), defaultVal)
			if err != nil {
				return fmt.Errorf("error setting default value for %q (%s)", key, err.Error())
			}
		}
	}

	return nil
}

// getTagName returns the user-specified tag name or defaults to "envvar" if none is specified.
func (b *Builder[T]) getTagName() string {
	if b.TagName == "" {
		return DefaultTagName
	}
	return b.TagName
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

func getTagKey(tagVal string) string {
	return strings.Split(tagVal, ",")[0]
}

func getTagAttribute(tagVal, attributeName string) (string, bool) {
	prefix := attributeName + "="
	for _, a := range strings.Split(tagVal, ",") {
		if a == attributeName {
			return "", true
		}
		if strings.HasPrefix(a, prefix) {
			return strings.TrimPrefix(a, prefix), true
		}
	}
	return "", false
}
