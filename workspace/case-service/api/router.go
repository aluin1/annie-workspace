package api

import (
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/resty.v1"
)

var headerAuth = "Authorization"
var clientRest *resty.Client

// ConfigRouter configure API router
func ConfigRouter() *echo.Echo {

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("2M"))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${method} ${uri} | status=${status} | ${latency_human}\n",
	}))
	allowedFrontEndUrl := os.Getenv("ALLOWED_FRONTEND_URL")
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{allowedFrontEndUrl},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	e.POST("/insertCase", HandleInsertCase)
	e.POST("/getToken", HandleToken)
	e.GET("/dataCase", HandleGetCase)
	e.POST("/editCaseStatus", HandleEditCase)
	e.GET("/ping", ping)
	e.POST("/uploadFiles", UploadFiles)
	e.POST("/validationTokenGmail", HandleValidationTokenGmail)
	e.POST("/GetDataAnnie", HandleGetDataAnnie)
	e.POST("/UploadFilesCloud", UploadFilesCloud)

	// config iso rest client
	clientTimeout, _ := strconv.Atoi(os.Getenv("ANNIE_TIMEOUT_SECOND"))
	clientRest = resty.New()
	clientRest.SetTimeout(time.Duration(clientTimeout) * time.Second)
	clientRest.SetDebug(true)

	return e
}

func ping(ctx echo.Context) (err error) {
	return nil
}
