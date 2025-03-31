package api

import (
	"case-service/model"
	"case-service/rc"
	"case-service/util"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"gopkg.in/asaskevich/govalidator.v8"
	"gorm.io/gorm"
)

// HandleValidationTokenGmail handles Gmail token validation
func HandleValidationTokenGmail(ctx echo.Context) error {
	var reqBody util.ReqTokenValidation
	if err := ctx.Bind(&reqBody); err != nil {
		log.Info().Msg(fmt.Sprintf("Invalid body from client: %s", err.Error()))

		errMessage := err.Error()
		errResp := util.ConstructGeneralErrorGetDataJSON(rc.Reject, errMessage)
		util.LogRespValidateAuthToken(errResp)
		return ctx.JSON(http.StatusBadRequest, errResp)
	}

	util.LogReqValidateAuthToken(reqBody)

	// validate request body
	_, err := govalidator.ValidateStruct(reqBody)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Invalid JSON values, %s ", err.Error()))
		errMessage := err.Error()
		errResp := util.ConstructGeneralErrorGetDataJSON(rc.Reject, errMessage)
		util.LogRespValidateAuthToken(errResp)
		return ctx.JSON(http.StatusBadRequest, errResp)

	}

	jwtToken := reqBody.Token
	respValidationToken, errValidateTokenAuth := util.ValidationTokenAuth(jwtToken)

	if errValidateTokenAuth != nil {
		errMessage := errValidateTokenAuth.Error()
		errResp := util.ConstructGeneralErrorGetDataJSON(rc.Reject, errMessage)
		util.LogRespValidateAuthToken(errResp)
		return ctx.JSON(http.StatusOK, errResp)
	}

	//get Data Case
	var originalDataUser model.DataUser
	db := util.GetDBConnection()
	ctxTimeout, cancel := util.SetContextTimeoutDatabase(context.Background())
	defer cancel()
	resultQry := db.WithContext(ctxTimeout).Debug().Where("email = ? AND active = ?", respValidationToken.Email, "1").Order("time_create DESC").First(&originalDataUser)
	if resultQry.Error == gorm.ErrRecordNotFound {

		errMessage := "User data is not allowed, please switch accounts."
		errResp := util.ConstructGeneralErrorGetDataJSON(rc.Reject, errMessage)
		util.LogRespValidateAuthToken(errResp)
		return ctx.JSON(http.StatusUnauthorized, errResp)
	}

	//Error get Data Case
	if resultQry.Error != nil {
		errMessage := "Something Wrong in Database, message error:" + resultQry.Error.Error()
		errResp := util.ConstructGeneralErrorGetDataJSON(rc.Reject, errMessage)
		util.LogRespValidateAuthToken(errResp)
		return ctx.JSON(http.StatusUnauthorized, errResp)
	}

	responseBody, errResponseBody := json.Marshal(respValidationToken)
	if errResponseBody != nil {

		log.Info().Msg(fmt.Sprintf("Error Marshal Response Validate Auth Token. error: %s", errResponseBody.Error()))
	}
	log.Info().Msg(fmt.Sprintf("Response Validate Auth Token : %s", responseBody))

	return ctx.JSON(http.StatusOK, respValidationToken)
}
