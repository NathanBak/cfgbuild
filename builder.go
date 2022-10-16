package cfgbuild

import (
	"encoding"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// A Builder is able to create and initialize a Config.  After creating a Builder, run the Build()
// method.
type Builder[T Config] struct {
	cfg          T
	instantiated bool
}

func (b *Builder[T]) Build() (T, error) {
	var err error

	err = b.instantiateCfg()
	if err != nil {
		return b.cfg, err
	}

	err = b.readEnvVars()
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
		key := structField.Tag.Get("envvar")
		if envVarVal, ok := os.LookupEnv(key); ok {
			err = setFieldValue(value.Field(i), envVarVal)
			if err != nil {
				return fmt.Errorf("error reading %q (%s)", key, err.Error())
			}
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
		return b.cfg.Init()
	}
	return nil
}

func setFieldValue(v reflect.Value, s string) error {

	if !v.CanAddr() {
		return errors.New("unable to obtain field address")
	}

	if !v.CanSet() {
		return errors.New("unable to set field value")
	}

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

	case reflect.TypeOf([]uint8{}): // we assume this to actually be []byte
		v.Set(reflect.ValueOf([]uint8(s)))

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
