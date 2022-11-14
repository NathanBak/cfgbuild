package cfgbuild

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringArrays(t *testing.T) {

	os.Setenv("MY_STRINGS", "one,two,three,four,five")

	cfg, err := NewConfig[*TestArrayConfig]()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, 5, len(cfg.MyStrings))
	assert.Equal(t, "one", cfg.MyStrings[0])
	assert.Equal(t, "two", cfg.MyStrings[1])
	assert.Equal(t, "three", cfg.MyStrings[2])
	assert.Equal(t, "four", cfg.MyStrings[3])
	assert.Equal(t, "five", cfg.MyStrings[4])
}

func TestIntegerArrays(t *testing.T) {

	os.Setenv("MY_INTS", " 1 ,2,3,4,5")
	os.Setenv("MY_INT8S", " 2,3,4,5,6")
	os.Setenv("MY_INT16S", "3,4,5,6,7")
	os.Setenv("MY_INT32S", "4,5,6,7,8")
	os.Setenv("MY_INT64S", "5,6,7,8,9")

	os.Setenv("MY_UINTS", "11,12,13,14,15")
	os.Setenv("MY_UINT8S", "12,13,14,15,16")
	os.Setenv("MY_UINT16S", "13,14,15,16,17")
	os.Setenv("MY_UINT32S", "14,15,16,17,18")
	os.Setenv("MY_UINT64S", "15,16,17,18,19")

	b := Builder[*TestArrayConfig]{Uint8Lists: true}

	cfg, err := b.Build()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, 5, len(cfg.MyInts))
	assert.Equal(t, 1, cfg.MyInts[0])
	assert.Equal(t, 2, cfg.MyInts[1])
	assert.Equal(t, 3, cfg.MyInts[2])
	assert.Equal(t, 4, cfg.MyInts[3])
	assert.Equal(t, 5, cfg.MyInts[4])

	assert.Equal(t, 5, len(cfg.MyInt8s))
	assert.EqualValues(t, 2, cfg.MyInt8s[0])
	assert.EqualValues(t, 3, cfg.MyInt8s[1])
	assert.EqualValues(t, 4, cfg.MyInt8s[2])
	assert.EqualValues(t, 5, cfg.MyInt8s[3])
	assert.EqualValues(t, 6, cfg.MyInt8s[4])

	assert.Equal(t, 5, len(cfg.MyInt16s))
	assert.EqualValues(t, 3, cfg.MyInt16s[0])
	assert.EqualValues(t, 4, cfg.MyInt16s[1])
	assert.EqualValues(t, 5, cfg.MyInt16s[2])
	assert.EqualValues(t, 6, cfg.MyInt16s[3])
	assert.EqualValues(t, 7, cfg.MyInt16s[4])

	assert.Equal(t, 5, len(cfg.MyInt32s))
	assert.EqualValues(t, 4, cfg.MyInt32s[0])
	assert.EqualValues(t, 5, cfg.MyInt32s[1])
	assert.EqualValues(t, 6, cfg.MyInt32s[2])
	assert.EqualValues(t, 7, cfg.MyInt32s[3])
	assert.EqualValues(t, 8, cfg.MyInt32s[4])

	assert.Equal(t, 5, len(cfg.MyInt64s))
	assert.EqualValues(t, 5, cfg.MyInt64s[0])
	assert.EqualValues(t, 6, cfg.MyInt64s[1])
	assert.EqualValues(t, 7, cfg.MyInt64s[2])
	assert.EqualValues(t, 8, cfg.MyInt64s[3])
	assert.EqualValues(t, 9, cfg.MyInt64s[4])

	assert.Equal(t, 5, len(cfg.MyUInts))
	assert.EqualValues(t, 11, cfg.MyUInts[0])
	assert.EqualValues(t, 12, cfg.MyUInts[1])
	assert.EqualValues(t, 13, cfg.MyUInts[2])
	assert.EqualValues(t, 14, cfg.MyUInts[3])
	assert.EqualValues(t, 15, cfg.MyUInts[4])

	assert.Equal(t, 5, len(cfg.MyUInt8s))
	assert.EqualValues(t, 12, cfg.MyUInt8s[0])
	assert.EqualValues(t, 13, cfg.MyUInt8s[1])
	assert.EqualValues(t, 14, cfg.MyUInt8s[2])
	assert.EqualValues(t, 15, cfg.MyUInt8s[3])
	assert.EqualValues(t, 16, cfg.MyUInt8s[4])

	assert.Equal(t, 5, len(cfg.MyUInt16s))
	assert.EqualValues(t, 13, cfg.MyUInt16s[0])
	assert.EqualValues(t, 14, cfg.MyUInt16s[1])
	assert.EqualValues(t, 15, cfg.MyUInt16s[2])
	assert.EqualValues(t, 16, cfg.MyUInt16s[3])
	assert.EqualValues(t, 17, cfg.MyUInt16s[4])

	assert.Equal(t, 5, len(cfg.MyUInt32s))
	assert.EqualValues(t, 14, cfg.MyUInt32s[0])
	assert.EqualValues(t, 15, cfg.MyUInt32s[1])
	assert.EqualValues(t, 16, cfg.MyUInt32s[2])
	assert.EqualValues(t, 17, cfg.MyUInt32s[3])
	assert.EqualValues(t, 18, cfg.MyUInt32s[4])

	assert.Equal(t, 5, len(cfg.MyUInt64s))
	assert.EqualValues(t, 15, cfg.MyUInt64s[0])
	assert.EqualValues(t, 16, cfg.MyUInt64s[1])
	assert.EqualValues(t, 17, cfg.MyUInt64s[2])
	assert.EqualValues(t, 18, cfg.MyUInt64s[3])
	assert.EqualValues(t, 19, cfg.MyUInt64s[4])
}

