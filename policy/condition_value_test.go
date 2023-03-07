package policy

import (
	"errors"
	"testing"
)

func TestNewConditionValue(t *testing.T) {
	cases := []struct {
		name string
		in   *ConditionValue
		want string
	}{
		{
			name: "SingularString",
			in:   NewConditionValueString(true, "test"),
			want: `"test"`,
		},
		{
			name: "SingularNumber",
			in:   NewConditionValueFloat(true, 123),
			want: `123`,
		},
		{
			name: "SingularBool",
			in:   NewConditionValueBool(true, true),
			want: `true`,
		},
		{
			name: "SliceString",
			in:   NewConditionValueString(false, "test"),
			want: `["test"]`,
		},
		{
			name: "SliceNumber",
			in:   NewConditionValueFloat(false, 123),
			want: `[123]`,
		},
		{
			name: "SliceBool",
			in:   NewConditionValueBool(false, true),
			want: `[true]`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.in.MarshalJSON()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("got '%s', want '%s'", string(got), tc.want)
			}
		})
	}
}

func TestConditionValueJSON(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "SingularString",
			in:   `"test"`,
			want: `"test"`,
		},
		{
			name: "SingularNumber",
			in:   `123`,
			want: `123`,
		},
		{
			name: "SingularBool",
			in:   `true`,
			want: `true`,
		},
		{
			name: "SliceString",
			in:   `["test"]`,
			want: `["test"]`,
		},
		{
			name: "SliceNumber",
			in:   `[123]`,
			want: `[123]`,
		},
		{
			name: "SliceBool",
			in:   `[true]`,
			want: `[true]`,
		},
		{
			name: "EmptySlice",
			in:   `[]`,
			want: `[]`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var cv ConditionValue
			err := cv.UnmarshalJSON([]byte(tc.in))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			got, err := cv.MarshalJSON()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("got '%s' %x, want '%s' %x", string(got), got, tc.want, tc.want)
			}
		})
	}
}

func TestInvalidConditionValueJSON(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want error
	}{
		{
			name: "NullSlice",
			in:   `[null]`,
			want: errors.New(ErrorInvalidConditionValueSlice),
		},
		{
			name: "Null",
			in:   `null`,
			want: errors.New(ErrorInvalidConditionValue),
		},
		{
			name: "InvalidType",
			in:   `{"test": "test"}`,
			want: errors.New(ErrorInvalidConditionValue),
		},
		{
			name: "InvalidSliceType",
			in:   `[{"test": "test"}]`,
			want: errors.New(ErrorInvalidConditionValueSlice),
		},
		{
			name: "InvalidJSON",
			in:   `{`,
			want: errors.New("unexpected end of JSON input"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var cv ConditionValue
			err := cv.UnmarshalJSON([]byte(tc.in))
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if err.Error() != tc.want.Error() {
				t.Errorf("got '%s', want '%s'", err.Error(), tc.want.Error())
			}
		})
	}
}

