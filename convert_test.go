package commonservice

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type dummy struct {
}

func ptr(value string) *string {
	return &value
}

func toStr(value *time.Time) string {
	return value.Format(layout)
}

func TestToTimePtr(t *testing.T) {
	var anyTime = time.Now().Add(-2500 * time.Hour) // Why not
	var anyTimeAsString = anyTime.Format(layout)
	assert.Nil(t, ToTimePtr(nil))
	assert.Nil(t, ToTimePtr(dummy{}))
	assert.Equal(t, toStr(&anyTime), toStr(ToTimePtr(anyTime)))
	assert.Equal(t, toStr(&anyTime), toStr(ToTimePtr(&anyTime)))
	assert.Equal(t, toStr(&anyTime), toStr(ToTimePtr(anyTimeAsString)))
	assert.Equal(t, toStr(&anyTime), toStr(ToTimePtr(&anyTimeAsString)))
	assert.Equal(t, toStr(&anyTime), toStr(ToTimePtr(anyTime)))
}

func TestToStringPtr(t *testing.T) {
	var str = "abc"
	assert.Nil(t, ToStringPtr(nil))
	assert.Equal(t, "2", *ToStringPtr(2))
	assert.Equal(t, "2", *ToStringPtr(2.0))
	assert.Equal(t, "abc", *ToStringPtr(str))
	assert.Equal(t, "abc", *ToStringPtr(&str))
}

func TestToInt(t *testing.T) {
	var invalid = -99999
	var vInt = 2
	var vInt32 = int32(2)
	var vInt64 = int64(2)
	var vFloat32 = float32(2.0)
	var vFloat64 = float64(2.0)
	assert.Equal(t, 2, ToInt(ptr("2"), invalid))
	assert.Equal(t, 2, ToInt("2", invalid))
	assert.Equal(t, invalid, ToInt("2x", invalid))
	assert.Equal(t, 2, ToInt(vInt, invalid))
	assert.Equal(t, 2, ToInt(&vInt, invalid))
	assert.Equal(t, 2, ToInt(vInt32, invalid))
	assert.Equal(t, 2, ToInt(&vInt32, invalid))
	assert.Equal(t, 2, ToInt(vInt64, invalid))
	assert.Equal(t, 2, ToInt(&vInt64, invalid))
	assert.Equal(t, 2, ToInt(vFloat32, invalid))
	assert.Equal(t, 2, ToInt(&vFloat32, invalid))
	assert.Equal(t, 2, ToInt(vFloat64, invalid))
	assert.Equal(t, 2, ToInt(&vFloat64, invalid))
	assert.Equal(t, invalid, ToInt(time.Now(), invalid))
}

func TestToFloat(t *testing.T) {
	var invalid = float64(-99999)
	var expected = float64(3.14)
	var vInt = 3
	var vInt32 = int32(3)
	var vInt64 = int64(3)
	var vFloat32 = float32(3.14)
	var vFloat64 = float64(3.14)
	assert.Equal(t, invalid, ToFloat("xxx", invalid))
	assert.Equal(t, float64(3), ToFloat(vInt, invalid))
	assert.Equal(t, float64(3), ToFloat(&vInt, invalid))
	assert.Equal(t, float64(3), ToFloat(vInt32, invalid))
	assert.Equal(t, float64(3), ToFloat(&vInt32, invalid))
	assert.Equal(t, float64(3), ToFloat(vInt64, invalid))
	assert.Equal(t, float64(3), ToFloat(&vInt64, invalid))
	var equals = func(a, b float64) bool {
		return int64(a*100.0) == int64(b*100)
	}
	assert.True(t, equals(expected, ToFloat(ptr("3.14"), invalid)))
	assert.True(t, equals(expected, ToFloat("3.14", invalid)))
	assert.True(t, equals(expected, ToFloat(vFloat32, invalid)))
	assert.True(t, equals(expected, ToFloat(&vFloat32, invalid)))
	assert.True(t, equals(expected, ToFloat(vFloat64, invalid)))
	assert.True(t, equals(expected, ToFloat(&vFloat64, invalid)))
	assert.Equal(t, invalid, ToFloat(time.Now(), invalid))
}
