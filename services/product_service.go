package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"bsnack/database/models"
)

type ProductService struct {
	DB    *gorm.DB
	redis *redis.Client
}

var ctx = context.Background()

func NewProductService(db *gorm.DB, rds *redis.Client) *ProductService {
	return &ProductService{
		DB:    db,
		redis: rds,
	}
}
func (s *ProductService) CreateProduct(payload *models.Product) error {
	payload.ID = uuid.New().String()
	now := time.Now()
	payload.CreatedAt = &now
	payload.UpdatedAt = time.Now()

	s.redis.Del(ctx, "products:all")
	s.redis.Del(ctx, "products:date")
	return s.DB.Create(payload).Error
}

func (s *ProductService) UpdateProduct(payload *models.Product) error {
	prd, err := s.GetByID(payload.ID)
	if err != nil {
		return err
	}

	if payload.CreatedAt == nil {
		payload.CreatedAt = prd.CreatedAt
	}
	payload.UpdatedAt = time.Now()

	s.redis.Del(ctx, "products:all")
	s.redis.Del(ctx, fmt.Sprintf("product:%s", payload.ID))
	return s.DB.Save(&payload).Error
}

func (s *ProductService) GetAll() ([]models.Product, error) {
	val, err := s.redis.Get(ctx, "products:all").Result()
	if err == nil {
		var cached []models.Product
		if err := json.Unmarshal([]byte(val), &cached); err == nil {
			return cached, nil
		}
	}

	var data []models.Product
	err = s.DB.Order("created_at DESC").Find(&data).Error

	jsonData, _ := json.Marshal(data)
	s.redis.Set(ctx, "products:all", jsonData, 30*time.Second)

	return data, err
}

func (s *ProductService) GetByID(id string) (*models.Product, error) {
	key := fmt.Sprintf("product:%s", id)
	val, err := s.redis.Get(ctx, key).Result()
	if err == nil {
		var cached models.Product
		if err := json.Unmarshal([]byte(val), &cached); err == nil {
			return &cached, nil
		}
	}

	var t models.Product
	err = s.DB.Where("id = ?", id).First(&t).Error
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	jsonData, _ := json.Marshal(t)
	s.redis.Set(ctx, key, jsonData, 30*time.Second)

	return &t, nil
}

func (s *ProductService) CheckStock(id string) (*models.Product, error) {
	var t models.Product
	err := s.DB.Where("id = ? AND stock > ?", id, 0).First(&t).Error
	if err != nil {
		return nil, fmt.Errorf("product not found or out of stock: %w", err)
	}
	return &t, nil
}

func (s *ProductService) GetByProductDate(date string) ([]models.Product, error) {
	layout := "2006-01-02T15:04:05Z07:00"

	t, err := time.Parse(layout, date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	end := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 0, 0, time.Local)

	val, err := s.redis.Get(ctx, "products:date").Result()
	if err == nil {
		var cached []models.Product
		if err := json.Unmarshal([]byte(val), &cached); err == nil {
			return cached, nil
		}
	}

	var data []models.Product

	err = s.DB.Where("created_at >= ? AND created_at <= ?", start, end).Find(&data).Error
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	jsonData, _ := json.Marshal(data)
	s.redis.Set(ctx, "products:date", jsonData, 30*time.Second)
	return data, nil
}

func (s *ProductService) GetBySize(size string) (*models.Product, error) {
	var t models.Product
	err := s.DB.Where("size = ? AND stock > ?", size, 0).First(&t).Error
	if err != nil {
		return nil, fmt.Errorf("product not found or out of stock: %w", err)
	}
	return &t, nil
}
