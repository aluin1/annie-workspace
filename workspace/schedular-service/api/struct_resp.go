package api

// Definisikan struct sesuai dengan respons JSON
type Case struct {
	TimePointId         string `json:"TimePointId"`
	CaseNumber          string `json:"CaseNumber"`
	IsPrint             int    `json:"IsPrint"`
	SubmitDate          string `json:"SubmitDate"`
	CompleteDate        string `json:"CompleteDate"`
	IsExtract           int    `json:"IsExtract"`
	SubmitExtractDate   string `json:"SubmitExtractDate"`
	CompleteExtractDate string `json:"CompleteExtractDate"`
	NumTryPrint         int    `json:"NumTryPrint"`
	NumTryExtract       int    `json:"NumTryExtract"`
	MachinePrint        string `json:"MachinePrint"`
	MachineExtract      string `json:"MachineExtract"`
}
