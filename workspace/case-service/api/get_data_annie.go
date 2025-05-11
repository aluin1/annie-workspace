package api

import (
	"case-service/rc"
	"case-service/util"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"gopkg.in/asaskevich/govalidator.v8"
	"gopkg.in/resty.v1"
)

// HandleGetDataAnnie
func HandleGetDataAnnie(ctx echo.Context) (err error) {
	log.Info().Msg(fmt.Sprintf("Incoming request Get Data Annie, %v", ctx.Request()))

	var reqBody util.GetDataAnnieRequest

	err = ctx.Bind(&reqBody)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Invalid JSON body from client, %s", err.Error()))
		errMessage := err.Error()
		ctx.JSON(http.StatusBadRequest, util.ConstructGeneralErrorGetDataJSON(rc.Reject, errMessage))
		return nil
	}

	// validate request body
	_, err = govalidator.ValidateStruct(reqBody)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Invalid JSON values, %s ", err.Error()))
		errMessage := err.Error()
		ctx.JSON(http.StatusBadRequest, util.ConstructGeneralErrorGetDataJSON(rc.Reject, errMessage))
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
		ctx.JSON(http.StatusBadRequest, util.ConstructGeneralErrorGetDataJSON(respCode, errMessage))
		return nil
	}

	splitToken := strings.Split(tokenValue, "Bearer ")
	token := splitToken[1]
	secretID := os.Getenv("CLIENT_SECRET")
	_, _, errToken := ValidateJWT(token, secretID)
	if errToken != nil {
		log.Info().Msg(fmt.Sprintf("Error Token, %s ", errToken.Error()))
		errMessage := errToken.Error()
		ctx.JSON(http.StatusUnauthorized, util.ConstructGeneralErrorGetDataJSON(rc.Reject, errMessage))
		return nil
	}

	// process Insert Case
	statusCode, respClient, err := GetDataAnnie(reqBody)
	if err != nil {
		log.Info().Msg(err.Error())
	}
	ctx.JSON(statusCode, respClient)

	return nil
}

func GetDataAnnie(reqBody util.GetDataAnnieRequest) (stsCode int, respClient interface{}, err error) {
	statusCode := http.StatusOK

	urlQuery := os.Getenv("CASE_DETAIL_URL")
	resp, err := clientRest.SetPreRequestHook(func(client *resty.Client, request *resty.Request) error {

		request.Header.Set("Content-Type", "application/json")

		return nil
	}).R().SetBody(reqBody).Post(urlQuery)
	if err != nil {
		errMessage := err.Error()
		log.Info().Msg(errMessage)

		return statusCode, util.ConstructGeneralErrorGetDataJSON(rc.Reject, errMessage), fmt.Errorf("error message: %v", errMessage)
	}

	respClientJSON := util.GetDataAnnieResponse{}
	err = json.Unmarshal(resp.Body(), &respClientJSON)
	if err != nil {
		errMessage := err.Error()
		log.Info().Msg(errMessage)
		return statusCode, util.ConstructGeneralErrorGetDataJSON(rc.Reject, err.Error()), fmt.Errorf("error message: %v", errMessage)

	}

	return statusCode, respClientJSON, err
}
