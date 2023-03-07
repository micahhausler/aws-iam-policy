package policy

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
)

const (
	ErrorInvalidStringOrSlice = "field neither slice of string or string"
	ErrorInvalidStringSlice   = "field not slice of string"
)

// NewStringOrSlice creates a new StringOrSlice. If singular is true and
// there is only one element, the structure will be marshaled as a string
// instead of a slice.
func NewStringOrSlice(singular bool, values ...string) *StringOrSlice {
	return &StringOrSlice{
		values:   values,
		singular: singular,
	}
}

// StringOrSlice is a type that can hold a string or a slice of strings.
// When unarshalling JSON, it will preserve whether the original value was
// a singular string or a slice of strings.
type StringOrSlice struct {
	values   []string
	singular bool
}

func (s *StringOrSlice) Add(value ...string) {
	s.values = append(s.values, value...)
	if len(s.values) != 1 {
		s.singular = false
	}
}

func (s *StringOrSlice) UnmarshalJSON(data []byte) error {
	var tmp interface{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	slice, ok := tmp.([]interface{})
	if ok {
		values := []string{}
		for _, item := range slice {
			if _, ok := item.(string); !ok {
				return errors.New(ErrorInvalidStringSlice)
			}
			values = append(values, item.(string))
			s.singular = false
		}
		s.values = values
		return nil
	}
	theString, ok := tmp.(string)
	if ok {
		s.values = []string{theString}
		s.singular = true
		return nil
	}
	return errors.New(ErrorInvalidStringOrSlice)
}

func (s *StringOrSlice) MarshalJSON() ([]byte, error) {
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)

	if s.singular && len(s.values) == 1 {
		err := enc.Encode(s.values[0])
		return []byte(strings.TrimSpace(buf.String())), err
	}
	err := enc.Encode(s.values)
	return []byte(strings.TrimSpace(buf.String())), err
}

// Values returns the values of the StringOrSlice.
func (s *StringOrSlice) Values() []string {
	return s.values
}

// Singular returns true if the StringOrSlice is a singular value and has zero
// or one value.
func (s *StringOrSlice) IsSingular() bool {
	return s.singular && len(s.values) <= 1
}
