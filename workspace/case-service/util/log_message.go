package util

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
)

func LogReqValidateAuthToken(reqBody ReqTokenValidation) {

	requesBody, errReqBody := json.Marshal(reqBody)
	if errReqBody != nil {

		log.Info().Msg(fmt.Sprintf("Error Marshal RequestValidate Auth Token. error: %s", errReqBody.Error()))
	}
	log.Info().Msg(fmt.Sprintf("Request Validate Auth Token : %s", requesBody))

}

func LogRespValidateAuthToken(respBody GetCaseResponse) {

	responseBody, errResponseBody := json.Marshal(respBody)
	if errResponseBody != nil {

		log.Info().Msg(fmt.Sprintf("Error Marshal Response Validate Auth Token. error: %s", errResponseBody.Error()))
	}
	log.Info().Msg(fmt.Sprintf("Response Validate Auth Token : %s", responseBody))

}
