package util

// EditCaseRequest xml data for request to case
type EditCaseRequest struct {
	CaseId         string `json:"case_id" valid:"required"`
	CustomerNumber string `json:"customer_number" valid:"required"`
	Comment        string `json:"comment"` //optional
	StatusCase     string `json:"status_case" valid:"required"`
}

// EditCaseResponse response from case
type EditCaseResponse struct {
	CustomerNumber  string `json:"customer_number"`
	ResponseCode    string `json:"response_code,omitempty"`
	ResponseMessage string `json:"response_message,omitempty"`
}