func TestConditionValueValues(t *testing.T) {
	cases := []struct {
		name         string
		in           *ConditionValue
		wantStr      []string
		wantFloat    []float64
		wantBool     []bool
		wantSingular bool
	}{
		{
			name:         "SingularString",
			in:           NewConditionValueString(true, "true"),
			wantStr:      []string{"true"},
			wantFloat:    []float64{},
			wantBool:     []bool{},
			wantSingular: true,
		},
		{
			name:         "SingularNumber",
			in:           NewConditionValueFloat(true, 123),
			wantStr:      []string{},
			wantFloat:    []float64{123},
			wantBool:     []bool{},
			wantSingular: true,
		},
		{
			name:         "SingularBool",
			in:           NewConditionValueBool(false, true),
			wantStr:      []string{},
			wantFloat:    []float64{},
			wantBool:     []bool{true},
			wantSingular: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotStr, gotBool, gotFloat := tc.in.Values()
			if len(gotStr) != len(tc.wantStr) {
				t.Errorf("got '%d', want '%d'", len(gotStr), len(tc.wantStr))
			}
			for i, v := range gotStr {
				if v != tc.wantStr[i] {
					t.Errorf("got '%s', want '%s'", v, tc.wantStr[i])
				}
			}

			if len(gotFloat) != len(tc.wantFloat) {
				t.Errorf("got '%d', want '%d'", len(gotFloat), len(tc.wantFloat))
			}
			for i, v := range gotFloat {
				if v != tc.wantFloat[i] {
					t.Errorf("got '%f', want '%f'", v, tc.wantFloat[i])
				}
			}

			if len(gotBool) != len(tc.wantBool) {
				t.Errorf("got '%d', want '%d'", len(gotBool), len(tc.wantBool))
			}
			for i, v := range gotBool {
				if v != tc.wantBool[i] {
					t.Errorf("got '%t', want '%t'", v, tc.wantBool[i])
				}
			}

			if tc.in.IsSingular() != tc.wantSingular {
				t.Errorf("got '%t', want '%t'", tc.in.IsSingular(), tc.wantSingular)
			}
		})
	}
}

func TestConditionValueAdd(t *testing.T) {
	cases := []struct {
		name     string
		in       *ConditionValue
		addStr   []string
		addFloat []float64
		addBool  []bool
		wantErr  error
	}{
		{
			name:    "AddStringSuccess",
			in:      NewConditionValueString(true, "test"),
			addStr:  []string{"test2"},
			wantErr: nil,
		},
		{
			name:    "AddStringFailHasBool",
			in:      NewConditionValueBool(true, true),
			addStr:  []string{"test2"},
			wantErr: errors.New("Cannot add strings, ConditionValue has existing bools"),
		},
		{
			name:    "AddStringFailHasFloat",
			in:      NewConditionValueFloat(true, 123),
			addStr:  []string{"test2"},
			wantErr: errors.New("Cannot add strings, ConditionValue has existing floats"),
		},
		{
			name:     "AddFloatSuccess",
			in:       NewConditionValueFloat(true, 123),
			addFloat: []float64{456},
			wantErr:  nil,
		},
		{
			name:     "AddFloatFailHasBool",
			in:       NewConditionValueBool(true, true),
			addFloat: []float64{456},
			wantErr:  errors.New("Cannot add floats, ConditionValue has existing bools"),
		},
		{
			name:     "AddFloatFailHasString",
			in:       NewConditionValueString(true, "test"),
			addFloat: []float64{456},
			wantErr:  errors.New("Cannot add floats, ConditionValue has existing strings"),
		},
		{
			name:    "AddBoolSuccess",
			in:      NewConditionValueBool(true, true),
			addBool: []bool{false},
			wantErr: nil,
		},
		{
			name:    "AddBoolFailHasFloat",
			in:      NewConditionValueFloat(true, 123),
			addBool: []bool{false},
			wantErr: errors.New("Cannot add bools, ConditionValue has existing floats"),
		},
		{
			name:    "AddBoolFailHasString",
			in:      NewConditionValueString(true, "test"),
			addBool: []bool{false},
			wantErr: errors.New("Cannot add bools, ConditionValue has existing strings"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			if len(tc.addStr) != 0 {
				err = tc.in.AddString(tc.addStr...)
			} else if len(tc.addFloat) != 0 {
				err = tc.in.AddFloat(tc.addFloat...)
			} else if len(tc.addBool) != 0 {
				err = tc.in.AddBool(tc.addBool...)
			}

			if err == nil && tc.wantErr != nil {
				t.Fatalf("expected error, got nil")
			}
			if err != nil && tc.wantErr == nil {
				t.Fatalf("got error, expected nil")
			}
			if err != nil && tc.wantErr != nil {
				if err.Error() != tc.wantErr.Error() {
					t.Errorf("got '%s', want '%s'", err.Error(), tc.wantErr.Error())
				}
			}
		})
	}

}
