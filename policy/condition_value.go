package policy

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
)

const (
	ErrorInvalidConditionValueSlice = "field not slice of string, bool or float64"
	ErrorInvalidConditionValue      = "field neither slice of string, bool, or float64 or string, bool or float64"
)

// NewConditionValueString creates a new ConditionValue. If singular is true and
// there is only one element, the structure will be marshaled as a string
// instead of a slice.
func NewConditionValueString(singular bool, values ...string) *ConditionValue {
	return &ConditionValue{
		strValues: values,
		singular:  singular,
	}
}

// NewConditionValueBool creates a new ConditionValue. If singular is true and
// there is only one element, the structure will be marshaled as a bool
// instead of a slice.
func NewConditionValueBool(singular bool, values ...bool) *ConditionValue {
	return &ConditionValue{
		boolValues: values,
		singular:   singular,
	}
}

// NewConditionValueFloat creates a new ConditionValue. If singular is true and
// there is only one element, the structure will be marshaled as an float64
// instead of a slice.
func NewConditionValueFloat(singular bool, values ...float64) *ConditionValue {
	return &ConditionValue{
		numValues: values,
		singular:  singular,
	}
}

// ConditionValue is a type that can hold an indivual or slice of string, bool or float64.
// When unarshalling JSON, it will preserve whether the original value was
// singular or a slice.
type ConditionValue struct {
	strValues  []string
	boolValues []bool
	numValues  []float64
	singular   bool
}

// AddStrings adds a slice of strings to the ConditionValue. If the
// ConditionValue already has bools or floats, an error is returned.
func (c *ConditionValue) AddString(values ...string) error {
	if len(c.boolValues) > 0 {
		return errors.New("Cannot add strings, ConditionValue has existing bools")
	} else if len(c.numValues) > 0 {
		return errors.New("Cannot add strings, ConditionValue has existing floats")
	}
	c.strValues = append(c.strValues, values...)
	if len(c.strValues) > 1 {
		c.singular = false
	}
	return nil
}

// AddBool adds a slice of bools to the ConditionValue. If the
// ConditionValue already has strings or floats, an error is returned.
func (c *ConditionValue) AddBool(values ...bool) error {
	if len(c.strValues) > 0 {
		return errors.New("Cannot add bools, ConditionValue has existing strings")
	} else if len(c.numValues) > 0 {
		return errors.New("Cannot add bools, ConditionValue has existing floats")
	}
	c.boolValues = append(c.boolValues, values...)
	if len(c.boolValues) > 1 {
		c.singular = false
	}
	return nil
}

// AddFloat adds a slice of floats to the ConditionValue. If the
// ConditionValue already has strings or bools, an error is returned.
func (c *ConditionValue) AddFloat(values ...float64) error {
	if len(c.strValues) > 0 {
		return errors.New("Cannot add floats, ConditionValue has existing strings")
	} else if len(c.boolValues) > 0 {
		return errors.New("Cannot add floats, ConditionValue has existing bools")
	}
	c.numValues = append(c.numValues, values...)
	if len(c.numValues) > 1 {
		c.singular = false
	}
	return nil
}

func (c *ConditionValue) UnmarshalJSON(data []byte) error {
	var tmp interface{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	theString, ok := tmp.(string)
	if ok {
		c.strValues = []string{theString}
		c.singular = true
		return nil
	}
	theBool, ok := tmp.(bool)
	if ok {
		c.boolValues = []bool{theBool}
		c.singular = true
		return nil
	}
	theFloat, ok := tmp.(float64)
	if ok {
		c.numValues = []float64{theFloat}
		c.singular = true
		return nil
	}

	slice, ok := tmp.([]interface{})
	if ok {
		strValues := []string{}
		boolValues := []bool{}
		numValues := []float64{}
		for _, item := range slice {
			switch item.(type) {
			case string:
				strValues = append(strValues, item.(string))
			case bool:
				boolValues = append(boolValues, item.(bool))
			case float64: // all numbers are float64
				numValues = append(numValues, item.(float64))
			default:
				return errors.New(ErrorInvalidConditionValueSlice)
			}
			c.singular = false
		}
		c.strValues = strValues
		c.boolValues = boolValues
		c.numValues = numValues
		return nil
	}

	return errors.New(ErrorInvalidConditionValue)
}

func (c *ConditionValue) MarshalJSON() ([]byte, error) {
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	var err error
	if c.singular {
		if len(c.strValues) > 0 {
			err = enc.Encode(c.strValues[0])
			goto eoCV
		}
		if len(c.boolValues) > 0 {
			err = enc.Encode(c.boolValues[0])
			goto eoCV
		}
		if len(c.numValues) > 0 {
			err = enc.Encode(c.numValues[0])
			goto eoCV
		}
	}
	if len(c.strValues) > 0 {
		err = enc.Encode(c.strValues)
		goto eoCV
	}
	if len(c.boolValues) > 0 {
		err = enc.Encode(c.boolValues)
		goto eoCV
	}
	err = enc.Encode(c.numValues)

eoCV:
	return []byte(strings.TrimSpace(buf.String())), err
}

// Values returns the values of the ConditionValue.
func (c *ConditionValue) Values() ([]string, []bool, []float64) {
	return c.strValues, c.boolValues, c.numValues
}

// IsSingular returns true if the ConditionValue is a singular value and has
// zero or one elements.
func (c *ConditionValue) IsSingular() bool {
	return c.singular && (len(c.strValues) <= 1) && (len(c.boolValues) <= 1) && (len(c.numValues) <= 1)
}
