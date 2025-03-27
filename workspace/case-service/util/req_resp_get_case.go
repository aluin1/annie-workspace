package util

// CaseList represents a paginated list of cases
type CaseList struct {
	Cases      []GetCaseData `json:"data_case"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalCount int           `json:"total_count"`
}

// GetCaseResponse
type GetCaseResponse struct {
	ResponseCode    string `json:"response_code,omitempty"`
	ResponseMessage string `json:"response_message,omitempty"`
}

// GetCaseData
type GetCaseData struct {
	CaseID             int    `json:"case_id"`
	CustomerNumber     string `json:"customer_number"`
	DoctorName         string `json:"doctor_name"`
	Email              string `json:"email"`
	PreviousCase       string `json:"previous_case"`
	PreviousCaseNumber string `json:"previous_case_number"`
	PatientName        string `json:"patient_name"`
	Dob                string `json:"dob"`
	HeightOfPatient    string `json:"height_of_patient"`
	Gender             string `json:"gender"`
	Race               string `json:"race"`
	PackageList        string `json:"package_list"`
	LateralXrayDate    string `json:"lateral_xray_date"`
	ConsultDate        string `json:"consult_date"`
	MissingTeeth       string `json:"missing_teeth"`
	AdenoidsRemoved    string `json:"adenoids_removed"`
	StatusCase         string `json:"status_case"`
	Comment            string `json:"comment"`

	LateralXrayImage   string `json:"lateral_xray_image"`   //optional
	FrontalXrayImage   string `json:"frontal_xray_image"`   //optional
	LowerArchImage     string `json:"lower_arch_image"`     //optional
	UpperArchImage     string `json:"upper_arch_image"`     //optional
	HandwristXrayImage string `json:"handwrist_xray_image"` //optional
	PanoramicXrayImage string `json:"panoramic_xray_image"` //optional
	AdditionalRecord_1 string `json:"additional_record_1"`  //optional
	AdditionalRecord_2 string `json:"additional_record_2"`  //optional
	AdditionalRecord_3 string `json:"additional_record_3"`  //optional
	AdditionalRecord_4 string `json:"additional_record_4"`  //optional
	AdditionalRecord_5 string `json:"additional_record_5"`  //optional

	TimeCreate string `json:"time_create"`
}
