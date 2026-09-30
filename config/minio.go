package config

import (
	"context"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var Client *minio.Client
var Ctx = context.Background()

func InitMinio() {
	var err error
	var secure bool

	if os.Getenv("MINIO_SECURE") == "true" {
		secure = true
	} else {
		secure = false
	}

	Client, err = minio.New(os.Getenv("MINIO_URL"), &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("MINIO_PUBLIC_KEY"), os.Getenv("MINIO_SECRET_KEY"), ""),
		Secure: secure,
	})

	if err != nil {
		log.Fatal("MinIO connection failed:", err)
	}
}
