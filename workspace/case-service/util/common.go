package util

import (
	"bytes"
	"case-service/model"
	"case-service/rc"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"net/http"
	"net/smtp"

	"github.com/golang-jwt/jwt/v4"

	"github.com/rs/zerolog/log"
)

type ValidationParamStruct struct {
	GrantType    string
	ClientID     string
	ClientSecret string
	Scope        string
	Token        string
}

// CheckValueParamToken
func CheckValueParamToken(validationProStruct ValidationParamStruct) (respCode string) {

	if validationProStruct.GrantType != "client_credentials" {
		respCode := rc.Reject
		return respCode
	}

	clientIdService := os.Getenv("CLIENT_ID")
	if validationProStruct.ClientID != clientIdService {
		respCode := rc.Reject
		return respCode
	}

	clientSecretService := os.Getenv("CLIENT_SECRET")
	if validationProStruct.ClientSecret != clientSecretService {
		respCode := rc.Reject
		return respCode
	}

	return ""

}

// CheckValueParamGeneral
func CheckValueParamGeneral(validationProStruct ValidationParamStruct) (respCode string) {

	if len(validationProStruct.Token) < 2 {
		respCode := rc.InvalidToken
		return respCode
	}

	return ""

}
func SendEmail(dataCase *model.DataCase) error {
	// Konfigurasi SMTP
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	senderEmail := os.Getenv("SENDER_EMAIL")
	appPassword := os.Getenv("SMTP_PASSWORD")

	if !dataCase.Email.Valid {
		log.Error().Msg("Recipient email is empty or invalid")
		return fmt.Errorf("recipient email is empty or invalid")
	}
	recipientEmail := dataCase.Email.String

	// Path file attachment
	// pathAttach := os.Getenv("PATH_ATTACH")
	// attachmentPath := filepath.Join(pathAttach, "image1.jpg")

	// Baca file attachment
	// fileData, err := ioReadFile(attachmentPath)
	// if err != nil {
	// 	log.Error().Err(err).Msg("Failed to read attachment")
	// 	return err
	// }

	// Encode file ke base64
	// encodedFile := base64.StdEncoding.EncodeToString(fileData)
	// filename := filepath.Base(attachmentPath)

	// Buat email message dengan attachment
	var emailBody bytes.Buffer
	writer := multipart.NewWriter(&emailBody)
	defer writer.Close() // Pastikan writer ditutup setelah digunakan

	boundary := writer.Boundary()
	subjectText := fmt.Sprintf("Ticket Number %s Annie VIP Order Form %s", strconv.Itoa(int(dataCase.CaseID.Int64)), dataCase.PatientName.String)
	headers := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: multipart/mixed; boundary=%s\r\n\r\n", senderEmail, recipientEmail, subjectText, boundary)
	emailBody.WriteString(headers)

	// Tambahkan teks email

	htmlContent := HtmlContent(dataCase)

	// textPart := fmt.Sprintf(
	// 	"--%s\r\nContent-Type: text/plain; charset=\"utf-8\"\r\n\r\nHello, this is a test email with an attachment!\r\n",
	// 	boundary,
	// )

	htmlPart := fmt.Sprintf(
		"--%s\r\nContent-Type: text/html; charset=\"utf-8\"\r\n\r\n%s\r\n",
		boundary,
		htmlContent,
	)

	// Gabungkan bagian teks dan HTML
	// textPartFull := textPart + htmlPart

	// textPart := fmt.Sprintf("--%s\r\nContent-Type: text/plain; charset=\"utf-8\"\r\n\r\nHello, this is a test email with an attachment!\r\n", boundary)
	emailBody.WriteString(htmlPart)

	// Tambahkan attachment
	// attachmentPart := fmt.Sprintf(
	// 	"--%s\r\nContent-Type: application/octet-stream\r\nContent-Transfer-Encoding: base64\r\nContent-Disposition: attachment; filename=\"%s\"\r\n\r\n%s\r\n",
	// 	boundary, filename, encodedFile)
	// emailBody.WriteString(attachmentPart)

	// Akhiri email dengan boundary akhir
	emailBody.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	// Autentikasi SMTP
	auth := smtp.PlainAuth("", senderEmail, appPassword, smtpHost)

	// Kirim email
	smtpAddr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	errSend := smtp.SendMail(smtpAddr, auth, senderEmail, []string{recipientEmail}, emailBody.Bytes())
	if errSend != nil {
		log.Error().Err(errSend).Msg("Failed to send email")
		return errSend
	}

	log.Info().Msg("Email sent successfully!")
	return nil
}

