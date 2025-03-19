package util

import (
	"bytes"
	"case-service/model"
	"case-service/rc"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"net/smtp"

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
	pathAttach := os.Getenv("PATH_ATTACH")
	attachmentPath := filepath.Join(pathAttach, "image1.jpg")

	// Baca file attachment
	fileData, err := ioutil.ReadFile(attachmentPath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read attachment")
		return err
	}

	// Encode file ke base64
	encodedFile := base64.StdEncoding.EncodeToString(fileData)
	filename := filepath.Base(attachmentPath)

	// Buat email message dengan attachment
	var emailBody bytes.Buffer
	writer := multipart.NewWriter(&emailBody)
	defer writer.Close() // Pastikan writer ditutup setelah digunakan

	boundary := writer.Boundary()
	headers := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: Email from FMBODS\r\nMIME-Version: 1.0\r\nContent-Type: multipart/mixed; boundary=%s\r\n\r\n", senderEmail, recipientEmail, boundary)
	emailBody.WriteString(headers)

	// Tambahkan teks email
	textPart := fmt.Sprintf("--%s\r\nContent-Type: text/plain; charset=\"utf-8\"\r\n\r\nHello, this is a test email with an attachment!\r\n", boundary)
	emailBody.WriteString(textPart)

	// Tambahkan attachment
	attachmentPart := fmt.Sprintf(
		"--%s\r\nContent-Type: application/octet-stream\r\nContent-Transfer-Encoding: base64\r\nContent-Disposition: attachment; filename=\"%s\"\r\n\r\n%s\r\n",
		boundary, filename, encodedFile)
	emailBody.WriteString(attachmentPart)

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
