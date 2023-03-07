package policy

const (
	// See https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_effect.html
	EffectAllow = "Allow"
	EffectDeny  = "Deny"

	// See https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_version.html
	Version2012_10_17 = "2012-10-17"
	Version2008_10_17 = "2008-10-17"
	VersionLatest     = Version2012_10_17
)

// Policy is a policy document.
type Policy struct {
	Id         string      `json:"Id,omitempty"`
	Statements []Statement `json:"Statement"`
	Version    string      `json:"Version"`
}

// Statement is a single statement in a policy document.
type Statement struct {
	Action       *StringOrSlice                        `json:"Action,omitempty"`
	Condition    map[string]map[string]*ConditionValue `json:"Condition,omitempty"`
	Effect       string                                `json:"Effect"`
	NotAction    *StringOrSlice                        `json:"NotAction,omitempty"`
	NotResource  *StringOrSlice                        `json:"NotResource,omitempty"`
	Principal    *Principal                            `json:"Principal,omitempty"`
	NotPrincipal *Principal                            `json:"NotPrincipal,omitempty"`
	Resource     *StringOrSlice                        `json:"Resource,omitempty"`
	Sid          string                                `json:"Sid,omitempty"`
}
