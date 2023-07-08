package domain

type QRCodeService interface {
	GenerateQRCode(shortURL string) (string, error)
}
