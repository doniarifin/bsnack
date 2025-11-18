package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"bsnack/database/models"
)

type CustomerService struct {
	DB     *gorm.DB
	redis  *redis.Client
	prdSrv *ProductService
}

func NewCustomerService(db *gorm.DB, rdc *redis.Client, prd *ProductService) *CustomerService {
	return &CustomerService{
		DB:     db,
		redis:  rdc,
		prdSrv: prd,
	}
}
func (s *CustomerService) CreateCustomer(payload *models.Customer) error {
	payload.ID = uuid.New().String()
	payload.CreatedAt = time.Now()
	payload.UpdatedAt = time.Now()

	s.redis.Del(ctx, "customers:all")
	return s.DB.Create(payload).Error
}
func (s *CustomerService) UpdateCustomer(payload *models.Customer) error {
	payload.UpdatedAt = time.Now()

	s.redis.Del(ctx, "customers:all")
	s.redis.Del(ctx, "customer:%s", payload.ID)
	return s.DB.Save(payload).Error
}

func (s *CustomerService) GetAll() ([]models.Customer, error) {
	val, err := s.redis.Get(ctx, "customers:all").Result()
	if err == nil {
		var cached []models.Customer
		if err := json.Unmarshal([]byte(val), &cached); err == nil {
			return cached, nil
		}
	}

	var data []models.Customer
	err = s.DB.Order("created_at DESC").Find(&data).Error

	jsonData, _ := json.Marshal(data)
	s.redis.Set(ctx, "customers:all", jsonData, 30*time.Second)

	return data, err
}

func (s *CustomerService) GetByID(id string) (*models.Customer, error) {
	key := fmt.Sprintf("customer:%s", id)
	val, err := s.redis.Get(ctx, key).Result()
	if err == nil {
		var cached models.Customer
		if err := json.Unmarshal([]byte(val), &cached); err == nil {
			return &cached, nil
		}
	}

	var t models.Customer
	err = s.DB.Where("id = ?", id).First(&t).Error
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	jsonData, _ := json.Marshal(t)
	s.redis.Set(ctx, key, jsonData, 30*time.Second)

	return &t, nil
}

func (s *CustomerService) ExchangePoint(customerId string, point int) error {
	size, err := getSizeByPoint(point)
	if err != nil {
		return err
	}

	customer, err := s.GetByID(customerId)
	if err != nil {
		return err
	}

	if point > customer.Points {
		return fmt.Errorf("point not enough, current point is %d", customer.Points)
	}

	product, err := s.prdSrv.GetBySize(size)
	if err != nil {
		return fmt.Errorf("error get product, %w", err)
	}

	//update product
	product.Stock -= 1

	s.redis.Del(ctx, "product:%s", product.ID)
	err = s.prdSrv.UpdateProduct(product)
	if err != nil {
		return fmt.Errorf("error update product, %w", err)
	}

	//update customer
	customer.Points -= point
	s.redis.Del(ctx, "customer:%s", customer.ID)
	return s.UpdateCustomer(customer)
}

func getSizeByPoint(point int) (string, error) {
	switch point {
	case 200:
		return "small", nil
	case 300:
		return "medium", nil
	case 500:
		return "large", nil
	default:
		return "", fmt.Errorf("point not enough,  minimal point is 200")
	}
}
