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

func TestStatementEqual(t *testing.T) {
	cases := []struct {
		name string
		a    *Statement
		b    *Statement
		want bool
	}{
		{
			name: "BothNil",
			a:    nil,
			b:    nil,
			want: true,
		},
		{
			name: "FirstNil",
			a:    nil,
			b:    &Statement{Effect: EffectAllow},
			want: false,
		},
		{
			name: "SecondNil",
			a:    &Statement{Effect: EffectAllow},
			b:    nil,
			want: false,
		},
		{
			name: "SameSimple",
			a: &Statement{
				Effect:   EffectAllow,
				Action:   NewStringOrSlice(true, "s3:GetObject"),
				Resource: NewStringOrSlice(true, "arn:aws:s3:::bucket/*"),
			},
			b: &Statement{
				Effect:   EffectAllow,
				Action:   NewStringOrSlice(true, "s3:GetObject"),
				Resource: NewStringOrSlice(true, "arn:aws:s3:::bucket/*"),
			},
			want: true,
		},
		{
			name: "DifferentEffect",
			a: &Statement{
				Effect: EffectAllow,
				Action: NewStringOrSlice(true, "s3:GetObject"),
			},
			b: &Statement{
				Effect: EffectDeny,
				Action: NewStringOrSlice(true, "s3:GetObject"),
			},
			want: false,
		},
		{
			name: "DifferentSid",
			a:    &Statement{Effect: EffectAllow, Sid: "1"},
			b:    &Statement{Effect: EffectAllow, Sid: "2"},
			want: false,
		},
		{
			name: "DifferentAction",
			a: &Statement{
				Effect: EffectAllow,
				Action: NewStringOrSlice(true, "s3:GetObject"),
			},
			b: &Statement{
				Effect: EffectAllow,
				Action: NewStringOrSlice(true, "s3:PutObject"),
			},
			want: false,
		},
		{
			name: "DifferentPrincipal",
			a: &Statement{
				Effect:    EffectAllow,
				Principal: NewAWSPrincipal("111122223333"),
			},
			b: &Statement{
				Effect:    EffectAllow,
				Principal: NewServicePrincipal("s3.amazonaws.com"),
			},
			want: false,
		},
		{
			name: "WithCondition",
			a: &Statement{
				Effect: EffectAllow,
				Condition: map[string]map[string]*ConditionValue{
					"StringEquals": {
						"aws:PrincipalOrgID": NewConditionValueString(true, "o-123456"),
					},
				},
			},
			b: &Statement{
				Effect: EffectAllow,
				Condition: map[string]map[string]*ConditionValue{
					"StringEquals": {
						"aws:PrincipalOrgID": NewConditionValueString(true, "o-123456"),
					},
				},
			},
			want: true,
		},
		{
			name: "DifferentConditionValue",
			a: &Statement{
				Effect: EffectAllow,
				Condition: map[string]map[string]*ConditionValue{
					"StringEquals": {
						"aws:PrincipalOrgID": NewConditionValueString(true, "o-111111"),
					},
				},
			},
			b: &Statement{
				Effect: EffectAllow,
				Condition: map[string]map[string]*ConditionValue{
					"StringEquals": {
						"aws:PrincipalOrgID": NewConditionValueString(true, "o-222222"),
					},
				},
			},
			want: false,
		},
		{
			name: "DifferentConditionKey",
			a: &Statement{
				Effect: EffectAllow,
				Condition: map[string]map[string]*ConditionValue{
					"StringEquals": {
						"aws:PrincipalOrgID": NewConditionValueString(true, "o-123456"),
					},
				},
			},
			b: &Statement{
				Effect: EffectAllow,
				Condition: map[string]map[string]*ConditionValue{
					"StringLike": {
						"aws:PrincipalOrgID": NewConditionValueString(true, "o-123456"),
					},
				},
			},
			want: false,
		},
		{
			name: "NilConditionVsEmpty",
			a:    &Statement{Effect: EffectAllow, Condition: nil},
			b:    &Statement{Effect: EffectAllow, Condition: map[string]map[string]*ConditionValue{}},
			want: true,
		},
		{
			name: "DifferentConditionLength",
			a: &Statement{
				Effect: EffectAllow,
				Condition: map[string]map[string]*ConditionValue{
					"StringEquals": {
						"aws:PrincipalOrgID": NewConditionValueString(true, "o-123456"),
					},
					"Bool": {
						"aws:SecureTransport": NewConditionValueString(true, "true"),
					},
				},
			},
			b: &Statement{
				Effect: EffectAllow,
				Condition: map[string]map[string]*ConditionValue{
					"StringEquals": {
						"aws:PrincipalOrgID": NewConditionValueString(true, "o-123456"),
					},
				},
			},
			want: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.a.Equal(tc.b)
			if got != tc.want {
				t.Errorf("got '%t', want '%t'", got, tc.want)
			}
		})
	}
}

