package services

import (
	"fmt"
	"github.com/skip2/go-qrcode"
	"os"
	"url-shortener/domain"
	"url-shortener/utils"
)

type QRCodeService struct {
	config *domain.Config
}

func NewQRCodeService(config *domain.Config) *QRCodeService {
	return &QRCodeService{
		config: config,
	}
}

func (s *QRCodeService) GenerateQRCode(shortURL string) (string, error) {
	// Generate a unique filename for the QR code image
	filename, err := utils.GenerateUniqueID()
	if err != nil {
		return "", err
	}

	filename = filename + ".png"

	// Check if the directory exists and create it if necessary
	err = os.MkdirAll(s.config.QRCodeDirectory, os.ModePerm)
	if err != nil {
		return "", err
	}

	// Generate the QR code
	err = qrcode.WriteFile(shortURL, qrcode.Medium, 256, fmt.Sprintf("%s/%s", s.config.QRCodeDirectory, filename))
	if err != nil {
		return "", err
	}

	// Return the URL of the saved QR code image
	return fmt.Sprintf("%s/%s", s.config.BaseURL+"/"+s.config.UrlImagePath, filename), nil
}
