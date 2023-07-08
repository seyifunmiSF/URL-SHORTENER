package repositories

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"time"
	"url-shortener/domain"
)

type RedisRepository struct {
	RedisClient *redis.Client
}

func NewRedisRepository(redisClient *redis.Client) *RedisRepository {
	return &RedisRepository{
		RedisClient: redisClient,
	}
}

const (
	redisURLKey = "shortened_urls"
)

func (r *RedisRepository) SaveURLData(urlData *domain.ShortenedURL) error {
	jsonData, err := json.Marshal(urlData)
	if err != nil {
		return err
	}

	err = r.RedisClient.RPush(r.RedisClient.Context(), redisURLKey, jsonData).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisRepository) GetAllURLData() ([]domain.ShortenedURL, error) {
	urlDataList := make([]domain.ShortenedURL, 0)

	urlDataJSONList, err := r.RedisClient.LRange(r.RedisClient.Context(), redisURLKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	for _, jsonData := range urlDataJSONList {
		var urlData domain.ShortenedURL
		err = json.Unmarshal([]byte(jsonData), &urlData)
		if err != nil {
			return nil, err
		}
		urlDataList = append(urlDataList, urlData)
	}

	return urlDataList, nil
}

func (r *RedisRepository) GetURLDataByAlias(alias string) (*domain.ShortenedURL, error) {
	urlDataJSON, err := r.RedisClient.LRange(r.RedisClient.Context(), redisURLKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	for _, jsonData := range urlDataJSON {
		var urlData domain.ShortenedURL
		err := json.Unmarshal([]byte(jsonData), &urlData)
		if err != nil {
			return nil, err
		}

		if urlData.CustomAlias == alias {
			return &urlData, nil
		}
	}

	// URL data not found in Redis cache
	return nil, nil
}

func (r *RedisRepository) PerformRateLimiting(alias string, duration time.Duration, maxRequests int64) bool {

	key := "rate_limit:" + alias

	// Increment the counter and get the updated count
	count := r.RedisClient.Incr(r.RedisClient.Context(), key).Val()

	// Set the expiration for the counter if it's a new key
	if count == 1 {
		r.RedisClient.Expire(r.RedisClient.Context(), key, duration)
	}

	// Check if the count exceeds the maximum allowed requests
	return count > maxRequests
}
