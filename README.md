# Go structures for AWS IAM Policy

[![Go Reference](https://pkg.go.dev/badge/github.com/micahhausler/aws-iam-policy.svg)](https://pkg.go.dev/github.com/micahhausler/aws-iam-policy)

Package policy implements types for [AWS's IAM policy grammar] and supports JSON serialization and deserialization.
No validation is performed on the policy, so it is possible to create invalid policies.

**Note**: This package is individually maintained and not supported by Amazon, AWS, or whoever employes the author.

[AWS's IAM policy grammar]: https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_grammar.html

```go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/micahhausler/aws-iam-policy/policy"
)

func main() {
	p := policy.Policy{
		Version: policy.VersionLatest,
		Id:      "CloudTrailBucketPolicy",
		Statements: []policy.Statement{
			{
				Sid:       "AWSCloudTrailWrite20150319",
				Effect:    policy.EffectAllow,
				Principal: policy.NewServicePrincipal("cloudtrail.amazonaws.com"),
				Action:    policy.NewStringOrSlice(false, "s3:PutObject"),
				Resource:  policy.NewStringOrSlice(false, "arn:aws:s3:::examplebucket/AWSLogs/123456789012/*"),
				Condition: map[string]map[string]*policy.ConditionValue{
					"StringEquals": {
						"s3:x-amz-acl": policy.NewConditionValueString(true, "bucket-owner-full-control"),
					},
				},
			},
			{
				Sid:       "AWSCloudTrailAclCheck20150319",
				Effect:    policy.EffectAllow,
				Principal: policy.NewServicePrincipal("cloudtrail.amazonaws.com"),
				Action:    policy.NewStringOrSlice(true, "s3:GetBucketAcl"),
				Resource:  policy.NewStringOrSlice(true, "arn:aws:s3:::examplebucket"),
			},
		},
	}
	out, _ := json.MarshalIndent(p, "", "\t")
	fmt.Println(string(out))
}
```
will output
```json
{
	"Id": "CloudTrailBucketPolicy",
	"Statement": [
		{
			"Action": [
				"s3:PutObject"
			],
			"Condition": {
				"StringEquals": {
					"s3:x-amz-acl": "bucket-owner-full-control"
				}
			},
			"Effect": "Allow",
			"Principal": {
				"Service": "cloudtrail.amazonaws.com"
			},
			"Resource": [
				"arn:aws:s3:::examplebucket/AWSLogs/123456789012/*"
			],
			"Sid": "AWSCloudTrailWrite20150319"
		},
		{
			"Action": "s3:GetBucketAcl",
			"Effect": "Allow",
			"Principal": {
				"Service": "cloudtrail.amazonaws.com"
			},
			"Resource": "arn:aws:s3:::examplebucket",
			"Sid": "AWSCloudTrailAclCheck20150319"
		}
	],
	"Version": "2012-10-17"
}
```

## Safety

You can use strict JSON decoding functionality to ensure that the JSON document
has a valid structure according to this package. Strict decoding functionality
is not enabled by default because it is possible that the IAM policy grammar
could be extended in the future to include new fields.

If you are modifying an existing policy document, or if you are using a policy
document that you did not create, you can enable strict decoding by setting the `DisallowUnknownFields` field on a custom `json.Decoder`.

The following example will fail because the JSON document has a hypothetical
new field `"Foo"` that is not part of the IAM policy statement grammar.

```go
invalidPolicyJSON := []byte(`{
  "Id": "CloudTrailBucketPolicy",
  "Statement": [
    {
      "Foo": "hypothetical new field",
      "Sid": "AWSCloudTrailWrite20150319",
      "Effect": "Allow",
      "Principal": {
        "Service": "cloudtrail.amazonaws.com"
      },
      "Action": "s3:PutObject",
      "Resource": "arn:aws:s3:::examplebucket/AWSLogs/123456789012/*"
    }
  ]
}`)
var p policy.Policy
decoder := json.NewDecoder(bytes.NewBuffer(invalidPolicyJSON))
decoder.DisallowUnknownFields()
err := decoder.Decode(&p)
if err != nil {
  fmt.Println(err)
}
// Output:
// json: unknown field "Foo"
```

## License

[MIT License](LICENSE)
