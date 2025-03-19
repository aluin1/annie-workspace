package util

// TokenRequestCase json data format for payment credit request from client
type TokenRequestCase struct {
	GrantType    string `form:"grant_type" query:"grant_type" valid:"required"`
	ClientID     string `form:"client_id" query:"client_id" valid:"required"`
	ClientSecret string `form:"client_secret" query:"client_secret" valid:"required"`
	Scope        string `form:"scope" query:"scope" valid:"required"`
}

// TokenResponseCase json data format for payment credit response from client
type TokenResponseCase struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	Scope       string `json:"scope"`
}

// TokenResponseCaseReject json data format for payment credit response from client
type TokenResponseCaseReject struct {
	ResponseCode    string `json:"response_code"`
	ResponseMessage string `json:"response_message"`
}
