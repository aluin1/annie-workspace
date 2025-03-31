package util

import (
	"case-service/constant"
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

	CaseData.StatusCase = null.StringFrom(constant.Submited)
	CaseData.LateralXrayImage = null.StringFrom(reqBody.LateralXrayImage)
	CaseData.FrontalXrayImage = null.StringFrom(reqBody.FrontalXrayImage)
	CaseData.LowerArchImage = null.StringFrom(reqBody.LowerArchImage)
	CaseData.UpperArchImage = null.StringFrom(reqBody.UpperArchImage)
	CaseData.HandwristXrayImage = null.StringFrom(reqBody.HandwristXrayImage)
	CaseData.PanoramicXrayImage = null.StringFrom(reqBody.PanoramicXrayImage)
	CaseData.AdditionalRecord_1 = null.StringFrom(reqBody.AdditionalRecord1)
	CaseData.AdditionalRecord_2 = null.StringFrom(reqBody.AdditionalRecord2)
	CaseData.AdditionalRecord_3 = null.StringFrom(reqBody.AdditionalRecord3)
	CaseData.AdditionalRecord_4 = null.StringFrom(reqBody.AdditionalRecord4)
	CaseData.AdditionalRecord_5 = null.StringFrom(reqBody.AdditionalRecord5)

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

// EditDataCase Insert DataQuery
func EditDataCase(reqBody EditCaseRequest, CaseData model.DataCase) (caseData *model.DataCase, err error) {

	trx := GetDBConnection().Begin()
	ctxTimeout, cancel := SetContextTimeoutDatabase(ctx)
	defer cancel()
	trxWithContext := trx.WithContext(ctxTimeout)

	CaseData.StatusCase = null.StringFrom(reqBody.StatusCase)
	CaseData.Comment = null.StringFrom(reqBody.Comment)
	err = trxWithContext.Save(&CaseData).Error
	if err != nil {
		trx.Rollback()
		return &CaseData, err
		// return err
	}

	trx.Commit()
	return &CaseData, err
	// return err
}
