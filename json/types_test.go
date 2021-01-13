package json

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestType struct {
	OptValue1 OptionalString `json:"value1"`
	OptValue2 OptionalString `json:"value2"`
}

func marshal(intf interface{}) string {
	var bytes, _ = json.Marshal(intf)
	return string(bytes)
}

func unmarshalTestType(JSON string) TestType {
	var value TestType
	_ = json.Unmarshal([]byte(JSON), &value)
	return value
}

func TestMarshalOptionalValue(t *testing.T) {
	t.Run("No value set", func(t *testing.T) {
		var testType = TestType{StringPtrToOptional(nil), StringPtrToOptional(nil)}
		var strTestType = marshal(testType)
		assert.Contains(t, strTestType, `"value1":null`)
		assert.Contains(t, strTestType, `"value2":null`)
	})

	t.Run("Defined true/false but Value is nil", func(t *testing.T) {
		var testType = TestType{OptValue1: OptionalString{Defined: false, Value: nil}, OptValue2: OptionalString{Defined: true, Value: nil}}
		var strTestType = marshal(testType)
		assert.Contains(t, strTestType, `"value1":null`)
		assert.Contains(t, strTestType, `"value2":null`)
	})

	var value = "value"
	t.Run("Defined true/false with a non nil value", func(t *testing.T) {
		var testType = TestType{OptValue1: OptionalString{Defined: false, Value: &value}, OptValue2: OptionalString{Defined: true, Value: &value}}
		var strTestType = marshal(testType)
		assert.Contains(t, strTestType, `"value1":null`)
		assert.Contains(t, strTestType, `"value2":"value"`)
	})

	t.Run("StringPtrToOptional vs StringToOptional", func(t *testing.T) {
		var testType = TestType{OptValue1: StringPtrToOptional(&value), OptValue2: StringToOptional(value)}
		var strTestType = marshal(testType)
		assert.Contains(t, strTestType, `"value1":"value"`)
		assert.Contains(t, strTestType, `"value2":"value"`)
	})
}

func TestUnmarshalOptionalValue(t *testing.T) {
	var testType = unmarshalTestType(`{"value1":null}`)
	assert.True(t, testType.OptValue1.Defined)
	assert.False(t, testType.OptValue2.Defined)

	testType = unmarshalTestType(`{"value1":"hello", "value2":"world"}`)
	assert.True(t, testType.OptValue1.Defined)
	assert.True(t, testType.OptValue2.Defined)
	assert.Equal(t, "hello", *testType.OptValue1.Value)
	assert.Equal(t, "world", *testType.OptValue2.Value)
}

func TestToValue(t *testing.T) {
	var value = "any value"
	var defaultValue = ""

	t.Run("Not defined", func(t *testing.T) {
		var notDefined = OptionalString{Defined: false, Value: &value}
		assert.Nil(t, notDefined.ToValue(defaultValue))
	})

	t.Run("Defined true but Value is nil", func(t *testing.T) {
		var definedAsNil = OptionalString{Defined: true, Value: nil}
		assert.Equal(t, &defaultValue, definedAsNil.ToValue(defaultValue))
	})

	t.Run("Defined true with non-nil Value", func(t *testing.T) {
		var defined = OptionalString{Defined: true, Value: &value}
		assert.Equal(t, &value, defined.ToValue(defaultValue))
	})
}
