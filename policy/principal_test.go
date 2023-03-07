package policy

import (
	"testing"
)

func TestNewPrincipal(t *testing.T) {
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

func TestPrincipalInvalidJSON(t *testing.T) {
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

func TestPrincipalAdd(t *testing.T) {
	cases := []struct {
		name         string
		in           *Principal
		addAws       []string
		addService   []string
		addFederated []string
		addCanonical []string
		want         string
	}{
		{
			name:   "AWS",
			in:     NewAWSPrincipal("111122223333"),
			addAws: []string{"222233334444"},
			want:   `{"AWS":["111122223333","222233334444"]}`,
		},
		{
			name:       "Service",
			in:         NewServicePrincipal("s3.amazonaws.com"),
			addService: []string{"ec2.amazonaws.com"},
			want:       `{"Service":["s3.amazonaws.com","ec2.amazonaws.com"]}`,
		},
		{
			name:         "Federated",
			in:           NewFederatedPrincipal("arn:aws:iam::123456789012:oidc-provider/oidc.eks.us-east-1.amazonaws.com/id/EXAMPLED539D4633E53DE1B716D3041E"),
			addFederated: []string{"arn:aws:iam::123456789012:oidc-provider/oidc.eks.us-west-2.amazonaws.com/id/EXAMPLED539D4633E53DE1B716D3041E"},
			want:         `{"Federated":["arn:aws:iam::123456789012:oidc-provider/oidc.eks.us-east-1.amazonaws.com/id/EXAMPLED539D4633E53DE1B716D3041E","arn:aws:iam::123456789012:oidc-provider/oidc.eks.us-west-2.amazonaws.com/id/EXAMPLED539D4633E53DE1B716D3041E"]}`,
		},
		{
			name:         "Canonical",
			in:           NewCanonicalUserPrincipal("e01ebb0e05f2b447b372b56ced947c1a89bfe77ba79896972ff49ddfdbd0ecdd"),
			addCanonical: []string{"e01ebb0e05f2b447b372b56ced947c1a89bfe77ba79896972ff49ddfdbd0ecdd"},
			want:         `{"CanonicalUser":["e01ebb0e05f2b447b372b56ced947c1a89bfe77ba79896972ff49ddfdbd0ecdd","e01ebb0e05f2b447b372b56ced947c1a89bfe77ba79896972ff49ddfdbd0ecdd"]}`,
		},
		{
			name:       "MixedAWSAndService",
			in:         NewAWSPrincipal("111122223333"),
			addService: []string{"ec2.amazonaws.com"},
			want:       `{"AWS":"111122223333","Service":"ec2.amazonaws.com"}`,
		},
		{
			name:         "MixedAWSAndFederated",
			in:           NewAWSPrincipal("111122223333"),
			addFederated: []string{"arn:aws:iam::123456789012:oidc-provider/oidc.eks.us-east-1.amazonaws.com/id/EXAMPLED539D4633E53DE1B716D3041E"},
			want:         `{"AWS":"111122223333","Federated":"arn:aws:iam::123456789012:oidc-provider/oidc.eks.us-east-1.amazonaws.com/id/EXAMPLED539D4633E53DE1B716D3041E"}`,
		},
		{
			name:         "MixedAWSAndCanonical",
			in:           NewAWSPrincipal("111122223333"),
			addCanonical: []string{"e01ebb0e05f2b447b372b56ced947c1a89bfe77ba79896972ff49ddfdbd0ecdd"},
			want:         `{"AWS":"111122223333","CanonicalUser":"e01ebb0e05f2b447b372b56ced947c1a89bfe77ba79896972ff49ddfdbd0ecdd"}`,
		},
		{
			name:   "MixedServiceAndAWS",
			in:     NewServicePrincipal("s3.amazonaws.com"),
			addAws: []string{"222233334444"},
			want:   `{"AWS":"222233334444","Service":"s3.amazonaws.com"}`,
		},
		{
			name:   "MixedAllandAWS",
			in:     NewGlobalPrincipal(),
			addAws: []string{"222233334444"},
			want:   `"*"`,
		},
		{
			name:         "MixedAllandFederated",
			in:           NewGlobalPrincipal(),
			addFederated: []string{"arn:aws:iam::123456789012:oidc-provider/oidc.eks.us-east-1.amazonaws.com/id/EXAMPLED539D4633E53DE1B716D3041E"},
			want:         `"*"`,
		},
		{
			name:         "MixedAllandCanonical",
			in:           NewGlobalPrincipal(),
			addCanonical: []string{"e01ebb0e05f2b447b372b56ced947c1a89bfe77ba79896972ff49ddfdbd0ecdd"},
			want:         `"*"`,
		},
		{
			name:       "MixedAllandService",
			in:         NewGlobalPrincipal(),
			addService: []string{"ec2.amazonaws.com"},
			want:       `"*"`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			for _, v := range tc.addAws {
				tc.in.AddAWS(v)
			}
			for _, v := range tc.addService {
				tc.in.AddService(v)
			}
			for _, v := range tc.addFederated {
				tc.in.AddFederated(v)
			}
			for _, v := range tc.addCanonical {
				tc.in.AddCanonicalUser(v)
			}
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