func TestStatementOrSliceMarshalJSON(t *testing.T) {
	cases := []struct {
		name         string
		in           *StatementOrSlice
		want         string
		wantSingular bool
	}{
		{
			name: "SingleStatement",
			in: NewSingularStatementOrSlice(Statement{
				Effect:   EffectAllow,
				Action:   NewStringOrSlice(true, "s3:GetObject"),
				Resource: NewStringOrSlice(true, "arn:aws:s3:::my_corporate_bucket/exampleobject.png"),
			}),
			want:         `{"Action":"s3:GetObject","Effect":"Allow","Resource":"arn:aws:s3:::my_corporate_bucket/exampleobject.png"}`,
			wantSingular: true,
		},
		{
			name: "SliceStatement",
			in: NewStatementOrSlice(Statement{
				Effect:   EffectAllow,
				Action:   NewStringOrSlice(true, "s3:GetObject"),
				Resource: NewStringOrSlice(true, "arn:aws:s3:::my_corporate_bucket/exampleobject.png"),
			}),
			want:         `[{"Action":"s3:GetObject","Effect":"Allow","Resource":"arn:aws:s3:::my_corporate_bucket/exampleobject.png"}]`,
			wantSingular: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := tc.in.MarshalJSON()
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			if string(b) != tc.want {
				t.Errorf("got '%s', want '%s'", string(b), tc.want)
			}
			if tc.wantSingular != tc.in.Singular() {
				t.Errorf("got '%t', want '%t'", tc.in.Singular(), tc.wantSingular)
			}
		})
	}
}

func TestStatementOrSliceEqual(t *testing.T) {
	cases := []struct {
		name string
		a    *StatementOrSlice
		b    *StatementOrSlice
		want bool
	}{
		{
			name: "BothNil",
			a:    nil,
			b:    nil,
			want: true,
		},
		{
			name: "FirstNil",
			a:    nil,
			b:    NewStatementOrSlice(Statement{Effect: EffectAllow}),
			want: false,
		},
		{
			name: "SecondNil",
			a:    NewStatementOrSlice(Statement{Effect: EffectAllow}),
			b:    nil,
			want: false,
		},
		{
			name: "SameStatements",
			a: NewStatementOrSlice(
				Statement{Effect: EffectAllow, Action: NewStringOrSlice(true, "s3:GetObject")},
			),
			b: NewStatementOrSlice(
				Statement{Effect: EffectAllow, Action: NewStringOrSlice(true, "s3:GetObject")},
			),
			want: true,
		},
		{
			name: "DifferentOrder",
			a: NewStatementOrSlice(
				Statement{Sid: "1", Effect: EffectAllow},
				Statement{Sid: "2", Effect: EffectDeny},
			),
			b: NewStatementOrSlice(
				Statement{Sid: "2", Effect: EffectDeny},
				Statement{Sid: "1", Effect: EffectAllow},
			),
			want: true,
		},
		{
			name: "DifferentLength",
			a: NewStatementOrSlice(
				Statement{Effect: EffectAllow},
			),
			b: NewStatementOrSlice(
				Statement{Effect: EffectAllow},
				Statement{Effect: EffectDeny},
			),
			want: false,
		},
		{
			name: "DifferentStatements",
			a: NewStatementOrSlice(
				Statement{Effect: EffectAllow, Action: NewStringOrSlice(true, "s3:GetObject")},
			),
			b: NewStatementOrSlice(
				Statement{Effect: EffectAllow, Action: NewStringOrSlice(true, "s3:PutObject")},
			),
			want: false,
		},
		{
			name: "DuplicateStatements",
			a: NewStatementOrSlice(
				Statement{Sid: "1", Effect: EffectAllow},
				Statement{Sid: "1", Effect: EffectAllow},
			),
			b: NewStatementOrSlice(
				Statement{Sid: "1", Effect: EffectAllow},
				Statement{Sid: "1", Effect: EffectAllow},
			),
			want: true,
		},
		{
			name: "DuplicateVsUnique",
			a: NewStatementOrSlice(
				Statement{Sid: "1", Effect: EffectAllow},
				Statement{Sid: "1", Effect: EffectAllow},
			),
			b: NewStatementOrSlice(
				Statement{Sid: "1", Effect: EffectAllow},
				Statement{Sid: "2", Effect: EffectAllow},
			),
			want: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.a.Equal(tc.b)
			if got != tc.want {
				t.Errorf("got '%t', want '%t'", got, tc.want)
			}
		})
	}
}

