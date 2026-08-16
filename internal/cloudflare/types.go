package cloudflare

// APIError represents the Cloudflare API error envelope.
type APIError struct {
	Code             int    `json:"code"`
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url,omitempty"`
}

// APIResponse is the common API payload wrapper used by Cloudflare.
type APIResponse struct {
	Success bool       `json:"success"`
	Result  any        `json:"result,omitempty"`
	Errors  []APIError `json:"errors,omitempty"`
}
