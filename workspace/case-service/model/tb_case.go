package model

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

type DataCase struct {
	CaseID             null.Int    `gorm:"column:case_id;primary_key"`
	CustomerNumber     null.String `gorm:"column:customer_number"`
	DoctorName         null.String `gorm:"column:doctor_name"`
	Email              null.String `gorm:"column:email"`
	PreviousCase       null.String `gorm:"column:previous_case"`
	PreviousCaseNumber null.String `gorm:"column:previous_case_number"`
	PatientName        null.String `gorm:"column:patient_name"`
	Dob                null.String `gorm:"column:dob"`
	HeightOfPatient    null.String `gorm:"column:height_of_patient"`
	Gender             null.String `gorm:"column:gender"`
	Race               null.String `gorm:"column:race"`
	PackageList        null.String `gorm:"column:package_list"`
	LateralXrayDate    null.String `gorm:"column:lateral_xray_date"`
	ConsultDate        null.String `gorm:"column:consult_date"`
	MissingTeeth       null.String `gorm:"column:missing_teeth"`
	AdenoidsRemoved    null.String `gorm:"column:adenoids_removed"`
	Comment            null.String `gorm:"column:comment"`
	StatusCase         null.String `gorm:"column:status_case"`

	LateralXrayImage   null.String `gorm:"lateral_xray_image"`   //optional
	FrontalXrayImage   null.String `gorm:"frontal_xray_image"`   //optional
	LowerArchImage     null.String `gorm:"lower_arch_image"`     //optional
	UpperArchImage     null.String `gorm:"upper_arch_image"`     //optional
	HandwristXrayImage null.String `gorm:"handwrist_xray_image"` //optional
	PanoramicXrayImage null.String `gorm:"panoramic_xray_image"` //optional
	AdditionalRecord_1 null.String `gorm:"additional_record_1"`  //optional
	AdditionalRecord_2 null.String `gorm:"additional_record_2"`  //optional
	AdditionalRecord_3 null.String `gorm:"additional_record_3"`  //optional
	AdditionalRecord_4 null.String `gorm:"additional_record_4"`  //optional
	AdditionalRecord_5 null.String `gorm:"additional_record_5"`  //optional

	TimeCreate null.Time `gorm:"column:time_create"`
}

// TableName sets the insert table name for this struct type
func (q *DataCase) TableName() string {
	return "tb_case"
}

// NewCase create new QR transaction model with predefined value
func NewCase(customerNumber string) *DataCase {
	return &DataCase{
		CustomerNumber: null.StringFrom(customerNumber),
		TimeCreate:     null.TimeFrom(time.Now()),
	}
}