func TestNilStatementOrSlice(t *testing.T) {
	cases := []struct {
		name string
		a    *StatementOrSlice
		want bool
	}{
		{
			name: "NotNilStatementOrSlice",
			a:    NewStatementOrSlice(),
			want: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.a.Values() == nil
			if got != tc.want {
				t.Errorf("got '%t', want '%t'", got, tc.want)
			}
		})
	}
}

func TestPolicyEqual(t *testing.T) {
	cases := []struct {
		name string
		a    *Policy
		b    *Policy
		want bool
	}{
		{
			name: "BothNil",
			a:    nil,
			b:    nil,
			want: true,
		},
		{
			name: "FirstNil",
			a:    nil,
			b:    &Policy{Version: VersionLatest},
			want: false,
		},
		{
			name: "SecondNil",
			a:    &Policy{Version: VersionLatest},
			b:    nil,
			want: false,
		},
		{
			name: "SamePolicy",
			a: &Policy{
				Version: VersionLatest,
				Statements: NewStatementOrSlice(
					Statement{Effect: EffectAllow, Action: NewStringOrSlice(true, "s3:GetObject")},
				),
			},
			b: &Policy{
				Version: VersionLatest,
				Statements: NewStatementOrSlice(
					Statement{Effect: EffectAllow, Action: NewStringOrSlice(true, "s3:GetObject")},
				),
			},
			want: true,
		},
		{
			name: "DifferentVersion",
			a: &Policy{
				Version:    Version2012_10_17,
				Statements: NewStatementOrSlice(Statement{Effect: EffectAllow}),
			},
			b: &Policy{
				Version:    Version2008_10_17,
				Statements: NewStatementOrSlice(Statement{Effect: EffectAllow}),
			},
			want: false,
		},
		{
			name: "DifferentId",
			a: &Policy{
				Id:         "policy-1",
				Version:    VersionLatest,
				Statements: NewStatementOrSlice(Statement{Effect: EffectAllow}),
			},
			b: &Policy{
				Id:         "policy-2",
				Version:    VersionLatest,
				Statements: NewStatementOrSlice(Statement{Effect: EffectAllow}),
			},
			want: false,
		},
		{
			name: "StatementsInDifferentOrder",
			a: &Policy{
				Version: VersionLatest,
				Statements: NewStatementOrSlice(
					Statement{Sid: "1", Effect: EffectAllow},
					Statement{Sid: "2", Effect: EffectDeny},
				),
			},
			b: &Policy{
				Version: VersionLatest,
				Statements: NewStatementOrSlice(
					Statement{Sid: "2", Effect: EffectDeny},
					Statement{Sid: "1", Effect: EffectAllow},
				),
			},
			want: true,
		},
		{
			name: "DifferentStatements",
			a: &Policy{
				Version: VersionLatest,
				Statements: NewStatementOrSlice(
					Statement{Effect: EffectAllow, Action: NewStringOrSlice(true, "s3:GetObject")},
				),
			},
			b: &Policy{
				Version: VersionLatest,
				Statements: NewStatementOrSlice(
					Statement{Effect: EffectAllow, Action: NewStringOrSlice(true, "s3:PutObject")},
				),
			},
			want: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.a.Equal(tc.b)
			if got != tc.want {
				t.Errorf("got '%t', want '%t'", got, tc.want)
			}
		})
	}
}
