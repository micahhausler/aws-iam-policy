package policy

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	fixtureDir = "./test_fixtures/managed-policies"
)

func TestManagedPolicies(t *testing.T) {
	fs.WalkDir(os.DirFS(fixtureDir), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			t.Fatal(err)
		}
		if d.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}
		t.Run(d.Name(), func(t *testing.T) {
			HelperTestPolicy(t, path)
		})
		return nil
	})
}

func HelperTestPolicy(t *testing.T, path string) {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(fixtureDir, path))
	if err != nil {
		t.Fatal(err)
	}

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

	if !bytes.Equal(bytes.TrimSpace(b), bytes.TrimSpace(newb)) {
		t.Fatalf("Serialized policy differed:\n%s", cmp.Diff(string(b), string(newb)))
	}

}
