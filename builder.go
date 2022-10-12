package cfgbuild

import (
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
				return err
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

	case reflect.TypeOf(time.Now()):
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(t))

	case reflect.TypeOf(time.Duration(3)):
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

		case reflect.String:
			v.SetString(s)

		default:
			return fmt.Errorf("unsupported type/kind \"%s/%s\"",
				v.Type().String(), v.Kind().String())
		}
	}
	return nil
}
