package policy

import (
	"testing"
)

func TestNewPolicyPrincipal(t *testing.T) {
	cases := []struct {
		name       string
		in         *Principal
		want       string
		wantKind   string
		wantValues []string
	}{
		{
			name:       "All",
			in:         NewGlobalPrincipal(),
			want:       `"*"`,
			wantKind:   PrincipalKindAll,
			wantValues: []string{PrincipalAll},
		},
		{
			name:       "AWS",
			in:         NewAWSPrincipal("arn:aws:iam::123456789012:root"),
			want:       `{"AWS":"arn:aws:iam::123456789012:root"}`,
			wantKind:   PrincipalKindAWS,
			wantValues: []string{"arn:aws:iam::123456789012:root"},
		},
		{
			name:       "CanonicalUser",
			in:         NewCanonicalUserPrincipal("e01ebb0e05f2b447b372b56ced947c1a89bfe77ba79896972ff49ddfdbd0ecdd"),
			want:       `{"CanonicalUser":"e01ebb0e05f2b447b372b56ced947c1a89bfe77ba79896972ff49ddfdbd0ecdd"}`,
			wantKind:   PrincipalKindCanonical,
			wantValues: []string{"e01ebb0e05f2b447b372b56ced947c1a89bfe77ba79896972ff49ddfdbd0ecdd"},
		},
		{
			name:       "Federated",
			in:         NewFederatedPrincipal("arn:aws:iam::123456789012:oidc-provider/oidc.eks.us-west-2.amazonaws.com/id/EXAMPLED539D4633E53DE1B716D3041E"),
			want:       `{"Federated":"arn:aws:iam::123456789012:oidc-provider/oidc.eks.us-west-2.amazonaws.com/id/EXAMPLED539D4633E53DE1B716D3041E"}`,
			wantKind:   PrincipalKindFederated,
			wantValues: []string{"arn:aws:iam::123456789012:oidc-provider/oidc.eks.us-west-2.amazonaws.com/id/EXAMPLED539D4633E53DE1B716D3041E"},
		},
		{
			name:       "Services",
			in:         NewServicePrincipal("s3.amazonaws.com"),
			want:       `{"Service":"s3.amazonaws.com"}`,
			wantKind:   PrincipalKindService,
			wantValues: []string{"s3.amazonaws.com"},
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
			if tc.in.Kind() != tc.wantKind {
				t.Errorf("got '%s', want '%s'", tc.in.Kind(), tc.wantKind)
			}
			if len(tc.in.Values()) != len(tc.wantValues) {
				t.Errorf("got '%d', want '%d'", len(tc.in.Values()), len(tc.wantValues))
			}
			for i, v := range tc.in.Values() {
				if v != tc.wantValues[i] {
					t.Errorf("got '%s', want '%s'", v, tc.wantValues[i])
				}
			}
		})
	}
}

func TestPolicyPrincipalInvalidJSON(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "Empty",
			in:   `true`,
			want: `json: cannot unmarshal bool into Go value of type policy.principal`,
		},
		{
			name: "InvalidJSON",
			in:   `{`,
			want: `unexpected end of JSON input`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got Principal
			err := got.UnmarshalJSON([]byte(tc.in))
			if err != nil && err.Error() != tc.want {
				t.Errorf("got '%s', want '%s'", err.Error(), tc.want)
			}
		})
	}
}
