package api

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"case-service/rc"
	"case-service/util"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"gopkg.in/asaskevich/govalidator.v8"
)

// HandleToken handle Token process
func HandleToken(ctx echo.Context) (err error) {
	log.Info().Msg(fmt.Sprintf("Incoming request Token, %v", ctx.Request()))

	// parse request body
	var reqBody util.TokenRequestCase

	err = ctx.Bind(&reqBody)
	if err != nil {
		errMessage := err.Error()
		log.Info().Msg(fmt.Sprintf("Invalid body from client, %s", errMessage))
		ctx.JSON(http.StatusBadRequest, util.ConstructErrorJSONToken(rc.Reject, errMessage, reqBody))
		return nil
	}

	// validate request body
	_, err = govalidator.ValidateStruct(reqBody)
	if err != nil {
		errMessage := err.Error()
		log.Info().Msg(fmt.Sprintf("Invalid Param values, %s ", errMessage))
		ctx.JSON(http.StatusBadRequest, util.ConstructErrorJSONToken(rc.ParamMandatoryNotsent, errMessage, reqBody))
		return nil
	}

	//Construct validate param
	validationParamStruct := util.ValidationParamStruct{
		GrantType:    reqBody.GrantType,
		ClientID:     reqBody.ClientID,
		ClientSecret: reqBody.ClientSecret,
		Scope:        reqBody.Scope,
	}

	respCode := util.CheckValueParamToken(validationParamStruct)
	if !govalidator.IsNull(respCode) {

		errMessage := rc.ClientResponseTextCase[respCode]

		log.Info().Msg(errMessage)
		ctx.JSON(http.StatusBadRequest, util.ConstructErrorJSONToken(respCode, errMessage, reqBody))
		return nil
	}

	// process Token
	respClient, err := processToken(reqBody)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("message error, %s", err.Error()))
	}

	ctx.JSON(http.StatusOK, respClient)

	return nil

}
func processToken(reqBody util.TokenRequestCase) (respClient interface{}, err error) {

	generateToken, expirationTime, errGenerateToken := GenerateJWT(reqBody.ClientID, reqBody.ClientSecret)
	if errGenerateToken != nil {
		errMessage := errGenerateToken.Error()
		tokenResponseCase := util.ConstructErrorJSONToken(rc.Reject, errMessage, reqBody)
		return tokenResponseCase, errGenerateToken
	}

	var tokenResponseCase util.TokenResponseCase
	tokenResponseCase.AccessToken = generateToken
	tokenResponseCase.TokenType = "Bearer"
	tokenResponseCase.ExpiresIn = expirationTime
	tokenResponseCase.Scope = reqBody.Scope

	return tokenResponseCase, err
}

// GenerateJWT generates a JWT token using client_id and secret_id
func GenerateJWT(clientID, secretID string) (string, int64, error) {

	expiredToken, errExpiredToken := strconv.Atoi(os.Getenv("EXPIRED_TOKEN_HOUR"))
	if errExpiredToken != nil {
		log.Info().Msg(fmt.Sprintf("error Get Data Expired Time: %s", errExpiredToken.Error()))
		return "", 0, errExpiredToken
	}

	expirationTime := time.Now().Add(time.Hour * time.Duration(expiredToken)).Unix() // Token berlaku satuan jam
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"client_id": clientID,
		"exp":       expirationTime,
	})

	currentTime := time.Now().Unix() // Waktu sekarang dalam UNIX timestamp

	remainingSeconds := expirationTime - currentTime // Selisih dalam detik

	tokenString, err := token.SignedString([]byte(secretID))
	if err != nil {
		log.Info().Msg(fmt.Sprintf("error Generate Token: %s", err.Error()))
		return "", 0, err
	}

	return tokenString, remainingSeconds, nil
}

// ValidateJWT validates a JWT token
func ValidateJWT(tokenString, secretID string) (bool, int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretID), nil
	})

	if err != nil {
		return false, 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp := int64(claims["exp"].(float64))
		if time.Now().Unix() > exp {
			return false, exp, fmt.Errorf("token expired")
		}
		return true, exp, nil
	}

	return false, 0, fmt.Errorf("invalid token")
}
