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
