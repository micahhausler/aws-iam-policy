package policy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseValidS3BucketPolicy(t *testing.T) {
	fixture := "./test_fixtures/valid_bucket_policies.json"

	reader, err := os.ReadFile(fixture)
	if err != nil {
		t.Fatal(err)
	}

	policies := []interface{}{}
	err = json.Unmarshal(reader, &policies)
	if err != nil {
		t.Fatal(err)
	}

	for i, policyIface := range policies {
		t.Run(fmt.Sprintf("Unmarshal S3 policy %d", i), func(t *testing.T) {
			// Yes we re-marshal the policy here, but this is the only way to
			// get the json package to use the MarshalJSON() method on each
			// Principal struct individually.
			b, err := json.MarshalIndent(policyIface, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
			// b, _ := json.Marshal(policyIface,)
			decoder := json.NewDecoder(bytes.NewBuffer(b))
			decoder.DisallowUnknownFields()
			policy := &Policy{}
			err = decoder.Decode(policy)
			if err != nil {
				testName := struct {
					Id string `json:"Id"`
				}{}
				uerr := json.Unmarshal(b, &testName)
				if uerr != nil {
					t.Fatal(uerr)
				}
				t.Fatalf("Error unmarshaling policy %s: %v", testName.Id, err)
			}

			// Now we re-marshal the policy to make sure it's the same
			newb, err := json.MarshalIndent(policy, "", "  ")
			// newb, err := json.Marshal(policy, )
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(b, newb) {
				t.Fatalf("Serialized policy differed:\n%s", cmp.Diff(string(b), string(newb)))
			}

		})
	}
}
