package api

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// Allowed file extensions
var allowedExtensions = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".pdf": true, ".docx": true,
}

// UploadFiles handles multiple file uploads
func UploadFiles(c echo.Context) error {
	uploadDir := os.Getenv("PATH_ATTACH")
	if uploadDir == "" {
		log.Error().Msg("Upload path is not set")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Upload path is not set"})
	}

	// Baca `customer_number` dari request
	customerNumber := c.FormValue("customer_number")
	if customerNumber == "" {
		log.Warn().Msg("Customer number is missing, using default directory")
		customerNumber = "default_dir"
	}

	// Buat sub-folder berdasarkan `customer_number`
	customerDir := filepath.Join(uploadDir, customerNumber)
	if err := os.MkdirAll(customerDir, os.ModePerm); err != nil {
		log.Error().Err(err).Msg("Failed to create customer directory")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create customer directory"})
	}

	form, err := c.MultipartForm()
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse form")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to parse form"})
	}

	files := form.File["file"]
	if len(files) == 0 {
		log.Warn().Msg("No files uploaded")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No files uploaded"})
	}

	var filePaths []string

	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file.Filename)) // Convert to lowercase for safety
		if ext == "" || !allowedExtensions[ext] {
			log.Warn().Str("filename", file.Filename).Msg("Invalid file type")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid file type"})
		}

		src, err := file.Open()
		if err != nil {
			log.Error().Err(err).Msg("Failed to open file")
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open file"})
		}

		// Generate unique file name
		newFileName := uuid.New().String() + "-" + time.Now().Format("20060102150405") + ext
		filePath := filepath.Join(customerDir, newFileName) // Perbaikan: Simpan file ke dalam direktori pelanggan
		filePath = filepath.Clean(filePath)                 // Ensure the path is safe

		// Create destination file
		dst, err := os.Create(filePath)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create file on server")
			src.Close()
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create file on server"})
		}

		// Copy file content
		if _, err = io.Copy(dst, src); err != nil {
			log.Error().Err(err).Msg("Failed to write file")
			src.Close()
			dst.Close()
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to write file"})
		}

		// Close file handles
		src.Close()
		dst.Close()

		// Return only the relative path
		filePaths = append(filePaths, uploadDir+"/"+customerNumber+"/"+newFileName)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":   "Files uploaded successfully",
		"filePaths": filePaths,
	})
}
