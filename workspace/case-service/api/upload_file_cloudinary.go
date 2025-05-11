package api

import (
	"case-service/rc"
	"case-service/util"
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"gopkg.in/asaskevich/govalidator.v8"
)

var allowedExtensionsCloud = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
}

func UploadFilesCloud(c echo.Context) error {
	tokenValue := c.Request().Header.Get("Authorization")
	validationParamStruct := util.ValidationParamStruct{Token: tokenValue}
	respCode := util.CheckValueParamGeneral(validationParamStruct)
	if !govalidator.IsNull(respCode) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": rc.ClientResponseTextCase[respCode]})
	}

	splitToken := strings.Split(tokenValue, "Bearer ")
	token := splitToken[1]
	secretID := os.Getenv("CLIENT_SECRET")
	_, _, errToken := ValidateJWT(token, secretID)
	if errToken != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": errToken.Error()})
	}

	customerNumber := c.FormValue("customer_number")
	if govalidator.IsNull(customerNumber) {
		customerNumber = "default"
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to parse form"})
	}

	files := form.File["file"]
	if len(files) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No files uploaded"})
	}

	folderCloud := os.Getenv("CLOUDINARY_CLOUD_FOLDER")
	// ✅ Init Cloudinary
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		log.Error().Err(err).Msg("Cloudinary config error")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Cloudinary configuration error"})
	}

	var fileURLs []string
	ctx := context.Background()

	for _, file := range files {
		// Log file information for debugging
		log.Info().Str("FileName", file.Filename).Str("FileSize", fmt.Sprintf("%d", file.Size)).Msg("Processing file")

		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !allowedExtensionsCloud[ext] {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid file type"})
		}

		src, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open file"})
		}
		defer src.Close()

		publicID := fmt.Sprintf("%s/%s-%s", folderCloud, uuid.New().String(), time.Now().Format("20060102150405"))

		// Log upload parameters
		log.Info().Str("PublicID", publicID).Str("Folder", folderCloud).Msg("Uploading to Cloudinary")

		uploadResult, err := cld.Upload.Upload(ctx, src, uploader.UploadParams{
			PublicID:     publicID,
			Folder:       folderCloud,
			ResourceType: "auto",
		})
		if err != nil {
			log.Error().Err(err).Str("PublicID", publicID).Msg("Upload to Cloudinary failed")
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Upload to Cloudinary failed"})
		}

		if uploadResult == nil {
			log.Error().Str("PublicID", publicID).Msg("No result returned from Cloudinary")
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "No result from Cloudinary upload"})
		}

		// Log successful upload
		log.Info().Interface("UploadResult", uploadResult).Msg("Upload result details")

		if !govalidator.IsNull(uploadResult.Error.Message) {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": uploadResult.Error.Message})
		}

		log.Info().Str("File URL", uploadResult.SecureURL).Msg("File uploaded successfully")

		fileURLs = append(fileURLs, uploadResult.SecureURL)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":   "Files uploaded successfully",
		"filePaths": fileURLs,
	})
}
