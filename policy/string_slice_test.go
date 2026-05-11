package policy

import (
	"testing"
)

func TestNewStringOrStringSlice(t *testing.T) {
	cases := []struct {
		name         string
		in           []string
		singular     bool
		want         string
		wantSingular bool
	}{
		{
			name:         "Singular",
			in:           []string{"arn:aws:iam::123456789012:root"},
			singular:     true,
			want:         `"arn:aws:iam::123456789012:root"`,
			wantSingular: true,
		},
		{
			name:         "SingleSlice",
			in:           []string{"arn:aws:iam::123456789012:root"},
			singular:     false,
			want:         `["arn:aws:iam::123456789012:root"]`,
			wantSingular: false,
		},
		{
			name:         "MultiSlice",
			in:           []string{"arn:aws:iam::111122223333:root", "arn:aws:iam::444455556666:root"},
			singular:     false,
			want:         `["arn:aws:iam::111122223333:root","arn:aws:iam::444455556666:root"]`,
			wantSingular: false,
		},
		{
			name:         "EmptySlice",
			in:           []string{},
			singular:     false,
			want:         `[]`,
			wantSingular: false,
		},
		{
			name:         "EmptyString",
			in:           []string{""},
			singular:     false,
			want:         `[""]`,
			wantSingular: false,
		},
		{
			name:         "EmptyStringSingular",
			in:           []string{""},
			singular:     true,
			want:         `""`,
			wantSingular: true,
		},
		{
			name:         "EmptyStringSlice",
			in:           []string{},
			singular:     true,
			want:         `[]`,
			wantSingular: true,
		},
		{
			name:         "IncorrectSingular",
			in:           []string{"arn:aws:iam::111122223333:root", "arn:aws:iam::444455556666:root"},
			singular:     true, // intentionally incorrect
			want:         `["arn:aws:iam::111122223333:root","arn:aws:iam::444455556666:root"]`,
			wantSingular: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ss := NewStringOrSlice(tc.singular, tc.in...)
			got, err := ss.MarshalJSON()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(got) != tc.want {
				t.Fatalf("got '%s', want '%s'", string(got), tc.want)
			}
			if ss.IsSingular() != tc.wantSingular {
				t.Fatalf("got '%t', want '%t'", ss.IsSingular(), tc.wantSingular)
			}
			if len(ss.Values()) != len(tc.in) {
				t.Fatalf("got '%d', want '%d'", len(ss.Values()), len(tc.in))
			}
			for i, v := range ss.Values() {
				if v != tc.in[i] {
					t.Fatalf("got '%s', want '%s'", v, tc.in[i])
				}
			}
		})
	}
}

func TestInvalidStringSliceJSON(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "NotSliceOfString",
			in:   `[{"foo": "bar"}]`,
			want: ErrorInvalidStringSlice,
		},
		{
			name: "InvalidString",
			in:   `123`,
			want: ErrorInvalidStringOrSlice,
		},
		{
			name: "InvalidJSON",
			in:   `{`,
			want: `unexpected end of JSON input`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var ss StringOrSlice
			err := ss.UnmarshalJSON([]byte(tc.in))
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if err.Error() != tc.want {
				t.Errorf("got '%s', want '%s'", err.Error(), tc.want)
			}
		})
	}
}

func TestStringOrSliceEqual(t *testing.T) {
	cases := []struct {
		name string
		a    *StringOrSlice
		b    *StringOrSlice
		want bool
	}{
		{
			name: "BothNil",
			a:    nil,
			b:    nil,
			want: true,
		},
		{
			name: "FirstNilSecondEmpty",
			a:    nil,
			b:    NewStringOrSlice(false),
			want: true,
		},
		{
			name: "FirstEmptySecondNil",
			a:    NewStringOrSlice(false),
			b:    nil,
			want: true,
		},
		{
			name: "FirstNilSecondNonEmpty",
			a:    nil,
			b:    NewStringOrSlice(true, "foo"),
			want: false,
		},
		{
			name: "FirstNonEmptySecondNil",
			a:    NewStringOrSlice(true, "foo"),
			b:    nil,
			want: false,
		},
		{
			name: "EqualSingular",
			a:    NewStringOrSlice(true, "s3:GetObject"),
			b:    NewStringOrSlice(true, "s3:GetObject"),
			want: true,
		},
		{
			name: "EqualSlice",
			a:    NewStringOrSlice(false, "s3:GetObject", "s3:PutObject"),
			b:    NewStringOrSlice(false, "s3:GetObject", "s3:PutObject"),
			want: true,
		},
		{
			name: "DifferentValues",
			a:    NewStringOrSlice(true, "s3:GetObject"),
			b:    NewStringOrSlice(true, "s3:PutObject"),
			want: false,
		},
		{
			name: "DifferentLength",
			a:    NewStringOrSlice(false, "s3:GetObject"),
			b:    NewStringOrSlice(false, "s3:GetObject", "s3:PutObject"),
			want: false,
		},
		{
			name: "BothEmpty",
			a:    NewStringOrSlice(false),
			b:    NewStringOrSlice(false),
			want: true,
		},
		{
			name: "DifferentOrder",
			a:    NewStringOrSlice(false, "a", "b"),
			b:    NewStringOrSlice(false, "b", "a"),
			want: true,
		},
		{
			name: "SingularEqualSlice",
			a:    NewStringOrSlice(true, "s3:GetObject"),
			b:    NewStringOrSlice(false, "s3:GetObject"),
			want: true,
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

func TestStringOrSliceAdd(t *testing.T) {
	cases := []struct {
		name string
		in   *StringOrSlice
		add  []string
		want []string
	}{

		{
			name: "Singular",
			in:   NewStringOrSlice(true, "arn:aws:iam::123456789012:root"),
			add:  []string{"arn:aws:iam::123456789012:root"},
			want: []string{"arn:aws:iam::123456789012:root", "arn:aws:iam::123456789012:root"},
		},
		{
			name: "Empty",
			in:   NewStringOrSlice(false, "arn:aws:iam::123456789012:root"),
			add:  []string{},
			want: []string{"arn:aws:iam::123456789012:root"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.in.Add(tc.add...)
			if len(tc.in.Values()) != len(tc.want) {
				t.Fatalf("got '%d', want '%d'", len(tc.in.Values()), len(tc.want))
			}
			for i, v := range tc.in.Values() {
				if v != tc.want[i] {
					t.Fatalf("got '%s', want '%s'", v, tc.want[i])
				}
			}
		})
	}
}
