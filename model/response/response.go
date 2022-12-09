package response

type ResponseContainer struct {
	Response Response `json:"RESPONSE"`
}

type Response struct {
	StatusCode      int           `json:"STATUS_CODE"`
	ErrorCode       *string       `json:"ERROR_CODE"`
	ResponseCode    *string       `json:"RESPONSE_CODE"`
	ResponseMessage *string       `json:"RESPONSE_MESSAGE"`
	Errors          []string      `json:"ERRORS"`
	Data            interface{}   `json:"DATA"`
	Info            *ResponseInfo `json:"INFO,omitempty"`
}

type ErrorContainer struct {
	Response Response `json:"RESPONSE"`
}

type ResponseInfo struct {
	Limit    int `json:"LIMIT"`
	Page     int `json:"PAGE"`
	PageSize int `json:"PAGE_SIZE"`
	Total    int `json:"TOTAL"`
}
