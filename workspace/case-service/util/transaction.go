package util

import (
	"case-service/model"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/guregu/null"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbInstance *dbUtil
	dbOnce     sync.Once
)

type dbUtil struct {
	db *gorm.DB
}

var ctx = context.Background()

func GetDBConnection() *gorm.DB {
	dbOnce.Do(func() {
		log.Info().Msg("Initializing database connection...")
		var dsn string
		var dialector gorm.Dialector

		dbType := strings.ToLower(os.Getenv("DATABASE_TYPE"))
		log.Info().Msgf("DATABASE_TYPE: %s", dbType)

		switch dbType {
		case "mysql":
			dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
				os.Getenv("USERNAME_DB"), os.Getenv("PASSWORD_DB"), os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_NAME"))
			dialector = mysql.Open(dsn)
		case "postgres":

			dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable search_path=%s",
				os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_PORT"), os.Getenv("USERNAME_DB"),
				os.Getenv("PASSWORD_DB"), os.Getenv("DATABASE_NAME"), os.Getenv("DATABASE_SCHEMA"))
			dialector = postgres.Open(dsn)
		default:
			log.Fatal().Msg("❌ Unsupported database type!")
			return
		}

		log.Info().Msgf("Connecting to database: %s", dsn)

		db, err := gorm.Open(dialector, &gorm.Config{})
		if err != nil {
			log.Fatal().Msgf("❌ Failed to connect to database. DSN: %s. Error: %v", dsn, err)
			panic(err)
			// return
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatal().Msg("❌ Failed to get *sql.DB object")
			return
		}

		// Set connection pool parameters
		sqlDB.SetMaxOpenConns(GetMaxOpenConns())
		sqlDB.SetMaxIdleConns(GetMaxIdleConns())
		sqlDB.SetConnMaxLifetime(time.Minute * GetConnMaxLifetime())

		log.Info().Msg("✅ Database has been initialized")
		dbInstance = &dbUtil{db: db}
	})

	return dbInstance.db
}

func SetContextTimeoutDatabase(ctx context.Context) (context.Context, context.CancelFunc) {
	databaseTimeout, _ := strconv.Atoi(os.Getenv("DATABASE_TIMEOUT_MILLISECOND"))
	contextTimeout, cancel := context.WithTimeout(ctx, time.Duration(databaseTimeout)*time.Millisecond)

	return contextTimeout, cancel
}

func GetMaxOpenConns() int {
	setMaxOpenConns, _ := strconv.Atoi(os.Getenv("DATABASE_SET_MAX_OPEN_CONNS"))
	log.Info().Msg(fmt.Sprintf("Set Max Open Conns: %v", setMaxOpenConns))
	return setMaxOpenConns
}

func GetMaxIdleConns() int {
	setMaxIdleConns, _ := strconv.Atoi(os.Getenv("DATABASE_SET_MAX_IDLE_CONNS"))
	log.Info().Msg(fmt.Sprintf("Set Max Idle Conns: %v", setMaxIdleConns))
	return setMaxIdleConns
}

func GetConnMaxLifetime() time.Duration {
	setConnMaxLifetime, err := time.ParseDuration(os.Getenv("DATABASE_SET_CONN_MAX_LIFETIME_MINUTE"))
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error: %s", err))
		return 0
	}
	log.Info().Msg(fmt.Sprintf("set Conn Max Lifetime: %v", setConnMaxLifetime))
	return setConnMaxLifetime
}

