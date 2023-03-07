package policy

import "encoding/json"

const (
	PrincipalAll = "*"

	PrincipalKindAll       = "All"
	PrincipalKindAWS       = "AWS"
	PrincipalKindCanonical = "CanonicalUser"
	PrincipalKindFederated = "Federated"
	PrincipalKindService   = "Service"
)

// Principal is a Principal in a policy document.
type Principal struct {
	principal *principal
	str       string
}

// AddService adds one or more services to the Principal.
func (p *Principal) AddService(service ...string) {
	if p.principal == nil {
		p.principal = &principal{}
	}
	p.principal.AddService(service...)
}

// AddAWS adds one or more AWS accounts to the Principal.
func (p *Principal) AddAWS(aws ...string) {
	if p.principal == nil {
		p.principal = &principal{}
	}
	p.principal.AddAWS(aws...)
}

// AddCanonicalUser adds one or more canonical users to the Principal.
func (p *Principal) AddCanonicalUser(canonicalUser ...string) {
	if p.principal == nil {
		p.principal = &principal{}
	}
	p.principal.AddCanonicalUser(canonicalUser...)
}

// AddFederated adds one or more federated users to the Principal.
func (p *Principal) AddFederated(federated ...string) {
	if p.principal == nil {
		p.principal = &principal{}
	}
	p.principal.AddFederated(federated...)
}

func newPrincipalFromString(s string) *Principal {
	return &Principal{
		str: s,
	}
}

// NewGlobalPrincipal creates a new Principal that matches all principals.
func NewGlobalPrincipal() *Principal {
	return newPrincipalFromString(PrincipalAll)
}

// NewServicePrincipal creates a new Principal that matches a service.
func NewServicePrincipal(service ...string) *Principal {
	return &Principal{
		principal: &principal{
			Service: NewStringOrSlice(true, service...),
		},
	}
}

// NewAWSPrincipal creates a new Principal that matches an AWS account.
func NewAWSPrincipal(aws ...string) *Principal {
	return &Principal{
		principal: &principal{
			AWS: NewStringOrSlice(true, aws...),
		},
	}
}

// NewCanonicalUserPrincipal creates a new Principal that matches a canonical user.
func NewCanonicalUserPrincipal(canonicalUser ...string) *Principal {
	return &Principal{
		principal: &principal{
			CanonicalUser: NewStringOrSlice(true, canonicalUser...),
		},
	}
}

// NewFederatedPrincipal creates a new Principal that matches a federated user.
func NewFederatedPrincipal(federated ...string) *Principal {
	return &Principal{
		principal: &principal{
			Federated: NewStringOrSlice(true, federated...),
		},
	}
}

func (p *Principal) UnmarshalJSON(data []byte) error {
	var tmp interface{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	str, ok := tmp.(string)
	if ok {
		p.str = str
		return nil
	}
	principal := &principal{}
	err = json.Unmarshal(data, principal)
	if err != nil {
		return err
	}
	p.principal = principal
	return nil
}

func (p *Principal) MarshalJSON() ([]byte, error) {
	if p.str != "" {
		return json.Marshal(p.str)
	}
	return json.Marshal(p.principal)
}

// Kinds returns the kinds of the Principal.
func (p *Principal) Kinds() []string {
	if p.str != "" {
		return []string{PrincipalKindAll}
	}
	resp := []string{}
	if p.principal.AWS != nil {
		resp = append(resp, PrincipalKindAWS)
	}
	if p.principal.CanonicalUser != nil {
		resp = append(resp, PrincipalKindCanonical)
	}
	if p.principal.Federated != nil {
		resp = append(resp, PrincipalKindFederated)
	}
	if p.principal.Service != nil {
		resp = append(resp, PrincipalKindService)
	}
	return resp
}

// AWS returns the AWS accounts of the Principal.
func (p *Principal) AWS() *StringOrSlice {
	if p.principal == nil {
		return nil
	}
	return p.principal.AWS
}

// CanonicalUser returns the canonical users of the Principal.
func (p *Principal) CanonicalUser() *StringOrSlice {
	if p.principal == nil {
		return nil
	}
	return p.principal.CanonicalUser
}

// Federated returns the federated users of the Principal.
func (p *Principal) Federated() *StringOrSlice {
	if p.principal == nil {
		return nil
	}
	return p.principal.Federated
}

// Service returns the services of the Principal.
func (p *Principal) Service() *StringOrSlice {
	if p.principal == nil {
		return nil
	}
	return p.principal.Service
}

// principal is a json-serializable type in a policy document.
type principal struct {
	AWS           *StringOrSlice `json:"AWS,omitempty"`
	CanonicalUser *StringOrSlice `json:"CanonicalUser,omitempty"`
	Federated     *StringOrSlice `json:"Federated,omitempty"`
	Service       *StringOrSlice `json:"Service,omitempty"`
}

func (p *principal) AddAWS(aws ...string) {
	if p.AWS == nil {
		p.AWS = NewStringOrSlice(true, aws...)
		return
	}
	p.AWS.Add(aws...)
}
func (p *principal) AddCanonicalUser(canonicalUser ...string) {
	if p.CanonicalUser == nil {
		p.CanonicalUser = NewStringOrSlice(true, canonicalUser...)
		return
	}
	p.CanonicalUser.Add(canonicalUser...)
}
func (p *principal) AddFederated(federated ...string) {
	if p.Federated == nil {
		p.Federated = NewStringOrSlice(true, federated...)
		return
	}
	p.Federated.Add(federated...)
}
func (p *principal) AddService(service ...string) {
	if p.Service == nil {
		p.Service = NewStringOrSlice(true, service...)
		return
	}
	p.Service.Add(service...)
}
