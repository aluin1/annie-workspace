package api

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"case-service/rc"
	"case-service/util"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"gopkg.in/asaskevich/govalidator.v8"
)

// HandleInsertCase handle Insert Case process
func HandleInsertCase(ctx echo.Context) (err error) {
	log.Info().Msg(fmt.Sprintf("Incoming request Insert Case, %v", ctx.Request()))

	// parse request body
	var reqBody util.InsertCaseRequest

	err = ctx.Bind(&reqBody)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Invalid JSON body from client, %s", err.Error()))
		errMessage := err.Error()
		ctx.JSON(http.StatusBadRequest, util.ConstructGeneralErrorJSON(rc.Reject, errMessage, reqBody))
		return nil
	}

	// validate request body
	_, err = govalidator.ValidateStruct(reqBody)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Invalid JSON values, %s ", err.Error()))
		errMessage := err.Error()
		ctx.JSON(http.StatusBadRequest, util.ConstructGeneralErrorJSON(rc.ParamMandatoryNotsent, errMessage, reqBody))
		return nil
	}

	tokenValue := ctx.Request().Header.Get(headerAuth)

	// Construct validate param
	validationParamStruct := util.ValidationParamStruct{
		Token: tokenValue,
	}

	respCode := util.CheckValueParamGeneral(validationParamStruct)
	if !govalidator.IsNull(respCode) {
		errMessage := rc.ClientResponseTextCase[respCode]
		ctx.JSON(http.StatusBadRequest, util.ConstructGeneralErrorJSON(respCode, errMessage, reqBody))
		return nil
	}

	splitToken := strings.Split(tokenValue, "Bearer ")
	token := splitToken[1]
	secretID := os.Getenv("CLIENT_SECRET")
	_, _, errToken := ValidateJWT(token, secretID)
	if errToken != nil {
		log.Info().Msg(fmt.Sprintf("Error Token, %s ", errToken.Error()))
		errMessage := errToken.Error()
		ctx.JSON(http.StatusUnauthorized, util.ConstructGeneralErrorJSON(rc.Reject, errMessage, reqBody))
		return nil
	}
	// process Insert Case
	statusCode, respClient, err := processCase(reqBody)
	if err != nil {
		log.Info().Msg(err.Error())
	}
	ctx.JSON(statusCode, respClient)

	return nil

}

func processCase(reqBody util.InsertCaseRequest) (stsCode int, respClient interface{}, err error) {
	statusCode := http.StatusOK
	dataCase, errQueryTran := util.InsertDataRequestCase(reqBody)
	if errQueryTran != nil {
		errMessage := fmt.Sprintf("Error Proses Data Case. \n customer_number: %s \n error message: %s", reqBody.CustomerNumber, errQueryTran.Error())
		log.Info().Msg(errMessage)
		statusCode = http.StatusInternalServerError
		jsonResponse := util.ConstructGeneralErrorJSON(rc.Reject, errMessage, reqBody)

		return statusCode, jsonResponse, errQueryTran
	}

	if os.Getenv("SEND_EMAIL") == "true" {

		errEmail := util.SendEmail(dataCase)
		if errEmail != nil {
			errMessage := fmt.Sprintf("message error send email: %s", errEmail.Error())
			log.Info().Msg(errMessage)
			statusCode = http.StatusInternalServerError
			jsonResponse := util.ConstructGeneralErrorJSON(rc.Reject, errMessage, reqBody)

			return statusCode, jsonResponse, errEmail
		}
	}

	var resp util.InsertCaseResponse
	resp.CustomerNumber = reqBody.CustomerNumber
	resp.ResponseCode = rc.Success
	resp.ResponseMessage = rc.ClientResponseTextCase[rc.Success]
	respClient = resp

	return statusCode, respClient, err
}