// InsertDataRequestCase Insert DataQuery
func InsertDataRequestCase(reqBody InsertCaseRequest) (caseData *model.DataCase, err error) {

	trx := GetDBConnection().Begin()
	ctxTimeout, cancel := SetContextTimeoutDatabase(ctx)
	defer cancel()
	trxWithContext := trx.WithContext(ctxTimeout)

	CaseData := model.NewCase(reqBody.CustomerNumber)
	CaseData.DoctorName = null.StringFrom(reqBody.DoctorName)
	CaseData.Email = null.StringFrom(reqBody.Email)
	CaseData.PreviousCase = null.StringFrom(reqBody.PreviousCase)
	CaseData.PreviousCaseNumber = null.StringFrom(reqBody.PreviousCaseNumber)
	CaseData.PatientName = null.StringFrom(reqBody.PatientName)
	CaseData.Dob = null.StringFrom(reqBody.Dob)
	CaseData.HeightOfPatient = null.StringFrom(reqBody.HeightOfPatient)
	CaseData.Gender = null.StringFrom(reqBody.Gender)
	CaseData.Race = null.StringFrom(reqBody.Race)
	CaseData.PackageList = null.StringFrom(reqBody.PackageList)
	CaseData.LateralXrayDate = null.StringFrom(reqBody.LateralXrayDate)
	CaseData.ConsultDate = null.StringFrom(reqBody.ConsultDate)
	CaseData.MissingTeeth = null.StringFrom(reqBody.MissingTeeth)
	CaseData.AdenoidsRemoved = null.StringFrom(reqBody.AdenoidsRemoved)
	CaseData.Comment = null.StringFrom(reqBody.Comment)

	//image upload insert DB
	// directoryName := reqBody.CustomerNumber
	// fileBase64LateralXrayImage := reqBody.LateralXrayImage
	// if !govalidator.IsNull(fileBase64LateralXrayImage) {
	// 	namingXray := directoryName + constant.NamingLateralXrayImage
	// 	pathLateralXrayImage, errUploadLateralXray := ConvertBase64(fileBase64LateralXrayImage, directoryName, namingXray)
	// 	if errUploadLateralXray != nil {
	// 		return CaseData, errUploadLateralXray
	// 	}

	CaseData.LateralXrayImage = null.StringFrom(reqBody.LateralXrayImage)
	// }

	// fileBase64FrontalXrayImage := reqBody.FrontalXrayImage
	// if !govalidator.IsNull(fileBase64FrontalXrayImage) {
	// 	namingFrontalXrayImage := directoryName + constant.NamingFrontalXrayImage
	// 	pathFrontalXrayImage, errFrontalXrayImage := ConvertBase64(fileBase64FrontalXrayImage, directoryName, namingFrontalXrayImage)
	// 	if errFrontalXrayImage != nil {
	// 		return CaseData, errFrontalXrayImage
	// 	}

	CaseData.FrontalXrayImage = null.StringFrom(reqBody.FrontalXrayImage)
	// }

	// fileBase64LowerArchImage := reqBody.LowerArchImage
	// if !govalidator.IsNull(fileBase64LowerArchImage) {
	// 	namingLowerArchImage := directoryName + constant.NamingLowerArchImage
	// 	pathLowerArchImage, errLowerArchImage := ConvertBase64(fileBase64LowerArchImage, directoryName, namingLowerArchImage)
	// 	if errLowerArchImage != nil {
	// 		return CaseData, errLowerArchImage
	// 	}

	CaseData.LowerArchImage = null.StringFrom(reqBody.LowerArchImage)
	// }

	// fileBase64UpperArchImage := reqBody.UpperArchImage
	// if !govalidator.IsNull(fileBase64UpperArchImage) {
	// 	namingUpperArchImage := directoryName + constant.NamingUpperArchImage
	// 	pathUpperArchImage, errUpperArchImage := ConvertBase64(fileBase64UpperArchImage, directoryName, namingUpperArchImage)
	// 	if errUpperArchImage != nil {
	// 		return CaseData, errUpperArchImage
	// 	}

	CaseData.UpperArchImage = null.StringFrom(reqBody.UpperArchImage)
	// }

	// fileBase64HandwristXrayImage := reqBody.HandwristXrayImage
	// if !govalidator.IsNull(fileBase64HandwristXrayImage) {
	// 	namingHandwristXrayImage := directoryName + constant.NamingHandwristXrayImage
	// 	pathHandwristXrayImage, errHandwristXrayImage := ConvertBase64(fileBase64HandwristXrayImage, directoryName, namingHandwristXrayImage)
	// 	if errHandwristXrayImage != nil {
	// 		return CaseData, errHandwristXrayImage
	// 	}

	CaseData.HandwristXrayImage = null.StringFrom(reqBody.HandwristXrayImage)
	// }

	// fileBase64PanoramicXrayImage := reqBody.PanoramicXrayImage
	// if !govalidator.IsNull(fileBase64PanoramicXrayImage) {
	// 	namingPanoramicXrayImage := directoryName + constant.NamingPanoramicXrayImage
	// 	pathPanoramicXrayImage, errPanoramicXrayImage := ConvertBase64(fileBase64PanoramicXrayImage, directoryName, namingPanoramicXrayImage)
	// 	if errPanoramicXrayImage != nil {
	// 		return CaseData, errPanoramicXrayImage
	// 	}

	CaseData.PanoramicXrayImage = null.StringFrom(reqBody.PanoramicXrayImage)
	// }

	// fileBase64AdditionalRecord1 := reqBody.AdditionalRecord1
	// if !govalidator.IsNull(fileBase64AdditionalRecord1) {
	// 	namingAdditionalRecord1 := directoryName + constant.NamingAdditionalRecord1
	// 	pathAdditionalRecord1, errAdditionalRecord1 := ConvertBase64(fileBase64AdditionalRecord1, directoryName, namingAdditionalRecord1)
	// 	if errAdditionalRecord1 != nil {
	// 		return CaseData, errAdditionalRecord1
	// 	}

	CaseData.AdditionalRecord_1 = null.StringFrom(reqBody.AdditionalRecord2)
	// }

	// fileBase64AdditionalRecord2 := reqBody.AdditionalRecord2
	// if !govalidator.IsNull(fileBase64AdditionalRecord2) {
	// 	namingAdditionalRecord2 := directoryName + constant.NamingAdditionalRecord2
	// 	pathAdditionalRecord2, errAdditionalRecord2 := ConvertBase64(fileBase64AdditionalRecord2, directoryName, namingAdditionalRecord2)
	// 	if errAdditionalRecord2 != nil {
	// 		return CaseData, errAdditionalRecord2
	// 	}

	CaseData.AdditionalRecord_2 = null.StringFrom(reqBody.AdditionalRecord2)
	// }

	// fileBase64AdditionalRecord3 := reqBody.AdditionalRecord3
	// if !govalidator.IsNull(fileBase64AdditionalRecord3) {
	// 	namingAdditionalRecord3 := directoryName + constant.NamingAdditionalRecord3
	// 	pathAdditionalRecord3, errAdditionalRecord3 := ConvertBase64(fileBase64AdditionalRecord3, directoryName, namingAdditionalRecord3)
	// 	if errAdditionalRecord3 != nil {
	// 		return CaseData, errAdditionalRecord3
	// 	}

	CaseData.AdditionalRecord_3 = null.StringFrom(reqBody.AdditionalRecord3)
	// }

	// fileBase64AdditionalRecord4 := reqBody.AdditionalRecord4
	// if !govalidator.IsNull(fileBase64AdditionalRecord4) {
	// 	namingAdditionalRecord4 := directoryName + constant.NamingAdditionalRecord4
	// 	pathAdditionalRecord4, errAdditionalRecord4 := ConvertBase64(fileBase64AdditionalRecord4, directoryName, namingAdditionalRecord4)
	// 	if errAdditionalRecord4 != nil {
	// 		return CaseData, errAdditionalRecord4
	// 	}

	CaseData.AdditionalRecord_4 = null.StringFrom(reqBody.AdditionalRecord4)
	// }

	// fileBase64AdditionalRecord5 := reqBody.AdditionalRecord5
	// if !govalidator.IsNull(fileBase64AdditionalRecord5) {
	// 	namingAdditionalRecord5 := directoryName + constant.NamingAdditionalRecord5
	// 	pathAdditionalRecord5, errAdditionalRecord5 := ConvertBase64(fileBase64AdditionalRecord5, directoryName, namingAdditionalRecord5)
	// 	if errAdditionalRecord5 != nil {
	// 		return CaseData, errAdditionalRecord5
	// 	}

	CaseData.AdditionalRecord_5 = null.StringFrom(reqBody.AdditionalRecord5)
	// }

	err = trxWithContext.Save(&CaseData).Error
	if err != nil {
		trx.Rollback()
		return CaseData, err
		// return err
	}

	trx.Commit()
	return CaseData, err
	// return err
}
