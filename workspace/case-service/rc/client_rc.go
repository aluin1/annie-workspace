package rc

// constants for Case response code
const (
	Success               = "00"
	Reject                = "01"
	ParamMandatoryNotsent = "02"
	InvalidToken          = "03"
)

// ClientResponseTextCase response code text
var ClientResponseTextCase = map[string]string{
	Success:               "Success",
	Reject:                "Reject",
	ParamMandatoryNotsent: "Mandatory Param Invalid",
	InvalidToken:          "Invalid Token",
}
