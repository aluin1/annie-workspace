package api

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/asaskevich/govalidator.v8"
	"gopkg.in/resty.v1"
)

var CaseNumber string

func init() {

	log.Infof("Configuring get case number...")
	reqTimeout, _ := strconv.Atoi(os.Getenv("CLIENT_TIMEOUT_SECOND"))
	caseRest = resty.New()
	caseRest.SetDebug(true)
	caseRest.SetTimeout(time.Duration(reqTimeout) * time.Second)
	caseRest.SetRetryCount(5)
	caseRest.SetRetryWaitTime(5 * time.Second)

	caseRest.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

}

// GetCase
func GetCase() (status bool) {

	resp, _ := caseRest.R().Post(os.Getenv("CASE_NUMBER_URL"))
	var cases []Case
	err := json.Unmarshal(resp.Body(), &cases)
	if err != nil {
		log.Fatalf("Error unmarshaling response: %v", err)
	}

	for _, c := range cases {
		CaseNumber = c.CaseNumber
	}

	log.Infof("value CaseNumber: %s", CaseNumber)
	if resp.IsSuccess() {
		// Get Response
		status = true
		// fmt.Println(runtime.GOOS)
		if !govalidator.IsNull(CaseNumber) {
			workDir := os.Getenv("DIR_APP") + "/01-DS_RUN_Log/"
			// cmdFile := workDir + "Extract-Report.CMD"

			// cmd := exec.Command(cmdFile, CaseNumber)
			// // Menyambungkan output dan error ke konsol
			// cmd.Stdout = os.Stdout
			// cmd.Stderr = os.Stderr

			// // Menjalankan perintah
			// err := cmd.Run()
			// if err != nil {
			// 	fmt.Printf("Error executing %s: %v\n", cmdFile, err)
			// 	return
			// }

			// fmt.Println("Command executed successfully!")

			// cmd := exec.Command("bash", "/app/Extract-Report.sh", CaseNumber)
			cmd := exec.Command("/bin/bash", "/app/Extract-Report.sh", CaseNumber)

			cmd.Dir = workDir

			output, err := cmd.CombinedOutput()
			if err != nil {
				log.Printf("Error: %v\n", err)
				return
			}

			log.Printf("message info: %s", string(output))
		}
		return
	} else {
		// response is ERROR
		status = false
		log.Error("Unable to get case number")
	}

	LogResty(resp.Request.Header, resp, "Log Execute Report: ")
	return status
}

func LogResty(headers http.Header, restResponse *resty.Response, prefix string) {
	jsonHeaderRequest, _ := json.Marshal(headers)
	jsonBodyRequest, _ := json.Marshal(restResponse.Request.Body)
	jsonHeaderResponse, _ := json.Marshal(restResponse.Header())

	log.Infof("%s Request Url: %s", prefix, restResponse.Request.URL)
	log.Infof("%s Request Header: %s", prefix, string(jsonHeaderRequest))
	log.Infof("%s Request Body: %s", prefix, string(jsonBodyRequest))
	log.Infof("%s Response Status: %s", prefix, restResponse.Status())
	log.Infof("%s Response Received At: %s", prefix, GetTimestampDate(restResponse.ReceivedAt()))
	log.Infof("%s Response Response Time: %s", prefix, restResponse.Time())
	log.Infof("%s Response Header: %s", prefix, string(jsonHeaderResponse))
	log.Infof("%s Response Body: %s", prefix, restResponse.String())
}

// GetTimestamp get current timestamp in ISO 8601 format
func GetTimestampDate(date time.Time) string {
	return fmt.Sprint(date.UTC().Format("2006-01-02 15:04:05.000Z"))
}
