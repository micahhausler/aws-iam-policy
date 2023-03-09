/*
Package policy implements types for [AWS's IAM policy grammar] and supports JSON serialization and deserialization.
No validation is performed on the policy, so it is possible to create invalid policies.

Here is an example that creates a policy document using this package.

	package main

	import (
		"encoding/json"
		"fmt"
		"github.com/micahhausler/aws-iam-policy/policy"
	)

	func main() {
		p := policy.Policy{
			Version: policy.VersionLatest,
			Statements: policy.NewStatementOrSlice([]policy.Statement{
				{
					Sid:       "S3Access",
					Effect:    policy.EffectAllow,
					Principal: policy.NewAWSPrincipal("arn:aws:iam::123456789012:role/my-role"),
					Action:    policy.NewStringOrSlice(true, "s3:ListBucket"),
					Resource:  policy.NewStringOrSlice(true, "arn:aws:s3:::examplebucket/AWSLogs/123456789012/*"),
				},
			},
		}...)
		b, err := json.MarshalIndent(p, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))
	}

[AWS's IAM policy grammar]: https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_grammar.html
*/
package policy
