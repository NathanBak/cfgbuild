package cfgbuild

import (
	"encoding"
	"errors"
	"fmt"
	"math/bits"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
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

// A Builder is able to create and initialize a Config.  After creating a Builder, run the Build()
// method.
type Builder[T interface{}] struct {
	cfg           T
	instantiated  bool
	setProps      map[string]bool
	Uint8Lists    bool
	ListSeparator string
	TagName       string
}

type initInterface interface {
	CfgBuildInit() error
}

type validateInterface interface {
	CfgBuildValidate() error
}

func (b *Builder[T]) Build() (cfg T, err error) {

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

	validator, ok := any(b.cfg).(validateInterface)
	if ok {
		err = validator.CfgBuildValidate()
	}
	return b.cfg, err
}

func (b *Builder[T]) readEnvVars() error {
	err := b.instantiateCfg()
	if err != nil {
		return err
	}

	typ := reflect.TypeOf(b.cfg).Elem()
	value := reflect.ValueOf(b.cfg).Elem()

	for i := 0; i < typ.NumField(); i++ {
		structField := typ.Field(i)
		tag := structField.Tag.Get(b.getTagName())
		split := strings.Split(tag, ",")
		key := "-"
		if len(split) > 0 {
			key = split[0]
		}
		if key == "-" {
			continue
		}

		if envVarVal, ok := os.LookupEnv(key); ok {
			err = b.setFieldValue(value.Field(i), envVarVal)
			if err != nil {
				return fmt.Errorf("error reading %q (%s)", key, err.Error())
			}
			b.setProps[key] = true
		}
	}

	return nil
}

func (b *Builder[T]) instantiateCfg() error {
	if !b.instantiated {
		typ := reflect.TypeOf(b.cfg)
		val := reflect.New(typ.Elem()).Interface().(T)
		b.cfg = val
		b.instantiated = true
	}
	return nil
}

func (b *Builder[T]) setFieldValue(v reflect.Value, s string) error {

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
	typ := reflect.TypeOf(b.cfg).Elem()
	missingRequired := []string{}

	for i := 0; i < typ.NumField(); i++ {
		structField := typ.Field(i)
		tag := structField.Tag.Get(b.getTagName())
		split := strings.Split(tag, ",")
		if len(split) > 0 {
			key := split[0]
			if key == "-" {
				continue
			}
			if strings.Contains(tag, ",required") && !b.setProps[key] {
				missingRequired = append(missingRequired, key)
			}
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

	typ := reflect.TypeOf(b.cfg).Elem()
	value := reflect.ValueOf(b.cfg).Elem()

	for i := 0; i < typ.NumField(); i++ {
		structField := typ.Field(i)
		tag := structField.Tag.Get(b.getTagName())
		split := strings.Split(tag, ",")

		if len(split) < 2 || split[0] == "-" {
			continue
		}

		for j := 1; j < len(split); j++ {
			if strings.HasPrefix(split[j], "default=") {
				defaultVal := strings.TrimPrefix(split[j], "default=")
				err := b.setFieldValue(value.Field(i), defaultVal)
				if err != nil {
					key := split[0]
					return fmt.Errorf("error setting default value for %q (%s)", key, err.Error())
				}
				break
			}
		}
	}

	return nil
}

// getTagName returns the user-specified tag name or defaults to "envvar" if none is specified.
func (b *Builder[T]) getTagName() string {
	if b.TagName == "" {
		return "envvar"
	}
	return b.TagName
}

func split(s, sep string) []string {
	if sep == "" {
		sep = ","
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