func ConvertBase64(base64String, directoryPath, fileName string) (string, error) {
	if base64String == "" {
		return "", fmt.Errorf("base64String is empty")
	}

	baseDir := os.Getenv("PATH_ATTACH")
	if baseDir == "" {
		baseDir = "/var/attach"
	}

	maxSizeStr := os.Getenv("MAX_SIZE_FILE")
	maxSize, err := strconv.Atoi(maxSizeStr)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error converting MAX_SIZE_FILE: %s", err))
		return "", fmt.Errorf("error converting MAX_SIZE_FILE: %w", err)
	}

	maxSizeBytes := maxSize * 1024 * 1024
	directoryFull := filepath.Clean(filepath.Join(baseDir, directoryPath))

	re := regexp.MustCompile(`^data:(.*?);base64,`)
	matches := re.FindStringSubmatch(base64String)
	if len(matches) < 2 {
		return "", fmt.Errorf("invalid Base64 format")
	}

	mimeType := matches[1]
	exts, err := mime.ExtensionsByType(mimeType)
	if err != nil || len(exts) == 0 {
		log.Warn().Msg(fmt.Sprintf("Unknown MIME type: %s, using default .bin", mimeType))
		exts = []string{".bin"}
	}
	ext := exts[0]

	base64Data := re.ReplaceAllString(base64String, "")
	fileData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("error decoding Base64: %v", err)
	}

	if len(fileData) > maxSizeBytes {
		return "", fmt.Errorf("file %s size exceeds the maximum allowed size (%d bytes)", fileName, maxSizeBytes)
	}

	absolutePath, err := filepath.Abs(directoryFull)
	if err != nil {
		return "", fmt.Errorf("error getting absolute path: %v", err)
	}

	if err := os.MkdirAll(absolutePath, os.ModePerm); err != nil {
		return "", fmt.Errorf("error creating directory: %v", err)
	}

	fileName = strings.Map(func(r rune) rune {
		if strings.ContainsRune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-", r) {
			return r
		}
		return '_'
	}, fileName)

	outputPath := filepath.Join(absolutePath, fileName+ext)
	if err := os.WriteFile(outputPath, fileData, 0644); err != nil {
		return "", fmt.Errorf("error writing file: %v", err)
	}

	fmt.Println("✅ File berhasil disimpan di:", outputPath)
	return outputPath, nil
}

