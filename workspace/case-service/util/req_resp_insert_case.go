package util

// InsertCaseRequest xml data for request to case
type InsertCaseRequest struct {
	CustomerNumber     string `json:"customer_number" valid:"required"`
	DoctorName         string `json:"doctor_name" valid:"required"`
	Email              string `json:"email" valid:"required"`
	PreviousCase       string `json:"previous_case"`        //optional
	PreviousCaseNumber string `json:"previous_case_number"` //optional
	PatientName        string `json:"patient_name" valid:"required"`
	Dob                string `json:"dob" valid:"required"`
	HeightOfPatient    string `json:"height_of_patient"` //optional
	Gender             string `json:"gender" valid:"required"`
	Race               string `json:"race" valid:"required"`
	PackageList        string `json:"package_list" valid:"required"`
	LateralXrayDate    string `json:"lateral_xray_date" valid:"required"`
	ConsultDate        string `json:"consult_date" valid:"required"`
	MissingTeeth       string `json:"missing_teeth"` //optional
	AdenoidsRemoved    string `json:"adenoids_removed" valid:"required"`
	Comment            string `json:"comment"` //optional

	LateralXrayImage   string `json:"lateral_xray_image"`   //optional
	FrontalXrayImage   string `json:"frontal_xray_image"`   //optional
	LowerArchImage     string `json:"lower_arch_image"`     //optional
	UpperArchImage     string `json:"upper_arch_image"`     //optional
	HandwristXrayImage string `json:"handwrist_xray_image"` //optional
	PanoramicXrayImage string `json:"panoramic_xray_image"` //optional
	AdditionalRecord1  string `json:"additional_record_1"`  //optional
	AdditionalRecord2  string `json:"additional_record_2"`  //optional
	AdditionalRecord3  string `json:"additional_record_3"`  //optional
	AdditionalRecord4  string `json:"additional_record_4"`  //optional
	AdditionalRecord5  string `json:"additional_record_5"`  //optional
}

// InsertCaseResponse response from case
type InsertCaseResponse struct {
	CustomerNumber  string `json:"customer_number"`
	ResponseCode    string `json:"response_code,omitempty"`
	ResponseMessage string `json:"response_message,omitempty"`
}
