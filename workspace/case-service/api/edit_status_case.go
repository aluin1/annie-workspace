package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"case-service/model"
	"case-service/rc"
	"case-service/util"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"gopkg.in/asaskevich/govalidator.v8"
	"gorm.io/gorm"
)

// HandleEditCase handle Insert Case process
func HandleEditCase(ctx echo.Context) (err error) {
	log.Info().Msg(fmt.Sprintf("Incoming request Insert Case, %v", ctx.Request()))

	// parse request body
	var reqBody util.EditCaseRequest

	err = ctx.Bind(&reqBody)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Invalid JSON body from client, %s", err.Error()))
		errMessage := err.Error()
		ctx.JSON(http.StatusBadRequest, util.ConstructEditCaseErrorJSON(rc.Reject, errMessage, reqBody))
		return nil
	}

	// validate request body
	_, err = govalidator.ValidateStruct(reqBody)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Invalid JSON values, %s ", err.Error()))
		errMessage := err.Error()
		ctx.JSON(http.StatusBadRequest, util.ConstructEditCaseErrorJSON(rc.ParamMandatoryNotsent, errMessage, reqBody))
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
		ctx.JSON(http.StatusBadRequest, util.ConstructEditCaseErrorJSON(respCode, errMessage, reqBody))
		return nil
	}

	splitToken := strings.Split(tokenValue, "Bearer ")
	token := splitToken[1]
	secretID := os.Getenv("CLIENT_SECRET")
	_, _, errToken := ValidateJWT(token, secretID)
	if errToken != nil {
		log.Info().Msg(fmt.Sprintf("Error Token, %s ", errToken.Error()))
		errMessage := errToken.Error()
		ctx.JSON(http.StatusUnauthorized, util.ConstructEditCaseErrorJSON(rc.Reject, errMessage, reqBody))
		return nil
	}
	// process Insert Case
	statusCode, respClient, err := processEditCase(reqBody)
	if err != nil {
		log.Info().Msg(err.Error())
	}
	ctx.JSON(statusCode, respClient)

	return nil

}

func processEditCase(reqBody util.EditCaseRequest) (stsCode int, respClient interface{}, err error) {
	statusCode := http.StatusOK

	//get Data Case
	var originalCase model.DataCase
	db := util.GetDBConnection()
	ctxTimeout, cancel := util.SetContextTimeoutDatabase(context.Background())
	defer cancel()
	resultQry := db.WithContext(ctxTimeout).Debug().Where("case_id = ? AND customer_number = ?", reqBody.CaseId, reqBody.CustomerNumber).Order("time_create DESC").First(&originalCase)
	if resultQry.Error == gorm.ErrRecordNotFound {

		errMessage := "Original Case not found"

		log.Info().Msg(errMessage)
		statusCode = http.StatusInternalServerError
		jsonResponse := util.ConstructEditCaseErrorJSON(rc.Reject, errMessage, reqBody)

		return statusCode, jsonResponse, resultQry.Error
	}

	//Error get Data Case
	if resultQry.Error != nil {
		errMessage := "Something Wrong in Database, message error:" + resultQry.Error.Error()
		log.Info().Msg(errMessage)
		statusCode = http.StatusInternalServerError
		jsonResponse := util.ConstructEditCaseErrorJSON(rc.Reject, errMessage, reqBody)

		return statusCode, jsonResponse, resultQry.Error
	}

	_, errQueryTran := util.EditDataCase(reqBody, originalCase)
	if errQueryTran != nil {
		errMessage := fmt.Sprintf("Error Proses Edit Data Case. \n customer_number: %s \n error message: %s", reqBody.CustomerNumber, errQueryTran.Error())
		log.Info().Msg(errMessage)
		statusCode = http.StatusInternalServerError
		jsonResponse := util.ConstructEditCaseErrorJSON(rc.Reject, errMessage, reqBody)

		return statusCode, jsonResponse, errQueryTran
	}

	var resp util.EditCaseResponse
	resp.CustomerNumber = reqBody.CustomerNumber
	resp.ResponseCode = rc.Success
	resp.ResponseMessage = rc.ClientResponseTextCase[rc.Success] + " - Edit Data Case"
	respClient = resp

	return statusCode, respClient, err
}