func HtmlContent(dataCase *model.DataCase) string {
	contentBody := `
	 
<html >
  <body>


    <h4 style="color:#ffff;padding: 10px; border-radius: 10px;font-size: 15px;  background: linear-gradient(to right, #099cd5, #191d7d)!important;">
    <b>   Information Ticket</b></h4>
    <table style="width: 100%; margin:2em;" >
    <tr style="font-size: 15px; ">
    <td style="width: 50%; "><b>Ticket Number:</b></td>
    <td><b>Customer Number:</b></td> 
    </tr>
    <tr style="font-size: 15px; "> 
    <td><span>` + strconv.Itoa(int(dataCase.CaseID.Int64)) + `</span></td>
    <td><span>` + dataCase.CustomerNumber.String + `</span></td>
    </tr>
    
    <tr style="font-size: 15px; "> 
    <td><b>Doctor Name:</b></td>
    <td><b>Email:</b></td> 
    </tr>
    <tr style="font-size: 15px; "> 
    <td><span>` + dataCase.DoctorName.String + `</span></td> 
    <td><span>` + dataCase.Email.String + `</span></td>
    </tr>

    <tr style="font-size: 15px; "> 
    <td><b>Previous Case:</b></td> 
    <td><b>Previous Case Number:</b></td> 
    </tr>
    <tr style="font-size: 15px; "> 
      <td><span>` + dataCase.PreviousCase.String + `</span></td>
      <td><span>` + dataCase.PreviousCaseNumber.String + `</span></td>
    </tr>

      </table>
        <hr>
        <h4 style="color:#ffff;padding: 10px; border-radius: 10px;font-size: 15px;  background: linear-gradient(to right, #099cd5, #191d7d)!important;">
          <b>Information Patient</b></h4>
          <table style="width: 100%; margin:2em;" >
          <tr style="font-size: 15px; ">
          <td style="width: 50%; "><b>Patient Name:</b></td>
          <td><b>DOB:</b></td> 
          </tr>
          <tr style="font-size: 15px; "> 
          <td><span>` + dataCase.PatientName.String + `</span></td>
          <td><span>` + dataCase.Dob.String + `</span></td>
          </tr>
          
          <tr style="font-size: 15px; "> 
          <td><b>Height of Patient:</b></td>
          <td><b>Gender:</b></td> 
          </tr>
          <tr style="font-size: 15px; "> 
          <td><span>` + dataCase.HeightOfPatient.String + `</span></td> 
          <td><span>` + dataCase.Gender.String + `</span></td>
          </tr>
      
          <tr style="font-size: 15px; "> 
          <td><b>Race:</b></td> 
          <td><b>Package List::</b></td> 
          </tr>
          <tr style="font-size: 15px; "> 
            <td><span>` + dataCase.Race.String + `</span></td>
            <td><span>` + dataCase.PackageList.String + `</span></td>
          </tr>
      
            </table> 
              <hr>  
              <h4 style="color:#ffff;padding: 10px; border-radius: 10px;font-size: 15px;  background: linear-gradient(to right, #099cd5, #191d7d)!important;">
                <b>Information Consultation</b></h4>
                <table style="width: 100%; margin:2em;" >
                <tr style="font-size: 15px; ">
                <td style="width: 50%; "><b> Lateral X-ray Date:</b></td>
                <td><b>	Consult Date:</b></td> 
                </tr>
                <tr style="font-size: 15px; "> 
                <td><span>` + dataCase.LateralXrayDate.String + `</span></td>
                <td><span>` + dataCase.ConsultDate.String + `</span></td>
                </tr>
                
                <tr style="font-size: 15px; "> 
                <td><b>Missing Teeth:</b></td>
                <td><b>	Adenoids Removed:</b></td> 
                </tr>
                <tr style="font-size: 15px; "> 
                <td><span>` + dataCase.MissingTeeth.String + `</span></td> 
                <td><span>` + dataCase.AdenoidsRemoved.String + `</span></td>
                </tr>
            
                <tr style="font-size: 15px; "> 
                <td><b>Comment:</b></td> 
                <td> </td> 
                </tr>
                <tr style="font-size: 15px; "> 
                  <td><span>` + dataCase.Comment.String + `</span></td>
                  <td><span> </span></td>
                </tr>
            
                  </table>
                    <hr>  
                          <h4 style="color:#ffff;padding: 10px; border-radius: 10px;font-size: 15px;  background: linear-gradient(to right, #099cd5, #191d7d)!important;">
                            <b>Information Consultation Detail</b></h4>
                            <table style="width: 100%; margin:2em;" >
                            <tr style="font-size: 15px; ">
                            <td style="width: 50%; "><b> Lateral X-Ray Image:</b></td>
                            <td><b>	Frontal X-Ray Image:</b></td> 
                            </tr>
                            <tr style="font-size: 15px; "> 
                            <td><span>` + dataCase.LateralXrayImage.String + `</span></td>
                            <td><span>` + dataCase.FrontalXrayImage.String + `</span></td>
                            </tr>
                            
                            <tr style="font-size: 15px; "> 
                            <td><b>Lower Arch Image:</b></td>
                            <td><b>	Upper Arch Image:</b></td> 
                            </tr>
                            <tr style="font-size: 15px; "> 
                            <td><span>` + dataCase.LowerArchImage.String + `</span></td> 
                            <td><span>` + dataCase.UpperArchImage.String + `</span></td>
                            </tr>

                            <tr style="font-size: 15px; "> 
                            <td><b>HandWrist X-Ray Image:</b></td>
                            <td><b>Panoramic X-Ray (Panorex) Image:</b></td> 
                            </tr>
                            <tr style="font-size: 15px; "> 
                            <td><span>` + dataCase.HandwristXrayImage.String + `</span></td> 
                            <td><span>` + dataCase.PanoramicXrayImage.String + `</span></td>
                            </tr>
                         
                            <tr style="font-size: 15px; "> 
                              <td><b>Additional Record 1:</b></td>
                              <td><b>Additional Record 2:</b></td> 
                              </tr>
                              <tr style="font-size: 15px; "> 
                              <td><span>` + dataCase.AdditionalRecord_1.String + `</span></td> 
                              <td><span>` + dataCase.AdditionalRecord_2.String + `</span></td>
                              </tr>
                         
                              <tr style="font-size: 15px; "> 
                                <td><b>Additional Record 3:</b></td>
                                <td><b>Additional Record 4:</b></td> 
                                </tr>
                                <tr style="font-size: 15px; "> 
                                <td><span>` + dataCase.AdditionalRecord_3.String + `</span></td> 
                                <td><span>` + dataCase.AdditionalRecord_4.String + `</span></td>
                                </tr>
                        
                                <tr style="font-size: 15px; "> 
                                  <td><b>Additional Record 5:</b></td>
                                  <td><b> </b></td> 
                                  </tr>
                                  <tr style="font-size: 15px; "> 
                                  <td><span>` + dataCase.AdditionalRecord_5.String + `</span></td> 
                                  <td><span> </span></td>
                                  </tr>
                              </table>
                                <hr>
								</body>
</html>

`
	return contentBody
}