func TestIntegerArraysWithSemiSeparator(t *testing.T) {

	os.Setenv("MY_INTS", " 1 ;2;3;4;5")
	os.Setenv("MY_INT8S", " 2;3;4;5;6")
	os.Setenv("MY_INT16S", "3;4;5;6;7")
	os.Setenv("MY_INT32S", "4;5;6;7;8")
	os.Setenv("MY_INT64S", "5;6;7;8;9")

	os.Setenv("MY_UINTS", "11;12;13;14;15")
	os.Setenv("MY_UINT8S", "12;13;14;15;16")
	os.Setenv("MY_UINT16S", "13;14;15;16;17")
	os.Setenv("MY_UINT32S", "14;15;16;17;18")
	os.Setenv("MY_UINT64S", "15;16;17;18;19")

	b := Builder[*TestArrayConfig]{Uint8Lists: true, ListSeparator: ";"}

	cfg, err := b.Build()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, 5, len(cfg.MyInts))
	assert.Equal(t, 1, cfg.MyInts[0])
	assert.Equal(t, 2, cfg.MyInts[1])
	assert.Equal(t, 3, cfg.MyInts[2])
	assert.Equal(t, 4, cfg.MyInts[3])
	assert.Equal(t, 5, cfg.MyInts[4])

	assert.Equal(t, 5, len(cfg.MyInt8s))
	assert.EqualValues(t, 2, cfg.MyInt8s[0])
	assert.EqualValues(t, 3, cfg.MyInt8s[1])
	assert.EqualValues(t, 4, cfg.MyInt8s[2])
	assert.EqualValues(t, 5, cfg.MyInt8s[3])
	assert.EqualValues(t, 6, cfg.MyInt8s[4])

	assert.Equal(t, 5, len(cfg.MyInt16s))
	assert.EqualValues(t, 3, cfg.MyInt16s[0])
	assert.EqualValues(t, 4, cfg.MyInt16s[1])
	assert.EqualValues(t, 5, cfg.MyInt16s[2])
	assert.EqualValues(t, 6, cfg.MyInt16s[3])
	assert.EqualValues(t, 7, cfg.MyInt16s[4])

	assert.Equal(t, 5, len(cfg.MyInt32s))
	assert.EqualValues(t, 4, cfg.MyInt32s[0])
	assert.EqualValues(t, 5, cfg.MyInt32s[1])
	assert.EqualValues(t, 6, cfg.MyInt32s[2])
	assert.EqualValues(t, 7, cfg.MyInt32s[3])
	assert.EqualValues(t, 8, cfg.MyInt32s[4])

	assert.Equal(t, 5, len(cfg.MyInt64s))
	assert.EqualValues(t, 5, cfg.MyInt64s[0])
	assert.EqualValues(t, 6, cfg.MyInt64s[1])
	assert.EqualValues(t, 7, cfg.MyInt64s[2])
	assert.EqualValues(t, 8, cfg.MyInt64s[3])
	assert.EqualValues(t, 9, cfg.MyInt64s[4])

	assert.Equal(t, 5, len(cfg.MyUInts))
	assert.EqualValues(t, 11, cfg.MyUInts[0])
	assert.EqualValues(t, 12, cfg.MyUInts[1])
	assert.EqualValues(t, 13, cfg.MyUInts[2])
	assert.EqualValues(t, 14, cfg.MyUInts[3])
	assert.EqualValues(t, 15, cfg.MyUInts[4])

	assert.Equal(t, 5, len(cfg.MyUInt8s))
	assert.EqualValues(t, 12, cfg.MyUInt8s[0])
	assert.EqualValues(t, 13, cfg.MyUInt8s[1])
	assert.EqualValues(t, 14, cfg.MyUInt8s[2])
	assert.EqualValues(t, 15, cfg.MyUInt8s[3])
	assert.EqualValues(t, 16, cfg.MyUInt8s[4])

	assert.Equal(t, 5, len(cfg.MyUInt16s))
	assert.EqualValues(t, 13, cfg.MyUInt16s[0])
	assert.EqualValues(t, 14, cfg.MyUInt16s[1])
	assert.EqualValues(t, 15, cfg.MyUInt16s[2])
	assert.EqualValues(t, 16, cfg.MyUInt16s[3])
	assert.EqualValues(t, 17, cfg.MyUInt16s[4])

	assert.Equal(t, 5, len(cfg.MyUInt32s))
	assert.EqualValues(t, 14, cfg.MyUInt32s[0])
	assert.EqualValues(t, 15, cfg.MyUInt32s[1])
	assert.EqualValues(t, 16, cfg.MyUInt32s[2])
	assert.EqualValues(t, 17, cfg.MyUInt32s[3])
	assert.EqualValues(t, 18, cfg.MyUInt32s[4])

	assert.Equal(t, 5, len(cfg.MyUInt64s))
	assert.EqualValues(t, 15, cfg.MyUInt64s[0])
	assert.EqualValues(t, 16, cfg.MyUInt64s[1])
	assert.EqualValues(t, 17, cfg.MyUInt64s[2])
	assert.EqualValues(t, 18, cfg.MyUInt64s[3])
	assert.EqualValues(t, 19, cfg.MyUInt64s[4])
}

