package brale

type APIResponse struct {
	Data     *APIData               `json:"data,omitempty"`
	Included []Included             `json:"included,omitempty"`
	Links    *APILinks              `json:"links,omitempty"`
	Errors   []ErrorDetail          `json:"errors,omitempty"`
	Meta     map[string]interface{} `json:"meta,omitempty"`
}

type APIData struct {
	Attributes    APIAttributes `json:"attributes"`
	ID            string        `json:"id"`
	Type          string        `json:"type"`
	Links         APILinks      `json:"links"`
	Relationships Relationships `json:"relationships"`
}

type APIAttributes struct {
	Created string `json:"created"`
	Status  string `json:"status"`
	Type    string `json:"type"`
	Updated string `json:"updated"`
}

type APILinks struct {
	AdditionalProp struct {
		Href string `json:"href"`
	} `json:"additionalProp"`
}

type Relationships struct {
	Transactions RelationshipData `json:"transactions"`
}

type RelationshipData struct {
	Data  []RelationshipDetail `json:"data"`
	Links struct {
		Related struct {
			Href string `json:"href"`
		} `json:"related"`
	} `json:"links"`
}

type RelationshipDetail struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type Included struct {
	Attributes    IncludedAttributes    `json:"attributes"`
	ID            string                `json:"id"`
	Type          string                `json:"type"`
	Links         APILinks              `json:"links"`
	Relationships IncludedRelationships `json:"relationships"`
}

type IncludedAttributes struct {
	Amount struct {
		Currency string `json:"currency"`
		Value    string `json:"value"`
	} `json:"amount"`
	Created string `json:"created"`
	Hash    string `json:"hash"`
	Note    string `json:"note"`
	Status  string `json:"status"`
	Type    string `json:"type"`
	Updated string `json:"updated"`
}

type IncludedRelationships struct {
	Deployment  RelationshipDetail `json:"deployment"`
	Destination RelationshipDetail `json:"destination"`
}

type ErrorDetail struct {
	Code   string            `json:"code"`
	Detail string            `json:"detail"`
	ID     string            `json:"id"`
	Links  map[string]string `json:"links"`
	Meta   map[string]string `json:"meta"`
	Source struct {
		Parameter string `json:"parameter"`
		Pointer   string `json:"pointer"`
	} `json:"source"`
	Status string `json:"status"`
	Title  string `json:"title"`
}