func ValidationTokenAuth(jwtToken string) (GoogleJWTClaims, error) {

	var googleJWTClaims GoogleJWTClaims
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return googleJWTClaims, fmt.Errorf("JWT does not contain a valid 'kid' field")
		}
		jwk, err := GetPublicKey(kid)
		if err != nil {
			return googleJWTClaims, err
		}
		return ParseRSAPublicKey(jwk)
	}, jwt.WithoutClaimsValidation()) // Tambahkan ini untuk menonaktifkan validasi bawaan

	if err != nil {
		log.Info().Msg(fmt.Sprintf("❌ JWT tidak valid: %v", err))
		return googleJWTClaims, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		log.Info().Msg("✅ JWT valid!")
		payloadJson, _ := json.MarshalIndent(claims, "", "  ")
		log.Info().Msg(fmt.Sprintf("✅ Payload: %s", string(payloadJson)))

		if err := json.Unmarshal(payloadJson, &googleJWTClaims); err != nil {
			log.Info().Msg(fmt.Sprintf("Error Unmarshal: %s", err.Error()))
		}

		log.Info().Msg(fmt.Sprintf("✅ Email: %v", googleJWTClaims.Email))
		return googleJWTClaims, err
	}

	log.Info().Msg("❌ JWT tidak valid!")
	return googleJWTClaims, err

}

// Convert base64 URL to big.Int
func Base64ToBigInt(encoded string) (*big.Int, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	bi := new(big.Int).SetBytes(decoded)
	return bi, nil
}

// Convert JSON Public Key to *rsa.PublicKey
func ParseRSAPublicKey(jwk JWK) (*rsa.PublicKey, error) {
	n, err := Base64ToBigInt(jwk.N)
	if err != nil {
		return nil, err
	}
	e, err := Base64ToBigInt(jwk.E)
	if err != nil {
		return nil, err
	}

	return &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}, nil
}

// Fetch public key from Google JWKS endpoint
func GetPublicKey(kid string) (JWK, error) {
	urlGetPublicKey := os.Getenv("URL_GET_PUBLIC_KEY")
	resp, err := http.Get(urlGetPublicKey)
	if err != nil {
		return JWK{}, fmt.Errorf("error fetching public keys: %w", err)
	}
	defer resp.Body.Close()

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return JWK{}, fmt.Errorf("error decoding JWKS: %w", err)
	}

	for _, key := range jwks.Keys {
		if key.Kid == kid {
			log.Info().Msg(fmt.Sprintf("✅ Using Key ID: %s", key.Kid))
			return key, nil
		}
	}

	return JWK{}, fmt.Errorf("no matching key found for kid: %s", kid)
}