func TestFloatArrays(t *testing.T) {
	os.Clearenv()
	os.Setenv("MY_FLOAT32S", "3.14,1.4142,0.7071,1.732")
	os.Setenv("MY_FLOAT64S", "1.4142,0.7071,1.732,2.71828182846")

	b := Builder[*TestArrayConfig]{}
	cfg, err := b.Build()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, 4, len(cfg.MyFloat32s))
	assert.EqualValues(t, 3.14, cfg.MyFloat32s[0])
	assert.EqualValues(t, 1.4142, cfg.MyFloat32s[1])
	assert.EqualValues(t, 0.7071, cfg.MyFloat32s[2])
	assert.EqualValues(t, 1.732, cfg.MyFloat32s[3])

	assert.Equal(t, 4, len(cfg.MyFloat64s))
	assert.EqualValues(t, 1.4142, cfg.MyFloat64s[0])
	assert.EqualValues(t, 0.7071, cfg.MyFloat64s[1])
	assert.EqualValues(t, 1.732, cfg.MyFloat64s[2])
	assert.EqualValues(t, 2.71828182846, cfg.MyFloat64s[3])
}

type TestArrayConfig struct {
	MyStrings []string `envvar:"MY_STRINGS"`

	MyInts   []int   `envvar:"MY_INTS"`
	MyInt8s  []int8  `envvar:"MY_INT8S"`
	MyInt16s []int16 `envvar:"MY_INT16S"`
	MyInt32s []int32 `envvar:"MY_INT32S"`
	MyInt64s []int64 `envvar:"MY_INT64S"`

	MyUInts   []uint   `envvar:"MY_UINTS"`
	MyUInt8s  []uint8  `envvar:"MY_UINT8S"`
	MyUInt16s []uint16 `envvar:"MY_UINT16S"`
	MyUInt32s []uint32 `envvar:"MY_UINT32S"`
	MyUInt64s []uint64 `envvar:"MY_UINT64S"`

	MyFloat32s []float32 `envvar:"MY_FLOAT32S"`
	MyFloat64s []float64 `envvar:"MY_FLOAT64S"`
}
