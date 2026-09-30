package controllers

import (
	"config-service/ApiResponse"
	"config-service/config"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

// GetGames godoc
// @Summary      upload files
// @Description  upload allowed files
// @Tags         Upload
// @Produce      json
// @Success      200 {object} ApiResponse.Response
// @Failure      500 {object} ApiResponse.Response
// @Router       /uploadfile [post]
func UploadToMinio(c *fiber.Ctx) error {

	// 1. Get file
	file, _ := c.FormFile("file")

	// 2. Open file
	src, err := file.Open()
	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	defer src.Close()

	// 3. Custom filename
	customName := buildCustomFileName(file.Filename)
	objectName := os.Getenv("FOLDER_PATH") + "/" + customName

	// 4. Upload to MinIO
	_, err = config.Client.PutObject(
		config.Ctx,
		os.Getenv("MINIO_BUCKET"),
		objectName,
		src,
		file.Size,
		minio.PutObjectOptions{
			ContentType: file.Header.Get("Content-Type"),
		},
	)

	if err != nil {
		return ApiResponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	// 5. Public URL
	fileURL := fmt.Sprintf("%s/%s",
		os.Getenv("MINIO_BUCKET"),
		objectName,
	)

	// 6. Return URL only
	return c.JSON(fiber.Map{
		"file_url": fileURL,
	})
}

// Build custom filename
func UploadToMinioplayrule(c *fiber.Ctx) error {

	// khan testing test
	// Recover from panic
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[UploadToMiniogamerule] Panic recovered: %v", r)
		}
	}()

	// =========================================
	// Request Parameters
	// =========================================
	playID := strings.TrimSpace(c.FormValue("game_id"))
	playType := strings.TrimSpace(c.FormValue("game_type"))
	agent := strings.TrimSpace(c.FormValue("agent"))
	folderName := "kh-game"

	// Accept-Language header
	acceptLanguage := strings.ToLower(strings.TrimSpace(c.Get("Accept-Language")))
	if acceptLanguage == "" {
		acceptLanguage = "en"
	}

	// Validate required fields
	if playID == "" || playType == "" || agent == "" {
		return ApiResponse.Error(
			c,
			fiber.StatusBadRequest,
			"game_id, game_type and agent are required",
		)
	}

	// =========================================
	// Get Uploaded File
	// =========================================
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("[UploadToMiniogamerule] FormFile error: %v", err)
		return ApiResponse.Error(
			c,
			fiber.StatusBadRequest,
			err.Error(),
		)
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		log.Printf("[UploadToMiniogamerule] File open error: %v", err)
		return ApiResponse.Error(
			c,
			fiber.StatusInternalServerError,
			err.Error(),
		)
	}
	defer src.Close()

	// =========================================
	// Build File Name
	// =========================================
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == "" {
		ext = ".bin"
	}

	customFileName := fmt.Sprintf(
		"rule_%s_%s_%s_%s%s",
		playID,
		playType,
		strings.ToUpper(agent),
		acceptLanguage,
		ext,
	)

	objectName := fmt.Sprintf("%s/%s", folderName, customFileName)

	bucketName := os.Getenv("MINIO_BUCKET")
	if bucketName == "" {
		return ApiResponse.Error(
			c,
			fiber.StatusInternalServerError,
			"MINIO_BUCKET is not configured",
		)
	}

	log.Printf("[UploadToMiniogamerule] Target object: %s", objectName)

	// =========================================
	// Delete Existing File (same exact name)
	// =========================================
	_, err = config.Client.StatObject(
		config.Ctx,
		bucketName,
		objectName,
		minio.StatObjectOptions{},
	)

	if err == nil {
		// File exists, delete it
		log.Printf(
			"[UploadToMiniogamerule] Existing file found, deleting: %s",
			objectName,
		)

		err = config.Client.RemoveObject(
			config.Ctx,
			bucketName,
			objectName,
			minio.RemoveObjectOptions{},
		)
		if err != nil {
			log.Printf(
				"[UploadToMiniogamerule] RemoveObject error: %v",
				err,
			)

			return ApiResponse.Error(
				c,
				fiber.StatusInternalServerError,
				"failed to delete existing file",
			)
		}

		log.Printf(
			"[UploadToMiniogamerule] Old file deleted successfully",
		)
	} else {
		// File does not exist (or StatObject returned not found)
		log.Printf(
			"[UploadToMiniogamerule] No existing file found: %s",
			objectName,
		)
	}

	// =========================================
	// Upload New File
	// =========================================
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	uploadInfo, err := config.Client.PutObject(
		config.Ctx,
		bucketName,
		objectName,
		src,
		file.Size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		log.Printf(
			"[UploadToMiniogamerule] Upload error: %v",
			err,
		)

		return ApiResponse.Error(
			c,
			fiber.StatusInternalServerError,
			err.Error(),
		)
	}

	// =========================================
	// Build File URL
	// =========================================
	endpoint := strings.TrimRight(os.Getenv("MINIO_PUBLIC_ENDPOINT"), "/")
	fileURL := fmt.Sprintf(
		"%s/%s/%s",
		endpoint,
		bucketName,
		objectName,
	)

	log.Printf(
		"[UploadToMiniogamerule] Upload success | Object: %s | Size: %d",
		uploadInfo.Key,
		uploadInfo.Size,
	)

	// =========================================
	// Response
	// =========================================
	return c.JSON(fiber.Map{
		"message":         "File uploaded successfully",
		"file_url":        fileURL,
		"file_name":       customFileName,
		"object_name":     objectName,
		"game_id":         playID,
		"game_type":       playType,
		"agent":           strings.ToUpper(agent),
		"accept_language": acceptLanguage,
		"size":            uploadInfo.Size,
	})
}
func buildCustomFileName(original string) string {
	ext := filepath.Ext(original)

	// filename without extension
	name := strings.TrimSuffix(original, ext)

	// make lowercase
	name = strings.ToLower(name)

	// remove unsafe characters (keep a-z, 0-9, -)
	reg := regexp.MustCompile(`[^a-z0-9-]`)
	name = reg.ReplaceAllString(name, "")

	// limit to 10 characters
	if len(name) > 10 {
		name = name[:10]
	}

	// fallback if name becomes empty
	if name == "" {
		name = "file"
	}

	return name + "-" + uuid.New().String() + ext
}
