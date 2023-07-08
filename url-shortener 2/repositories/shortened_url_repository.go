package repositories

import (
	"github.com/jinzhu/gorm"
	"url-shortener/domain"
)

type ShortenedURLRepository struct {
	DB *gorm.DB
}

func NewShortenedURLRepository(db *gorm.DB) *ShortenedURLRepository {
	return &ShortenedURLRepository{
		DB: db,
	}
}

func (r *ShortenedURLRepository) SaveURLData(urlData *domain.ShortenedURL) error {
	err := r.DB.Create(urlData).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *ShortenedURLRepository) GetAllURLData() ([]domain.ShortenedURL, error) {
	var urlDataList []domain.ShortenedURL
	err := r.DB.Find(&urlDataList).Error
	if err != nil {
		return nil, err
	}
	return urlDataList, nil
}

func (r *ShortenedURLRepository) CheckCustomAliasExists(alias string) (bool, error) {
	var urlData domain.ShortenedURL
	err := r.DB.Where("custom_alias = ?", alias).First(&urlData).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *ShortenedURLRepository) GetURLDataByAlias(alias string) (*domain.ShortenedURL, error) {
	var urlData domain.ShortenedURL
	if err := r.DB.Where("custom_alias = ?", alias).First(&urlData).Error; err != nil {
		return nil, err
	}
	return &urlData, nil
}

func (r *ShortenedURLRepository) UpdateURLData(urlData *domain.ShortenedURL) error {
	if err := r.DB.Save(urlData).Error; err != nil {
		return err
	}
	return nil
}
