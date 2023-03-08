package policy

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestStatementOrSliceConstructor(t *testing.T) {
	cases := []struct {
		name         string
		in           *StatementOrSlice
		add          []Statement
		want         []Statement
		wantSingular bool
	}{
		{
			name:         "SingleStatement",
			in:           NewSingularStatementOrSlice(Statement{Sid: "1", Effect: EffectAllow}),
			add:          []Statement{{Sid: "2", Effect: EffectDeny}},
			want:         []Statement{{Sid: "1", Effect: EffectAllow}, {Sid: "2", Effect: EffectDeny}},
			wantSingular: false,
		},
		{
			name:         "SingleStatement",
			in:           NewSingularStatementOrSlice(Statement{Sid: "1", Effect: EffectAllow}),
			add:          []Statement{},
			want:         []Statement{{Sid: "1", Effect: EffectAllow}},
			wantSingular: true,
		},
		{
			name:         "SliceStatement",
			in:           NewStatementOrSlice([]Statement{{Sid: "1", Effect: EffectAllow}, {Sid: "2", Effect: EffectDeny}}...),
			add:          []Statement{{Sid: "3", Effect: EffectAllow}},
			want:         []Statement{{Sid: "1", Effect: EffectAllow}, {Sid: "2", Effect: EffectDeny}, {Sid: "3", Effect: EffectAllow}},
			wantSingular: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, s := range tc.add {
				tc.in.Add(s)
			}
			if len(tc.want) != len(tc.in.Values()) {
				t.Errorf("got '%d', want '%d'", len(tc.in.Values()), len(tc.want))
				return
			}
			if !cmp.Equal(tc.want, tc.in.Values()) {
				t.Errorf("%s", cmp.Diff(tc.want, tc.in.Values()))
				return
			}
			if tc.wantSingular != tc.in.Singular() {
				t.Errorf("got '%t', want '%t'", tc.in.Singular(), tc.wantSingular)
			}
		})
	}

}

func TestDisallowUnknownFields(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		wantErr string
	}{
		{
			name: "AllowUnknownFieldsInPolicy",
			in: `{
				"Version": "2012-10-17",
				"NewField": "NewValue",
				"Statement": [
					{
						"Effect": "Allow",
						"Action": "s3:GetObject",
						"Resource": "arn:aws:s3:::my_corporate_bucket/exampleobject.png"
					}
				]
			}`,
			wantErr: `json: unknown field "NewField"`,
		},
		{
			name: "AllowUnknownFieldsInStatement",
			in: `{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Effect": "Allow",
						"Action": "s3:GetObject",
						"Resource": "arn:aws:s3:::my_corporate_bucket/exampleobject.png",
						"NewField": "NewValue"
					}
				]
			}`,
			wantErr: `StatementOrSlice is not a slice of statements: json: unknown field "NewField"`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var p Policy
			decoder := json.NewDecoder(bytes.NewBufferString(tc.in))
			decoder.DisallowUnknownFields()
			err := decoder.Decode(&p)
			if err == nil {
				t.Fatalf("expect error, got none")
			}
			if err.Error() != tc.wantErr {
				t.Fatalf("expect error %q, got %q", tc.wantErr, err)
			}
		})
	}
}

func TestStatementOrSliceUnmarshalJSON(t *testing.T) {
	cases := []struct {
		name         string
		in           string
		want         []Statement
		wantSingular bool
		wantErr      string
	}{
		{
			name: "SingleStatement",
			in: `{
				"Effect": "Allow",
				"Action": "s3:GetObject",
				"Resource": "arn:aws:s3:::my_corporate_bucket/exampleobject.png",
				"Principal": {
					"AWS": "123456789012"
				}
			}`,
			want: []Statement{
				{
					Effect:    EffectAllow,
					Action:    NewStringOrSlice(true, "s3:GetObject"),
					Resource:  NewStringOrSlice(true, "arn:aws:s3:::my_corporate_bucket/exampleobject.png"),
					Principal: NewAWSPrincipal("123456789012"),
				},
			},
			wantSingular: true,
		},
		{
			name: "SliceStatement",
			in: `[
				{
					"Effect": "Allow",
					"Action": "s3:GetObject",
					"Resource": "arn:aws:s3:::my_corporate_bucket/exampleobject.png",
					"Principal": {
						"AWS": "123456789012"
					}
				}
			]`,
			want: []Statement{
				{
					Effect:    EffectAllow,
					Action:    NewStringOrSlice(true, "s3:GetObject"),
					Resource:  NewStringOrSlice(true, "arn:aws:s3:::my_corporate_bucket/exampleobject.png"),
					Principal: NewAWSPrincipal("123456789012"),
				},
			},
			wantSingular: false,
		},
		{
			name: "InvalidJSON",
			in: `{
				"Effect": "Allow",
				"Action": "s3:GetObject",
				`,
			wantErr:      "unexpected end of JSON input",
			wantSingular: false,
		},
		{
			name:    "BooleanJSON",
			in:      `true`,
			wantErr: ErrorInvalidStatementOrSlice,
		},
		{
			name:    "BadJSON",
			in:      `{`,
			wantErr: `unexpected end of JSON input`,
		},
		{
			name: "InvalidList",
			in: `[
				{
					"Effect": "Allow",
					"NotAField": "s3:GetObject"
				}
			]`,
			wantErr:      `StatementOrSlice is not a slice of statements: json: unknown field "NotAField"`,
			wantSingular: false,
		},
		{
			name: "InvalidStatement",
			in: `{
				"Effect": "Allow",
				"NotAField": "s3:GetObject"
			}`,
			wantErr:      `StatementOrSlice must be a single Statement or a slice of Statements: json: unknown field "NotAField"`,
			wantSingular: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var s StatementOrSlice
			err := s.UnmarshalJSON([]byte(tc.in))
			if err != nil {
				if tc.wantErr == "" {
					t.Fatalf("expect no error, got %v", err)
				}
				if err.Error() != tc.wantErr {
					t.Fatalf("expect error %q, got %q", tc.wantErr, err)
				}
				return
			}
			if len(tc.want) != len(s.Values()) {
				t.Errorf("got '%d', want '%d'", len(s.Values()), len(tc.want))
				return
			}
			if !cmp.Equal(tc.want, s.Values(), cmpopts.IgnoreUnexported(StringOrSlice{}, Principal{})) {
				t.Errorf("%s", cmp.Diff(tc.want, s.Values(), cmpopts.IgnoreUnexported(StringOrSlice{}, Principal{})))
				return
			}
			if tc.wantSingular != s.Singular() {
				t.Errorf("got '%t', want '%t'", s.Singular(), tc.wantSingular)
			}
		})
	}
}
