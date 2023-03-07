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
	principal principal
	str       string
}

// AddService adds one or more services to the Principal.
func (p *Principal) AddService(service ...string) {
	p.principal.Service.Add(service...)
}

// AddAWS adds one or more AWS accounts to the Principal.
func (p *Principal) AddAWS(aws ...string) {
	p.principal.AWS.Add(aws...)
}

// AddCanonicalUser adds one or more canonical users to the Principal.
func (p *Principal) AddCanonicalUser(canonicalUser ...string) {
	p.principal.CanonicalUser.Add(canonicalUser...)
}

// AddFederated adds one or more federated users to the Principal.
func (p *Principal) AddFederated(federated ...string) {
	p.principal.Federated.Add(federated...)
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
		principal: principal{
			Service: NewStringOrSlice(true, service...),
		},
	}
}

// NewAWSPrincipal creates a new Principal that matches an AWS account.
func NewAWSPrincipal(aws ...string) *Principal {
	return &Principal{
		principal: principal{
			AWS: NewStringOrSlice(true, aws...),
		},
	}
}

// NewCanonicalUserPrincipal creates a new Principal that matches a canonical user.
func NewCanonicalUserPrincipal(canonicalUser ...string) *Principal {
	return &Principal{
		principal: principal{
			CanonicalUser: NewStringOrSlice(true, canonicalUser...),
		},
	}
}

// NewFederatedPrincipal creates a new Principal that matches a federated user.
func NewFederatedPrincipal(federated ...string) *Principal {
	return &Principal{
		principal: principal{
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
	principal := principal{}
	err = json.Unmarshal(data, &principal)
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

// Kind returns the kind of the Principal.
func (p *Principal) Kind() string {
	if p.str != "" {
		return PrincipalKindAll
	}
	if p.principal.AWS != nil {
		return PrincipalKindAWS
	} else if p.principal.CanonicalUser != nil {
		return PrincipalKindCanonical
	} else if p.principal.Federated != nil {
		return PrincipalKindFederated
	}
	return PrincipalKindService
}

// Values returns the values of the Principal.
func (p *Principal) Values() []string {
	if p.str != "" {
		return []string{p.str}
	}
	if p.principal.AWS != nil {
		return p.principal.AWS.Values()
	} else if p.principal.CanonicalUser != nil {
		return p.principal.CanonicalUser.Values()
	} else if p.principal.Federated != nil {
		return p.principal.Federated.Values()
	}
	return p.principal.Service.Values()
}

// principal is a json-serializable type in a policy document.
type principal struct {
	AWS           *StringOrSlice `json:"AWS,omitempty"`
	CanonicalUser *StringOrSlice `json:"CanonicalUser,omitempty"`
	Federated     *StringOrSlice `json:"Federated,omitempty"`
	Service       *StringOrSlice `json:"Service,omitempty"`
}

func (p *principal) AddAWS(aws ...string)                     { p.AWS.Add(aws...) }
func (p *principal) AddCanonicalUser(canonicalUser ...string) { p.CanonicalUser.Add(canonicalUser...) }
func (p *principal) AddFederated(federated ...string)         { p.Federated.Add(federated...) }
func (p *principal) AddService(service ...string)             { p.Service.Add(service...) }
