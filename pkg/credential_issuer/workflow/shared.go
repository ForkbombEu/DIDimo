package workflow

const EbsiIssuersUrl = "https://api-conformance.ebsi.eu/trusted-issuers-registry/v5/issuers"

type FetchIssuersActivityResponse struct { Issuers []string }

type ApiResponse struct {
	Items    []Item          `json:"items"`
	Links    Links `json:"links"` 
	PageSize int             `json:"pageSize"`
	Self     string          `json:"self"`
	Total    int             `json:"total"`
}

type Item struct {
	Did  string `json:"did"`
	Href string `json:"href"`
}

type Links struct {
	First string `json:"first"`
	Last  string `json:"last"`
	Next  string `json:"next"`
	Prev  string `json:"prev"`
}


type CredentialWorkflowInput struct {
	BaseURL  string // Base URL for the credential issuer
	IssuerID string // ID of the credentials issuer from PB
}

type CredentialWorkflowResponse struct {
	Message string
}

const FetchIssuersTaskQueue = "FetchIssuersTaskQueue"