package util

// ConstructGeneralErrorJSON
func ConstructGeneralErrorJSON(responseCode, errMessage string, reqBody InsertCaseRequest) InsertCaseResponse {
	err := InsertCaseResponse{}
	err.CustomerNumber = reqBody.CustomerNumber
	err.ResponseCode = responseCode
	err.ResponseMessage = errMessage

	return err
}

// ConstructErrorJSONToken
func ConstructErrorJSONToken(responseCode, errMessage string, reqBody TokenRequestCase) TokenResponseCaseReject {
	err := TokenResponseCaseReject{}
	err.ResponseCode = responseCode
	err.ResponseMessage = errMessage

	return err
}

// ConstructGeneralErrorGetDataJSON
func ConstructGeneralErrorGetDataJSON(responseCode, errMessage string) GetCaseResponse {
	err := GetCaseResponse{}
	err.ResponseCode = responseCode
	err.ResponseMessage = errMessage

	return err
}
