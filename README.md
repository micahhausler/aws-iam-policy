# Go structures for AWS IAM Policy

[![Go Reference](https://pkg.go.dev/badge/github.com/micahhausler/aws-iam-policy.svg)](https://pkg.go.dev/github.com/micahhausler/aws-iam-policy)

Package policy implements types for [AWS's IAM policy grammar] and supports JSON serialization and deserialization.
No validation is performed on the policy, so it is possible to create invalid policies.

[AWS's IAM policy grammar]: https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_grammar.html
