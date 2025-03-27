package api

import (
	"case-service/model"
	"case-service/rc"
	"case-service/util"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/guregu/null"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"gopkg.in/asaskevich/govalidator.v8"
)

// HandleGetCase handle Data Case process
func HandleGetCase(ctx echo.Context) (err error) {
	log.Info().Msg(fmt.Sprintf("Incoming request Data Case, %v", ctx.Request()))

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
	// process Data Case
	return getCases(ctx)
}

func getCases(c echo.Context) error {
	db := util.GetDBConnection()

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var cases []model.DataCase
	var totalCount int64
	db.Model(&model.DataCase{}).Count(&totalCount)
	db.Order("time_create DESC").Limit(pageSize).Offset(offset).Find(&cases)

	var responseCases []util.GetCaseData
	for _, c := range cases {
		responseCases = append(responseCases, util.GetCaseData{
			CustomerNumber:     getNullString(c.CustomerNumber),
			CaseID:             int(c.CaseID.Int64),
			DoctorName:         getNullString(c.DoctorName),
			Email:              getNullString(c.Email),
			PreviousCase:       getNullString(c.PreviousCase),
			PreviousCaseNumber: getNullString(c.PreviousCaseNumber),
			PatientName:        getNullString(c.PatientName),
			Dob:                getNullString(c.Dob),
			HeightOfPatient:    getNullString(c.HeightOfPatient),
			Gender:             getNullString(c.Gender),
			Race:               getNullString(c.Race),
			PackageList:        getNullString(c.PackageList),
			LateralXrayDate:    getNullString(c.LateralXrayDate),
			ConsultDate:        getNullString(c.ConsultDate),
			MissingTeeth:       getNullString(c.MissingTeeth),
			AdenoidsRemoved:    getNullString(c.AdenoidsRemoved),
			Comment:            getNullString(c.Comment),
			StatusCase:         getNullString(c.StatusCase),
			LateralXrayImage:   getNullString(c.LateralXrayImage),
			FrontalXrayImage:   getNullString(c.FrontalXrayImage),
			LowerArchImage:     getNullString(c.LowerArchImage),
			UpperArchImage:     getNullString(c.UpperArchImage),
			HandwristXrayImage: getNullString(c.HandwristXrayImage),
			PanoramicXrayImage: getNullString(c.PanoramicXrayImage),
			AdditionalRecord_1: getNullString(c.AdditionalRecord_1),
			AdditionalRecord_2: getNullString(c.AdditionalRecord_2),
			AdditionalRecord_3: getNullString(c.AdditionalRecord_3),
			AdditionalRecord_4: getNullString(c.AdditionalRecord_4),
			AdditionalRecord_5: getNullString(c.AdditionalRecord_5),
			TimeCreate:         getNullTime(c.TimeCreate),
		})
	}

	return c.JSON(http.StatusOK, util.CaseList{
		Cases:      responseCases,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: int(totalCount),
	})
}

// Helper untuk menangani null.String -> string
func getNullString(s null.String) string {
	if s.Valid {
		return s.String
	}
	return ""
}

// Helper untuk menangani null.Time -> string
func getNullTime(t null.Time) string {
	if t.Valid {
		return t.Time.Format("2006-01-02 15:04:05")
	}
	return ""
}
